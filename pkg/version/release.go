package release

import (
	"context"
	"fmt"

	"github.com/google/go-github/v24/github"
)

const (
	owner = "weaveworks"
	repo  = "footloose"
)

// FindLastRelease searches latest release of the project
func FindLastRelease() (*github.RepositoryRelease, error) {
	githubclient := github.NewClient(nil)
	repoRelease, _, err := githubclient.Repositories.GetLatestRelease(context.Background(), owner, repo)
	if err != nil {
		return nil, fmt.Errorf("Failed to get latest release information")
	}
	return repoRelease, nil
}
