package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

func publicKey(privateKeyPath string) (*ssh.PublicKeys, error) {
	var publicKey *ssh.PublicKeys
	sshKey, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("ioutil.ReadFile: %w", err)
	}
	publicKey, err = ssh.NewPublicKeys("git", sshKey, "")
	if err != nil {
		return nil, fmt.Errorf("ssh.NewPublicKeys: %w", err)
	}
	return publicKey, err
}

func gitAuth() (*ssh.PublicKeys, error) {
	return publicKey(path.Clean(*privateKeyFile))
}

func gitCloneOrGetRepo() (*git.Repository, error) {
	if _, err := os.Stat(path.Join(*repoDir, ".git")); !os.IsNotExist(err) {
		// Already exists
		repo, err := git.PlainOpen(*repoDir)
		if err != nil {
			return nil, err
		}
		return repo, nil
	}
	auth, err := gitAuth()
	if err != nil {
		return nil, err
	}
	log.Println("Cloning git repo")
	repo, err := git.PlainClone(*repoDir, false, &git.CloneOptions{
		URL:           *gitRepo,
		Auth:          auth,
		Progress:      os.Stdout,
		Depth:         1,
		SingleBranch:  true,
		ReferenceName: plumbing.NewBranchReferenceName(*gitBranch),
		Tags:          git.NoTags,
		NoCheckout:    true,
	})
	if err != nil {
		return nil, err
	}
	return repo, nil
}

// gitSync returns true if any changes were found
func gitSync() (bool, error) {
	repo, err := gitCloneOrGetRepo()
	if err != nil {
		return false, fmt.Errorf("gitCloneOrGetRepo: %w", err)
	}
	remote, err := repo.Remote("origin")
	if err != nil {
		return false, fmt.Errorf("repo.Remote: %w", err)
	}
	auth, err := gitAuth()
	if err != nil {
		return false, fmt.Errorf("gitAuth: %w", err)
	}
	log.Println("Doing git fetch")
	err = remote.Fetch(&git.FetchOptions{
		RemoteName: "origin",
		Auth:       auth,
		Depth:      1,
		Force:      true,
		Progress:   os.Stdout,
		Tags:       git.NoTags,
	})
	if err != git.NoErrAlreadyUpToDate && err != nil {
		return false, fmt.Errorf("remote.Fetch: %w", err)
	}
	wt, err := repo.Worktree()
	if err != nil {
		return false, fmt.Errorf("repo.Worktree: %w", err)
	}
	initialHead, err := repo.Head()
	if err != nil {
		return false, fmt.Errorf("repo.Head: %w", err)
	}
	remoteRef, err := repo.Reference(plumbing.NewRemoteReferenceName("origin", *gitBranch), true)
	if err != nil {
		return false, fmt.Errorf("repo.Reference: %w", err)
	}
	err = wt.Reset(&git.ResetOptions{
		Commit: remoteRef.Hash(),
		Mode:   git.HardReset,
	})
	if err != nil {
		return false, fmt.Errorf("wt.Reset: %w", err)
	}
	newHead, err := repo.Head()
	if err != nil {
		return false, fmt.Errorf("repo.Head: %w", err)
	}
	return *initialHead != *newHead, nil
}
