package report

import (
	"encoding/json"
	"fmt"
	"github.com/cldcvr/api-automation-framework/model"
	"github.com/gavv/httpexpect"
	"github.com/slack-go/slack"

	"os"
	"testing"
)

var (
	module      = "SlackNotification"
	httpRequest *model.HttpRequest
)

var section []*slack.SectionBlock

// Body to create slack body
type Body struct {
	Blocks slack.Blocks `json:"blocks,omitempty"`
}

// SendSlackMessage create slack messge template and send it using SLACK_WEB_HOOK_URL
func SendSlackMessage(t *testing.T, report *Report, url string) {
	httpRequest1 := model.HttpRequest{
		HttpUrl: model.HttpUrl{
			BaseUrl: os.Getenv("SLACK_WEB_HOOK_URL")},
		Method: "POST",
	}
	httpRequest = &httpRequest1
	str := createMessage(report, url)
	body := Body{
		Blocks: str,
	}

	byt, _ := json.Marshal(body)
	var jsonData map[string]interface{}
	json.Unmarshal(byt, &jsonData)
	httpRequest.Body = jsonData
	e := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  httpRequest.HttpUrl.BaseUrl,
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewDebugPrinter(t, true),
		},
	})
	request := e.POST(httpRequest.HttpUrl.BaseUrl)
	request = request.WithJSON(httpRequest.Body)
	request.Expect()
}

func addTextBlock(msg string) *slack.TextBlockObject {
	return slack.NewTextBlockObject("mrkdwn", msg, false, false)
}

func createMessage(report *Report, url string) slack.Blocks {
	textInfo := addTextBlock(fmt.Sprintf("*%s*  ", report.Project))
	envField := addTextBlock(fmt.Sprintf("*Environment*\n %s", report.Environment))
	var statusField *slack.TextBlockObject
	if report.Totalfail > 0 || report.TotalSkip > 0 {
		statusField = addTextBlock(fmt.Sprintf("*Status*\n %s", "Fail"))
	} else {
		statusField = addTextBlock(fmt.Sprintf("*Status*\n %s", "Pass"))
	}
	countField := addTextBlock(fmt.Sprintf("*Total Run:* %d, *Pass:* %d,  *Fail:* %d *Skip:* %d             %s",
		report.TotalPass+report.Totalfail+report.TotalSkip, report.TotalPass, report.Totalfail, report.TotalSkip, fmt.Sprintf("*<%s|see full results>*", url)))

	sectionBlock1 := slack.NewSectionBlock(textInfo, []*slack.TextBlockObject{envField, statusField}, nil, slack.SectionBlockOptionBlockID("test_block23"))
	sectionBlock2 := slack.NewSectionBlock(countField, nil, nil)
	section = append(section, sectionBlock1, sectionBlock2)
	var arr = make([]slack.Block, len(section))
	for index, val := range section {
		arr[index] = val
	}
	return slack.Blocks{arr}
}
