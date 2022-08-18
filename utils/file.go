// Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package utils

import (
	"bytes"
	"context"
	"strings"

	"github.com/google/go-github/v45/github"
	"github.com/reviewpad/reviewpad/v3"
	reviewpad_gh "github.com/reviewpad/reviewpad/v3/codehost/github"
	"github.com/reviewpad/reviewpad/v3/engine"
)

const pullRequestFileLimit = 50

func LoadReviewpadFile(ctx context.Context, githubClient *reviewpad_gh.GithubClient, filePath string, branch *github.PullRequestBranch) (*engine.ReviewpadFile, error) {
	reviewpadFileContent, err := githubClient.DownloadContents(ctx, filePath, branch)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(reviewpadFileContent)

	return reviewpad.Load(buf)
}

// reviewpad-an: experimental
// ReviewpadFileChanges checks if a file path was changed in a pull request.
// The way this is done depends on the number of files changed in the pull request.
// If the number of files changed is greater than pullRequestFileLimit,
// then we download both files using the filePath and check their contents.
// This strategy assumes that the file path exists in the head branch.
// Otherwise, we download the pull request files and check the filePath exists in them.
func ReviewpadFileChanged(ctx context.Context, githubClient *reviewpad_gh.GithubClient, filePath string, pullRequest *github.PullRequest) (bool, error) {
	if *pullRequest.ChangedFiles > pullRequestFileLimit {
		rawHeadFile, err := githubClient.DownloadContents(ctx, filePath, pullRequest.Head)
		if err != nil {
			return false, err
		}

		rawBaseFile, err := githubClient.DownloadContents(ctx, filePath, pullRequest.Base)
		if err != nil {
			if strings.HasPrefix(err.Error(), "no file named") {
				return true, nil
			}
			return false, err
		}

		// FIXME: use the hashes of the files
		return !bytes.Equal(rawBaseFile, rawHeadFile), nil
	}

	branchRepoOwner := *pullRequest.Base.Repo.Owner.Login
	branchRepoName := *pullRequest.Base.Repo.Name

	files, err := githubClient.GetPullRequestFiles(ctx, branchRepoOwner, branchRepoName, *pullRequest.Number)
	if err != nil {
		return false, err
	}

	for _, file := range files {
		if filePath == *file.Filename {
			return true, nil
		}
	}

	return false, nil
}
