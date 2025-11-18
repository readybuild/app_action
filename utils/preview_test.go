package utils

import (
	"testing"

	"github.com/digitalocean/godo"
	gha "github.com/sethvargo/go-githubactions"
	"github.com/stretchr/testify/require"
)

func TestSanitizeSpecForPullRequestPreview(t *testing.T) {
	spec := &godo.AppSpec{
		Name:    "foo",
		Domains: []*godo.AppDomainSpec{{Domain: "foo.com"}},
		Alerts:  []*godo.AppAlertSpec{{Value: 80}},
		Services: []*godo.AppServiceSpec{{
			Name: "web",
			GitHub: &godo.GitHubSourceSpec{
				Repo:         "foo/bar",
				Branch:       "main",
				DeployOnPush: true,
			},
		}, {
			Name: "web2",
			GitHub: &godo.GitHubSourceSpec{
				Repo:         "another/repo",
				Branch:       "main",
				DeployOnPush: true,
			},
		}},
		Workers: []*godo.AppWorkerSpec{{
			Name: "worker",
			GitHub: &godo.GitHubSourceSpec{
				Repo:         "foo/bar",
				Branch:       "main",
				DeployOnPush: true,
			},
		}},
		Jobs: []*godo.AppJobSpec{{
			Name: "job",
			GitHub: &godo.GitHubSourceSpec{
				Repo:         "foo/bar",
				Branch:       "main",
				DeployOnPush: true,
			},
		}},
		Functions: []*godo.AppFunctionsSpec{{
			Name: "function",
			GitHub: &godo.GitHubSourceSpec{
				Repo:         "foo/bar",
				Branch:       "main",
				DeployOnPush: true,
			},
		}},
	}

	ghCtx := &gha.GitHubContext{
		Repository: "foo/bar",
		HeadRef:    "feature-branch",
		Event: map[string]any{
			"pull_request": map[string]any{
				"number": float64(3),
			},
		},
	}

	err := SanitizeSpecForPullRequestPreview(spec, ghCtx, false)
	require.NoError(t, err)

	expected := &godo.AppSpec{
		Name: "feature-branch", // Name got generated.
		// Domains and alerts got removed.
		Services: []*godo.AppServiceSpec{{
			Name: "web",
			GitHub: &godo.GitHubSourceSpec{
				Repo:         "foo/bar",
				Branch:       "feature-branch", // Branch got updated.
				DeployOnPush: false,            // DeployOnPush got set to false.
			},
		}, {
			Name: "web2",
			GitHub: &godo.GitHubSourceSpec{
				Repo:         "another/repo", // No change.
				Branch:       "main",
				DeployOnPush: true,
			},
		}},
		Workers: []*godo.AppWorkerSpec{{
			Name: "worker",
			GitHub: &godo.GitHubSourceSpec{
				Repo:         "foo/bar",
				Branch:       "feature-branch", // Branch got updated.
				DeployOnPush: false,            // DeployOnPush got set to false.
			},
		}},
		Jobs: []*godo.AppJobSpec{{
			Name: "job",
			GitHub: &godo.GitHubSourceSpec{
				Repo:         "foo/bar",
				Branch:       "feature-branch", // Branch got updated.
				DeployOnPush: false,            // DeployOnPush got set to false.
			},
		}},
		Functions: []*godo.AppFunctionsSpec{{
			Name: "function",
			GitHub: &godo.GitHubSourceSpec{
				Repo:         "foo/bar",
				Branch:       "feature-branch", // Branch got updated.
				DeployOnPush: false,            // DeployOnPush got set to false.
			},
		}},
	}

	require.Equal(t, expected, spec)
}

func TestGenerateAppName(t *testing.T) {
	tests := []struct {
		name       string
		repoOwner  string
		repo       string
		branchName string
		expected   string
	}{{
		name:       "success",
		repoOwner:  "foo",
		repo:       "bar",
		branchName: "feature-test-do-deploy2",
		expected:   "feature-test-do-deploy2",
	}, {
		name:       "branch with slashes",
		repoOwner:  "foo",
		repo:       "bar",
		branchName: "feature/test",
		expected:   "feature-test",
	}, {
		name:       "branch with underscores and dots",
		repoOwner:  "foo",
		repo:       "bar",
		branchName: "feature_test.v2",
		expected:   "feature-test-v2",
	}, {
		name:       "long branch name",
		repoOwner:  "foo",
		repo:       "bar",
		branchName: "this-is-an-extremely-long-branch-name-that-exceeds-the-limit",
		expected:   "this-is-an-extremely-long-branch",
	}, {
		name:       "long branch with truncation at hyphen",
		repoOwner:  "foo",
		repo:       "bar",
		branchName: "feature-with-a-very-long-name-",
		expected:   "feature-with-a-very-long-name",
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := GenerateAppName(test.repoOwner, test.repo, test.branchName)
			require.Equal(t, test.expected, got)
		})
	}
}
