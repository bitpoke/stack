package git

import (
	"regexp"

	"github.com/pkg/errors"
)

var regex = regexp.MustCompile(`(github\.com\/[a-zA-Z0-9-_]+\/[a-zA-Z0-9-_]+)`)

func NormalizeGitRepoUrl(repoUrl string) (string, error) {
	submatch := regex.FindStringSubmatch(repoUrl)
	if len(submatch) == 0 {
		return "", errors.New("invalid git repository url")
	}
	return submatch[0], nil
}

func GitRepoFollowedRef(repoUrl, branch string) (string, error) {
	gitUrl, err := NormalizeGitRepoUrl(repoUrl)
	if err != nil {
		return "", err
	}

	return gitUrl + "." + branch, nil
}
