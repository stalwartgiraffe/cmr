package gitlab

import (
	"os"
	"path/filepath"

	//"github.com/stalwartgiraffe/cmr/internal/gitlab"
	"gopkg.in/yaml.v3"
)

type ProjectModel struct {
	ID                   int      `json:"id" yaml:"id"`
	Description          string   `json:"description" yaml:"description"`
	DefaultBranch        string   `json:"default_branch" yaml:"default_branch"`
	Topics               []string `json:"topics" yaml:"topics"`
	Archived             bool     `json:"archived" yaml:"archived"`
	Visibility           string   `json:"visibility" yaml:"visibility"`
	SSHURLToRepo         string   `json:"ssh_url_to_repo" yaml:"ssh_url_to_repo"`
	HTTPURLToRepo        string   `json:"http_url_to_repo" yaml:"http_url_to_repo"`
	WebURL               string   `json:"web_url" yaml:"web_url"`
	Name                 string   `json:"name" yaml:"name"`
	NameWithNamespace    string   `json:"name_with_namespace" yaml:"name_with_namespace"`
	Path                 string   `json:"path" yaml:"path"`
	PathWithNamespace    string   `json:"path_with_namespace" yam:"path_with_namespace"`
	IssuesEnabled        bool     `json:"issues_enabled" yaml:"issues_enabled"`
	MergeRequestsEnabled bool     `json:"merge_requests_enabled" yaml:"merge_requests_enabled"`
	WikiEnabled          bool     `json:"wiki_enabled" yaml:"wiki_enabled"`
	JobsEnabled          bool     `json:"jobs_enabled" yaml:"jobs_enabled"`
	SnippetsEnabled      bool     `json:"snippets_enabled" yaml:"snippets_enabled"`
	CreatedAt            Time     `json:"created_at" yaml:"created_at"`
	LastActivityAt       Time     `json:"last_activity_at" yaml:"last_activity_at"`
	SharedRunnersEnabled bool     `json:"shared_runners_enabled" yaml:"shared_runners_enabled"`
	CreatorID            int      `json:"creator_id" yaml:"creator_id"`
	Namespace            struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Path string `json:"path"`
		Kind string `json:"kind"`
	} `json:"namespace"`
	AvatarURL            string       `json:"avatar_url" yaml:"avatar_url"`
	StarCount            int          `json:"star_count" yaml:"star_count"`
	ForksCount           int          `json:"forks_count" yaml:"forks_count"`
	OpenIssuesCount      int          `json:"open_issues_count" yaml:"open_issues_count"`
	PublicJobs           bool         `json:"public_jobs" yaml:"public_jobs"`
	SharedWithGroups     []GroupShare `json:"shared_with_groups" yaml:"shared_with_groups"`
	RequestAccessEnabled bool         `json:"request_access_enabled" yaml:"request_access_enabled"`
}

type GroupShare struct {
	GroupID          int    `json:"group_id" yaml:"group_id"`
	GroupName        string `json:"group_name" yaml:"group_name"`
	GroupFullPath    string `json:"group_full_path" yaml:"group_full_path"`
	GroupAccessLevel int    `json:"group_access_level" yaml:"group_access_level"`
	ExpiresAt        Time   `json:"expires_at" yaml:"expires_at"`
}

func RepoFilePath(home string, root string, project ProjectModel) string {
	var start string
	if len(root) == 0 {
		start = home
	} else if root[0] == '/' {
		start = root
	} else {
		start = filepath.Join(home, root)
	}
	dir := project.PathWithNamespace
	return filepath.Join(start, "repos", dir)
}

func ReadProjectsSlice(filepath string) ([]ProjectModel, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var projects []ProjectModel
	err = yaml.NewDecoder(file).Decode(&projects)
	if err != nil {
		return nil, err
	}
	return projects, nil
}

func ReadProjects() (map[int]ProjectModel, error) {
	projectsSlice, err := ReadProjectsSlice("ignore/projects.yaml")
	if err != nil {
		return nil, err
	}
	return MakeProjectMap(projectsSlice), nil
}

func MakeProjectMap(projectsSlice []ProjectModel) map[int]ProjectModel {
	projects := make(map[int]ProjectModel)
	for _, p := range projectsSlice {
		projects[p.ID] = p
	}
	return projects
}
