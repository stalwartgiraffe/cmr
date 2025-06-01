package gitlab

// see https://mholt.github.io/json-to-go/
type GroupModel struct {
	ID                             int    `json:"id"`
	Name                           string `json:"name"`
	Path                           string `json:"path"`
	Description                    string `json:"description"`
	Visibility                     string `json:"visibility"`
	ShareWithGroupLock             bool   `json:"share_with_group_lock"`
	RequireTwoFactorAuthentication bool   `json:"require_two_factor_authentication"`
	TwoFactorGracePeriod           int    `json:"two_factor_grace_period"`
	ProjectCreationLevel           string `json:"project_creation_level"`
	AutoDevopsEnabled              bool   `json:"auto_devops_enabled"`
	SubgroupCreationLevel          string `json:"subgroup_creation_level"`
	EmailsDisabled                 bool   `json:"emails_disabled"`
	EmailsEnabled                  bool   `json:"emails_enabled"`
	MentionsDisabled               bool   `json:"mentions_disabled"`
	LfsEnabled                     bool   `json:"lfs_enabled"`
	DefaultBranchProtection        int    `json:"default_branch_protection"`
	AvatarURL                      string `json:"avatar_url"`
	WebURL                         string `json:"web_url"`
	RequestAccessEnabled           bool   `json:"request_access_enabled"`
	RepositoryStorage              string `json:"repository_storage"`
	FullName                       string `json:"full_name"`
	FullPath                       string `json:"full_path"`
	FileTemplateProjectID          int    `json:"file_template_project_id"`
	ParentID                       int    `json:"parent_id"`
	CreatedAt                      Time   `json:"created_at"`
	IPRestrictionRanges            string `json:"ip_restriction_ranges"`
}
