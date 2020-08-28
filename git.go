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
	if _, err := os.Stat(privateKeyPath); os.IsNotExist(err) {
		// Key doesn't exist - create new one
		err := generateRsaKeyPairs(*keysDir)
		if err != nil {
			return nil, err
		}
		pubKey, err := getPublicKey()
		if err != nil {
			return nil, err
		}
		log.Printf("Created new key pair:\n%s", pubKey)
	}
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
	return publicKey(path.Join(*keysDir, privateKeyFile))
}

func gitCloneOrGetRepo() (*git.Repository, error) {
	if _, err := os.Stat(*repoDir); !os.IsNotExist(err) {
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
		return false, err
	}
	wt, err := repo.Worktree()
	if err != nil {
		return false, err
	}
	auth, err := gitAuth()
	if err != nil {
		return false, err
	}
	initialHead, err := repo.Head()
	if err != nil {
		return false, err
	}
	log.Println("Doing git pull")
	err = wt.Pull(&git.PullOptions{
		RemoteName:    "origin",
		Auth:          auth,
		Progress:      os.Stdout,
		Depth:         1,
		SingleBranch:  true,
		Force:         true,
		ReferenceName: plumbing.NewBranchReferenceName(*gitBranch),
	})
	newHead, err := repo.Head()
	if err != nil {
		return false, err
	}
	return *initialHead != *newHead, nil
}
