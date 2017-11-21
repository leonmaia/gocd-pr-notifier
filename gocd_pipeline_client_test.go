package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_shouldReturnTrueWhenPipelineIsAvailable(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		js, _ := json.Marshal(GoCDPipelineResponse{Schedulable: true})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(js)
	}))
	defer ts.Close()

	if result, err := isPipelineAvailable("pipeName", ts.URL, "blah"); result != true && err != nil {
		t.Error("Decode didn't happen correctly, action could not be retrieved")
	}
}

func Test_shouldReturnFalseWhenPipelineIsUnavailable(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		js, _ := json.Marshal(GoCDPipelineResponse{Schedulable: false})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(js)
	}))
	defer ts.Close()

	if result, err := isPipelineAvailable("pipeName", ts.URL, "blah"); result != false && err != nil {
		t.Error("Decode didn't happen correctly, action could not be retrieved")
	}
}

func Test_shouldReturnFalseWithErrorWhenPipelineCheckGoesWrong(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}))
	defer ts.Close()

	if result, err := isPipelineAvailable("pipeName", ts.URL, "blah"); result != false && err == nil {
		t.Error("Decode didn't happen correctly, action could not be retrieved")
	}
}
