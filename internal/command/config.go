package command

import (
	"fmt"

	"github.com/isacikgoz/gitbatch/internal/git"
)

// ConfigOptions defines the rules for commit operation
type ConfigOptions struct {
	// Section
	Section string
	// Option
	Option string
	// Site should be Global or Local
	Site ConfigSite
	// Mode is the command mode
	CommandMode Mode
}

// ConfigSite defines a string type for the site.
type ConfigSite string

const (
	// ConfigSiteLocal defines a local config.
	ConfigSiteLocal ConfigSite = "local"

	// ConfigSiteGlobal defines a global config.
	ConfigSiteGlobal ConfigSite = "global"
)

// Config adds or reads config of a repository
func Config(r *git.Repository, o *ConfigOptions) (value string, err error) {
	// here we configure config operation

	switch o.CommandMode {
	case ModeLegacy:
		return configWithGit(r, o)
	case ModeNative:
		return configWithGoGit(r, o)
	}
	return value, fmt.Errorf("unhandled config operation")
}

// configWithGit is simply a bare git config --site <option>.<section> command which is flexible
func configWithGit(r *git.Repository, options *ConfigOptions) (value string, err error) {
	args := make([]string, 0)
	args = append(args, "config")
	if len(string(options.Site)) > 0 {
		args = append(args, "--"+string(options.Site))
	}
	args = append(args, "--get")
	args = append(args, options.Section+"."+options.Option)
	// parse options to command line arguments
	out, err := Run(r.AbsPath, "git", args)
	if err != nil {
		return out, err
	}
	// till this step everything should be ok
	return out, nil
}

// commitWithGoGit is the primary commit method
func configWithGoGit(r *git.Repository, options *ConfigOptions) (value string, err error) {
	// TODO: add global search
	config, err := r.Repo.Config()
	if err != nil {
		return value, err
	}
	return config.Raw.Section(options.Section).Option(options.Option), nil
}

// AddConfig adds an entry on the ConfigOptions field.
func AddConfig(r *git.Repository, options *ConfigOptions, value string) (err error) {
	return addConfigWithGit(r, options, value)

}

// addConfigWithGit is simply a bare git config --add <option> command which is flexible
func addConfigWithGit(r *git.Repository, options *ConfigOptions, value string) (err error) {
	args := make([]string, 0)
	args = append(args, "config")
	if len(string(options.Site)) > 0 {
		args = append(args, "--"+string(options.Site))
	}
	args = append(args, "--add")
	args = append(args, options.Section+"."+options.Option)
	if len(value) > 0 {
		args = append(args, value)
	}
	if _, err := Run(r.AbsPath, "git", args); err != nil {
		return err
	}
	// till this step everything should be ok
	return r.Refresh()
}
