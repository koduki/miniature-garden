package main

import (
    "encoding/json"
    "fmt"
    "os"
    "os/exec"
    "io/ioutil"
    "github.com/robfig/cron"
    "time"
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

func addJobs(c *cron.Cron, dirname string) {
   fileInfos, err := ioutil.ReadDir(dirname)
   if err != nil {
     fmt.Println(err)
     os.Exit(1)
   }
   
   for _, fi := range fileInfos {
      filepath := dirname + "/" + fi.Name()

      job, err := readJob(filepath)
      if err != nil {
	fmt.Println(err)
	os.Exit(1)
      }
      c.AddFunc(job.Schedule, func(){ executeJob(job) })
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

func main() {

   c := cron.New()
   addJobs(c, "jobs")
   c.Start()

   for {
     time.Sleep(10000000000000)
     fmt.Println("sleep")
   } 
}

