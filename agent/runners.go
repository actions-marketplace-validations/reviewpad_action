// Copyright (C) 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package agent

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/go-github/v48/github"
	"github.com/reviewpad/action/v3/utils"
	"github.com/reviewpad/reviewpad/v3"
	reviewpad_gh "github.com/reviewpad/reviewpad/v3/codehost/github"
	"github.com/reviewpad/reviewpad/v3/collector"
	"github.com/reviewpad/reviewpad/v3/engine"
	"github.com/reviewpad/reviewpad/v3/handler"
)

// reviewpad-an: critical
func runReviewpad(entity *handler.TargetEntity, e *handler.ActionEvent, mixpanelToken, filePath, fileUrl string) {
	repo := *e.Repository
	splittedRepo := strings.Split(repo, "/")
	repoOwner := splittedRepo[0]
	repoName := splittedRepo[1]

	var err error
	var eventPayload interface{}

	// schedule events do not have a payload.
	if *e.EventName != "schedule" {
		eventPayload, err = github.ParseWebHook(*e.EventName, *e.EventPayload)
		if err != nil {
			utils.LogFatalErr(err)
		}
	}

	ctx, canc := context.WithTimeout(context.Background(), time.Minute*10)
	defer canc()

	githubClient := reviewpad_gh.NewGithubClientFromToken(ctx, *e.Token)

	var dryRun bool
	var reviewpadFileChanged bool
	var reviewpadFile *engine.ReviewpadFile
	var pullRequest *github.PullRequest

	githubUrl := fmt.Sprintf("github.com/%v/%v/%v/%v", entity.Owner, entity.Repo, entity.Kind, entity.Number)
	collectorClient := collector.NewCollector(mixpanelToken, entity.Owner, string(entity.Kind), githubUrl, "gh_action")

	if entity.Kind == handler.PullRequest {
		pullRequest, _, err = githubClient.GetPullRequest(ctx, repoOwner, repoName, entity.Number)
		if err != nil {
			utils.LogFatalErr(err)
		}

		if pullRequest.Merged != nil && *pullRequest.Merged {
			log.Print("skip execution for merged pull requests")
			return
		}

		if err := utils.ValidateBranch(pullRequest.Base); err != nil {
			utils.LogFatalErr(err)
		}

		if err := utils.ValidateBranch(pullRequest.Head); err != nil {
			utils.LogFatalErr(err)
		}
	}

	if fileUrl != "" {
		log.Printf("using remote config file %s", fileUrl)
		branch, filePath, err := utils.ValidateUrl(fileUrl)
		if err != nil {
			utils.LogFatalErr(err)
		}
		if reviewpadFile, err = utils.LoadReviewpadFile(ctx, githubClient, filePath, branch); err != nil {
			utils.LogFatalErr(err)
		}
	} else {
		log.Printf("using local config file %s", filePath)
		if entity.Kind == handler.PullRequest {
			reviewpadFileChanged, err = utils.ReviewpadFileChanged(ctx, githubClient, filePath, pullRequest)
			if err != nil {
				utils.LogFatalErr(err)
			}

			if reviewpadFileChanged {
				if reviewpadFile, err = utils.LoadReviewpadFile(ctx, githubClient, filePath, pullRequest.Head); err != nil {
					utils.LogFatalErr(err)
				}
			} else {
				if reviewpadFile, err = utils.LoadReviewpadFile(ctx, githubClient, filePath, pullRequest.Base); err != nil {
					utils.LogFatalErr(err)
				}
			}
		} else {
			reviewpadFileChanged = false
			defaultBranchName, err := githubClient.GetDefaultRepositoryBranch(ctx, repoOwner, repoName)
			if err != nil {
				utils.LogFatalErr(err)
			}

			defaultBranch := &github.PullRequestBranch{
				Repo: &github.Repository{
					Owner: &github.User{
						Login: github.String(repoOwner),
					},
					Name: github.String(repoName),
				},
				Ref: github.String(defaultBranchName),
			}

			if reviewpadFile, err = utils.LoadReviewpadFile(ctx, githubClient, filePath, defaultBranch); err != nil {
				utils.LogFatalErr(err)
			}
		}
	}

	dryRun = reviewpadFileChanged

	exitStatus, err := reviewpad.Run(ctx, githubClient, collectorClient, entity, eventPayload, reviewpadFile, dryRun, reviewpadFileChanged)
	if err != nil {
		if reviewpadFile.IgnoreErrors {
			utils.LogErr(err)
			return
		}

		utils.LogFatalErr(err)
	}

	if exitStatus == engine.ExitStatusFailure {
		log.Fatal("The executed program exited with code 1")
	}
}

// reviewpad-an: critical
func RunAction(githubToken, mixpanelToken, rawEvent, file, fileUrl string) {
	event, err := handler.ParseEvent(rawEvent)

	if err != nil {
		log.Fatalf("error parsing event: %v", err)
	}

	targetEntities, err := handler.ProcessEvent(event)
	if err != nil {
		log.Fatalf("error parsing event: %v", err)
	}

	event.Token = &githubToken

	for _, targetEntity := range targetEntities {
		runReviewpad(targetEntity, event, mixpanelToken, file, fileUrl)
	}
}
