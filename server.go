package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
)

var (
	httpAddr       = flag.String("http", "127.0.0.1", "gocd address")
	authentication = flag.String("auth", "", "gocd authentication")
	pipelineName   = flag.String("pipeline", "", "pipeline name to trigger build")
	materialName   = flag.String("material", "", "material name to trigger build")
)

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	// Save a copy of this request for debugging.
	requestDump, _ := httputil.DumpRequest(r, true)
	fmt.Println(string(requestDump))

	ghPayload := GitHubPullRequestPayload{}
	err := json.NewDecoder(r.Body).Decode(&ghPayload)
	if err != nil {
		http.Error(w, "Please send a valid request body", 400)
		return
	}
	fmt.Println(ghPayload)

	statusCheckURL := fmt.Sprintf("http://%s:8153/go/api/pipelines/%s/status", *httpAddr, *pipelineName)
	notifyURL := fmt.Sprintf("http://%s:8153/go/api/pipelines/%s/schedule", *httpAddr, *pipelineName)

	if err := notifyGoCDOfChangeInPR(*pipelineName, *materialName, notifyURL, statusCheckURL, *authentication, ghPayload); err != nil {
		http.Error(w, err.Error(), 500)
	}
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
