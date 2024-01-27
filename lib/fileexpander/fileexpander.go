package fileexpander

import (
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"

	"github.com/valyala/fasttemplate"
)

func Expand(exp string) (string, error) {
	return expand(exp)
}

func readContentFromFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	truncData := strings.TrimRightFunc(string(data), unicode.IsSpace)
	return truncData, nil
}

func expand(s string) (string, error) {
	if !strings.Contains(s, "$__file{") {
		return s, nil
	}
	result, err := fasttemplate.ExecuteFuncStringWithErr(s, "$__file{", "}", func(w io.Writer, tag string) (int, error) {
		data, err := readContentFromFile(tag)
		if err != nil {
			return 0, fmt.Errorf("unable to read the file %q", tag)
		}
		return fmt.Fprintf(w, "%s", data)
	})
	if err != nil {
		return "", err
	}
	return result, nil
}
