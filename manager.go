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

type Job struct {
    Name  string
    Script  string
}

func executeJob(filename string) {
    file, err := ioutil.ReadFile(filename)
    var job Job

    if err != nil {
	fmt.Println(err)
	os.Exit(1)
    }
  
    err = json.Unmarshal(file, &job) 

    fmt.Println("EVENT:JOB-START\tJOB-NAME:" + job.Name)
    out, err := exec.Command("sh", "-c", job.Script).Output()

    if err != nil {
	fmt.Println(err)
	os.Exit(1)
    }

    fmt.Println(string(out))
    fmt.Println("EVENT:JOB-END\tJOB-NAME:" + job.Name)
}

func main() {
   c := cron.New()
   c.AddFunc("*/1 * * * * *", func() {executeJob("jobs/crawle.json")})
   c.Start()

   for {
     time.Sleep(10000000000000)
     fmt.Println("sleep")
   } 
}

