package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/go-git/go-git/v5"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var (
	gitlab_host  string
	gitlab_token string
	code_path    string
)

func init() {
	flag.StringVar(&gitlab_host, "host", "https://git.renzhikeji.com", "gitlab_host 请设置")
	flag.StringVar(&gitlab_token, "token", "eJAfTmyyyVYxxxtH5jsse", "gitlab_token 请设置")
	flag.StringVar(&code_path, "path", "/Users/renzhikeji/logs", "path 代码下载路径请设置")
}

func main() {
	flag.Parse()

	groups, err := allGroups()
	if err != nil {
		panic(err)
	}

	for _, group := range groups {
		marshal, _ := json.Marshal(&group)
		fmt.Println()
		fmt.Println()
		fmt.Println()
		Info("=========项目组信息=================")
		fmt.Printf("%s", marshal)
		fmt.Println()
	}

	fmt.Println()
	fmt.Println()
	Info("=========项目信息收集中，请稍等=================")
	allProjectInfos, err := allProjectInfo(groups)
	if err != nil {
		panic(err)
	}

	fmt.Println()
	fmt.Println()
	fmt.Println()
	fmt.Println()
	Info("=========项目信息=================")
	Info("项目总数: ", len(allProjectInfos))

	fmt.Println()
	fmt.Println()

	fmt.Println()
	fmt.Println()

	for _, info := range allProjectInfos {
		fmt.Println()

		dir := strings.Join([]string{code_path, info.Name}, "/")
		Info("git clone %s   %s  ......................", info.HTTPURLToRepo, dir)

		_, err := git.PlainClone(dir, false, &git.CloneOptions{
			// The intended use of a GitHub personal access token is in replace of your password
			// because access tokens can easily be revoked.
			// https://help.github.com/articles/creating-a-personal-access-token-for-the-command-line/
			Auth: &githttp.BasicAuth{
				Username: "abc123", // yes, this can be anything except an empty string
				Password: gitlab_token,
			},
			URL:      info.HTTPURLToRepo,
			Progress: os.Stdout,
		})

		Error(err)
	}
}

type GitGropInfo struct {
	AvatarURL            string      `json:"avatar_url"`
	Description          string      `json:"description"`
	FullName             string      `json:"full_name"`
	FullPath             string      `json:"full_path"`
	ID                   int64       `json:"id"`
	LdapAccess           interface{} `json:"ldap_access"`
	LdapCn               interface{} `json:"ldap_cn"`
	LfsEnabled           bool        `json:"lfs_enabled"`
	Name                 string      `json:"name"`
	ParentID             int64       `json:"parent_id"`
	Path                 string      `json:"path"`
	RequestAccessEnabled bool        `json:"request_access_enabled"`
	Visibility           string      `json:"visibility"`
	WebURL               string      `json:"web_url"`
}

func allGroups() ([]GitGropInfo, error) {

	request, err := http.NewRequest(http.MethodGet, strings.Join([]string{gitlab_host, "api/v4/groups"}, "/"), nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("PRIVATE-TOKEN", gitlab_token)
	query := request.URL.Query()
	query.Set("per_page", "50")

	allGroups := make([]GitGropInfo, 0, 100)
	pageNum := 1

	for {
		query.Set("page", strconv.Itoa(pageNum))
		request.URL.RawQuery = query.Encode()

		response, err := http.DefaultClient.Do(request)
		pageNum++

		if err != nil {
			return nil, err
		}

		if response.StatusCode != 200 {
			Warning("返回code 不是 200", request)
			continue
		}

		body, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}

		var gitGropInfo []GitGropInfo
		err = json.Unmarshal(body, &gitGropInfo)

		if len(gitGropInfo) == 0 {
			break
		}

		allGroups = append(allGroups, gitGropInfo...)

		if err != nil {
			return nil, err
		}

		response.Body.Close()

	}

	return allGroups, nil
}

