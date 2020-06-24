package main

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/creativeprojects/resticprofile/config"
	"github.com/creativeprojects/resticprofile/systemd"
)

type ownCommand struct {
	name        string
	description string
	action      func(commandLineFlags, []string) error
}

var (
	ownCommands = []ownCommand{
		{
			name:        "profiles",
			description: "display profile names from the configuration file",
			action:      displayProfilesCommand,
		},
		{
			name:        "self-update",
			description: "update resticprofile to latest version (does not update restic)",
			action:      selfUpdate,
		},
		{
			name:        "systemd-unit",
			description: "create a user systemd timer",
			action:      createSystemdTimer,
		},
	}
)

func displayOwnCommands() {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	for _, command := range ownCommands {
		_, _ = fmt.Fprintf(w, "\t%s\t%s\n", command.name, command.description)
	}
	_ = w.Flush()
}

func isOwnCommand(command string) bool {
	for _, commandDef := range ownCommands {
		if commandDef.name == command {
			return true
		}
	}
	return false
}

func runOwnCommand(command string, flags commandLineFlags, args []string) error {
	for _, commandDef := range ownCommands {
		if commandDef.name == command {
			return commandDef.action(flags, args)
		}
	}
	return fmt.Errorf("command not found: %v", command)
}

func displayProfilesCommand(commandLineFlags, []string) error {
	displayProfiles()
	displayGroups()
	return nil
}

func displayProfiles() {
	profileSections := config.ProfileSections()
	if profileSections == nil || len(profileSections) == 0 {
		fmt.Println("\nThere's no available profile in the configuration")
	} else {
		fmt.Println("\nProfiles available:")
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		for name, sections := range profileSections {
			if sections == nil || len(sections) == 0 {
				_, _ = fmt.Fprintf(w, "\t%s:\t(n/a)\n", name)
			} else {
				_, _ = fmt.Fprintf(w, "\t%s:\t(%s)\n", name, strings.Join(sections, ", "))
			}
		}
		_ = w.Flush()
	}
	fmt.Println("")
}

func displayGroups() {
	groups := config.ProfileGroups()
	if groups == nil || len(groups) == 0 {
		return
	}
	fmt.Println("Groups available:")
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	for name, groupList := range groups {
		_, _ = fmt.Fprintf(w, "\t%s:\t%s\n", name, strings.Join(groupList, ", "))
	}
	_ = w.Flush()
	fmt.Println("")
}

func selfUpdate(flags commandLineFlags, args []string) error {
	err := confirmAndSelfUpdate(flags.verbose)
	if err != nil {
		return err
	}
	return nil
}

func createSystemdTimer(flags commandLineFlags, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("OnCalendar argument required")
	}
	systemd.Generate(flags.name, args[0])
	return nil
}