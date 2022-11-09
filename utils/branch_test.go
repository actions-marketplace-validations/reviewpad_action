// Copyright (C) 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package utils

import (
	"testing"

	"github.com/google/go-github/v48/github"
)

func TestValidateBranch(t *testing.T) {
	tests := []struct {
		name    string
		arg     *github.PullRequestBranch
		wantErr error
	}{
		{
			name: "no error",
			arg: &github.PullRequestBranch{
				Ref: github.String("refs/heads/test"),
				Repo: &github.Repository{
					Name: github.String("reviewpad"),
					Owner: &github.User{
						Login: github.String("reviewpad"),
					},
				},
			},
			wantErr: nil,
		},
		{
			name:    "branch is nil",
			arg:     nil,
			wantErr: ErrBranchNil,
		},
		{
			name: "branch ref is nil",
			arg: &github.PullRequestBranch{
				Ref: nil,
			},
			wantErr: ErrBranchRefNil,
		},
		{
			name: "branch repo is nil",
			arg: &github.PullRequestBranch{
				Ref:  github.String("refs/heads/test"),
				Repo: nil,
			},
			wantErr: ErrBranchRepoNil,
		},
		{
			name: "branch repo name is nil",
			arg: &github.PullRequestBranch{
				Ref: github.String("refs/heads/test"),
				Repo: &github.Repository{
					Name: nil,
				},
			},
			wantErr: ErrBranchRepoNameNil,
		},
		{
			name: "branch repo owner is nil",
			arg: &github.PullRequestBranch{
				Ref: github.String("refs/heads/test"),
				Repo: &github.Repository{
					Name:  github.String("reviewpad"),
					Owner: nil,
				},
			},
			wantErr: ErrBranchRepoOwnerNil,
		},
		{
			name: "branch repo owner login is nil",
			arg: &github.PullRequestBranch{
				Ref: github.String("refs/heads/test"),
				Repo: &github.Repository{
					Name: github.String("reviewpad"),
					Owner: &github.User{
						Login: nil,
					},
				},
			},
			wantErr: ErrBranchRepoOwnerLoginNil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ValidateBranch(test.arg)
			if err != test.wantErr {
				t.Fatalf("got error %v, want %v", err, test.wantErr)
			}
		})
	}
}