type ProjectInfo struct {
	Links struct {
		Events        string `json:"events"`
		Issues        string `json:"issues"`
		Labels        string `json:"labels"`
		Members       string `json:"members"`
		MergeRequests string `json:"merge_requests"`
		RepoBranches  string `json:"repo_branches"`
		Self          string `json:"self"`
	} `json:"_links"`
	ApprovalsBeforeMerge       int64       `json:"approvals_before_merge"`
	Archived                   bool        `json:"archived"`
	AutoCancelPendingPipelines string      `json:"auto_cancel_pending_pipelines"`
	AutoDevopsDeployStrategy   string      `json:"auto_devops_deploy_strategy"`
	AutoDevopsEnabled          bool        `json:"auto_devops_enabled"`
	AvatarURL                  interface{} `json:"avatar_url"`
	BuildCoverageRegex         interface{} `json:"build_coverage_regex"`
	BuildTimeout               int64       `json:"build_timeout"`
	BuildsAccessLevel          string      `json:"builds_access_level"`
	CiConfigPath               interface{} `json:"ci_config_path"`
	CiDefaultGitDepth          interface{} `json:"ci_default_git_depth"`
	ContainerRegistryEnabled   bool        `json:"container_registry_enabled"`
	CreatedAt                  string      `json:"created_at"`
	CreatorID                  int64       `json:"creator_id"`
	DefaultBranch              string      `json:"default_branch"`
	Description                string      `json:"description"`
	EmptyRepo                  bool        `json:"empty_repo"`
	ForksCount                 int64       `json:"forks_count"`
	HTTPURLToRepo              string      `json:"http_url_to_repo"`
	ID                         int64       `json:"id"`
	ImportStatus               string      `json:"import_status"`
	IssuesAccessLevel          string      `json:"issues_access_level"`
	IssuesEnabled              bool        `json:"issues_enabled"`
	JobsEnabled                bool        `json:"jobs_enabled"`
	LastActivityAt             string      `json:"last_activity_at"`
	LfsEnabled                 bool        `json:"lfs_enabled"`
	MergeMethod                string      `json:"merge_method"`
	MergeRequestsAccessLevel   string      `json:"merge_requests_access_level"`
	MergeRequestsEnabled       bool        `json:"merge_requests_enabled"`
	Mirror                     bool        `json:"mirror"`
	Name                       string      `json:"name"`
	NameWithNamespace          string      `json:"name_with_namespace"`
	Namespace                  struct {
		AvatarURL interface{} `json:"avatar_url"`
		FullPath  string      `json:"full_path"`
		ID        int64       `json:"id"`
		Kind      string      `json:"kind"`
		Name      string      `json:"name"`
		ParentID  interface{} `json:"parent_id"`
		Path      string      `json:"path"`
		WebURL    string      `json:"web_url"`
	} `json:"namespace"`
	OnlyAllowMergeIfAllDiscussionsAreResolved bool   `json:"only_allow_merge_if_all_discussions_are_resolved"`
	OnlyAllowMergeIfPipelineSucceeds          bool   `json:"only_allow_merge_if_pipeline_succeeds"`
	OpenIssuesCount                           int64  `json:"open_issues_count"`
	PackagesEnabled                           bool   `json:"packages_enabled"`
	Path                                      string `json:"path"`
	PathWithNamespace                         string `json:"path_with_namespace"`
	PrintingMergeRequestLinkEnabled           bool   `json:"printing_merge_request_link_enabled"`
	PublicJobs                                bool   `json:"public_jobs"`
	ReadmeURL                                 string `json:"readme_url"`
	RepositoryAccessLevel                     string `json:"repository_access_level"`
	RequestAccessEnabled                      bool   `json:"request_access_enabled"`
	ResolveOutdatedDiffDiscussions            bool   `json:"resolve_outdated_diff_discussions"`
	SharedRunnersEnabled                      bool   `json:"shared_runners_enabled"`
	SharedWithGroups                          []struct {
		ExpiresAt        interface{} `json:"expires_at"`
		GroupAccessLevel int64       `json:"group_access_level"`
		GroupFullPath    string      `json:"group_full_path"`
		GroupID          int64       `json:"group_id"`
		GroupName        string      `json:"group_name"`
	} `json:"shared_with_groups"`
	SnippetsAccessLevel string        `json:"snippets_access_level"`
	SnippetsEnabled     bool          `json:"snippets_enabled"`
	SSHURLToRepo        string        `json:"ssh_url_to_repo"`
	StarCount           int64         `json:"star_count"`
	TagList             []interface{} `json:"tag_list"`
	Visibility          string        `json:"visibility"`
	WebURL              string        `json:"web_url"`
	WikiAccessLevel     string        `json:"wiki_access_level"`
	WikiEnabled         bool          `json:"wiki_enabled"`
}

func allProjectInfo(allGroups []GitGropInfo) ([]ProjectInfo, error) {
	allProjectInfos := make([]ProjectInfo, 0, 200)
	if len(allGroups) == 0 {
		return allProjectInfos, nil
	}

	for _, group := range allGroups {

		request, err := http.NewRequest(http.MethodGet, strings.Join([]string{gitlab_host, "api/v4/groups", strconv.FormatInt(group.ID, 10), "projects"}, "/"), nil)
		if err != nil {
			return nil, err
		}

		request.Header.Set("PRIVATE-TOKEN", gitlab_token)
		query := request.URL.Query()
		query.Set("per_page", "50")

		pageNum := 1
		for {

			query.Set("page", strconv.Itoa(pageNum))
			request.URL.RawQuery = query.Encode()

			response, err := http.DefaultClient.Do(request)
			if err != nil {
				return nil, err
			}

			if response.StatusCode != 200 {
				panic("error")
			}

			body, err := io.ReadAll(response.Body)
			response.Body.Close()
			if err != nil {
				return nil, err
			}

			var projectInfos []ProjectInfo

			err = json.Unmarshal(body, &projectInfos)
			if err != nil {
				return nil, err
			}

			if len(projectInfos) == 0 {
				break
			}

			allProjectInfos = append(allProjectInfos, projectInfos...)

			pageNum++
		}
	}

	return allProjectInfos, nil
}

func CheckIfError(err error) {
	if err == nil {
		return
	}

	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	panic(err)
}

func Error(err error) {
	if err == nil {
		return
	}

	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
}

func Info(format string, args ...interface{}) {
	fmt.Printf("\x1b[34;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}

func Warning(format string, args ...interface{}) {
	fmt.Printf("\x1b[36;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}
