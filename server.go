package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
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

type GitHubPayload struct {
	after string
}

func isPipelineAvailable() bool {
	status := GoCDPipelineResponse{}
	var url = fmt.Sprintf("http://%s:8153/go/api/pipelines/orgs-service-pr/status", *HTTPAddr)
	req, _ := requests.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", *Authentication)
	resp, err := req.Do()
	if err != nil {
		log.Fatal("Error while getting pipeline status: ", err)
	}
	json.NewDecoder(resp.Body).Decode(&status)
	return status.Schedulable
}

func notifyGoCDOfChangeInPR(w http.ResponseWriter, r *http.Request) {
	ghRequest := GitHubPayload{}
	var apiURL = fmt.Sprintf("http://%s:8153/go/api/pipelines/orgs-service-pr/schedule", *HTTPAddr)

	form := url.Values{}
	json.NewDecoder(r.Body).Decode(&ghRequest)
	form.Add("materials[pr-material]", ghRequest.after)

	req, err := requests.NewRequest("POST", apiURL, strings.NewReader(form.Encode()))
	if err != nil {
		log.Fatal("error while creating request: ", err)
	}
	req.Header.Set("Authorization", *Authentication)
	req.Header.Set("Confirm", "true")

	if isPipelineAvailable() {
		_, err := req.Do()
		if err != nil {
			log.Fatal("Error while scheduling: ", err)
		}
	} else {
		work := WorkRequest{Request: req, Delay: 10 * time.Second}
		Collector(work)
	}
}

func main() {
	flag.Parse()
	StartDispatcher(100)
	http.HandleFunc("/github-webhook", notifyGoCDOfChangeInPR)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
