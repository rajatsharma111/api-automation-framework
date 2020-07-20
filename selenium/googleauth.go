package selenium

import (
	"github.com/tebeka/selenium"

	"errors"
	"net/url"
	"regexp"
	"strings"
)

const (
	googleUsernameField      = "identifierId"
	googleUsernameNextButton = "identifierNext"
	googlePasswordField      = "//input[@name='password']"
	googlePasswordNextButton = "passwordNext"
	googleAllow              = "#submit_approve_access > span"
)

// EnterUserName enter email/username on google sign in page
func enterUserName(userName string) error {
	webelement, err := FindElement(selenium.ByID, googleUsernameField)
	if err != nil {
		return err
	}
	return SendKeys(userName, webelement)
}

// EnterPassword enter passord on google sign in page
func enterPassword(password string) error {
	webelement, err := FindElement(selenium.ByXPATH, googlePasswordField)
	if err != nil {
		return err
	}
	return SendKeys(password, webelement)
}

// ClickUserNextButton click next button after enter username on google sign in page
func clickUserNextButton() error {
	webelement, err := FindElement(selenium.ByID, googleUsernameNextButton)
	if err != nil {
		return err
	}
	return Click(webelement)
}

// ClickPasswordNextButton click next button after enter password on google sign in page
func clickPasswordNextButton() error {
	webelement, err := FindElement(selenium.ByID, googlePasswordNextButton)
	if err != nil {
		return err
	}
	return Click(webelement)
}

// ClickAllow for allowing google access to strato app
func clickAllow() error {
	webelement, err := FindElement(selenium.ByCSSSelector, googleAllow)
	if err != nil {
		return err
	}
	return Click(webelement)
}

// GetGoogleAuthCode hit the given oauth url, enter credentials and allow access for strato
func GetGoogleAuthCode(url, userName, password string) (string, string, error) {
	NewDriver("chrome")
	defer StopService()
	defer QuitRemote()

	// open browser with the url
	err := Get(url)
	if err != nil {
		return "", "", err
	}

	err = WaitForTitle("Sign in - Google Accounts")
	if err != nil {
		return "", "", err
	}

	// Enter user name and click next
	err = enterUserName(userName)
	if err != nil {
		return "", "", err
	}

	// Enter password name and click next
	err = enterPassword(password)
	if err != nil {
		return "", "", err
	}

	err = WaitForTitle("Strato - One dash to rule them all!")
	if err != nil {
		return "", "", err
	}

	err = WaitForURL("code=")
	if err != nil {
		return "", "", err
	}

	currentURL, errURL := GetCurrentURL()
	if errURL != nil {
		return "", "", errURL
	}

	conditon := strings.Contains(currentURL, "code=")
	if !conditon {
		return "", "", errors.New("Auth Code is not generated")
	}

	code, errCode := getURLData(currentURL, "code")
	if errCode != nil {
		return "", "", errCode
	}
	state, errState := getURLData(currentURL, "state")
	if errState != nil {
		return code, "", errState
	}

	return code, state, nil
}

func getURLData(urlStr string, key string) (string, error) {
	url1, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}
	s, errQueryUnescape := url.QueryUnescape(url1.String())
	if err != nil {
		return "", errQueryUnescape
	}
	path := "(" + key + "=)(.*?)(\\&)"
	re := regexp.MustCompile(path)
	if len(re.FindStringSubmatch(s)) >= 2 {
		return re.FindStringSubmatch(s)[2], nil
	}
	return "", errors.New("Error while extracting " + key + " from the URL, " + url1.String())
}
