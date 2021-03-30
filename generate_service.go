package main

import (
	"go/format"
	"strings"

	"go.uber.org/zap"
)

type GenerateService struct {
	logger *zap.SugaredLogger
}

type GenerateServiceInput struct {
	PackageName string
	ModuleName  string
	ModelName   string
	FirstChar   string
}

func NewGenerateService(logger *zap.SugaredLogger) ICommand {
	return GenerateService{
		logger: logger,
	}
}

func (m GenerateService) CommandDescription() string {
	return "generate service"
}

func (m GenerateService) Run() {

	packageName := NewStringPrompt(m.logger, "Package Name", "")
	moduleName := NewStringPrompt(m.logger, "Module Name", "")
	modelName := NewStringPrompt(m.logger, "Model Name", "")

	firstChar := strings.ToLower(string(modelName[0]))

	input := GenerateServiceInput{
		PackageName: packageName,
		ModelName:   modelName,
		FirstChar:   firstChar,
		ModuleName:  moduleName,
	}

	output, err := ParseTemplate("generate_service.txt", input)
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
