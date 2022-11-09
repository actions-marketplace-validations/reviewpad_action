// Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package utils

import "github.com/google/go-github/v48/github"

func ValidateBranch(branch *github.PullRequestBranch) error {
	if branch == nil {
		return ErrBranchNil
	}

	if branch.Ref == nil {
		return ErrBranchRefNil
	}

	baseRepo := branch.Repo
	if baseRepo == nil {
		return ErrBranchRepoNil
	}

	if baseRepo.Name == nil {
		return ErrBranchRepoNameNil
	}

	if baseRepo.Owner == nil {
		return ErrBranchRepoOwnerNil
	}

	if baseRepo.Owner.Login == nil {
		return ErrBranchRepoOwnerLoginNil
	}

	return nil
}
