package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
)

var (
	httpAddr       = flag.String("http", "127.0.0.1", "gocd address")
	authentication = flag.String("auth", "", "gocd authentication")
	pipelineName   = flag.String("pipeline_name", "", "pipeline name to trigger build")
	materialName   = flag.String("material_name", "", "material name to trigger build")
)

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	ghRequest := GitHubPayload{}
	err := json.NewDecoder(r.Body).Decode(&ghRequest)
	if err != nil {
		log.Fatal("Error while decoding request: ", err)
		return
	}

	statusCheckURL := fmt.Sprintf("http://%s:8153/go/api/pipelines/%s/status", *httpAddr, *pipelineName)
	notifyURL := fmt.Sprintf("http://%s:8153/go/api/pipelines/%s/schedule", *httpAddr, *pipelineName)

	notifyGoCDOfChangeInPR(*pipelineName, *materialName, notifyURL, statusCheckURL, *authentication, ghRequest)
}

func main() {
	flag.Parse()
	StartDispatcher(100)
	http.HandleFunc("/github-webhook", webhookHandler)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
