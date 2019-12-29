package indeed

import (
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

// GetJobs searches for links on glassdoor
// based on the keyword (k), location (l), and age (a)
// the user inputs. It then writes the data to a json file.
func GetJobs(k, l, a string) {
	jsonJobData, templateList, existingURLs := jobdata.GetOld()

	var urlQuery string

	// Gets keyword to search from user input
	// and prepares it to be used in the get request
	// then does the same for location
	keyword := k
	keyword = url.QueryEscape(strings.TrimSpace(keyword))

	location := l
	location = url.QueryEscape(location)

	// Age is how old the jobs you want to collect are
	// converted to an int to use as comparison later
	fromAge := a
	fromAge = strings.TrimSuffix(fromAge, "\r\n")
	intAge, _ := strconv.Atoi(fromAge)
	// builds the url for the get request
	urlQuery += "q=" + keyword + "&"
	urlQuery += "l=" + location + "&"
	urlQuery += "sort=date"
	u := &url.URL{
		Scheme:   "https",
		Host:     "indeed.com",
		Path:     "jobs",
		RawQuery: urlQuery,
	}

	var jobList jobsList
	getJobURLs(u.String(), &jobList, intAge)
	parseJobData(&jobList, &templateList, existingURLs)

	jobdata.WriteNew(&templateList, jsonJobData)
}

func getJobURLs(uurl string, jl *jobsList, age int) {
	i := 0
	for {
		nurl := uurl + "&start=" + strconv.Itoa(i)
		res, err := http.Get(nurl)
		if err != nil {
			log.Fatal(err)
		}
		if res.StatusCode == 404 {
			break
		}
		defer res.Body.Close()

		data, _ := ioutil.ReadAll(res.Body)
		bodyStr := string(data)

		// Either of these strings being found means
		// there are no jobs on the page so break
		noJobs := "did not match any jobs</h2>"
		if strings.Index(bodyStr, noJobs) != -1 {
			break
		}
		noJobs = "could not be found.</h2>"
		if strings.Index(bodyStr, noJobs) != -1 {
			break
		}
		// Added links is used to determine whether or not to
		// keep going through pages since indeed doesn't completely sort
		// their results.
		addedLinks := 0
		for {
			// jk:' is a key that is searched for in the response
			// body. ' is the end of the key
			link, err := srchMovBody("jk:'", "'", &bodyStr)
			if err != nil {
				break
			}

			// datajk is equal to the key stored in link.
			// In order to get the dates from each job I search for
			// the job key in the response and then find the first date
			// the appeaers which is the date that job was posted
			datajk := link
			dateEst := strings.Index(bodyStr, datajk)
			if dateEst == -1 {
				break
			}
			estDateStart := bodyStr[dateEst:]
			dateStr := "<span class=\"date\">"
			endDateStr := "</span>"
			date, err := srchBody(dateStr, endDateStr, &estDateStart)
			if err != nil {
				break
			}

			// Splits the date so we just have the number or first word
			date = strings.Split(date, " ")[0]
			date = strings.Trim(date, "+")

			// Always add jobs with today as the date to the list
			// then compares the date posted to the age you chose
			// and adds all jobs that were posted on or before then
			intDate, _ := strconv.Atoi(date)
			if date == "Today" || date == "Just" {
				jl.jobs = append(jl.jobs, link)
				addedLinks++
			} else if intDate <= age {
				jl.jobs = append(jl.jobs, link)
				addedLinks++
			} else {
				continue
			}
		}
		// The first time you go through a page and don't add any links
		// meaning all the links on the page were out of the date range
		// stop searching
		if addedLinks == 0 {
			break
		}
		// If the next page arrow is gone stop searching as there are no more pages
		noJobs = "<span class=np>Next&nbsp"
		if strings.Index(bodyStr, noJobs) == -1 {
			break
		}
		i += 10
	}
}

func parseJobData(jl *jobsList, tl *jobdata.Templates, existingURLs map[string]int) {
	for _, job := range jl.jobs {
		job = "https://www.indeed.com/viewjob?jk=" + job
		res, err := http.Get(job)
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()

		for {
			// if the URL already exists in the file
			// adding it is skipped to prevent duplicates
			_, ok := existingURLs[job]
			if ok {
				break
			}
			data, _ := ioutil.ReadAll(res.Body)
			bodyStr := string(data)

			// Searches for the title
			titleStr := "jobsearch-JobInfoHeader-title\">"
			endTitleStr := "</h3>"
			title, err := srchBody(titleStr, endTitleStr, &bodyStr)
			if err != nil {
				break
			}

			// Searches for the company
			compStr := "companyrating\"><div class=\"icl-u-lg-mr--sm icl-u-xs-mr--xs\">"
			endCompStr := "</div>"
			company, err := srchBody(compStr, endCompStr, &bodyStr)
			if err != nil {
				break
			}

			// Searches for the apply link and returns a link back
			// to the page if it is a job that you apply directly from indeed
			applyStr := "class=\"icl-u-lg-hide\"><a class=\"icl-Button icl-Button--primary icl-Button--block\" href=\""
			endApplyStr := "\""
			apply, err := srchBody(applyStr, endApplyStr, &bodyStr)
			if err != nil {
				applyStr = "jobsearch-IndeedApplyButton-buttonWrapper icl-u-lg-block icl-u-xs-hide"
				endApplyStr = "\""
				_, err = srchBody(applyStr, endApplyStr, &bodyStr)
				if err != nil {
					break
				} else {
					apply = job
				}
			}

			// Searches for the job description
			descStr := "class=\"jobsearch-JobComponent-description icl-u-xs-mt--md\">"
			endDescStr := "<div class=\"jobsearch-JobDescriptionTab-content\">"
			description, err := srchBody(descStr, endDescStr, &bodyStr)
			if err != nil {
				break
			}

			// Gets the time so we can use the
			// date as the date we added each URL
			currentTime := time.Now()
			t := jobdata.TemplateData{
				Key:            len(tl.TemplateData) + 1,
				Date:           currentTime.Format("02-01-2006"),
				URL:            job,
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
func srchMovBody(start, end string, body *string) (string, error) {
	strStart := strings.Index(*body, start)
	if strStart == -1 {
		return "", errors.New("Error retrieving str")
	}
	beginStr := (*body)[strStart+len(start):]
	strEnd := strings.Index(beginStr, end)
	if strEnd == -1 {
		return "", errors.New("Error retrieving str")
	}
	str := beginStr[:strEnd]
	*body = (*body)[strStart+strEnd:]
	return str, nil
}

// Searches for the start string and end string and returns whats between them
// from the end of start to the beginning of end.
func srchBody(start, end string, body *string) (string, error) {
	strStart := strings.Index(*body, start)
	if strStart == -1 {
		return "", errors.New("Error retrieving str")
	}
	beginStr := (*body)[strStart+len(start):]
	strEnd := strings.Index(beginStr, end)
	if strEnd == -1 {
		return "", errors.New("Error retrieving str")
	}
	str := beginStr[:strEnd]
	return str, nil
}
