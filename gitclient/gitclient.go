package gitclient

import (
	"context"
	"fmt"

	"github.com/google/go-github/v66/github"
	"golang.org/x/oauth2"
)

// GitClientInterface defines the interface for Git operations.
type GitClientInterface interface {
	Fork(ctx context.Context, org, name string) error
}

// GitClientRepo implements GitClientInterface for GitHub.
type GitClientRepo struct {
	client *github.Client
	owner  string
	repo   string
}

// NewGitClientRepo initializes a GitClientRepo with a GitHub client.
func NewGitClientRepo(ctx context.Context, token string, owner string, repo string) GitClientInterface {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	return &GitClientRepo{
		client: client,
		owner:  owner,
		repo:   repo,
	}
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
