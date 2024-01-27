package fileexpander

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileExpander(t *testing.T) {
	_, err := os.Getwd()
	if err != nil {
		t.Error("failed to get the working directory")
	}
	var testParams = []struct {
		inputStr  string
		resultStr string
		errorStr  string
	}{
		{"--httpAuth.username=myadmin", "--httpAuth.username=myadmin", ""},
		{"--httpAuth.username=$__file{./testdata/username-1.txt}", "--httpAuth.username=test-user", ""},
		{"$__file{./testdata/username-1.txt}", "test-user", ""},
		{"--httpAuth.password=$__file{./testdata/password-1.txt}", "--httpAuth.password=test-password", ""},
		{"--httpAuth.password=$__file{./testdata/non-existing-file}", "", "unable to read the file \"./testdata/non-existing-file\""},
		{"--httpAuth.username=$__file{./testdata/username-1.txt} --httpAuth.password=$__file{./testdata/password-1.txt}", "--httpAuth.username=test-user --httpAuth.password=test-password", ""},
	}
	for _, param := range testParams {
		res, err := Expand(param.inputStr)
		if err != nil {
			assert.Equal(t, param.errorStr, err.Error())
			assert.Equal(t, param.resultStr, res)
		}
		assert.Equal(t, param.resultStr, res)
	}
}
