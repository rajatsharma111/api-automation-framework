package helpers

import (
	"encoding/json"
	"fmt"
	"github.com/cldcvr/api-automation-framework/report"
	"github.com/gavv/httpexpect"
	"github.com/stretchr/testify/assert"

	"io/ioutil"
	"os"
	"strings"
	"testing"
)

// Equal is a custome assertion which assert the condition and report them if fails or skips
func Equal(t *testing.T, actual, expected interface{}, msgAndArgs ...interface{}) {
	if assert.Equal(t, expected, actual, msgAndArgs...) {
		return
	}
	message := fmt.Sprintf("Not equal: \n"+
		"expected: %#v\n"+
		"actual  : %#v", expected, actual)
	report.Fail(t, strings.Join(assert.CallerInfo(), "\n\t\t\t"), messageFromMsgAndArgs(msgAndArgs), message)
	t.FailNow()
}

// NotNil is a custome assertion which assert the condition and report them if fails or skips
func NotNil(t *testing.T, expected interface{}, msgAndArgs ...interface{}) {
	if assert.NotNil(t, expected, msgAndArgs...) {
		return
	}
	message := fmt.Sprintf("Not nil: \n %#v", expected)
	report.Fail(t, strings.Join(assert.CallerInfo(), "\n\t\t\t"), messageFromMsgAndArgs(msgAndArgs), message)
	t.FailNow()
}

// Nil is a custome assertion which assert the condition and report them if fails or skips
func Nil(t *testing.T, expected interface{}, msgAndArgs ...interface{}) {
	if assert.True(t, (expected == nil || expected == ""), msgAndArgs...) {
		return
	}
	message := fmt.Sprintf("Nil: \n %#v", expected)
	report.Fail(t, strings.Join(assert.CallerInfo(), "\n\t\t\t"), messageFromMsgAndArgs(msgAndArgs), message)
	t.FailNow()
}

// Contains is a custome assertion which assert the condition and report them if fails or skips
func Contains(t *testing.T, actual, expected interface{}, msgAndArgs ...interface{}) {
	if assert.Contains(t, actual, expected, msgAndArgs...) {
		return
	}
	message := fmt.Sprintf("Not contains: \n"+
		"expected: %#v\n"+
		"actual  : %#v", expected, actual)
	report.Fail(t, strings.Join(assert.CallerInfo(), "\n\t\t\t"), messageFromMsgAndArgs(msgAndArgs), message)
	t.FailNow()
}

// False is a custome assertion which assert the condition and report them if fails or skips
func False(t *testing.T, actual bool, msgAndArgs ...interface{}) {
	if assert.False(t, actual, msgAndArgs...) {
		return
	}
	message := fmt.Sprintf("Not false: \n"+
		"actual  : %#v", actual)
	report.Fail(t, strings.Join(assert.CallerInfo(), "\n\t\t\t"), messageFromMsgAndArgs(msgAndArgs), message)
	t.FailNow()
}

// True is a custome assertion which assert the condition and report them if fails or skips
func True(t *testing.T, actual bool, msgAndArgs ...interface{}) {
	if assert.True(t, actual, msgAndArgs...) {
		return
	}
	message := fmt.Sprintf("Not false: \n"+
		"actual  : %#v", actual)
	report.Fail(t, strings.Join(assert.CallerInfo(), "\n\t\t\t"), messageFromMsgAndArgs(msgAndArgs), message)
	t.FailNow()
}

func messageFromMsgAndArgs(msgAndArgs ...interface{}) string {
	if len(msgAndArgs) == 0 || msgAndArgs == nil {
		return ""
	}
	if len(msgAndArgs) == 1 {
		msg := msgAndArgs[0]
		if msgAsStr, ok := msg.(string); ok {
			return msgAsStr
		}
		return fmt.Sprintf("%+v", msg)
	}
	if len(msgAndArgs) > 1 {
		return fmt.Sprintf(msgAndArgs[0].(string), msgAndArgs[1:]...)
	}
	return ""
}

// ValidateSchema validate schema for the given response
func ValidateSchema(response *httpexpect.Response, schemaFilePath string) {
	var schema interface{}
	jsonFile, err := os.Open(schemaFilePath)
	if err != nil {
		fmt.Println(err)
	}
	data, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		LogError("Error loading " + schemaFilePath + " file")
	}
	json.Unmarshal([]byte(data), &schema)
	response.JSON().Schema(schema)
}
