package main

import (
	"encoding/json"
	"testing"
)

func TestShouldBeAbleToDecodeGoCDResponse(t *testing.T) {
	ghPayload := GitHubPullRequestPayload{}
	sha := "d21783b6f67c9a0c19425f19305a2a35149c2e45"
	payload := []byte(`{
		"action": "synchronize",
		"pull_request": {
		  "head": {
			"label": "Tradeshift:ahhdh",
			"ref": "ahhdh",
			"sha": "d21783b6f67c9a0c19425f19305a2a35149c2e45"
		  }
		},
		"before": "21fa6915f780fc7c730c7ae42799cafdcf57a55c",
		"after": "d21783b6f67c9a0c19425f19305a2a35149c2e45"
	  }`)

	json.Unmarshal(payload, &ghPayload)
	if ghPayload.PullRequest.Head.Sha != sha {
		t.Error("Decode didn't happen correctly, SHA could not be retrieved")
	}
	if ghPayload.Action != "synchronize" {
		t.Error("Decode didn't happen correctly, action could not be retrieved")
	}
	if ghPayload.After != sha {
		t.Error("Decode didn't happen correctly, action could not be retrieved")
	}
}
