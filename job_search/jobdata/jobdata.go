package jobdata

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

/*
Templates is a list of TemplateData
*/
type Templates struct {
	TemplateData []TemplateData
}

/*
TemplateData is the structure of the data retrieved from
each URL stored in json so it can be displayed to webpage
*/
type TemplateData struct {
	Key            int    `json:"key"`
	Date           string `json:"dateAdded"`
	URL            string `json:"URL"`
	JobTitle       string `json:"jobTitle"`
	CompanyName    string `json:"companyName"`
	JobDescription string `json:"jobDescription"`
	ApplyLink      string `json:"applyLink"`
	AppliedStatus  int    `json:"appliedStatus"`
}

/*
UpdateStatus gets the updates the status in the json file.
1 denotes applied, 2 denotes neutral, and 3 denotes not interested
*/
func UpdateStatus(key, value int) {
	jsonJobData, templateList, _ := GetOld()

	templateList.TemplateData[key-1].AppliedStatus = value

	// Marshal the data into j then write
	// it all back to the open file
	WriteNew(&templateList, jsonJobData)
}

/*
GetOld opens the old json file and gets the information from it.
It then unmarshals the data into the Template and creates a map of URLs.
It returns the open file, the Template, and the map of URLs.
*/
func GetOld() (*os.File, Templates, map[string]int) {
	jd, err := os.OpenFile("jobData.json", os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}

	// Get the information already stored in the json file
	data, err := ioutil.ReadFile("jobData.json")
	if err != nil {
		log.Fatal(err)
	}

	// Don't unmarshal the data if the file is empty
	var templateList Templates
	if len(data) != 0 {
		err = json.Unmarshal(data, &templateList.TemplateData)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Create a map of all URLs from the file to quickly check
	// if a link already exists so duplicates are not added
	existingURLs := make(map[string]int)
	for i, url := range templateList.TemplateData {
		existingURLs[url.URL] = i
	}

	return jd, templateList, existingURLs
}

/*
WriteNew marshals the data from the Template and then writes it back to the file.
*/
func WriteNew(tl *Templates, jd *os.File) {
	j, err := json.MarshalIndent(tl.TemplateData, "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	_, err = jd.Write(j)
	if err != nil {
		log.Fatal(err)
	}
}
