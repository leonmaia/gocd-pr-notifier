package main

import (
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

func notifyGoCDOfChangeInPR(w http.ResponseWriter, r *http.Request) {
	var url = fmt.Sprintf("http://%s:8153/go/api/pipelines/orgs-service-pr/schedule", *HTTPAddr)
	req, _ := requests.NewRequest("POST", url, nil)

	req.Header.Set("Authorization", *Authentication)
	req.Header.Set("Confirm", "true")
	response, err := req.Do()
	if err != nil {
		fmt.Println(err)
	}
	if response.StatusCode == 409 {
		work := WorkRequest{Request: req, Delay: 90 * time.Second}
		Collector(work)
		return
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
