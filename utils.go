package main

import (
	"bytes"
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/manifoldco/promptui"
	"go.uber.org/zap"
)

// NewStringPrompt creates new string prompt
func NewStringPrompt(
	logger *zap.SugaredLogger,
	label string,
	defaultValue string,
) string {
	prompt := promptui.Prompt{
		Label:   label,
		Default: defaultValue,
	}
	result, err := prompt.Run()
	if err != nil {
		logger.Error("prompt failed")
		return ""
	}

	return result
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

// ToSnakeCase convers pascal case to snake case
// taken from: https://gist.github.com/stoewer/fbe273b711e6a06315d19552dd4d33e6
func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

//ParseTemplate to parse the template with given data
func ParseTemplate(templateFileName string, data interface{}) (string, error) {
	templatePath, err := filepath.Abs(fmt.Sprintf("templates/%s", templateFileName))
	if err != nil {
		return "", errors.New("invalid template name")
	}
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return "", err
	}
	body := buf.String()
	return body, nil
}
