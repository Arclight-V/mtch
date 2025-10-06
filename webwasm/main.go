package main

import (
	"encoding/json"
	"strconv"
	"strings"
	"syscall/js"
	"time"
)

var document js.Value

const projectName = "MTCH"

func main() {
	document = js.Global().Get("document")
	buildPage()
	select {}
}

func buildPage() {
	body := document.Get("body")
	body.Set("className", "fb-body")
	body.Set("innerHTML", "")

	page := createElement("div", "page")
	body.Call("appendChild", page)

	header := createElement("header", "page__header")
	brand := createElement("div", "brand")
	brand.Set("textContent", projectName)
	header.Call("appendChild", brand)

	loginLink := createElement("a", "page__login-link")
	loginLink.Set("href", "#")
	loginLink.Set("textContent", "Already have an account?")
	header.Call("appendChild", loginLink)

	page.Call("appendChild", header)

	main := createElement("main", "page__main")
	intro := createElement("section", "intro")
	introTitle := createElement("h1", "intro__title")
	introTitle.Set("textContent", "Create a new account")
	introSubtitle := createElement("p", "intro__subtitle")
	introSubtitle.Set("textContent", "It's quick and easy.")
	appendChildren(intro, introTitle, introSubtitle)

	card := createElement("section", "signup-card")
	cardTitle := createElement("h2", "signup-card__title")
	cardTitle.Set("textContent", "Sign Up")
	cardSubtitle := createElement("p", "signup-card__subtitle")
	cardSubtitle.Set("textContent", "It's free and always will be.")

	form := createElement("form", "signup-form")
	form.Set("autocomplete", "off")

	nameRow := createElement("div", "signup-form__row", "signup-form__row--split")
	firstName := createInput("text", "First name")
	firstName.Set("name", "firstname")
	lastName := createInput("text", "Surname")
	lastName.Set("name", "lastname")
	appendChildren(nameRow, firstName, lastName)

	contact := createInput("text", "Mobile number or email address")
	contact.Set("name", "contact")
	password := createInput("password", "New password")
	password.Set("name", "password")

	dobGroup := createElement("div", "signup-form__group")
	dobLabel := createElement("span", "signup-form__label")
	dobLabel.Set("textContent", "Birthday")
	dobRow := createElement("div", "signup-form__row", "signup-form__row--split")

	monthSelect := createSelect()
	daySelect := createSelect()
	yearSelect := createSelect()

	populateMonthOptions(monthSelect)
	populateDayOptions(daySelect)
	populateYearOptions(yearSelect)

	appendChildren(dobRow, monthSelect, daySelect, yearSelect)
	appendChildren(dobGroup, dobLabel, dobRow)

	genderGroup := createElement("div", "signup-form__group")
	genderLabel := createElement("span", "signup-form__label")
	genderLabel.Set("textContent", "Gender")
	genderChoices := createElement("div", "signup-form__row", "signup-form__row--choices")

	for _, option := range []struct {
		label string
		value string
	}{
		{label: "Female", value: "female"},
		{label: "Male", value: "male"},
		{label: "Custom", value: "custom"},
	} {
		genderChoices.Call("appendChild", createRadioChoice(option.label, option.value))
	}

	appendChildren(genderGroup, genderLabel, genderChoices)

	helperText := createElement("p", "signup-form__helper")
	helperText.Set("textContent", "People who use our service may have uploaded your contact information to "+projectName+".")

	policyText := createElement("p", "signup-form__policy")
	policyText.Set("innerHTML", "By clicking Sign Up, you agree to our <a href='#'>Terms</a>, <a href='#'>Privacy Policy</a> and <a href='#'>Cookies Policy</a>.")

	submit := createElement("button", "signup-form__submit")
	submit.Set("type", "submit")
	submit.Set("textContent", "Sign Up")

	status := createElement("p", "signup-form__status")
	status.Call("setAttribute", "role", "status")

	appendChildren(form, nameRow, contact, password, dobGroup, genderGroup, helperText, policyText, submit, status)
	appendChildren(card, cardTitle, cardSubtitle, form)
	appendChildren(main, intro, card)
	page.Call("appendChild", main)

	footer := createElement("footer", "page__footer")
	footerNote := createElement("p", "page__footer-note")
	footerNote.Set("textContent", projectName+" helps you connect and share with the people in your life.")
	footer.Call("appendChild", footerNote)
	page.Call("appendChild", footer)

	submitHandler := js.FuncOf(func(this js.Value, args []js.Value) any {
		event := args[0]
		event.Call("preventDefault")

		status.Set("textContent", "Submitting…")
		status.Get("classList").Call("add", "signup-form__status--visible")

		payload := map[string]string{
			"first_name":  strings.TrimSpace(firstName.Get("value").String()),
			"last_name":   strings.TrimSpace(lastName.Get("value").String()),
			"contact":     strings.TrimSpace(contact.Get("value").String()),
			"password":    password.Get("value").String(),
			"birth_month": monthSelect.Get("value").String(),
			"birth_day":   daySelect.Get("value").String(),
			"birth_year":  yearSelect.Get("value").String(),
		}

		if gender := getSelectedGender(form); gender != "" {
			payload["gender"] = gender
		}

		jsonBody, err := json.Marshal(payload)
		if err != nil {
			status.Set("textContent", "Failed to prepare request")
			return nil
		}

		init := map[string]any{
			"method": "POST",
			"headers": map[string]any{
				"Content-Type": "application/json",
			},
			"body": string(jsonBody),
		}

		fetchPromise := js.Global().Call("fetch", "/api/v1/auth/register", js.ValueOf(init))

		successHandler := js.FuncOf(func(this js.Value, args []js.Value) any {
			resp := args[0]
			if resp.Get("ok").Bool() {
				status.Set("textContent", "Registration submitted successfully!")
				return nil
			}
			status.Set("textContent", "Registration failed. Please try again.")
			return nil
		})

		errorHandler := js.FuncOf(func(this js.Value, args []js.Value) any {
			status.Set("textContent", "Network error. Please try again.")
			return nil
		})

		fetchPromise.Call("then", successHandler).Call("catch", errorHandler)

		return nil
	})

	form.Call("addEventListener", "submit", submitHandler)
}

