package config

const (
	configField = `%s %s`

	configGetter = `
	func (c *Config) %s() %s {
		return c.%s
	}`

	defIf = `
	if !f.%s.set {
		var err = f.%s.Set("%s")
		if err != nil {
			// _probably_ just for now
			panic(err)
		}
	}`

	flagField = `%s %sFlag`

	flagVar = `flag.Var(&f.%s, "%s", "%s")`

	mapper = `%s: f.%s.value,`

	reqIf = `
	if !f.%s.set {
		// just for now
		panic("%s is a required flag")
	}`

	validator = `Validate%s(%s %s) error`
)
