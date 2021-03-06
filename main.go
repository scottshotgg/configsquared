package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/common/log"
	"github.com/scottshotgg/configsquared/assets"
)

const (
	rootAsset     = "assets/"
	templateAsset = rootAsset + "templates/"
	typesAsset    = rootAsset + "types/"
	zeroString    = "0"
	falseString   = "false"
	emptyString   = ""
)

type (
	baseValue struct {
		Type        *string
		Description string
		Required    bool
		Validate    bool
		Layout      string
		Format      string
		Attributes  map[string]*baseValue
		Example     string
		// Extra       map[string]string
		// Example interface{}
		// Format string?
		// Qualifiers? regex, lt, gt, etc

		Default interface{}
	}

	statements struct {
		libTypes map[string]struct{}

		configFields  []string
		configGetters []string

		flagFields []string
		flagVars   []string

		requiredIfs []string
		defaultIfs  []string

		mappers []string

		validators    []string
		validateCalls []string

		extraFields []string
	}
)

var (
	pwd  string
	root string

	allowedExtraFields = map[string][]string{
		"time": []string{
			"format",
			"layout",
		},
	}

	timeFormats = map[string]string{
		"ansic":       "Mon Jan _2 15:04:05 2006",
		"unixdate":    "Mon Jan _2 15:04:05 MST 2006",
		"rubydate":    "Mon Jan 02 15:04:05 -0700 2006",
		"rfc822":      "02 Jan 06 15:04 MST",
		"rc822z":      "02 Jan 06 15:04 -0700",
		"rfc850":      "Monday, 02-Jan-06 15:04:05 MST",
		"rfc1123":     "Mon, 02 Jan 2006 15:04:05 MST",
		"rfc1123Z":    "Mon, 02 Jan 2006 15:04:05 -0700",
		"rfc3339":     "2006-01-02T15:04:05Z07:00",
		"rfc3339nano": "2006-01-02T15:04:05.999999999Z07:00",
		"kitchen":     "3:04PM",
		// Handy time stamps.
		"stamp":      "Jan _2 15:04:05",
		"stampmilli": "Jan _2 15:04:05.000",
		"stampmicro": "Jan _2 15:04:05.000000",
		"stampnano":  "Jan _2 15:04:05.000000000",
	}

	types = map[string]string{
		"int":      zeroString,
		"int32":    zeroString,
		"int64":    zeroString,
		"uint":     zeroString,
		"uint32":   zeroString,
		"uint64":   zeroString,
		"string":   emptyString,
		"bool":     falseString,
		"url":      emptyString,
		"time":     emptyString,
		"duration": emptyString,
		"ip":       emptyString,
		"ipv4":     emptyString,
		"ipv6":     emptyString,
	}
)

func main() {
	var err error

	pwd, err = os.Getwd()
	if err != nil {
		panic(err)
	}

	pwd += "/"

	// TODO: make this an arg
	contents, err := ioutil.ReadFile(pwd + "config.json")
	if err != nil {
		panic(err)
	}

	var config map[string]*baseValue

	err = json.Unmarshal(contents, &config)
	if err != nil {
		panic(err)
	}

	// Wipe the generated folder
	err = removeConfigDir()
	if err != nil {
		panic(err)
	}

	err = createConfigDir()
	if err != nil {
		panic(err)
	}

	// fmt.Printf("config: %+v\n", config)

	var stmts = statements{
		libTypes: map[string]struct{}{},
	}

	for k := range config {
		var v = config[k]

		fmt.Printf("k, v: %+v\n %+v\n", k, v)

		if v.Type == nil {
			// TODO: No type provided; assume its an object
			// not sure if I want to do this
			continue
		}

		switch *v.Type {
		case "object":
			stmts.parseStruct(k, v)

		default:
			stmts.parseBase("", k, v)
		}
	}

	fmt.Println("Generating package ...")

	// If there is a validator function, that means something requested validation.
	// In this case:
	//	- generate the validator.go file that contains the interface
	// This is to allow configs without validation to not concern themselves with the interface entirely
	if len(stmts.validators) > 0 {
		err = stmts.createValidatorFile()
		if err != nil {
			log.Fatalln("err", err)
		}
	}

	err = stmts.createConfigFile()
	if err != nil {
		log.Fatalln("err", err)
	}

	err = stmts.createFlagsFile()
	if err != nil {
		log.Fatalln("err", err)
	}

	err = stmts.copyLibFiles()
	if err != nil {
		log.Fatalln("err", err)
	}

	err = importAndFormat()
	if err != nil {
		log.Fatalln("err", err)
	}

	fmt.Println("Done!")
}

