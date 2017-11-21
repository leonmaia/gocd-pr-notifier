package main

import (
	"encoding/json"
	"testing"
)

func Test_SchedulableShouldBeTrue(t *testing.T) {
	gocdPayload := GoCDPipelineResponse{}
	payload := []byte(`{
		"pausedCause": "",
		"pausedBy": "",
		"paused": false,
		"schedulable": true,
		"locked": false
	  }`)

	json.Unmarshal(payload, &gocdPayload)
	if gocdPayload.Schedulable != true {
		t.Error("Decode didn't happen correctly")
	}
}

func Test_SchedulableShouldBeFalse(t *testing.T) {
	gocdPayload := GoCDPipelineResponse{}
	payload := []byte(`{
		"pausedCause": "",
		"pausedBy": "",
		"paused": false,
		"schedulable": false,
		"locked": false
	  }`)

	json.Unmarshal(payload, &gocdPayload)
	if gocdPayload.Schedulable != false {
		t.Error("Decode didn't happen correctly")
	}
}
