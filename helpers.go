package main

import (
	"fmt"
	"strings"
)

func makeFuncName(k string) string {
	if len(k) < 2 {
		panic(fmt.Errorf("string less than length"))
	}

	return strings.ToUpper(string(k[0])) + k[1:]
}

var realType = map[string]string{
	"unix":        "time.Time",
	"time":        "time.Time",
	"duration":    "time.Duration",
	"url":         "url.URL",
	"any":         "string",
	"stringArray": "[]string",
}

func makeConfigField(configName, configType string) string {
	var rt = realType[configType]
	if len(rt) > 0 {
		configType = rt
	}

	return fmt.Sprintf(configField, configName, configType)
}

func makeConfigGetter(k, configName, configType string) string {
	var rt = realType[configType]
	if len(rt) > 0 {
		configType = rt
	}

	return fmt.Sprintf(configGetter, makeFuncName(k), configType, configName)
}

func makeFlagField(configName, configType string) string {
	return fmt.Sprintf(flagField, configName, configType)
}

func makeFlagVar(configName, desc string) string {
	return fmt.Sprintf(flagVar, configName, configName, desc)
}

func makeRequiredIf(configName string) string {
	return fmt.Sprintf(reqIf, configName, configName)
}

func makeMapper(configName string) string {
	return fmt.Sprintf(mapper, configName, configName)
}

func makeDefaultIf(configName, configType string, def interface{}) string {
	return fmt.Sprintf(defIf, configName, configName, makeDefaultString(configType, def))
}

func makeValidator(k, configName, configType string) string {
	return fmt.Sprintf(validator, makeFuncName(k), configName, configType)
}

func makeValidateCall(k, configName string) string {
	return fmt.Sprintf(validateCall, makeFuncName(k), configName)
}

func makeExtraField(configName, fieldName, fieldValue string) string {
	return fmt.Sprintf(extraField, configName, fieldName, fieldValue)
}
