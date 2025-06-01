package prompts

import (
	"fmt"
	"regexp"

	"github.com/go-git/go-git/v5"
	"github.com/rivo/tview"
)

// To step through the debugger in a TUI,
// must start dlv hosting the app in headless mode
//     cd <dir with main.go>
//     dlv debug --headless --listen=:2345 --api-version=2  ./main.go
//
// in vs code add this to launch.json
//
//     {
//         "name": "Connect to headless",
//         "type": "go",
//         "debugAdapter": "legacy",
//         "request": "attach",
//         "mode": "remote",
//         "port": 2345,
//           // note this could also be a server or a container
//         "host": "127.0.0.1"
//     }
//
// alternatively just log it...
// in main.go
//
// file, err := os.OpenFile("./mylog.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
// if err != nil {
//   log.Fatal(err)
// }
// log.SetOutput(file)
// ..
// log.Println(spew.Sdump(okButton))
// and then tail the log file in another terminal

// these operations are typical of conventional commits.
var operations []string = []string{
	"chore", "build", "ci", "docs", "feat", "fix", "perf", "refactor", "revert", "style", "test"}

func CommitToLocal(
	statuses git.Status,
	issue string,
	description string) (
	[]string, string, string, string, error) {
	app := tview.NewApplication()
	var statusField *tview.TextArea
	var okButton *tview.Button

	// FIXME implement custom primitives to do validation warning highlights.
	// FIXME implement custom buttons to customize the disabled state - can not persistently set style
	// The default Form over rides items styles.
	// https://github.com/rivo/tview/issues/931
	// we would need a custom form with its own drawing to implement more customized item styles
	form := tview.NewForm().
		AddDropDown(operationLabel, operations, 0, nil)

	filePaths, lines := toFileLines(statuses, ToStagingStatus)
	addLinesText(form, "Files", lines)
	form.AddInputField(issueLabel, issue, 20,
		func(textToCheck string, lastChar rune) bool {
			onValidate(form, statusField, okButton)
			return true
		},
		nil)

	form.AddInputField(descriptionLabel, description, 60,
		func(txt string, lastChar rune) bool {
			onValidate(form, statusField, okButton)
			return true
		},
		nil)

	isOk := false
	form.
		AddButton(cancelLabel, func() {
			app.Stop()
		})
	form.
		AddButton(okLabel, func() {
			isOk = true
			app.Stop()
		})
	okButton = form.GetButton(form.GetButtonIndex(okLabel))

	form.AddTextArea(statusLabel, okStatus, 60, 20, 500, nil)
	statusField = form.GetFormItemByLabel(statusLabel).(*tview.TextArea)

	form.SetBorder(true).SetTitle("Enter some data").SetTitleAlign(tview.AlignLeft)
	if err := app.SetRoot(form, true).EnableMouse(true).Run(); err != nil { // block
		return nil, "", "", "", err
	}
	if !isOk {
		return nil, "", "", "", nil
	}

	opTxt, descriptionTxt, issueTxt, err := validateFields(form)
	return filePaths, opTxt, descriptionTxt, issueTxt, err
}

// onValidate will write validation state to form .
func onValidate(form *tview.Form, statusField *tview.TextArea, okButton *tview.Button) {
	if form == nil ||
		statusField == nil ||
		okButton == nil {
		return
	}
	if _, _, _, err := validateFields(form); err != nil {
		okButton.SetDisabled(true)
		statusField.SetText(err.Error(), true)
	} else {
		okButton.SetDisabled(false)
		statusField.SetText(okStatus, true)
	}
}

const operationLabel = "Operation"
const issueLabel = "Issue"
const descriptionLabel = "Description"
const statusLabel = "Status"
const okLabel = "Ok"
const cancelLabel = "Cancel"
const badOperationStatus = "Operation does not match conventional commit rules"
const badIssueStatus = "Issue does not match jira naming rules."
const badDescriptionStatus = "Description does not match conventional commit rules."
const okStatus = "-----------"

func validateFields(form *tview.Form) (string, string, string, error) {
	_, opTxt := form.GetFormItemByLabel(operationLabel).(*tview.DropDown).GetCurrentOption()
	if len(opTxt) < 1 {
		return "", "", "", fmt.Errorf(badOperationStatus)
	}

	issueTxt := form.GetFormItemByLabel(issueLabel).(*tview.InputField).GetText()
	if !jiraIssueRE.Match([]byte(issueTxt)) {
		return "", "", "", fmt.Errorf(badIssueStatus)
	}
	descriptionTxt := form.GetFormItemByLabel(descriptionLabel).(*tview.InputField).GetText()
	if !isConventionalCommitDescription(descriptionTxt) {
		return "", "", "", fmt.Errorf(badDescriptionStatus)
	}
	return opTxt, descriptionTxt, issueTxt, nil
}

// var partialIssueRE = regexp.MustCompile(`^([A-Z]{1,5})?(-)?([0-9]{1,5})?$`)
var jiraIssueRE = regexp.MustCompile(`^([A-Z]{1,5})(-)([0-9]{1,5})$`)

/*
// these are the default case rules for conventional commitlint
'sentence-case', // Sentence case
  'start-case'.    // Start Case
  'pascal-case',   // PascalCase
  'upper-case',    // UPPERCASE
*/

// Common punctuation: Periods, commas, semicolons, colons, hyphens, question marks, exclamation marks, parentheses, brackets
// Some special characters: Underscores, ampersands, at signs, plus signs, equal signs, asterisks, percent signs, slashes (forward and backward)
// var specialCh = `\.\,\;\:\-\?\!\(\)\_\&\@\+\=\*\%\\\/`
var specialCh = `\:`
var digitSpecialCh = `0-9` + specialCh
var lowerCh = `a-z` + digitSpecialCh
var upperCh = `A-Z` + digitSpecialCh

// 'sentence-case', // Sentence case
var sentenceCaseTxt = `^([A-Z]?[` + lowerCh + `]*)(\s+[` + lowerCh + `]+)*$`
var sentenceCaseRE = regexp.MustCompile(sentenceCaseTxt)

var startWord = `[A-Z][` + lowerCh + `]*`
var digitSpecialWord = `[` + digitSpecialCh + `]+`

var startOrDigitSpecialWord = startWord + `|` + digitSpecialWord

// 'start-case'.    // Start Case
var startTxt = `^(` + startOrDigitSpecialWord + `)(\s+(` + startOrDigitSpecialWord + `))*$`
var startRe = regexp.MustCompile(startTxt)

// 'pascal-case',   // PascalCase
// are space allowed?
var pascalWord = `[A-Z][` + lowerCh + `]+`
var pascalTxt = `^(` + pascalWord + `)+$`
var pascalRe = regexp.MustCompile(pascalTxt)

var upperCaseTxt = `^([` + upperCh + `]+)(\s+[` + upperCh + `]+)*$`
var upperRe = regexp.MustCompile(upperCaseTxt)

func isConventionalCommitDescription(txt string) bool {
	b := []byte(txt)
	return sentenceCaseRE.Match(b) ||
		startRe.Match(b) ||
		pascalRe.Match(b) ||
		upperRe.Match(b)
}
