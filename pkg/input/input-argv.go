package input

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func NewArgvInput(argv []string) *ArgvInput {
	input := new(ArgvInput)

	if nil == argv {
		input.tokens = os.Args[1:]
	} else {
		input.tokens = argv[1:]
	}

	input.doParse = input.Parse

	return input
}

type ArgvInput struct {
	abstractInput
	tokens []string
	parsed []string
}

func (i *ArgvInput) GetFirstArgument() string {
	for _, token := range i.tokens {
		if "" != token && '-' == token[0] {
			continue
		}

		return token
	}

	panic(errors.New("first argument not found"))
}

func (i *ArgvInput) HasParameterOption(values []string, onlyParams bool) bool {
	panic("implement me")
}

func (i *ArgvInput) GetParameterOption(values []string, defaultValue string, onlyParams bool) {
	panic("implement me")
}

func (i *ArgvInput) Parse() {
	parseOptions := true
	i.parsed = i.tokens

	longOptionRegex := regexp.MustCompile("^--")

	for {
		if 0 == len(i.parsed) {
			break
		}

		token := i.parsed[0]
		i.parsed = i.parsed[1:]

		if parseOptions && "" == token {
			i.parseArgument(token)
		} else if parseOptions && "--" == token {
			parseOptions = false
		} else if parseOptions && longOptionRegex.MatchString(token) {
			i.parseLongOption(token)
		} else if parseOptions && '-' == token[0] && "-" != token {
			i.parseShortOption(token)
		} else {
			i.parseArgument(token)
		}

	}
}

//
// internal
//

func (i *ArgvInput) parseShortOption(token string) {
	name := token[1:]

	if len(name) > 1 {
		if i.definition.HasShortcut(name[0:0]) && i.definition.GetOptionForShortcut(name[0:0]).AcceptValue() {
			// an option with a value (with no space)
			i.addShortOption(name[0:0], name[1:])
		} else {
			i.parseShortOptionSet(name)
		}
	} else {
		i.addShortOption(name, "")
	}
}

func (i *ArgvInput) parseShortOptionSet(name string) {
	length := len(name)

	for index := 0; index < length; index++ {
		if !i.definition.HasShortcut(name[index:index]) {
			panic(errors.New(fmt.Sprintf("the '-%s' option does not exist", name[index:index])))
		}

		opt := i.definition.GetOptionForShortcut(name[index:index])

		if opt.AcceptValue() {
			if index == length - 1 {
				i.addLongOption(opt.GetName(), "")
			} else {
				i.addLongOption(opt.GetName(), name[index+1:])
			}

			break
		} else {
			i.addLongOption(opt.GetName(), "")
		}
	}
}

func (i *ArgvInput) parseLongOption(token string) {
	name := token[2:]
	pos := strings.Index(name, "=")

	if pos != -1 {
		value := name[pos+1:]

		if 0 == len(value) {
			i.parsed = append([]string{value}, i.parsed...)
		}

		i.addLongOption(name[0:pos], value)
	} else {
		i.addLongOption(name, "")
	}
}

func (i *ArgvInput) parseArgument(token string) {
	keys := i.definition.GetArgumentsOrder()

	if len(keys) > 0 && i.definition.HasArgument(keys[len(keys)-1]) {
		arg := i.definition.GetArgument(keys[len(keys)-1])

		if arg.IsArray() {
			i.argumentArrays[arg.GetName()] = []string{token}
		} else {
			i.arguments[arg.GetName()] = token
		}
	} else if len(keys) > 1 &&
		i.definition.HasArgument(keys[len(keys)-2]) &&
		i.definition.GetArgument(keys[len(keys)-2]).IsArray() {
		arg := i.definition.GetArgument(keys[len(keys)-2])
		i.argumentArrays[arg.GetName()] = append(i.argumentArrays[arg.GetName()], token)
	}
}

func (i *ArgvInput) addShortOption(shortcut string, value string) {
	if !i.definition.HasShortcut(shortcut) {
		panic(errors.New(fmt.Sprintf("the '-%s' option does not exist", shortcut)))
	}

	opt := i.definition.GetOptionForShortcut(shortcut)

	i.addLongOption(opt.GetName(), value)
}

func (i *ArgvInput) addLongOption(name string, value string) {
	if !i.definition.HasOption(name) {
		panic(errors.New(fmt.Sprintf("the '--%s' option does not exist", name)))
	}

	opt := i.definition.GetOption(name)

	if "" != value && !opt.AcceptValue() {
		panic(errors.New(fmt.Sprintf("the '--%s' option does not accept a value", name)))
	}

	acceptValue := opt.AcceptValue()
	size := len(i.parsed)

	if "" == value && acceptValue && size > 0 {
		// if option accepts an optional or mandatory argument
		// let's see if there is one provided
		next := i.parsed[0]
		i.parsed = i.parsed[1:]

		if len(next) > 0 || "" == next {
			value = next
		} else {
			i.parsed = append([]string{next}, i.parsed...)
		}
	}

	if "" == value {
		if opt.IsValueRequired() {
			panic(errors.New(fmt.Sprintf("the '--%s' option requires a value", name)))
		}
	}

	if opt.IsArray() {
		i.optionArrays[name] = append(i.optionArrays[name], value)
	} else {
		i.options[name] = value
	}
}
