package fileexpander

import (
	"flag"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
)

var dynamicFlagsList []string = []string{}

func AppendToDynamicFlagList(name string) {
	dynamicFlagsList = append(dynamicFlagsList, name)
}

func Parse() {
	ParseFlagSet(flag.CommandLine, os.Args[1:])
}

// ParseFlagSet parses the given args into the given fs.
func ParseFlagSet(fs *flag.FlagSet, args []string) {
	fmt.Printf("Args Before: %v\n\n", args)
	args = expandArgs(args)
	fmt.Println(args)
	if err := fs.Parse(args); err != nil {
		// Do not use lib/logger here, since it is uninitialized yet.
		log.Fatalf("cannot parse flags %q: %s", args, err)
	}
	if fs.NArg() > 0 {
		// See https://github.com/VictoriaMetrics/VictoriaMetrics/issues/4845
		log.Fatalf("unprocessed command-line args left: %s; the most likely reason is missing `=` between boolean flag name and value; "+
			"see https://pkg.go.dev/flag#hdr-Command_line_flag_syntax", fs.Args())
	}
}

func expandArgs(args []string) []string {
	for i, arg := range args {
		splittedArg := strings.Split(arg, "=")
		if slices.Contains(dynamicFlagsList, splittedArg[0]) {
			continue
		}
		expandedArg, err := Expand(arg)
		if err != nil {
			log.Fatalf("Ran into the error: %s", err.Error())
		}
		args[i] = string(expandedArg)
	}
	return args
}
