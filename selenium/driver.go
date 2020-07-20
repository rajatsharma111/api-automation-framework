package selenium

import (
	"fmt"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"os"
	"stratoApiAutomation/helpers"
)

var headless = false
var driver selenium.WebDriver
var seleniumService *selenium.Service

// Need to change these path once got the cldcvr test account
const (
	seleniumPath = "/selenium-server-standalone-3.141.59.jar"
	driverPath   = "/chromedriver"
	port         = 8080
)

// StartService Starts selenium server
func StartService() {
	opts := []selenium.ServiceOption{
		selenium.ChromeDriver(helpers.AutomationPath + driverPath), // Specify the path to GeckoDriver in order to use Firefox.
		selenium.Output(os.Stderr),                                 // Output debug information to STDERR.
	}
	selenium.SetDebug(false)
	selenium, err := selenium.NewSeleniumService(helpers.AutomationPath+seleniumPath, port, opts...)
	if err != nil {
		panic(err)
	}
	seleniumService = selenium
}

// StopService stops selenium server
func StopService() {
	seleniumService.Stop()
}

// NewDriver create new browser driver
func NewDriver(browser string) selenium.WebDriver {
	StartService()
	caps := selenium.Capabilities{"browserName": browser}

	switch browser {
	case "chrome":
		chrCaps := chrome.Capabilities{
			Args: []string{
				"--no-sandbox",
			},
			W3C: true,
		}
		if headless {
			chrCaps.Args = append(chrCaps.Args, "--headless")
		}
		caps.AddChrome(chrCaps)

	case "htmlunit":
		caps["javascriptEnabled"] = true
	}
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	if err != nil {
	}
	driver = wd
	return wd
}
