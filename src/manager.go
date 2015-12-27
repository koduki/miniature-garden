package main

import (
    "encoding/json"
    "fmt"
    "os"
    "os/exec"
    "io/ioutil"
    "github.com/robfig/cron"
    "net/http"
    "strings"
    "log"
)

type FileInfos []os.FileInfo

type Job struct {
    Name string
    Script string
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

func readJobs(dirname string)([]Job) {
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

func addJobs(c *cron.Cron, jobs []Job) {
   for _, job := range jobs {
      c.AddFunc(job.Schedule, func(){ executeJob(job) })
      fmt.Println("EVENT:JOB-REGIST\tJOB-NAME:" + job.Name)
   }
}

func executeJob(job Job) {
    fmt.Println("EVENT:JOB-START\tJOB-NAME:" + job.Name)
    out, err := exec.Command("sh", "-c", job.Script).Output()

    if err != nil {
	fmt.Println(err)
	os.Exit(1)
    }

    fmt.Print(string(out))
    fmt.Println("EVENT:JOB-END\tJOB-NAME:" + job.Name)
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

func handlerJobs(w http.ResponseWriter, r *http.Request)(string) {
    var html string
    jobs := readJobs("jobs")


    html += "<h1>Miniature Garden</h1>"
    html += "<h2>Miniature Garden</h2>"
    html += "<ul>"
    for _, job := range jobs {
      html += "<li><a href=''>" + job.Name + "</a></li>"
    }
    html += "</ul>"

    return html
}

func handler2(w http.ResponseWriter, r *http.Request) {
    path := r.URL.Path

    var html string

    switch path {
    case "/":
        html = handlerJobs(w, r)
    default:
        html = "404 not found."
    }

    fmt.Fprintf(w, html)
}


func main() {

    c := cron.New()
    jobs := readJobs("jobs")
    addJobs(c, jobs)
//    c.Start()

    http.HandleFunc("/", handler2)
    fmt.Println("Run server, http://localhost:9090/")

    err := http.ListenAndServe(":9090", nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}

