package git

import (
	"errors"

	log "github.com/sirupsen/logrus"
)

var (
	configCmdMode string

	configCommand       = "config"
	configCmdModeLegacy = "git"
	configCmdModeNative = "go-git"
)

// ConfigOptions defines the rules for commit operation
type ConfigOptions struct {
	// Section
	Section string
	// Option
	Option string
	// Site should be Global or Local
	Site ConfigSite
}

// ConfigSite defines a string type for the site.
type ConfigSite string

const (
	// ConfigSiteLocal defines a local config.
	ConfigSiteLocal ConfigSite = "local"

	// ConfgiSiteGlobal defines a global config.
	ConfgiSiteGlobal ConfigSite = "global"
)

// Config
func Config(e *RepoEntity, options ConfigOptions) (value string, err error) {
	// here we configure config operation
	// default mode is go-git (this may be configured)
	configCmdMode = configCmdModeLegacy

	switch configCmdMode {
	case configCmdModeLegacy:
		return configWithGit(e, options)
	case configCmdModeNative:
		return configWithGoGit(e, options)
	}
	return value, errors.New("Unhandled config operation")
}

// configWithGit is simply a bare git commit -m <msg> command which is flexible
func configWithGit(e *RepoEntity, options ConfigOptions) (value string, err error) {
	args := make([]string, 0)
	args = append(args, configCommand)
	if len(string(options.Site)) > 0 {
		args = append(args, "--"+string(options.Site))
	}
	args = append(args, "--get")
	args = append(args, options.Section+"."+options.Option)
	// parse options to command line arguments
	out, err := GenericGitCommandWithOutput(e.AbsPath, args)
	if err != nil {
		return out, err
	}
	// till this step everything should be ok
	return out, nil
}

// commitWithGoGit is the primary commit method
func configWithGoGit(e *RepoEntity, options ConfigOptions) (value string, err error) {
	// TODO: add global search
	config, err := e.Repository.Config()
	if err != nil {
		return value, err
	}
	return config.Raw.Section(options.Section).Option(options.Option), nil
}

// AddConfig adds an entry on the ConfigOptions field.
func AddConfig(e *RepoEntity, options ConfigOptions, value string) (err error) {
	return addConfigWithGit(e, options, value)

}

// addConfigWithGit is simply a bare git config --add <option> command which is flexible
func addConfigWithGit(e *RepoEntity, options ConfigOptions, value string) (err error) {
	args := make([]string, 0)
	args = append(args, configCommand)
	if len(string(options.Site)) > 0 {
		args = append(args, "--"+string(options.Site))
	}
	args = append(args, "--add")
	args = append(args, options.Section+"."+options.Option)
	if len(value) > 0 {
		args = append(args, value)
	}
	if err := GenericGitCommand(e.AbsPath, args); err != nil {
		log.Warn("Error at git command (config)")
		return err
	}
	// till this step everything should be ok
	return e.Refresh()
}
