package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type baseValue struct {
	Type        *string
	Description string
	Required    bool
	Validate    bool
	Layout      string
	Format      string
	// Extra       map[string]string
	// Example interface{}
	// Format string?
	// Qualifiers? regex, lt, gt, etc

	Default interface{}
}

var (
	pwd  string
	root string
)

func main() {
	var err error

	pwd, err = os.Getwd()
	if err != nil {
		panic(err)
	}

	pwd += "/"

	root, err = filepath.Abs(os.Args[0])
	if err != nil {
		panic(err)
	}

	root = filepath.Dir(root) + "/"

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

	// fmt.Printf("config: %+v\n", config)

	var stmts = statements{
		libTypes: map[string]struct{}{},
	}

	for k := range config {
		var v = config[k]

		// fmt.Println("k, v", k, v)

		if v.Type == nil {
			// TODO: No type provided; assume its an object
			// skip struct config objects for now
			fmt.Printf("struct not implemented: %s - %v\n", k, v)

			continue
		}

		stmts.parseBase(k, v)
	}

	fmt.Println("Generating package ...")

	// Wipe the generated folder
	err = removeDir()
	if err != nil {
		panic(err)
	}

	err = createDir()
	if err != nil {
		panic(err)
	}

	// If there is a validator function, that means something requested validation.
	// In this case:
	//	- generate the validator.go file that contains the interface
	// This is to allow configs without validation to not concern themselves with the interface entirely
	if len(stmts.validators) > 0 {
		err = stmts.createValidatorFile()
		if err != nil {
			panic(err)
		}
	}

	err = stmts.createConfigFile()
	if err != nil {
		panic(err)
	}

	err = stmts.createFlagsFile()
	if err != nil {
		panic(err)
	}

	err = stmts.copyLibFiles()
	if err != nil {
		panic(err)
	}

	err = importAndFormat()
	if err != nil {
		panic(err)
	}

	fmt.Println("Done!")
}

func importAndFormat() error {
	var err = exec.Command("goimports", "-w", pwd+"config").Run()
	if err != nil {
		return err
	}

	return exec.Command("gofmt", "-w", pwd+"config").Run()
}

func removeDir() error {
	return os.RemoveAll(pwd + "config")
}

func createDir() error {
	return os.MkdirAll(pwd+"config", 0777)
}

func (s *statements) copyLibFiles() error {
	var err error

	// For not just copy the example files
	for k := range s.libTypes {
		// fmt.Println("copying \"" + root + "assets/lib/" + k + ".go\" to " + pwd + "config/")
		err = exec.Command("cp", root+"assets/lib/"+k+".go", pwd+"config/").Run()
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *statements) createConfigFile() error {
	var configTemplate, err = ioutil.ReadFile(root + "assets/templates/config.go")
	if err != nil {
		panic(err)
	}

	var cf = string(configTemplate)

	cf = strings.Replace(cf, "// configFields",
		strings.Join(s.configFields, "\n"), 1)

	cf = strings.Replace(cf, "// configGetters",
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

func (s *statements) createFlagsFile() error {
	var flagTemplate, err = ioutil.ReadFile(root + "assets/templates/flags.go")
	if err != nil {
		panic(err)
	}

	var cf = string(flagTemplate)

	cf = strings.Replace(cf, "// flagFields",
		strings.Join(s.flagFields, "\n"), 1)

	cf = strings.Replace(cf, "// flagVars",
		strings.Join(s.flagVars, "\n"), 1)

	cf = strings.Replace(cf, "// requiredIfs",
		strings.Join(s.requiredIfs, "\n"), 1)

	cf = strings.Replace(cf, "// defaultIfs",
		strings.Join(s.defaultIfs, "\n"), 1)

	cf = strings.Replace(cf, "// mappers",
		strings.Join(s.mappers, "\n"), 1)

	cf = strings.Replace(cf, "// extraFields",
		strings.Join(s.extraFields, "\n"), 1)

	return ioutil.WriteFile(pwd+"config/flags.go", []byte(cf), 0666)
}

func (s *statements) createValidatorFile() error {
	var validateTemplate, err = ioutil.ReadFile(root + "assets/templates/validator.go")
	if err != nil {
		panic(err)
	}

	var cf = strings.Replace(string(validateTemplate), "// validators",
		strings.Join(s.validators, "\n"), 1)

	return ioutil.WriteFile(pwd+"config/validator.go", []byte(cf), 0666)
}

type statements struct {
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

func (s *statements) parseBase(k string, v *baseValue) {
	var (
		configName = strings.ToLower(k)
		configType = *v.Type
	)

	s.libTypes[configType] = struct{}{}

	// Add config fields
	s.configFields = append(s.configFields, makeConfigField(configName, configType))

	// Add config getters
	s.configGetters = append(s.configGetters, makeConfigGetter(k, configName, configType))

	// Add the flag fields
	s.flagFields = append(s.flagFields, makeFlagField(configName, configType))

	// TODO: this is one that needs to change if we are in a struct
	// Create the flag.Var statements
	s.flagVars = append(s.flagVars, makeFlagVar(configName, v.Description))

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
	s.mappers = append(s.mappers, makeMapper(configName))

	if v.Validate {
		s.validators = append(s.validators, makeValidator(k, configName, configType))
		s.validateCalls = append(s.validateCalls, makeValidateCall(k, configName))
	}

	// If there are defined extra fields then we will add them
	if len(allowedExtraFields[configType]) > 0 {
		s.extraFields = append(s.extraFields, resolveExtraFields(configType, configName, v)...)
	}
}

var allowedExtraFields = map[string][]string{
	"time": []string{
		"format",
		"layout",
	},
}

var timeFormats = map[string]string{
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

func resolveExtraFields(configType, configName string, v *baseValue) []string {
	var fields []string

	// If we have extra fields, we need to use them
	if configType == "time" {
		// Get the extra fields
		// somehow check the fields?
		// going to have to change "rfc3339" -> "time.RFC3339"

		// We don't have a layout
		if v.Layout == "" {
			if v.Format == "" {
				// If both are empty, we will default to time.RFC3339
				fields = append(fields, makeExtraField(configName, "layout", time.RFC3339))
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
	}

	return fields
}

// makeDefaultString takes an interface makes and either:
// makes a stringified default default (heh) value in the case of a nil default
// or
// returns the stringified version of the provided default
func makeDefaultString(typeOf string, d interface{}) string {
	// If the default provided was nil (not specified), give it the Go default
	if d == nil {
		switch typeOf {
		case "string":
			return ""

		case "int":
			return "0"

		case "bool":
			return "false"

		default:
			panic(fmt.Errorf("invalid type: %s", typeOf))
		}
	}

	switch typeOf {
	case "string":
		var dd, ok = d.(string)
		if ok {
			return dd
		}

	case "int":
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
