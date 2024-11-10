package main

import (
	"envious/tools"
	"flag"
	"os"
)

func main() {
	var tokens []tools.Token = tools.ParseIniFile("my.ini")
	var profiles []tools.Profile = tools.BuildProfiles(tokens)

	enviousLsCmd := flag.NewFlagSet("ls", flag.ExitOnError)
	detailed := enviousLsCmd.Bool("details", false, "details")

	enviousUseCmd := flag.NewFlagSet("use", flag.ExitOnError)

	if len(os.Args) == 1 {
		tools.UseDefaultProfile(profiles)
		return
	}

	switch os.Args[1] {
	case "ls":
		enviousLsCmd.Parse(os.Args[2:])
		tools.ListProfiles(profiles, detailed)

	case "use":
		enviousUseCmd.Parse(os.Args[2:])
		tools.UseProfile(profiles, os.Args[2])
	}
}
