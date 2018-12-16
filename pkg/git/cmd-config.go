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
func Config(entity *RepoEntity, options ConfigOptions) (value string, err error) {
	// here we configure config operation
	// default mode is go-git (this may be configured)
	configCmdMode = configCmdModeLegacy

	switch configCmdMode {
	case configCmdModeLegacy:
		value, err = configWithGit(entity, options)
		return value, err
	case configCmdModeNative:
		value, err = configWithGoGit(entity, options)
		return value, err
	}
	return value, errors.New("Unhandled config operation")
}

// configWithGit is simply a bare git commit -m <msg> command which is flexible
func configWithGit(entity *RepoEntity, options ConfigOptions) (value string, err error) {
	args := make([]string, 0)
	args = append(args, configCommand)
	if len(string(options.Site)) > 0 {
		args = append(args, "--"+string(options.Site))
	}
	args = append(args, "--get")
	args = append(args, options.Section+"."+options.Option)
	// parse options to command line arguments
	out, err := GenericGitCommandWithOutput(entity.AbsPath, args)
	if err != nil {
		return out, err
	}
	// till this step everything should be ok
	return out, nil
}

// commitWithGoGit is the primary commit method
func configWithGoGit(entity *RepoEntity, options ConfigOptions) (value string, err error) {
	// TODO: add global search
	config, err := entity.Repository.Config()
	if err != nil {
		return value, err
	}
	value = config.Raw.Section(options.Section).Option(options.Option)
	return value, nil
}

// AddConfig
func AddConfig(entity *RepoEntity, options ConfigOptions, value string) (err error) {
	err = addConfigWithGit(entity, options, value)
	return err

}

// addConfigWithGit is simply a bare git config --add <option> command which is flexible
func addConfigWithGit(entity *RepoEntity, options ConfigOptions, value string) (err error) {
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
	if err := GenericGitCommand(entity.AbsPath, args); err != nil {
		log.Warn("Error at git command (config)")
		return err
	}
	// till this step everything should be ok
	return entity.Refresh()
}
