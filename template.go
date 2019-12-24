package main

var (
	configField = `%s %s`

	configGetter = `func (c *Config) %s() %s {
	return c.%s
}`

	flagField = `%s %sFlag`

	flagVar = `flag.Var(&f.%s, "%s", "%s")`

	reqIf = `if !f.%s.set {
	// just for now
	panic("%s is a required flag")
}`

	defIf = `if !f.%s.set {
	var err = f.%s.Set("%s")
	if err != nil {
		// _probably_ just for now
		panic(err)
	}
}`

	mapper = `%s: f.%s.value,`

	validator = `Validate%s(%s %s) error`

	validateFunc = `func (c *Config) validate(v Validator) error {
		var err error

		%s

		return nil
	}`

	validateCall = `err = v.Validate%s(c.%s)
	if err != nil {
		return err
}`

	parseValidate = `c = f.toConfig()
	
	return &c, c.validate(v)`

	extraField = `f.%s.%s = "%s"`
)

// TODO: not sure if I want to do this for now
type templateInfo struct {
	configField  string
	configGetter string
	flagField    string
	flagVar      string
	reqIf        string
	defIf        string
	mapper       string
	validator    string
}

func (t *templateInfo) Init() error {
	// ioutil.ReadFile("templates")
	return nil
}
