package report

import (
	"log"
	"os"
	"testing"
	"text/template"
	"time"
)

var pass = 1
var fail = 0
var skip = -1

type ReportRequest interface {
	ReqestDetail() map[string]string
	GetModuleName() string
	ResponseDetail() map[string]string
}

var HTTPRequest ReportRequest

const storeObject = "automation/api/results"

// Report is the parent struct for the report generation which consist(Project name, Environment name, StartTime, Module info,
// TotalPass, Totalfail,TotalSkip)
type Report struct {
	Project     string
	Environment string
	StartTime   string
	Module      map[string]Module
	TotalPass   int
	Totalfail   int
	TotalSkip   int
}

// Module consists module wise test info
type Module struct {
	TotalPass int
	Totalfail int
	TotalSkip int
	Test      []Test
}

// Test consists test info(Status, TestNamne, Exception, ErrorTrace, AssertDetail, Request, Response)
type Test struct {
	Status       int
	TestNamne    string
	Exception    string
	ErrorTrace   string
	AssertDetail string
	Request      map[string]string
	Response     map[string]string
}

// Exception consists test case exception detail(ErrorTrace, Expected, Actual and Messages)
type Exception struct {
	ErrorTrace string
	Expected   string
	Actual     string
	Messages   string
}

// HtmlReport singleton instance for report
var HtmlReport *Report

// NewReport creates singelton report object with required information
func NewReport(projectName, environment string) *Report {
	if HtmlReport == nil {
		HtmlReport = &Report{
			Project:     projectName,
			Environment: environment,
			StartTime:   time.Now().Format("Mon, Jan 2 15:04:05"),
			Module:      map[string]Module{},
			TotalPass:   0,
			Totalfail:   0,
			TotalSkip:   0,
		}
	}
	return HtmlReport
}

// Fail sets the failed test detail can be called once test case got failed
func Fail(t *testing.T, errorTrace string, message string, assert string) {
	test := Test{
		Status:       fail,
		TestNamne:    t.Name(),
		Exception:    message,
		ErrorTrace:   errorTrace,
		AssertDetail: assert,
		Request:      HTTPRequest.ReqestDetail(),
		Response:     HTTPRequest.ResponseDetail(),
	}
	HtmlReport.setModuleData(HTTPRequest.GetModuleName(), test)
}

// Set test module data to the module struct
func (rep *Report) setModuleData(module string, test Test) {
	var mod Module
	if val, ok := rep.Module[module]; ok {
		mod = val
		mod.Test = append(mod.Test, test)
	} else {
		mod = Module{
			Test: []Test{test},
		}
	}
	mod = rep.updateModuleStatus(mod, test.Status)
	rep.Module[module] = mod
}

// SetModulePass for a module and test
func (rep *Report) SetModulePass() {
	var mod Module
	if val, ok := rep.Module[HTTPRequest.GetModuleName()]; ok {
		mod = val
	} else {
		mod = Module{}
	}
	mod = rep.updateModuleStatus(mod, pass)
	rep.Module[HTTPRequest.GetModuleName()] = mod
}

// Increase pass count for a module and test
func (rep *Report) updateModuleStatus(mod Module, status int) Module {
	switch status {
	case fail:
		mod.Totalfail = mod.Totalfail + 1
		rep.Totalfail = rep.Totalfail + 1
	case pass:
		mod.TotalPass = mod.TotalPass + 1
		rep.TotalPass = rep.TotalPass + 1
	case skip:
		mod.TotalSkip = mod.TotalSkip + 1
		rep.TotalSkip = rep.TotalSkip + 1
	}
	return mod
}

// CreateReport using the '/report.gohtml' template
func (rep *Report) CreateReport(htmlTemplatePath, reportPath string) string {
	tpl, err := template.ParseFiles(htmlTemplatePath)
	if err != nil {
		log.Println(time.Now().Format("02-Jan-2006"), "---- ", err)
	}
	file, _ := os.Create(reportPath)
	err = tpl.Execute(file, rep)
	if err != nil {
		panic(err)
	}
	return reportPath
}

// UploadReport create a final report and upload it to the cloud storage
func (rep *Report) UploadReport(htmlTemplatePath, reportPath string) (string, error) {
	path := rep.CreateReport(htmlTemplatePath, reportPath)
	client, err := GetStorageObject()
	if err != nil {
		log.Println(time.Now().Format("02-Jan-2006"), "---- ", err)
	}
	bucketName := os.Getenv("BUCKET_NAME")
	if bucketName == "" {
		log.Println(time.Now().Format("02-Jan-2006"), "---- ", "Cloud storage BUCKET_NAME environment not set")
	}
	url, err := client.UploadFile(bucketName,
		storeObject+"/"+time.Now().Format("02-Jan-2006")+"/"+time.Now().Format("15:04:05"), path)
	if err != nil {
		return "", err
	}
	return url, nil
}

// GetReport return the current report object
func (rep *Report) GetReport() *Report {
	return rep
}
