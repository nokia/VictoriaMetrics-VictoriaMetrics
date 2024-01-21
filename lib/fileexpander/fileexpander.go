package fileexpander

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"
	"unicode"
)

const FILE_PROVIDER_REGEX = `\$__file\{(.*)\}`

var regexVar = regexp.MustCompile(FILE_PROVIDER_REGEX)

func Expand(exp string) (string, error) {
	if !strings.Contains(exp, "$__file{") {
		return exp, nil
	}
	if regexVar.MatchString(exp) {
		submatches := regexVar.FindAllStringSubmatch(exp, -1)
		if len(submatches) > 1 {
			return "", fmt.Errorf("more than one match found for file provider expression ($__file{}) in the arguments.")
		}
		if submatches != nil && len(submatches[0]) > 1 && submatches[0][1] != "" {
			a, err := readContentFromFile(submatches[0][1])
			if err != nil {
				return "", fmt.Errorf("error reading the file %s. reason:%s", submatches[0][1], err.Error())
			}
			data := []byte(exp)
			data = bytes.Replace(data, []byte(submatches[0][0]), []byte(a), -1)
			return string(data), nil
		}
	}
	return exp, nil
}

func readContentFromFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	pass := strings.TrimRightFunc(string(data), unicode.IsSpace)
	return pass, nil
}
