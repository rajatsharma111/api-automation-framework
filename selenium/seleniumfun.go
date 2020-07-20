package selenium

import (
	"github.com/tebeka/selenium"
	"strings"
	"time"
)

// QuitRemote close current driver session
func QuitRemote() error {
	if err := driver.Quit(); err != nil {
		return err
	}
	return nil
}

// PageLoadTimeout wait for page load
func PageLoadTimeout() error {
	if err := driver.SetPageLoadTimeout(200); err != nil {
		return err
	}
	return nil
}

// Get the url
func Get(url string) error {
	if err := driver.Get(url); err != nil {
		return err
	}
	return nil
}

// GetTitle reterive the current title
func GetTitle() (string, error) {
	title, err := driver.Title()
	if err != nil {
		return "", err
	}
	return title, nil
}

// GetCurrentURL reterive the current URL
func GetCurrentURL() (string, error) {
	title, err := driver.CurrentURL()
	if err != nil {
		return "", err
	}
	return title, nil
}

// FindElement wait and find element
func FindElement(by, query string) (selenium.WebElement, error) {
	waitErr := WaitForElement(by, query)
	if waitErr != nil {
		return nil, waitErr
	}
	elem, err := driver.FindElement(by, query)
	return elem, err
}

// SendKeys enter the text with a enter
func SendKeys(input string, element selenium.WebElement) error {
	if err := element.SendKeys(input + selenium.EnterKey); err != nil {
		return err
	}
	return nil
}

// Click click on the given webelement
func Click(element selenium.WebElement) error {
	if err := element.Click(); err != nil {
		return err
	}
	return nil
}

// WaitForElement wait for elemnt to be displayed
func WaitForElement(by, query string) error {
	time.Sleep(5 * time.Second)
	titleChangeCondition := func(wd selenium.WebDriver) (bool, error) {
		elem, err := wd.FindElement(by, query)
		if err != nil {
			return false, err
		}
		isDisplayed, _ := elem.IsDisplayed()
		return isDisplayed, nil
	}
	if err := driver.Wait(titleChangeCondition); err == nil {
		return err
	}
	return nil
}

// WaitForTitle wait for  the expected title
func WaitForTitle(expected string) error {
	titleChangeCondition := func(wd selenium.WebDriver) (bool, error) {
		title, err := wd.Title()
		if err != nil {
			return false, err
		}

		return strings.Contains(title, expected), nil
	}
	if err := driver.WaitWithTimeoutAndInterval(titleChangeCondition, 10*time.Second, 200*time.Millisecond); err == nil {
		return err
	}
	return nil
}

// WaitForURL wait for  the current URL
func WaitForURL(expected string) error {
	titleChangeCondition := func(wd selenium.WebDriver) (bool, error) {
		title, err := wd.CurrentURL()
		if err != nil {
			return false, err
		}
		return strings.Contains(title, expected), nil
	}
	if err := driver.WaitWithTimeoutAndInterval(titleChangeCondition, 20*time.Second, 200*time.Millisecond); err == nil {
		return err
	}
	return nil
}
