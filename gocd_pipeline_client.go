package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/leonmaia/requests"
)

func isPipelineAvailable(pipelineName, statusCheckURL, auth string) bool {
	status := GoCDPipelineResponse{}
	req, _ := requests.NewRequest("GET", statusCheckURL, nil)
	req.Header.Set("Authorization", auth)
	resp, err := req.Do()
	if err != nil {
		log.Fatal("Error while getting pipeline status: ", err)
		return false
	}
	err = json.NewDecoder(resp.Body).Decode(&status)
	if err != nil {
		log.Fatal("Error while decoding response: ", err)
		return false
	}
	return status.Schedulable
}

func notifyGoCDOfChangeInPR(pipelineName, materialName, notifyURL, statusCheckURL, auth string, ghPayload GitHubPayload) {
	form := url.Values{}
	form.Add(fmt.Sprintf("materials[%s]", materialName), ghPayload.after)

	req, err := requests.NewRequest("POST", notifyURL, strings.NewReader(form.Encode()))
	if err != nil {
		log.Fatal("error while creating request: ", err)
	}
	req.Header.Set("Authorization", auth)
	req.Header.Set("Confirm", "true")

	if isPipelineAvailable(pipelineName, statusCheckURL, auth) {
		_, err = req.Do()
		if err != nil {
			log.Fatal("Error while scheduling: ", err)
		}
	} else {
		work := WorkRequest{Request: req, Delay: 10 * time.Second, PipelineName: pipelineName, StatusCheckURL: statusCheckURL, Auth: auth}
		Collector(work)
	}
}
