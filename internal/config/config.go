package config

import (
	"fmt"
	"io"
	"os"
	"path"
	"regexp"

	"github.com/spf13/viper"
)

type Config struct {
	Repos    MyRepos   `yaml:"repos"`
	Projects []Project `yaml:"projects"`
}

type MyRepos struct {
	Root string `yaml:"root"`
}

type Project struct {
	Name    string   `yaml:"name"`
	Linters []Linter `yaml:"linters,omitempty"`
}
type Linter struct {
	Name    string   `yaml:"name"`
	Args    string   `yaml:"args"`
	CmdArgs []string `yaml:"-"`
}

func LoadConfigFile(filepath string) (*Config, error) {
	if filepath == "" {
		// Find home directory.
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		filepath = path.Join(home, ".cmr.yaml")
	}
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	return LoadConfig(file)
}

// LoadConfig reads in config file and ENV variables if set.
func LoadConfig(file io.Reader) (*Config, error) {
	viper.AutomaticEnv() // read in environment variables that match

	// must tell viper what to unmarshal before we read the file
	viper.SetConfigType("yaml")

	var err error
	if err = viper.ReadConfig(file); err != nil {
		return nil, err
	}

	cfg := &Config{}
	err = viper.Unmarshal(cfg)
	if err != nil {
		return nil, err
	}

	// put the args into
	if err := cfg.parse(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) parse() error {
	for i := range c.Projects {
		if err := c.Projects[i].parse(); err != nil {
			return err
		}
	}
	return nil
}

func (p *Project) parse() error {

	if len(p.Name) < 1 {
		return fmt.Errorf("Project has empty name")
	}
	for i := range p.Linters {
		if err := p.Linters[i].Parse(); err != nil {
			return err
		}
	}
	return nil
}

func (l *Linter) Parse() error {
	if len(l.Args) < 1 {
		return nil
	}

	l.CmdArgs = splitCmdArgs(l.Args)

	for _, a := range l.CmdArgs {
		if err := verifyQuote(a); err != nil {
			return err
		}
	}
	return nil
}

var argsRE = regexp.MustCompile(`("[^"]*")|[^\s]+`)

func splitCmdArgs(s string) []string {
	return argsRE.FindAllString(s, -1)
}

var quotedRE = regexp.MustCompile(`^"[^"]*"$`)
var noQuoteRE = regexp.MustCompile(`^[^"]+$`)

func verifyQuote(s string) error {
	if len(s) < 1 {
		return nil
	}
	if quotedRE.MatchString(s) {
		return nil
	}
	if noQuoteRE.MatchString(s) {
		return nil
	}
	return fmt.Errorf("argument is improperly quoted:%s", s)
}
