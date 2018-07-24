package webdriver

import "fmt"

type WebForm struct {
	WebElement
}

func (f *WebForm) GetAllElements() ([]WebElement, error) {
	return f.FindElements(ByXPATH, ".//input | .//textarea | .//select | .//button")
}

func (f *WebForm) Get(elem string) (WebElement, error) {
	return f.FindElement(ByName, elem)
}

func (f *WebForm) GetValue(we WebElement) (string, error) {
	return we.GetAttribute("value")
}

func (f *WebForm) SetValue(elem string, value string) error {
	field, err := f.FindElement(ByName, elem)
	if err != nil {
		return fmt.Errorf("%q: can't find form element %s", err, elem)
	}
	err = field.SendKeys(value)
	if err != nil {
		return fmt.Errorf("%q: can't set form element %s value %s", err, elem, value)
	}
	return nil
}
