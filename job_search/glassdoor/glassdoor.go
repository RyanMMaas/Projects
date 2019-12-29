package glassdoor

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/RyanMMaas/projects/job_search/jobdata"
)

// Stores a list of URLs of jobs
type jobsList struct {
	jobs []string
}

type locData struct {
	CompoundID   string `json:"compundId"`
	CountryName  string `json:"countryName"`
	ID           string `json:"id"`
	Label        string `json:"label"`
	LocationID   int    `json:"locationId"`
	LocationType string `json:"locationType"`
	LongName     string `json:"longName"`
	RealID       int    `json:"realId"`
}

// GetJobs searches for links on glassdoor
// based on the keyword (k), location (l), and age (a)
// the user inputs. It then writes the data to a json file.
func GetJobs(k, l, a string) {
	jsonJobData, templateList, existingURLs := jobdata.GetOld()

	var urlQuery string

	// Leave these for now
	suggestCount := "0"
	suggestChosen := "false"
	clickSource := "searchBtn"

	// Prepares keyword to be used in the get request
	// then does the same for location
	typedKeyword := k
	typedKeyword = url.QueryEscape(strings.TrimSpace(typedKeyword))
	sckeyword := url.QueryEscape(getscKeyword(typedKeyword))

	location := l
	location = url.QueryEscape(location)
	locT, locID := getLocationID(location)
	jobType := ""

	// How old the jobs you want to search for are.
	fromAge := a
	fromAge = strings.TrimSuffix(fromAge, "\r\n")

	// builds the url for the get request
	urlQuery += "suggestCount=" + suggestCount + "&"
	urlQuery += "suggestChosen=" + suggestChosen + "&"
	urlQuery += "clickSource=" + clickSource + "&"
	urlQuery += "typedKeyword=" + typedKeyword + "&"
	urlQuery += "sc.keyword=" + sckeyword + "&"
	urlQuery += "locT=" + locT + "&"
	urlQuery += "locId=" + locID + "&"
	urlQuery += "jobType=" + jobType + "&"
	urlQuery += "fromAge=" + fromAge
	u := &url.URL{
		Scheme:   "https",
		Host:     "glassdoor.com",
		Path:     "Job/jobs.htm",
		RawQuery: urlQuery,
	}

	// Send get request to url
	response, err := http.Get(u.String())
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	// Gets the body of the response and turns it into a string
	buf := new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	bodyStr := buf.String()

	// URL used to go to different pages
	untranslatedURLStart := strings.Index(bodyStr, "untranslatedUrl")
	untranslatedURLStart += len("untranslatedUrl") + 5
	untranslatedURLEnd := strings.Index(bodyStr[untranslatedURLStart:], "'")
	untranslatedURL := bodyStr[untranslatedURLStart : untranslatedURLStart+untranslatedURLEnd]
	untranslatedURL += "?fromAge=" + fromAge

	var jobList jobsList
	getJobURLs(untranslatedURL, &jobList)
	getJobRedirects(&jobList)

	parseJobData(&jobList, &templateList, existingURLs)
	jobdata.WriteNew(&templateList, jsonJobData)
}

