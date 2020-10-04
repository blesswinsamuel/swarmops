package git

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	log "github.com/sirupsen/logrus"
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

type Git struct {
	keys   *ssh.PublicKeys
	repo   *git.Repository
	wt     *git.Worktree
	branch string
}

func NewGit(repoUrl, branch, repoDir, privateKeyFile string) (*Git, error) {
	keys, err := publicKey(path.Clean(privateKeyFile))
	if err != nil {
		return nil, fmt.Errorf("gitAuth: %w", err)
	}
	repo, err := gitCloneOrGetRepo(repoUrl, branch, repoDir, keys)
	if err != nil {
		return nil, fmt.Errorf("gitCloneOrGetRepo: %w", err)
	}
	wt, err := repo.Worktree()
	if err != nil {
		return nil, fmt.Errorf("repo.Worktree: %w", err)
	}
	return &Git{
		keys:   keys,
		repo:   repo,
		wt:     wt,
		branch: branch,
	}, nil
}

func gitCloneOrGetRepo(gitRepo, gitBranch, repoDir string, keys *ssh.PublicKeys) (*git.Repository, error) {
	repo, err := git.PlainOpen(repoDir)
	if err != nil {
		if err == git.ErrRepositoryNotExists {
			// Repo doesn't exist
			log.Infoln("Cloning git repo")
			repo, err := git.PlainClone(repoDir, false, &git.CloneOptions{
				URL:      gitRepo,
				Auth:     keys,
				Progress: os.Stdout,
				// Depth:         1,
				SingleBranch:  true,
				RemoteName:    "origin",
				ReferenceName: plumbing.NewBranchReferenceName(gitBranch),
				Tags:          git.NoTags,
			})
			if err != nil {
				return nil, fmt.Errorf("git.PlainClone: %w", err)
			}
			return repo, nil
		}
		return nil, fmt.Errorf("git.PlainOpen: %w", err)
	}
	return repo, nil
}

// gitSync returns true if any changes were found
func (g *Git) Sync() (bool, error) {
	remoteRef, err := g.gitFetch()
	if err != nil {
		return false, fmt.Errorf("g.gitFetch: %w", err)
	}
	outOfSync, err := g.gitIsOutOfSync(remoteRef)
	if err != nil {
		return false, fmt.Errorf("g.gitFetch: %w", err)
	}
	if !outOfSync {
		return false, nil
	}
	// repo out of sync
	err = g.gitReset(remoteRef.Hash())
	if err != nil {
		return false, fmt.Errorf("g.gitReset: %w", err)
	}
	return true, nil
}

// gitFetch returns does git fetch and returns new remote ref
func (g *Git) gitFetch() (*plumbing.Reference, error) {
	remote, err := g.repo.Remote("origin")
	if err != nil {
		return nil, fmt.Errorf("repo.Remote: %w", err)
	}
	log.Infoln("Doing git fetch")
	err = remote.Fetch(&git.FetchOptions{
		RemoteName: "origin",
		Auth:       g.keys,
		Force:      true,
		Progress:   os.Stdout,
		Tags:       git.NoTags,
		// Depth:      1,
	})
	if err != git.NoErrAlreadyUpToDate && err != nil {
		return nil, fmt.Errorf("remote.Fetch: %w", err)
	}
	remoteRef, err := g.repo.Reference(plumbing.NewRemoteReferenceName("origin", g.branch), true)
	if err != nil {
		return nil, fmt.Errorf("repo.Reference: %w", err)
	}
	return remoteRef, nil
}

// gitIsOutOfSync returns true if any changes were found
func (g *Git) gitIsOutOfSync(remoteRef *plumbing.Reference) (bool, error) {
	headRef, err := g.repo.Head()
	if err != nil {
		return false, fmt.Errorf("repo.Head: %w", err)
	}
	return headRef.Hash() != remoteRef.Hash(), nil
}

// gitReset returns true if any changes were found
func (g *Git) gitReset(hash plumbing.Hash) error {
	log.Infoln("Doing git reset")
	err := g.wt.Reset(&git.ResetOptions{
		Commit: hash,
		Mode:   git.HardReset,
	})
	if err != nil {
		return fmt.Errorf("wt.Reset: %w", err)
	}
	return nil
}
