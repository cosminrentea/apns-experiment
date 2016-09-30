package main

import (
	"sourcegraph.com/sourcegraph/go-selenium"
	"testing"
	"fmt"
)

var executorURL = "http://localhost:4444/wd/hub"

// An example test using the WebDriverT and WebElementT interfaces. If you use the non-*T
// interfaces, you must perform error checking that is tangential to what you are testing,
// and you have to destructure results from method calls.
func TestWithT(t *testing.T) {
	var webDriver selenium.WebDriver
	var err error
	caps := selenium.Capabilities(map[string]interface{}{"browserName": "firefox"})
	if webDriver, err = selenium.NewRemote(caps, executorURL); err != nil {
		fmt.Printf("Failed to open session: %s\n", err)
		return
	}
	defer webDriver.Quit()

	err = webDriver.Get("https://google.com")
	if err != nil {
		fmt.Printf("Failed to load page: %s\n", err)
		return
	}

	if title, err := webDriver.Title(); err == nil {
		fmt.Printf("Page title: %s\n", title)
	} else {
		fmt.Printf("Failed to get page title: %s", err)
		return
	}

	var elem selenium.WebElement
	elem, err = webDriver.FindElement(selenium.ByCSSSelector, ".repo .name")
	if err != nil {
		fmt.Printf("Failed to find element: %s\n", err)
		return
	}

	if text, err := elem.Text(); err == nil {
		fmt.Printf("Repository: %s\n", text)
	} else {
		fmt.Printf("Failed to get text of element: %s\n", err)
		return
	}
}