// Used to get the location id
func getLocationID(l string) (string, string) {
	locAddr := "https://www.glassdoor.com/findPopularLocationAjax.htm?term=" + l + "&maxLocationsToReturn=10"
	locResp, err := http.Get(locAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer locResp.Body.Close()
	body, err := ioutil.ReadAll(locResp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var ld []locData
	err = json.Unmarshal(body, &ld)
	if err != nil {
		log.Fatal(err)
	}
	return ld[0].LocationType + "&", strconv.Itoa(ld[0].RealID) + "&"
}

type keywordData struct {
	Suggestion string  `json:"suggestion"`
	Source     string  `json:"source"`
	Confidence float32 `json:"confidence"`
	Version    string  `json:"version"`
	Category   string  `json:"category"`
	EmployerID string  `json:"employerId"`
}

// Used to get the keyword id
func getscKeyword(k string) string {
	sckeyword := "https://www.glassdoor.com/searchsuggest/typeahead?source=Jobs&version=New&input=" + k + "&rf=full"
	kwResp, err := http.Get(sckeyword)
	if err != nil {
		log.Fatal(err)
	}
	defer kwResp.Body.Close()
	body, err := ioutil.ReadAll(kwResp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var kwd []keywordData
	err = json.Unmarshal(body, &kwd)
	if err != nil {
		log.Fatal(err)
	}

	return kwd[0].Suggestion
}

func getJobURLs(uurl string, jl *jobsList) error {
	i := 1
	nURL1 := uurl[:len(uurl)-14]
	nURL2 := uurl[len(uurl)-14:]
	for {
		// Builds the url with _IP{} where
		// the number after IP is the page
		uurl = nURL1 + "_IP" + strconv.Itoa(i) + nURL2
		res, err := http.Get(uurl)
		if err != nil {
			log.Fatal(err)
		}
		if res.StatusCode == 404 {
			break
		}
		defer res.Body.Close()

		data, _ := ioutil.ReadAll(res.Body)
		bodyStr := string(data)

		noJobs := "<h4 class='bold'>Your search for"
		if strings.Index(bodyStr, noJobs) != -1 {
			break
		}
		noJobs = "<h4>Your filtered search does not"
		if strings.Index(bodyStr, noJobs) != -1 {
			break
		}
		for {
			// Searches for all the links on one age
			searchStr := "<div><a href='/partner/jobListing.htm"
			endStr := "'"
			str, err := srchMovBody(searchStr, endStr, -23, 0, &bodyStr)
			if err != nil {
				break
			}

			jl.jobs = append(jl.jobs, str)
		}
		i++
	}
	return nil
}

// Loop through all URLs in jobslist and get the redirect URl
// Then empty the jobs and refill the struct with the correct URLs
func getJobRedirects(jl *jobsList) {
	temp := jobsList{}
	for _, url := range jl.jobs {
		res, err := http.Get("https://www.glassdoor.com" + url)
		if err != nil {
			log.Fatal(err)
		}
		temp.jobs = append(temp.jobs, res.Request.URL.String())
	}
	jl.jobs = nil
	for _, url := range temp.jobs {
		jl.jobs = append(jl.jobs, url)
	}
}

func parseJobData(jl *jobsList, tl *jobdata.Templates, existingURLs map[string]int) {
	for _, job := range jl.jobs {
		res, err := http.Get(job)
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()

		for {
			// if the URL already exists in the file
			// adding it is skipped to prevent duplicates
			_, ok := existingURLs[strings.Split(job, "&")[0]]
			if ok {
				break
			}
			data, _ := ioutil.ReadAll(res.Body)
			bodyStr := string(data)

			// Searches for the title
			titleStr := "jobViewJobTitleWrap'>"
			endTitleStr := "</h2>"
			title, err := srchMovBody(titleStr, endTitleStr, 39, 0, &bodyStr)
			if err != nil {
				break
			}

			// Searches for the company
			compStr := "class='strong ib'>"
			endCompStr := "</span>"
			company, err := srchMovBody(compStr, endCompStr, 0, 0, &bodyStr)
			if err != nil {
				break
			}

			// Searches for the apply link
			applyStr := "class=\"regToApplyArrowBoxContainer\">"
			endApplyStr := "'"
			apply, err := srchMovBody(applyStr, endApplyStr, 9, 0, &bodyStr)
			if err != nil {
				break
			}

			// Searches for the job description
			descStr := "class='jobDescriptionContent desc module pad noMargBot'>"
			endDescStr := "</section>"
			description, err := srchMovBody(descStr, endDescStr, 0, 0, &bodyStr)
			if err != nil {
				break
			}

			// Gets the time so we can use the
			// date as the date we added each URL
			currentTime := time.Now()
			t := jobdata.TemplateData{
				Key:            len(tl.TemplateData) + 1,
				Date:           currentTime.Format("02-01-2006"),
				URL:            strings.Split(job, "&")[0],
				JobTitle:       title,
				CompanyName:    company,
				JobDescription: description,
				ApplyLink:      apply,
				AppliedStatus:  2}

			tl.TemplateData = append(tl.TemplateData, t)
		}
	}
}

// Searches for the start string and end string and returns whats between them
// from the end of start to the beginning of end. Then it moves the slice
// so it holds less and searching will be quicker.
// Works best on well structured pages
// so and eo are start offset and end offset. They are used in cases where finding
// a string requires searching for part of it in order to be accurate.
func srchMovBody(start, end string, so, eo int, body *string) (string, error) {
	strStart := strings.Index(*body, start)
	if strStart == -1 {
		return "", errors.New("Error retrieving str")
	}
	beginStr := (*body)[strStart+len(start)+so:]
	strEnd := strings.Index(beginStr, end)
	if strEnd == -1 {
		return "", errors.New("Error retrieving str")
	}
	str := beginStr[:strEnd]
	*body = (*body)[strStart+strEnd:]
	return str, nil
}
