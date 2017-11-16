package main

import (
	"time"

	"github.com/leonmaia/requests"
)

type WorkRequest struct {
	Request        *requests.Request
	Delay          time.Duration
	PipelineName   string
	StatusCheckURL string
	Auth           string
}
