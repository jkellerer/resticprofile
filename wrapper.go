package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/creativeprojects/resticprofile/constants"

	"github.com/creativeprojects/resticprofile/clog"
	"github.com/creativeprojects/resticprofile/config"
)

type resticWrapper struct {
	resticBinary string
	profile      *config.Profile
	moreArgs     []string
}

func newResticWrapper(resticBinary string, profile *config.Profile, moreArgs []string) *resticWrapper {
	return &resticWrapper{
		resticBinary: resticBinary,
		profile:      profile,
		moreArgs:     moreArgs,
	}
}

func (r *resticWrapper) runInitialize() error {
	rCommand := r.prepareCommand(constants.CommandInit)
	rCommand.displayStderr = false
	return runCommand(rCommand)
}

func (r *resticWrapper) runCleanup() {

}

func (r *resticWrapper) runCheck() {

}

func (r *resticWrapper) runCommand(command string) error {
	rCommand := r.prepareCommand(command)
	return runCommand(rCommand)
}

func (r *resticWrapper) prepareCommand(command string) commandDefinition {
	// place the restic command first, there are some flags not recognized otherwise (like --stdin)
	arguments := append([]string{command}, convertIntoArgs(r.profile.GetCommandFlags(command))...)

	if r.moreArgs != nil && len(r.moreArgs) > 0 {
		arguments = append(arguments, r.moreArgs...)
	}

	// Special case for backup command
	if command == constants.CommandBackup {
		arguments = append(arguments, r.profile.GetBackupSource()...)
	}

	env := append(os.Environ(), r.getEnvironment()...)

	clog.Debugf("Starting command: %s %s", r.resticBinary, strings.Join(arguments, " "))
	rCommand := newCommand(r.resticBinary, arguments, env)

	if command == constants.CommandBackup && r.profile.Backup.UseStdin {
		clog.Debug("Redirecting stdin to the backup")
		rCommand.useStdin = true
	}
	return rCommand
}

func (r *resticWrapper) getEnvironment() []string {
	if r.profile.Environment == nil || len(r.profile.Environment) == 0 {
		return nil
	}
	env := make([]string, len(r.profile.Environment))
	i := 0
	for key, value := range r.profile.Environment {
		// env variables are always uppercase
		key = strings.ToUpper(key)
		clog.Debugf("Setting up environment variable '%s'", key)
		env[i] = fmt.Sprintf("%s=%s", key, value)
		i++
	}
	return env
}

func convertIntoArgs(flags map[string][]string) []string {
	args := make([]string, 0)

	if flags == nil || len(flags) == 0 {
		return args
	}

	for key, values := range flags {
		if values == nil {
			continue
		}
		if len(values) == 0 {
			args = append(args, fmt.Sprintf("--%s", key))
			continue
		}
		for _, value := range values {
			args = append(args, fmt.Sprintf("--%s", key))
			if value != "" {
				if strings.Contains(value, " ") {
					// quote the string containing spaces
					value = fmt.Sprintf(`"%s"`, value)
				}
				args = append(args, value)
			}
		}
	}
	return args
}