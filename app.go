package main

import (
	"bytes"
	"encoding/json"
	"flag"
	grace "gopkg.in/paytm/grace.v1"
	logging "gopkg.in/tokopedia/logging.v1"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type (
	Config struct {
		App AppConfig
	}

	AppConfig struct {
		MattermostUrl      string
	}
)

var config Config

/****************************************************************************
 *
 * following is the callback data that runscope sends to custom webhook.
 *
 *****************************************************************************
{
   "test_id": "76598752-cbda-4e1d-820f-6274a62f74ff",
   "test_name": "Buckets Test",
   "test_run_id": "9c15aa62-21f0-48f2-a819-c99bdf8e4543",
   "team_id": "6b9c7f65-9e11-4f77-85ad-e6ee7a28232d",
   "team_name": "Acme Inc.",
   "environment_uuid": "98290cfc-a008-4ab7-9ea4-8906f12b228f",
   "environment_name": "Staging Settings",
   "bucket_name": "Rocket Sled",
   "bucket_key": "7xzcnsgbwox2",
   "test_url": "https://www.runscope.com/radar/7xzcnsgbwox2/76598752-cbda-4e1d-820f-6274a62f74ff",
   "test_run_url": "https://www.runscope.com/radar/7xzcnsgbwox2/76598752-cbda-4e1d-820f-6274a62f74ff/results/9c15aa62-21f0-48f2-a819-c99bdf8e4543",
   "trigger_url": "https://api.runscope.com/radar/09039249-fdfd-4e1d-820f-6274a62f74ff/trigger",
   "result": "fail",
   "started_at": 1384281308.548077,
   "finished_at": 1384281310.680218,
   "agent": null,
   "region": "us1",
   "region_name": "US East - Northern Virginia",
   "initial_variables": {},
   "requests": [{
     "url": "https://api.runscope.com/",
     "variables": {
        "fail": 0,
        "total": 1,
        "pass": 1
     },
     "assertions": {
        "fail": 0,
        "total": 2,
        "pass": 2
     },
     "scripts": {
        "fail": 0,
        "total": 1,
        "pass": 1
     },
     "result": "pass",
     "method": "GET",
     "response_time_ms": 123,
     "response_size_bytes": 2048,
     "response_status_code": 200,
     "note": "Root URL"
  }]
}
*/

type RunscopeRequestData struct {
	Url          string         `json:"url"`
	Variables    map[string]int `json:"variables"`
	Assertions   map[string]int `json:"assertions"`
	Scripts      map[string]int `json:"scipts"`
	Result       string         `json:"result"`
	ResponseTime int            `json:"response_time_ms"`
}

type RunscopeData struct {
	TestId          string  `json:"test_id"`
	TestName        string  `json:"test_name"`
	TestRunId       string  `json:"test_run_id"`
	TeamId          string  `json:"team_id"`
	TeamName        string  `json:"team_name"`
	EnvironmentUuid string  `json:"environment_uuid"`
	EnvironmentName string  `json:"environment_name"`
	BucketName      string  `json:"bucket_name"`
	BucketKey       string  `json:"bucket_key"`
	TestUrl         string  `json:"test_url"`
	TestRunUrl      string  `json:"test_run_url"`
	TriggerUrl      string  `json:"trigger_url"`
	Result          string  `json:"result"`
	StartedAt       float64 `json:"started_at"`
	FinishedAt      float64 `json:"finished_at"`
	Agent           string  `json:"agent"`
	Region          string  `json:"region"`
	RegionName      string  `json:"region_name"`

	/**
	 * other fields left. if you need please include.
	 */
	Requests []RunscopeRequestData `json:"requests"`
}

func runscopeHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("err", "reading the body from request", err)
		return
	}

	var data RunscopeData
	if err := json.Unmarshal(body, &data); err != nil {
		log.Println("err", "unmarshaling the body into data structure.", err)
		return
	}

    if data.Result == "pass" {
        log.Println("info:", "skipping for successful result.")
        return
    }

	status := "API test run completed. <" + data.TestRunUrl + "|View Results>  - <" + data.TestUrl + "|Edit Test>\\n"
	for _, req := range data.Requests {
		status = status + "> " + "<" + req.Url + "|" + data.TestName + ">" + "\\n"
		status = status + "> " + "**Status:** " + req.Result + "\t" + "**Requests Executed:** " + strconv.Itoa(req.Variables["total"]) + "\\n"
		status = status + "> " + "**Environment:** " + data.EnvironmentName + "\t" + "**Assertions Passed:** " + strconv.Itoa(req.Assertions["pass"]) + " of " + strconv.Itoa(req.Assertions["total"]) + "\\n"
		status = status + "> " + "**Location:** " + data.RegionName + "\t" + "**Scripts Passed:** " + strconv.Itoa(req.Scripts["pass"]) + " of " + strconv.Itoa(req.Scripts["total"]) + "\\n"
		status = status + "> " + "**Total Response Time:** " + strconv.Itoa(req.ResponseTime) + "\\n"
	}

	var jsonStr = []byte("{\"text\": \"" + status + "\"}")
	resp, err := http.Post(config.App.MattermostUrl, "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Println("err", "sending HTTP POST to mattermost", err)
		return
	}

	log.Println(resp.Status)
	defer resp.Body.Close()
}

func main() {
	flag.Parse()

	ok := logging.ReadModuleConfig(&config, ".", "app")
	if !ok {
		log.Fatal("Could not find configuration files.")
	}

	logging.LogInit()
	http.HandleFunc("/v1/runscope", runscopeHandler)
	port := ":8008"
	log.Fatal(grace.Serve(port, nil))
}
