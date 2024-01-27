package fileexpander

import (
	"fmt"
	"io"
	"strings"
	"unicode"

	"github.com/VictoriaMetrics/VictoriaMetrics/lib/fs"
	"github.com/valyala/fasttemplate"
)

func Expand(exp string) (string, error) {
	return expand(exp)
}

func readContentFromFile(path string) (string, error) {
	data, err := fs.ReadFileOrHTTP(path)
	if err != nil {
		return "", err
	}
	pass := strings.TrimRightFunc(string(data), unicode.IsSpace)
	return pass, nil
}

func expand(s string) (string, error) {
	if !strings.Contains(s, "$__file{") {
		return s, nil
	}
	result, err := fasttemplate.ExecuteFuncStringWithErr(s, "$__file{", "}", func(w io.Writer, tag string) (int, error) {
		if tag == "" {
			return 0, fmt.Errorf("file path cannot be empty under file provider expression %q", s)
		}
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
