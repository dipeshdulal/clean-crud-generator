package main

import (
	"go/format"
	"strings"

	"go.uber.org/zap"
)

type GenerateController struct {
	logger *zap.SugaredLogger
}

type GenerateControllerInput struct {
	PackageName string
	ModuleName  string
	ModelName   string
	FirstChar   string
}

func NewGenerateController(logger *zap.SugaredLogger) ICommand {
	return GenerateController{
		logger: logger,
	}
}

func (m GenerateController) CommandDescription() string {
	return "generate controller"
}

func (m GenerateController) Run() {

	packageName := NewStringPrompt(m.logger, "Package Name", "")
	moduleName := NewStringPrompt(m.logger, "Module Name", "")
	modelName := NewStringPrompt(m.logger, "Model Name", "")

	firstChar := strings.ToLower(string(modelName[0]))

	input := GenerateControllerInput{
		PackageName: packageName,
		ModelName:   modelName,
		FirstChar:   firstChar,
		ModuleName:  moduleName,
	}

	output, err := ParseTemplate("generate_controller.txt", input)
	if err != nil {
		m.logger.Error("parsing error: ", err.Error())
		return
	}

	formatedOutput, err := format.Source([]byte(output))
	if err != nil {
		m.logger.Error("go formating error: ", err.Error())
		return
	}

	m.logger.Infof("%+v \n", string(formatedOutput))
}
