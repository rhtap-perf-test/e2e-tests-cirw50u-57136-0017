package build

import (
	"fmt"
	"os"
	"strconv"

	"github.com/konflux-ci/e2e-tests/pkg/clients/github"
	"github.com/konflux-ci/e2e-tests/pkg/constants"
	"github.com/konflux-ci/e2e-tests/pkg/utils"
)

// resolve the git url and revision from a pull request. If not found, return a default
// that is set from environment variables.
func ResolveGitDetails(repoUrlENV, repoRevisionENV string) (string, string, error) {
	defaultGitURL := fmt.Sprintf("https://github.com/%s/%s", constants.DEFAULT_GITHUB_BUILD_ORG, constants.DEFAULT_GITHUB_BUILD_REPO)
	defaultGitRevision := "main"
	// If we are testing the changes from a pull request, APP_SUFFIX may contain the
	// pull request ID. If it looks like an ID, then fetch information about the pull
	// request and use it to determine which git URL and revision to use for the EC
	// pipelines. NOTE: This is a workaround until Pipeline as Code supports passing
	// the source repo URL: https://issues.redhat.com/browse/SRVKP-3427. Once that's
	// implemented, remove the APP_SUFFIX support below and simply rely on the other
	// environment variables to set the git revision and URL directly.
	appSuffix := os.Getenv("APP_SUFFIX")
	if pullRequestID, err := strconv.ParseInt(appSuffix, 10, 64); err == nil {
		gh, err := github.NewGithubClient(utils.GetEnv(constants.GITHUB_TOKEN_ENV, ""), constants.DEFAULT_GITHUB_BUILD_ORG)
		if err != nil {
			return "", "", err
		}
		return gh.GetPRDetails(constants.DEFAULT_GITHUB_BUILD_REPO, int(pullRequestID))

	}
	return utils.GetEnv(repoUrlENV, defaultGitURL), utils.GetEnv(repoRevisionENV, defaultGitRevision), nil
}
