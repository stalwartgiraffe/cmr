package gitlab

import (
	"fmt"

	"github.com/stalwartgiraffe/cmr/internal/utils"
	"github.com/stalwartgiraffe/cmr/kam"
)

type Group struct {
	Name     string
	FullPath string
	GroupID  int
}

func getVal[T any](m kam.Map, k string, t *T) bool {
	anyv, ok := m[k]
	if ok {
		v, ok := anyv.(T)
		if ok {
			*t = v
			return true
		}
	}
	return false
}

func makeGroup(m kam.Map) (Group, error) {
	g := Group{}
	if !getVal(m, "name", &g.Name) ||
		!getVal(m, "id", &g.GroupID) ||
		!getVal(m, "full_path", &g.FullPath) {
		return Group{}, fmt.Errorf("make missing expected group attribute")
	}

	return g, nil
}

type GetGroupsError struct {
	err  error
	page kam.JSONValue
}

func (e *GetGroupsError) Error() string {
	return fmt.Sprintf("%s\n%s",
		e.err.Error(),
		utils.YamlString(e.page),
	)
}

func GetGroupsFromPage(page kam.JSONValue) ([]Group, error) {
	if page.AnyVal == nil {
		return nil, &GetGroupsError{fmt.Errorf(" page has empty body"), page}
	}

	ar, ok := page.AnyVal.(kam.Array)
	if !ok {
		return nil, &GetGroupsError{fmt.Errorf("page missing expected array"), page}
	}

	groups := []Group{}
	for _, a := range ar {
		m, ok := a.(map[string]any)
		if !ok {
			return nil, &GetGroupsError{fmt.Errorf("array element not a kam.Map"), page}
		}
		g, err := makeGroup(m)
		if err != nil {
			return nil, &GetGroupsError{err, page}
		}
		groups = append(groups, g)
	}
	return groups, nil
}

type GetProjectError struct {
	err  error
	page kam.JSONValue
}

func (e *GetProjectError) Error() string {
	return fmt.Sprintf("%s\n%s",
		e.err.Error(),
		utils.YamlString(e.page),
	)
}

func GetProjectFromPage(page kam.JSONValue) ([]Project, error) {
	if page.AnyVal == nil {
		return nil, &GetProjectError{fmt.Errorf("page has empty body"), page}
	}

	ar, ok := page.AnyVal.(kam.Array)
	if !ok {
		return nil, &GetProjectError{fmt.Errorf("page missing expected array"), page}
	}

	projects := []Project{}
	for _, a := range ar {
		m, ok := a.(map[string]any)
		if !ok {
			return nil, &GetProjectError{fmt.Errorf("array element not a kam.Map"), page}
		}
		p, err := makeProject(m)
		if err != nil {
			return nil, &GetProjectError{err, page}
		}
		projects = append(projects, p)
	}
	return projects, nil
}

type ProjectNamespace struct {
	AvatarUrl string
	ID        int
	Kind      string
	ParentID  int
	FullPath  string
	Path      string
	Name      string
	WebUrl    string
}

func makeProjectNamespace(m kam.Map) (ProjectNamespace, error) {
	n := ProjectNamespace{}
	if !getVal(m, "id", &n.ID) {
		return ProjectNamespace{}, fmt.Errorf("namespace is missing id")
	}
	if !getVal(m, "name", &n.Name) {
		return ProjectNamespace{}, fmt.Errorf("namespace is missing name")
	}

	getVal(m, "avatar_url", &n.AvatarUrl)
	getVal(m, "full_path", &n.FullPath)
	getVal(m, "kind", &n.Kind)
	getVal(m, "parent_id", &n.ParentID)
	getVal(m, "path", &n.Path)
	getVal(m, "web_url", &n.WebUrl)
	return n, nil
}

type Project struct {
	WebUrl        string
	SshUrlToRepo  string
	HttpUrlToRepo string

	Namespace ProjectNamespace
}

func makeProject(m kam.Map) (Project, error) {
	p := Project{}
	if !getVal(m, "web_url", &p.WebUrl) ||
		!getVal(m, "ssh_url_to_repo", &p.SshUrlToRepo) ||
		!getVal(m, "http_url_to_repo", &p.HttpUrlToRepo) {

		return Project{}, fmt.Errorf("make missing expected group attribute")
	}
	var namespaceMap map[string]any //  kam.Map does not work as expected with type assertion
	if getVal(m, "namespace", &namespaceMap) {
		var err error
		p.Namespace, err = makeProjectNamespace(namespaceMap)
		if err != nil {
			fmt.Println("make missing group.namespace attribute")
		}
	} else {
		fmt.Println("no namespace")
	}

	return p, nil
}
