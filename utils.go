package main

import (
	"fmt"
	"net/url"
	"strings"
)

// extract owner and repo name from github url 
// For eg 
// url =  https://github.com/{owner}/{repo_name}


func ExtractOwnerRepo(githubURL string) (string, string, error) {
	u, err := url.Parse(githubURL)
	if err != nil {
		return "", "", fmt.Errorf("invalid URL: %v", err)
	}

	if u.Host != "github.com" {
		return "", "", fmt.Errorf("not a GitHub URL")
	}

	pathParts := strings.Split(u.Path, "/")
	if len(pathParts) < 3 { // Expecting at least /owner/repo
		return "", "", fmt.Errorf("invalid GitHub path format")
	}

	owner := pathParts[1]
	repo := pathParts[2]

	// Handle potential .git suffix
	repo = strings.TrimSuffix(repo, ".git")

	return owner, repo, nil
}
