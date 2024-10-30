package gitclient

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/google/go-github/v66/github"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
)

const (
	baseApiRepoOwner string = "darwishdev"
	baseApiRepoName  string = "devkit-api"
)

// GitClientInterface defines the interface for Git operations.
type GitClientInterface interface {
	Fork(ctx context.Context, org, name string) error
	Clone(ctx context.Context, org string, name string) error
}

// GitClientRepo implements GitClientInterface for GitHub.
type GitClientRepo struct {
	client *github.Client
	owner  string
	repo   string
}

// NewGitClientRepo initializes a GitClientRepo with a GitHub client.
func NewGitClientRepo(ctx context.Context, token string) GitClientInterface {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	log.Debug().Str("token", token).Msg("get token")
	client := github.NewClient(tc)

	return &GitClientRepo{
		client: client,
		owner:  baseApiRepoOwner,
		repo:   baseApiRepoName,
	}
}

func (g *GitClientRepo) Clone(ctx context.Context, org string, name string) error {
	gitCloneCmd := exec.Command("git", "clone", fmt.Sprintf(fmt.Sprintf("git@github.com:%s/%s.git", org, name)))
	gitCloneCmd.Stdout = os.Stdout
	gitCloneCmd.Stderr = os.Stderr
	err := gitCloneCmd.Run()
	if err != nil {
		return err
	}
	return nil

}

// Fork forks the base repository into a new one in the specified organization with the given name.
func (g *GitClientRepo) Fork(ctx context.Context, org string, name string) error {
	opts := &github.RepositoryCreateForkOptions{
		Organization: org,
		Name:         name,
	}
	fork, resp, err := g.client.Repositories.CreateFork(ctx, g.owner, g.repo, opts)
	if resp.StatusCode == 202 {
		fmt.Println("Fork is being created in a background task. Check again later.")
	} else if err != nil {
		return fmt.Errorf("failed to create fork: %w", err)
	}

	fmt.Printf("Repository forked to %s with URL: %s\n", fork.GetFullName(), fork.GetHTMLURL())
	return nil
}
