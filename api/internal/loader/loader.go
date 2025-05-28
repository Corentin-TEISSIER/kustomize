// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

// Package loader has a data loading interface and various implementations.
package loader

import (
	"os"
	"path/filepath"

	"sigs.k8s.io/kustomize/api/ifc"
	"sigs.k8s.io/kustomize/api/internal/git"
	"sigs.k8s.io/kustomize/kyaml/errors"
	"sigs.k8s.io/kustomize/kyaml/filesys"
)

// NewLoader returns a Loader pointed at the given target.
// If the target is remote, the loader will be restricted
// to the root and below only.  If the target is local, the
// loader will have the restrictions passed in.  Regardless,
// if a local target attempts to transitively load remote bases,
// the remote bases will all be root-only restricted.
func NewLoader(
	lr LoadRestrictorFunc,
	target string, fSys filesys.FileSystem) (ifc.Loader, error) {
	repoSpec, err := git.NewRepoSpecFromURL(target)
	if err == nil {
		// The target qualifies as a remote git target.
		return newLoaderAtGitClone(
			repoSpec, fSys, nil, git.ClonerUsingGitExec)
	}

	root, err := filesys.ConfirmDir(fSys, target)
	if err == nil {
		return newLoaderAtConfirmedDir(
			lr, root, fSys, nil, git.ClonerUsingGitExec), nil
	}

	info, err := os.Stat(target)
	if err != nil {
		return nil, errors.WrapPrefixf(err, ErrRtNotDirNorFile.Error()) //nolint:govet
	}

	fileName := info.Name()

	absPath, err := filepath.Abs(target)
	if err != nil {
		return nil, err
	}

	root, err = filesys.ConfirmDir(fSys, filepath.Dir(absPath))
	if err != nil {
		return nil, err
	}

	return newLoaderAtConfirmedPath(lr, root, fileName, fSys, nil, git.ClonerUsingGitExec), nil

}
