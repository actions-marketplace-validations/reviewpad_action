// Copyright (C) 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package main

import (
	"log"
	"os"

	"github.com/reviewpad/action/v3/agent"
)

var (
	rawEvent,
	file,
	fileUrl,
	gitHubToken,
	MixpanelToken string
)

func init() {
	rawEvent = os.Getenv("INPUT_EVENT")
	if rawEvent == "" {
		log.Fatal("missing variable INPUT_EVENT")
	}

	gitHubToken = os.Getenv("INPUT_TOKEN")
	if gitHubToken == "" {
		log.Fatal("missing variable INPUT_TOKEN")
	}

	file = os.Getenv("INPUT_FILE")
	fileUrl = os.Getenv("INPUT_FILE_URL")
	if file == "" && fileUrl == "" {
		log.Fatal("missing variable INPUT_FILE or INPUT_FILE_URL")
	}
}

func main() {
	agent.RunAction(gitHubToken, MixpanelToken, rawEvent, file, fileUrl)
}
