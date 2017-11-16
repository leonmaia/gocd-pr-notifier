package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/leonmaia/requests"
)

var (
	HTTPAddr       = flag.String("http", "127.0.0.1", "gocd address")
	Authentication = flag.String("auth", "", "gocd authentication")
)

type GoCDPipelineResponse struct {
	Schedulable bool
}

func isPipelineAvailable() bool {
	status := GoCDPipelineResponse{}
	var url = fmt.Sprintf("http://%s:8153/go/api/pipelines/orgs-service-pr/status", *HTTPAddr)
	req, _ := requests.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", *Authentication)
	resp, _ := req.Do()
	json.NewDecoder(resp.Body).Decode(&status)
	return status.Schedulable
}

func notifyGoCDOfChangeInPR(w http.ResponseWriter, r *http.Request) {
	flag.Parse()
	var url = fmt.Sprintf("http://%s:8153/go/api/pipelines/orgs-service-pr/schedule", *HTTPAddr)
	req, _ := requests.NewRequest("POST", url, nil)
	req.Header.Set("Authorization", *Authentication)
	req.Header.Set("Confirm", "true")

	if isPipelineAvailable() {
		req.Do()
	} else {
		work := WorkRequest{Request: req, Delay: 10 * time.Second}
		Collector(work)
	}
}

func main() {
	StartDispatcher(100)
	http.HandleFunc("/github-webhook", notifyGoCDOfChangeInPR)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
