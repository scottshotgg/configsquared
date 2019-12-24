# configsquared

ConfigSquared: A configurable code generator to easily produce type-safe config packages for your application

`bin/example` will regenerate the `examples/config` directory; it can also be used as a starting scratchpad to play around.

<br>

## Example:

_Let's walk through a simple example from the `examples` directory:_

To start off, let's look at the contained `config.json` file (`examples/config.json`). We will see that the file is very similar to a JSON Schema; this is by design. Your config should be thought of as a contract between you and the user of your application, whoever it may be - Kubernetes, Docker, yourself, "the client", Spock - and this should be reflected in it's design and external behavior.

<br>

### Definition

Every `config.json` file contains a single object in which the keys identify the flags you intend to utilize in your application with their respective value representing intrinsic attributes about their value, such as the variable `type`, whether or not is is `required` for the user to provide this, if it has a `default` value, and a `description` that will appear in the help menu. In addition, the user can also specify whether or not they would like to require the library to `validate` the config values at run-time. <br>
The _only_ required value to have is `type`; everything else is optional.

\*_Currently, the only types implemented are `bool`, `int`, and `string`._

<br>

### Generation

After `config.json` has been sufficiently defined, the tool can be run from the same directory to produce a `config` package with the needed utilities. Taking a look inside, we can see the `config.go` file; the root of the package. In this file lies the external `Config`
structure you will work work with in your code.

<br>

### Usage

To utilize the generated package, start by importing it and then just call `config.Parse()`. That's all! This will **always** return a `*Config`, but will also return an `error` if validation has been required.

_\*Take special note here, you do have to supply a `Validator` if you have required run-time validation._

<br>

### Validation

Examining `validator.go`, we can see how to implement the validation for our flags - by supplying a `Validator` which will ensure that values marked for validation are checked as we have specified. See `examples/main.go` for a simple demo implementation.

<br>
<br>

Try out the different scripts in the `examples/bin` folder to get a feel for what you can expect at run-time.

- `bin/valid` - simple script to display how values are passed
- `bin/invalid` - discovering an invalid flag via a validation function
- `bin/required` - displaying what happens if a required flag is not passed
