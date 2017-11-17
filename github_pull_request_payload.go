package main

// GitHubPullRequestPayload is the type of payload received for when an action happens inside a pull request
type GitHubPullRequestPayload struct {
	Action, Before, After string
	PullRequest           struct {
		Head struct {
			Sha string
		} `json:"head"`
	} `json:"pull_request"`
}