func getSelectedGender(form js.Value) string {
	selected := form.Call("querySelector", "input[name='gender']:checked")
	if selected.IsUndefined() || selected.IsNull() {
		return ""
	}
	return selected.Get("value").String()
}

func createElement(tag string, classes ...string) js.Value {
	el := document.Call("createElement", tag)
	filtered := filterEmpty(classes)
	if len(filtered) > 0 {
		el.Set("className", strings.Join(filtered, " "))
	}
	return el
}

func createInput(inputType, placeholder string, classes ...string) js.Value {
	classList := append([]string{"signup-form__input"}, classes...)
	input := createElement("input", classList...)
	input.Set("type", inputType)
	if placeholder != "" {
		input.Set("placeholder", placeholder)
	}
	return input
}

func createSelect(classes ...string) js.Value {
	classList := append([]string{"signup-form__input", "signup-form__select"}, classes...)
	return createElement("select", classList...)
}

func createRadioChoice(label, value string) js.Value {
	wrapper := createElement("label", "signup-form__choice")
	input := createElement("input", "signup-form__choice-input")
	input.Set("type", "radio")
	input.Set("name", "gender")
	input.Set("value", value)
	text := createElement("span", "signup-form__choice-label")
	text.Set("textContent", label)
	appendChildren(wrapper, input, text)
	return wrapper
}

func appendChildren(parent js.Value, children ...js.Value) {
	for _, child := range children {
		parent.Call("appendChild", child)
	}
}

func populateMonthOptions(selectEl js.Value) {
	months := []string{
		"January", "February", "March", "April", "May", "June",
		"July", "August", "September", "October", "November", "December",
	}
	for index, month := range months {
		option := createOption(strconv.Itoa(index+1), month)
		selectEl.Call("appendChild", option)
	}
}

func populateDayOptions(selectEl js.Value) {
	for day := 1; day <= 31; day++ {
		option := createOption(strconv.Itoa(day), strconv.Itoa(day))
		selectEl.Call("appendChild", option)
	}
}

func populateYearOptions(selectEl js.Value) {
	currentYear := time.Now().Year()
	for year := currentYear; year >= currentYear-110; year-- {
		option := createOption(strconv.Itoa(year), strconv.Itoa(year))
		selectEl.Call("appendChild", option)
	}
}

func createOption(value, label string) js.Value {
	option := document.Call("createElement", "option")
	option.Set("value", value)
	option.Set("textContent", label)
	return option
}

func filterEmpty(values []string) []string {
	filtered := values[:0]
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			filtered = append(filtered, value)
		}
	}
	return filtered
}
