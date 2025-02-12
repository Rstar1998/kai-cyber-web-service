package main

import (
	"strings"
	"testing"
)

func TestExtractOwnerRepo(t *testing.T) {
	testCases := []struct {
		name       string
		githubURL  string
		wantOwner  string
		wantRepo   string
		wantErr    bool
		wantErrMsg string // Optional: for more specific error checking
	}{
		{
			name:      "Valid GitHub URL",
			githubURL: "https://github.com/owner/repo",
			wantOwner: "owner",
			wantRepo:  "repo",
			wantErr:   false,
		},
		{
			name:      "Valid GitHub URL with .git suffix",
			githubURL: "https://github.com/owner/repo.git",
			wantOwner: "owner",
			wantRepo:  "repo",
			wantErr:   false,
		},
		{
			name:      "Valid GitHub URL with trailing slash",
			githubURL: "https://github.com/owner/repo/",
			wantOwner: "owner",
			wantRepo:  "repo",
			wantErr:   false,
		},
		{
			name:       "Invalid URL format",
			githubURL:  "invalid-url",
			wantOwner:  "",
			wantRepo:   "",
			wantErr:    true,
			wantErrMsg: "not a GitHub URL", // check for specific error message
		},
		{
			name:       "Not a GitHub URL",
			githubURL:  "https://example.com/owner/repo",
			wantOwner:  "",
			wantRepo:   "",
			wantErr:    true,
			wantErrMsg: "not a GitHub URL",
		},
		{
			name:       "Empty URL",
			githubURL:  "",
			wantOwner:  "",
			wantRepo:   "",
			wantErr:    true,
			wantErrMsg: "not a GitHub URL",
		},
		{
			name:      "GitHub URL with query parameters",
			githubURL: "https://github.com/owner/repo?query=param",
			wantOwner: "owner",
			wantRepo:  "repo",
			wantErr:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			owner, repo, err := ExtractOwnerRepo(tc.githubURL)

			if (err != nil) != tc.wantErr {
				t.Errorf("ExtractOwnerRepo(%q) error = %v, wantErr %v", tc.githubURL, err, tc.wantErr)
				return // Stop here if error expectation doesn't match
			}

			if tc.wantErr && tc.wantErrMsg != "" && !strings.Contains(err.Error(), tc.wantErrMsg) {
				t.Errorf("ExtractOwnerRepo(%q) error message = %v, wantErrMsg %v", tc.githubURL, err.Error(), tc.wantErrMsg)
				return // Stop here if error message expectation doesn't match

			}

			if owner != tc.wantOwner {
				t.Errorf("ExtractOwnerRepo(%q) owner = %q, want %q", tc.githubURL, owner, tc.wantOwner)
			}
			if repo != tc.wantRepo {
				t.Errorf("ExtractOwnerRepo(%q) repo = %q, want %q", tc.githubURL, repo, tc.wantRepo)
			}
		})
	}
}
