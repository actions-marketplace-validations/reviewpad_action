// Copyright (C) 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package utils

import (
	"errors"
	"testing"

	"github.com/google/go-github/v48/github"
	"github.com/stretchr/testify/assert"
)

func TestValidateUrl(t *testing.T) {
	tests := []struct {
		name       string
		arg        string
		wantBranch *github.PullRequestBranch
		wantFile   string
		wantErr    error
	}{
		{
			name: "no error",
			arg:  "https://github.com/reviewpad/reviewpad/blob/main/reviewpad.yml",
			wantBranch: &github.PullRequestBranch{
				Ref: github.String("main"),
				Repo: &github.Repository{
					Name: github.String("reviewpad"),
					Owner: &github.User{
						Login: github.String("reviewpad"),
					},
				},
			},
			wantFile: "reviewpad.yml",
			wantErr:  nil,
		},
		{
			name:       "error on invalid url",
			arg:        "https://gitlab.com/reviewpad/reviewpad/blob/main/reviewpad.yml",
			wantBranch: nil,
			wantFile:   "",
			wantErr:    errors.New("fatal: url must be a link to a GitHub blob, e.g. https://github.com/reviewpad/action/blob/main/main.go"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			wantBranch, wantFile, err := ValidateUrl(test.arg)

			assert.Equal(t, test.wantBranch, wantBranch)
			assert.Equal(t, test.wantFile, wantFile)
			assert.Equal(t, test.wantErr, err)
		})
	}
}
