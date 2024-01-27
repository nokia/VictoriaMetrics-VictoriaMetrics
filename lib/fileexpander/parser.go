package fileexpander

import (
	"flag"
	"log"
	"slices"
	"strings"
)

var enableFileExpander = flag.Bool("enableFileExpander", false, "enables all the cli flags set of file expander format `$__file{<path>}` to be expanded with the contents of the respective flag.")

var reloadableFlagsList []string = []string{}

// AppendToReloadableFlagList appends the name of the flag to reloadableFlagsList
func AppendToReloadableFlagList(name string) {
	reloadableFlagsList = append(reloadableFlagsList, name)
}

func Parse(args []string) []string {
	return expandArgs(args)
}

func expandArgs(args []string) []string {
	expandedArgs := make([]string, len(args))
	for _, arg := range args {
		splittedArg := strings.Split(arg, "=")
		// do not expand right away in case the flag is reloadable
		if slices.Contains(reloadableFlagsList, splittedArg[0]) {
			expandedArgs = append(expandedArgs, arg)
			continue
		}
		expArg, err := Expand(arg)
		if err != nil {
			log.Fatalf("Ran into the error: %s", err.Error())
		}
		expandedArgs = append(expandedArgs, string(expArg))
	}
	return expandedArgs
}
