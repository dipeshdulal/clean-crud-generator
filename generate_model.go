package main

import (
	"fmt"
	"go/format"
	"strings"

	"go.uber.org/zap"
)

type GenerateModel struct {
	logger *zap.SugaredLogger
}

type GenerateModelInput struct {
	PackageName string
	ModelName   string
	DataTypes   string
	FirstChar   string
	TableName   string
	MapData     string
}

func NewGenerateModel(logger *zap.SugaredLogger) ICommand {
	return GenerateModel{
		logger: logger,
	}
}

func (m GenerateModel) CommandDescription() string {
	return "generate model with mapping functions"
}

func (m GenerateModel) Run() {

	packageName := NewStringPrompt(m.logger, "Package Name", "")
	modelName := NewStringPrompt(m.logger, "Model Name", "")
	tableName := NewStringPrompt(m.logger, "Table Name", ToSnakeCase(modelName)+"s")

	m.logger.Info("For columns info we require columns in format")
	m.logger.Info("<PascalCaseColumnName>:<datatype>")
	m.logger.Info("eg: ")
	m.logger.Info("UserID:string PhoneNo:uint Password:string")

	columnsString := NewStringPrompt(m.logger, "Enter columns list", "")
	firstChar := strings.ToLower(string(modelName[0]))

	columns := strings.Split(columnsString, " ")
	dataTypes := ""
	mapData := []string{}
	for _, column := range columns {
		trimedColumn := strings.Trim(column, " ")
		tokens := strings.Split(trimedColumn, ":")
		if len(tokens) == 2 {
			tag := fmt.Sprintf("`json:\"%s\"`", ToSnakeCase(tokens[0]))
			datatype := fmt.Sprintf("%s %s %s \n", tokens[0], tokens[1], tag)

			md := fmt.Sprintf("\"%s\": %s.%s", ToSnakeCase(tokens[0]), firstChar, tokens[0])
			mapData = append(mapData, md)
			dataTypes = dataTypes + datatype
		}
	}

	input := GenerateModelInput{
		PackageName: packageName,
		ModelName:   modelName,
		TableName:   tableName,
		DataTypes:   dataTypes,
		FirstChar:   firstChar,
		MapData:     strings.Join(mapData, ",\n") + ",",
	}

	output, err := ParseTemplate("generate_model.txt", input)
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
