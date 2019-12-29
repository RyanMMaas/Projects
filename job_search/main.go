package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/RyanMMaas/projects/job_search/glassdoor"
	"github.com/RyanMMaas/projects/job_search/indeed"
	"github.com/RyanMMaas/projects/job_search/jobdata"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}

func main() {
	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/post", postHandler)
	http.HandleFunc("/updateStatus", updateStatusHandler)

	http.Handle("/data/", http.StripPrefix("/data", http.FileServer(http.Dir("./data"))))
	http.Handle("/public/", http.StripPrefix("/public", http.FileServer(http.Dir("./pub"))))
	http.ListenAndServe(":9000", nil)
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Title string
	}{"Search job sites"}
	err := tpl.ExecuteTemplate(w, "index.gohtml", data)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.NotFound(w, r)
	}
	keyword := r.FormValue("job-title")
	location := r.FormValue("location")
	age := r.FormValue("age")
	gdSearch := r.FormValue("glassdoor")
	indSearch := r.FormValue("indeed")

	if gdSearch == "Glassdoor" {
		fmt.Println("starting glassdoor search")
		glassdoor.GetJobs(keyword, location, age)
		fmt.Println("glassdoor search complete")
	}
	if indSearch == "Indeed" {
		fmt.Println("starting indeed search")
		indeed.GetJobs(keyword, location, age)
		fmt.Println("indeed search complete")
	}

	// Get the information already stored in the json file
	data, err := ioutil.ReadFile("jobData.json")
	if err != nil {
		log.Fatal(err)
	}
	w.Write(data)
}

func updateStatusHandler(w http.ResponseWriter, r *http.Request) {
	updateKey := r.FormValue("key")
	updateValue := r.FormValue("status")
	key, _ := strconv.Atoi(updateKey)
	val, _ := strconv.Atoi(updateValue)

	jobdata.UpdateStatus(key, val)
}
