package main

import (
	"go/format"
	"strings"

	"go.uber.org/zap"
)

type GenerateRepository struct {
	logger *zap.SugaredLogger
}

type GenerateRepositoryInput struct {
	PackageName string
	ModuleName  string
	ModelName   string
	FirstChar   string
}

func NewGenerateRepository(logger *zap.SugaredLogger) ICommand {
	return GenerateRepository{
		logger: logger,
	}
}

func (m GenerateRepository) CommandDescription() string {
	return "generate gorm repository"
}

func (m GenerateRepository) Run() {

	packageName := NewStringPrompt(m.logger, "Package Name", "")
	moduleName := NewStringPrompt(m.logger, "Module Name", "")
	modelName := NewStringPrompt(m.logger, "Model Name", "")

	firstChar := strings.ToLower(string(modelName[0]))

	input := GenerateRepositoryInput{
		PackageName: packageName,
		ModelName:   modelName,
		FirstChar:   firstChar,
		ModuleName:  moduleName,
	}

	output, err := ParseTemplate("generate_repository.txt", input)
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