// TODO: should probably check if they have this and if they dont then throw an error or offer to install it for them
func importAndFormat() error {
	var err = exec.Command("goimports", "-w", pwd+"config").Run()
	if err != nil {
		return err
	}

	return exec.Command("gofmt", "-w", pwd+"config").Run()
}

func removeConfigDir() error {
	return os.RemoveAll(pwd + "config")
}

func createConfigDir() error {
	return os.MkdirAll(pwd+"config", 0777)
}

func (s *statements) copyLibFiles() error {
	// For not just copy the example files
	for k := range s.libTypes {
		var (
			fileName  = k + ".go"
			assetName = typesAsset + fileName
		)

		a, err := assets.Asset(assetName)
		if err != nil {
			return errors.New("could not find asset: " + assetName)
		}

		err = ioutil.WriteFile(pwd+"config/"+fileName, a, 0777)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *statements) createConfigFile() error {
	var configTemplate, err = assets.Asset(templateAsset + "config.go")
	if err != nil {
		panic(err)
	}

	var cf = string(configTemplate)

	cf = strings.Replace(cf, "// {{ configFields }}",
		strings.Join(s.configFields, "\n"), 1)

	cf = strings.Replace(cf, "// {{ configGetters }}",
		strings.Join(s.configGetters, "\n"), 1)

	if len(s.validateCalls) > 0 {
		// TODO: this could be done more efficiently, but this works for now
		cf = strings.Replace(cf, "Parse() *Config", "Parse(v Validator) (*Config, error)", 1)
		cf = strings.Replace(cf, "return &c", "return &c, nil", 1)
		cf = strings.Replace(cf, "c = f.toConfig()", parseValidate, 1)

		// APPEND the `validate()` function with the injected validate calls to the config file
		cf += "\n" + fmt.Sprintf(validateFunc, strings.Join(s.validateCalls, "\n"))
	}

	return ioutil.WriteFile(pwd+"config/config.go", []byte(cf), 0666)
}

// TODO: find some way to keep the values standard in the replace
func (s *statements) createFlagsFile() error {
	var flagTemplate, err = assets.Asset(templateAsset + "flags.go")
	if err != nil {
		panic(err)
	}

	var cf = string(flagTemplate)

	cf = strings.Replace(cf, "// {{ flagFields }}",
		strings.Join(s.flagFields, "\n"), 1)

	cf = strings.Replace(cf, "// {{ flagVars }}",
		strings.Join(s.flagVars, "\n"), 1)

	cf = strings.Replace(cf, "// {{ requiredIfs }}",
		strings.Join(s.requiredIfs, "\n"), 1)

	cf = strings.Replace(cf, "// {{ defaultIfs }}",
		strings.Join(s.defaultIfs, "\n"), 1)

	cf = strings.Replace(cf, "// {{ mappers }}",
		strings.Join(s.mappers, "\n"), 1)

	cf = strings.Replace(cf, "// {{ extraFields }}",
		strings.Join(s.extraFields, "\n"), 1)

	return ioutil.WriteFile(pwd+"config/flags.go", []byte(cf), 0666)
}

func (s *statements) createValidatorFile() error {
	var validateTemplate, err = assets.Asset(templateAsset + "validator.go")
	if err != nil {
		panic(err)
	}

	var cf = strings.Replace(string(validateTemplate), "// {{ validators }}",
		strings.Join(s.validators, "\n"), 1)

	return ioutil.WriteFile(pwd+"config/validator.go", []byte(cf), 0666)
}

func (s *statements) parseArray(k string, v *baseValue) {}

var (
	nestedStructTemplate = `
		package config

		type %s struct {
	`
)

func (s *statements) parseStruct(k string, v *baseValue) {
	fmt.Printf("i be parsing struct: %s, %+v\n", k, *v)

	var s1 = statements{
		libTypes: s.libTypes,
	}

	for k1, attr := range v.Attributes {
		if attr.Type == nil {
			panic("wtf nil type; fix me")
		}

		switch *attr.Type {
		case "object":
			fmt.Println("k1", k1)
			s1.parseStruct(k+"."+k1, attr)

		default:
			fmt.Println("k11", k1)
			s1.parseBase(k, k1, attr)
		}
	}

	var (
		configName = strings.ToLower(k)
		// configType = *v.Type
	)

	// Add base struct to outer config variable like you would a regular config
	s.configFields = append(s.configFields, makeConfigField(k, strings.ToUpper(k[0:1])+k[1:]))
	s.configGetters = append(s.configGetters, makeNestedStructConfigGetter("Config", k, configName, strings.ToUpper(k[0:1])+k[1:]))

	// Make the typedef header for file
	var header = fmt.Sprintf(nestedStructTemplate, strings.ToUpper(k[0:1])+k[1:])

	// Take all fields and add them here
	header += "\n" + strings.Join(s1.configFields, "\n")
	header += "\n}\n" + strings.Join(s1.configGetters, "\n")

	// Stringify the config write it to a file

	// // fold everything in
	// s.configFields = append(s.configFields, s1.configFields...)
	// s.configGetters = append(s.configGetters, s1.configGetters...)
	s.flagFields = append(s.flagFields, s1.flagFields...)
	s.flagVars = append(s.flagVars, s1.flagVars...)

	// // Check if they have marked it required
	// if v.Required {
	// 	s.requiredIfs = append(s.requiredIfs, makeRequiredIf(configName))
	// }

	// // Check if we need to make a default assertion
	// if v.Default != nil {
	// 	s.defaultIfs = append(s.defaultIfs, makeDefaultIf(configName, configType, v.Default))
	// }

	// TODO: this is one that needs to change if we are in a struct
	// Create the mapper from flags to config
	var mappers = `
		%s: %s{
			%s
		},
	`

	// fmt.Println("wat", strings.ToUpper(k[0:1])+strings.ToLower(k[1:]))
	// fmt.Println("wat2", k)

	fmt.Println("s1.mappers", s1.mappers)

	s.mappers = append(s.mappers, fmt.Sprintf(mappers, k, strings.ToUpper(k[0:1])+strings.ToLower(k[1:]), strings.Join(s1.mappers, "\n")))

	fmt.Println("s.mappers", s.mappers)

	// if v.Validate {
	// 	s.validators = append(s.validators, makeValidator(k, configName, configType))
	// 	s.validateCalls = append(s.validateCalls, makeValidateCall(k, configName))
	// }

	// // If there are defined extra fields then we will add them
	// if len(allowedExtraFields[configType]) > 0 {
	// 	s.extraFields = append(s.extraFields, resolveExtraFields(configType, configName, v)...)
	// }

	var err = ioutil.WriteFile(pwd+"config/"+k+".go", []byte(header), 0666)
	fmt.Println("err", err)
	if err != nil {
		fmt.Println("err:", err.Error())
		panic("err:" + err.Error())
		os.Exit(9)
	}
	/*
		flag: needs to be mongo.port, mongo.addr, mongo.etc...
		config: value needs to be a custom struct
	*/

	// add parent struct type to makeConfigGetter
	// <<< ULTIMATELY add parent struct name to parseBase >>>
	// s.configFields = append(s.configFields, makeConfigField(configName, configType))
}

func (s *statements) parseBase(base, k string, v *baseValue) {
	// It doesn't make sense to require the value but then also provide a default... whats the point?
	if v.Default != nil && v.Required {
		panic(fmt.Errorf("cannot apply a default to a required value - { name: %s, default: %+v }", k, v.Default))
	}

	var (
		configName = strings.ToLower(k)
		configType = *v.Type
	)

	// Determine if we have an array
	if strings.Contains(configType, "[") || strings.Contains(configType, "]") {
		configType = strings.Replace(configType, "[", "", 1)
		configType = strings.Replace(configType, "]", "", 1)
		fmt.Println("configTypedStripped", configType)

		configType = configType + "Array"
	}

	// Might need to be something here for the "real type"

	// // We have an array
	// if v.Items != nil {
	// 	// TODO: implement array parsing
	// 	panic(fmt.Errorf("arrays are not implemented: %s", k))
	// 	// s.parseArray(k, v)
	// }

	// If it is a `time` type and they want `unix` format then we need
	// to technically change the type since `unix` is an integer
	// instead of a string and it will be parsed differently
	if configType == "time" && v.Format == "unix" {
		configType = v.Format
	}

	s.libTypes[configType] = struct{}{}

	// Add config fields
	s.configFields = append(s.configFields, makeConfigField(configName, configType))

	// Add config getters
	if base == "" {
		s.configGetters = append(s.configGetters, makeConfigGetter(k, configName, configType))
	} else {

		if len(v.Attributes) > 0 {
			s.configGetters = append(s.configGetters, makeNestedStructConfigGetter(strings.ToUpper(base[0:1])+base[1:], k, configName, configType))
		} else {
			s.configGetters = append(s.configGetters, makeNestedConfigGetter(strings.ToUpper(base[0:1])+base[1:], k, configName, configType))
		}
	}

	// Add the flag fields
	s.flagFields = append(s.flagFields, makeFlagField(base, configName, configType))

	// TODO: this is one that needs to change if we are in a struct
	// Create the flag.Var statements
	s.flagVars = append(s.flagVars, makeFlagVar(base, configName, v.Description))

	// Check if they have marked it required
	if v.Required {
		s.requiredIfs = append(s.requiredIfs, makeRequiredIf(configName))
	}

	// Check if we need to make a default assertion
	if v.Default != nil {
		s.defaultIfs = append(s.defaultIfs, makeDefaultIf(configName, configType, v.Default))
	}

	// TODO: this is one that needs to change if we are in a struct
	// Create the mapper from flags to config
	if base == "" {
		s.mappers = append(s.mappers, makeMapper(configName))
	} else {
		s.mappers = append(s.mappers, makeNestedMapper(configName, base+strings.ToUpper(configName[0:1])+configName[1:]))
	}

	if v.Validate {
		s.validators = append(s.validators, makeValidator(k, configName, configType))
		s.validateCalls = append(s.validateCalls, makeValidateCall(k, configName))
	}

	// If there are defined extra fields then we will add them
	if len(allowedExtraFields[configType]) > 0 {
		s.extraFields = append(s.extraFields, resolveExtraFields(configType, configName, v)...)
	}
}

func timeFields(configName string, v *baseValue) []string {
	var fields []string

	// Get the extra fields
	// somehow check the fields?
	// going to have to change "rfc3339" -> "time.RFC3339"

	// We don't have a layout
	if v.Layout == "" {
		if v.Format == "" {
			// If both are empty, we will default to time.RFC3339
			fields = append(fields, makeExtraField(configName, "layout", time.RFC3339))
		} else if v.Format == "unix" {
			// panic("wtf")
		} else {
			// Ensure that the time format is valid (one of Go's predefined time formats)
			var f, ok = timeFormats[strings.ToLower(v.Format)]
			if !ok {
				panic(fmt.Errorf("invalid time format: %s", v.Format))
			}

			fields = append(fields, makeExtraField(configName, "layout", f))
		}
		// We have a layout
	} else {
		if v.Format != "" {
			// Invalid - You cannot have a "layout" and a "format", it must be one of the other because a "format" ultimately produces a "layout"
			panic(errors.New("cannot specify both a format and a layout; a format produces a layout"))
		}

		// First make sure they have not accidentally-ed the format into the layout
		// Ensure that the time format is valid (one of Go's predefined time formats)
		var _, ok = timeFormats[strings.ToLower(v.Layout)]
		if ok {
			panic(fmt.Errorf("please use the \"format\" attribute for a specific format: %s", v.Layout))
		}

		// Use their layout unconditionally
		fmt.Println("Using provided layout:", v.Layout)

		fields = append(fields, makeExtraField(configName, "layout", v.Layout))
	}

	return fields
}

func resolveExtraFields(configType, configName string, v *baseValue) []string {
	// If we have extra fields, we need to use them
	switch configType {
	case "time":
		return timeFields(configName, v)
	}

	return []string{}
}

// makeDefaultString takes an interface makes and either:
// makes a stringified default default (heh) value in the case of a nil default
// or
// returns the stringified version of the provided default
func makeDefaultString(typeOf string, d interface{}) string {
	// If the default provided was nil (not specified), give it the Go default
	if d == nil {
		var def, ok = types[typeOf]
		if !ok {
			panic(fmt.Errorf("invalid type: %s", typeOf))
		}

		return def
	}

	switch typeOf {
	case "url",
		"time",
		"duration",
		"ip",
		"ipv4",
		"ipv6",
		"string":
		var dd, ok = d.(string)
		if ok {
			return dd
		}

	case "uint",
		"uint32",
		"uint64",
		"int",
		"int32",
		"int64":
		// Go will unmarshal int fields to float64 due to the json spec
		var dd, ok = d.(float64)
		if ok {
			return strconv.FormatInt(int64(dd), 10)
		}

	case "bool":
		var dd, ok = d.(bool)
		if ok {
			return strconv.FormatBool(dd)
		}

	default:
		panic(fmt.Errorf("invalid type: %s", typeOf))
	}

	panic(invalidDefaultForType(typeOf, d))
}

func invalidDefaultForType(typeOf string, d interface{}) error {
	return fmt.Errorf("invalid default for type: %v for %s", d, typeOf)
}
