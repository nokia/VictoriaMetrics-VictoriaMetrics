package fileexpander

import (
	"log"
	"slices"
	"strings"

	"github.com/VictoriaMetrics/VictoriaMetrics/lib/flagutil"
)

func Parse(args []string) []string {
	return expandArgs(args)
}

func expandArgs(args []string) []string {
	expandedArgs := make([]string, len(args))
	for _, arg := range args {
		splittedArg := strings.Split(arg, "=")
		// do not expand right away in case the flag is reloadable
		if slices.Contains(flagutil.ReloadableFlagsList, strings.TrimLeft(splittedArg[0], "-")) {
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
