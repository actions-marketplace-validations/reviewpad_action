// Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package utils

import (
	"errors"
)

var (
	ErrBranchNil               = errors.New("branch is nil")
	ErrBranchRefNil            = errors.New("branch ref is nil")
	ErrBranchRepoNil           = errors.New("branch repo is nil")
	ErrBranchRepoNameNil       = errors.New("branch repo name is nil")
	ErrBranchRepoOwnerNil      = errors.New("branch repo owner is nil")
	ErrBranchRepoOwnerLoginNil = errors.New("branch repo owner login is nil")
)
