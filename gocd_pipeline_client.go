package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/leonmaia/requests"
)

func isPipelineAvailable(pipelineName, statusCheckURL, auth string) (bool, error) {
	status := GoCDPipelineResponse{}
	req, _ := createGoCDRequest("GET", statusCheckURL, auth, strings.NewReader(""))
	resp, err := req.Do()
	if err != nil {
		return false, fmt.Errorf("error while doing check pipeline request to gocd: %s", err.Error())
	}
	err = json.NewDecoder(resp.Body).Decode(&status)
	if err != nil {
		return false, fmt.Errorf("error while decoding gocd response: %s", err.Error())
	}

	return status.Schedulable, nil
}

func notifyGoCDOfChangeInPR(pipelineName, materialName, notifyURL, statusCheckURL, auth string, ghPayload GitHubPullRequestPayload) error {
	form := url.Values{}
	form.Add(fmt.Sprintf("materials[%s]", materialName), ghPayload.After)
	req, err := createGoCDRequest("POST", notifyURL, auth, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Confirm", "true")

	if result, err := isPipelineAvailable(pipelineName, statusCheckURL, auth); result && err == nil {
		_, err = req.Do()
		if err != nil {
			return fmt.Errorf("error while doing schedule request to gocd: %s", err.Error())
		}
	} else {
		work := WorkRequest{Request: req, Delay: 10 * time.Second, PipelineName: pipelineName, StatusCheckURL: statusCheckURL, Auth: auth}
		Collector(work)
	}

	return nil
}

func createGoCDRequest(method, url, auth string, body *strings.Reader) (*requests.Request, error) {
	req, err := requests.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("error creating gocd request: %s", err.Error())
	}
	req.Header.Set("Authorization", auth)

	return req, nil
}
