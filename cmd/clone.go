package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/stalwartgiraffe/cmr/internal/elog"
	"github.com/stalwartgiraffe/cmr/internal/gitlab"
	"github.com/stalwartgiraffe/cmr/internal/gitutil"
)

func NewCloneCommand(cfg *CmdConfig) *cobra.Command {
	return &cobra.Command{
		Use:   "clone",
		Short: "clone git repos",
		Long:  `clone well known repos using bulk commands`,

		Args: func(cmd *cobra.Command, args []string) error {
			if 0 < len(args) {
				return fmt.Errorf("unexpected args %v", args)
			} else {
				return nil
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			logger := elog.New()

			logger.Println("so cloney")

			file, err := os.Open("ignore/projects.yaml")
			if err != nil {
				fmt.Println(err)
				return
			}
			defer file.Close()

			//var projects map[int]gitlab.ProjectModel
			var projects []gitlab.ProjectModel
			err = yaml.NewDecoder(file).Decode(&projects)
			if err != nil {
				fmt.Println(err)
				return
			}

			home, err := os.UserHomeDir()
			if err != nil {
				fmt.Println(err)
				return
			}
			// Now you can use the groups variable
			token, err := loadGitlabAuthToken()
			if err != nil {
				fmt.Println(err)
				return
			}
			known := wellKnownProjects()
			for _, project := range projects {
				if _, ok := known[project.PathWithNamespace]; !ok {
					continue
				}

				if err := Clone(cfg, project, home, token); err != nil {
					fmt.Println("ERROR", err)
				}
			}
		},
	}
}

func Clone(cfg *CmdConfig, project gitlab.ProjectModel, home string, token string) error {
	path := gitlab.RepoFilePath(
		home,
		cfg.Config.Repos.Root,
		project,
	)

	fmt.Println("cloning ", path)
	dot := filepath.Join(path, ".git")
	_, err := os.Stat(dot)
	if err == nil {
		fmt.Println(path, "exists. Skipping...")
		return nil
	}

	err = os.MkdirAll(path, os.ModeDir|0755)
	if err != nil {
		return err
	}

	return gitutil.Clone(path, project.HTTPURLToRepo, token, os.Stdout)
}

// to find dirs that contain git
// fd -H --no-ignore -t d '.git'  | grep -F -v '.gitlab' | sort
func wellKnownProjects() map[string]struct{} {
	return map[string]struct{}{
		"ad-registration/schema":                                {},
		"adserving/datagrid/profileservice":                     {},
		"app/phoenix/tool/deals/routetest":                      {},
		"build-images":                                          {},
		"build-images/debian-slim-git-curl-jq/":                 {},
		"continuous-integration/build-images":                   {},
		"exchange-node/billing":                                 {},
		"exchange-node/buyer-traffic-optimizer-lib":             {},
		"exchange-node/dash-hunter":                             {},
		"exchange-node/demand":                                  {},
		"exchange-node/deployment/demand-deployment":            {},
		"exchange-node/deployment/exchange-node-argocd-apps":    {},
		"exchange-node/deployment/exchange-node-deployment":     {},
		"exchange-node/deployment/exchange-pod-deployment":      {},
		"exchange-node/deployment/gauntlet-deployment":          {},
		"exchange-node/deployment/grafana-pyroscope-deployment": {},
		"exchange-node/encoding":                                {},
		"exchange-node/feature-toggles/features":                {},
		"exchange-node/feature-toggles/features-api":            {},
		"exchange-node/feature-toggles/features-deploy":         {},
		"exchange-node/feature-toggles/featureslib":             {},
		"exchange-node/gauntlet":                                {},
		"exchange-node/gitlab-ci-modules":                       {},
		"exchange-node/impression":                              {},
		"exchange-node/ixlib":                                   {},
		"exchange-node/kit/moneylib":                            {},
		"exchange-node/load-balancer-info":                      {},
		"exchange-node/machine-learning/rivr-go-library":        {},
		"exchange-node/privacy-sandbox-reporting":               {},
		"exchange-node/rules-api":                               {},
		"exchange-node/rules-lib":                               {},
		"exchange-node/rules-updater":                           {},
		"exchange-node/rules/monorepo":                          {},
		"exchange-node/schema":                                  {},
		"exchange-node/signal-management-lib":                   {},
		"exchange-node/supply":                                  {},
		"exchange-node/system-tests":                            {},
		"exchange-node/telemetry":                               {},
		"exchange-node/third-party/protobuf-go":                 {},
		"gauntlet/grafana-dashboards":                           {},
		"gauntlet/local-gauntlet":                               {},
		"ix/engineering/interviewing":                           {},
		"kit/moneylib":                                          {},
		"m8s/awd":                                               {},
		"machine-learning-optimization/inference":               {},
		"marvel-reference/bracelet":                             {},
		"marvel-reference/bracelet-env":                         {},
		"nomix/arc3":                                            {},
		"nomix/engineering":                                     {},
		"nomix/glossary":                                        {},
		"observability/":                                        {},
		"observability/grafana":                                 {},
		"observability/grafana-automation":                      {},
		"observability/grafana-automation/grafana-alerts":       {},
		"observability/grafana-automation/grafana-dashboards":   {},
		"observability/mimir-rules":                             {},
		"operations/ansible":                                    {},
		"operations/inventory":                                  {},
		"operations/k8s/argocd-resources":                       {},
		"operations/k8s/core-components":                        {},
		"operations/k8s/k8s-data-argocd-resources":              {},
		"operations/k8s/namespaces":                             {},
		"orchestration/charts":                                  {},
		"orchestration/pandora":                                 {},
		"pipelines/base-pipelines/all-deployment-repositories":  {},
		"pipelines/base-pipelines/ci-module-repositories":       {},
		"pipelines/base-pipelines/go/services":                  {},
		"pipelines/gitlab-ci-modules/mkdocs-gpages":             {},
		"platform-dev/AS":                                       {},
		"platform-dev/operations/inventory":                     {},
		"platform-test/AS":                                      {},
		"platform-test/lib":                                     {},
		"platform-tools/experiment-booker":                      {},
		"regulations/privacy-lib":                               {},
		"regulations/schema":                                    {},
		"tools/linting":                                         {},
	}
}
