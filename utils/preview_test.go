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
		Name: "bar-feature-branch", // Name got generated.
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
		expected:   "bar-feature-test-do-deploy2",
	}, {
		name:       "branch with slashes",
		repoOwner:  "foo",
		repo:       "bar",
		branchName: "feature/test",
		expected:   "bar-feature-test",
	}, {
		name:       "long repo owner",
		repoOwner:  "thisisanextremelylongrepohostname",
		repo:       "bar",
		branchName: "feature-branch",
		expected:   "bar-feature-branch",
	}, {
		name:       "long repo",
		repoOwner:  "foo",
		repo:       "thisisanextremelylongreponame",
		branchName: "feature-branch",
		expected:   "thisisanextremelylongreponame-feature-branch",
	}, {
		name:       "repo with hostname",
		repoOwner:  "foo",
		repo:       "my.domain.com",
		branchName: "feature-branch",
		expected:   "my-domain-com-feature-branch",
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := GenerateAppName(test.repoOwner, test.repo, test.branchName)
			require.Equal(t, test.expected, got)
		})
	}
}
