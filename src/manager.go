package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/robfig/cron"
)

type FileInfos []os.FileInfo

type Job struct {
	Name     string
	Script   string
	Schedule string
}

func readJob(filename string) (Job, error) {
	file, err := ioutil.ReadFile(filename)
	var job Job

	if err != nil {
		return job, err
	}

	err = json.Unmarshal(file, &job)

	if err != nil {
		return job, err
	}

	return job, nil

}

func readJobs(dirname string) []Job {
	fis, err := ioutil.ReadDir(dirname)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	jobs := make([]Job, len(fis))
	for i, fi := range fis {
		filepath := dirname + "/" + fi.Name()

		job, err := readJob(filepath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		jobs[i] = job
	}

	return jobs
}

func addJobs(jobs []Job) {
	c := cron.New()
	for _, j := range jobs {
		job := j
		c.AddFunc(job.Schedule, func() { executeJob(job) })
		fmt.Println("EVENT:JOB-REGIST\tJOB-NAME:" + job.Name)
	}
	c.Start()
}

func executeJob(job Job) {
	path := "logs/" + job.Name + ".log"

	writeLogln(path, "EVENT:JOB-START\tJOB-NAME:"+job.Name)
	fmt.Println(job.Script)

	log.Info("Name: ", job.Name, " Script: ", job.Script)

	out, _ := exec.Command("sh", "-xc", job.Script).Output()

	writeLog(path, string(out))
	writeLogln(path, "EVENT:JOB-END\tJOB-NAME:"+job.Name)
}

func handler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	fmt.Println(r.Form)
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "Hello astaxie!")
}

func renderJobs(w http.ResponseWriter, r *http.Request) string {
	var html string
	jobs := readJobs("jobs")

	html += "<h1>Miniature Garden</h1>"
	html += "<h2>Job List</h2>"
	html += "<ul>"
	for _, job := range jobs {
		html += "<li><a href='/" + job.Name + "'>" + job.Name + "</a></li>"
	}
	html += "</ul>"

	return html
}

func renderJob(w http.ResponseWriter, r *http.Request, job Job) string {
	filename := "logs/" + job.Name + ".log"
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
	}
	text := string(file)
	lines := regexp.MustCompile(`\s*\n\s*`).Split(text, -1)
	sort.Sort(sort.Reverse(sort.StringSlice(lines)))

	var html string

	html += "<h1>Miniature Garden</h1>"
	html += "<h2>Job - " + job.Name + "</h2>"
	for _, l := range lines {
		html += "<p>" + l + "</p>"
	}
	return html
}

func writeLog(path string, message string) {
	var writer *bufio.Writer

	file, _ := os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
	writer = bufio.NewWriter(file)
	writer.WriteString(message)
	writer.Flush()
}

func writeLogln(path string, message string) {
	t := time.Now()
	writeLog(path, message+"\tTIMESTAMP:"+t.Format("2006-01-02 15:04:05")+"\n")
}

func handler2(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	var html string

	switch path {
	case "/":
		html = renderJobs(w, r)
	default:
		html = "404 not found."

		jobs := readJobs("jobs")
		for _, job := range jobs {
			if "/"+job.Name == path {
				html = renderJob(w, r, job)
			}
		}
	}

	fmt.Fprintf(w, html)
}

func main() {

	jobs := readJobs("jobs")
	addJobs(jobs)

	http.HandleFunc("/", handler2)
	fmt.Println("Run server, http://localhost:9090/")

	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
