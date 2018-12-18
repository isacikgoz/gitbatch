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

// CommitOptions defines the rules for commit operation
type ConfigOptions struct {
	// Section
	Section string
	// Option
	Option string
	// Site should be Global or Local
	Site ConfigSite
}

type ConfigSite string

const (
	// ConfigStieLocal
	ConfigSiteLocal ConfigSite = "local"
	// ConfgiSiteGlobal
	ConfgiSiteGlobal ConfigSite = "global"
)

// Config
func Config(e *RepoEntity, options ConfigOptions) (value string, err error) {
	// here we configure config operation
	// default mode is go-git (this may be configured)
	configCmdMode = configCmdModeLegacy

	switch configCmdMode {
	case configCmdModeLegacy:
		value, err = configWithGit(e, options)
		return value, err
	case configCmdModeNative:
		value, err = configWithGoGit(e, options)
		return value, err
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

// AddConfig
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
