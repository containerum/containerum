# Flags based on structures. [![GoDoc](https://godoc.org/github.com/octago/sflags?status.svg)](http://godoc.org/github.com/octago/sflags) [![Build Status](https://travis-ci.org/octago/sflags.svg?branch=master)](https://travis-ci.org/octago/sflags)  [![codecov](https://codecov.io/gh/octago/sflags/branch/master/graph/badge.svg)](https://codecov.io/gh/octago/sflags)  [![Go Report Card](https://goreportcard.com/badge/github.com/octago/sflags)](https://goreportcard.com/report/github.com/octago/sflags)

The sflags package uses structs, reflection and struct field tags
to allow you specify command line options. It supports [different types](#supported-types-in-structures) and [features](#features).

An example:

```golang
type HTTPConfig struct {
	Host    string `desc:"HTTP host"`
	Port    int `flag:"port p" desc:"some port"`
	SSL     bool `env:"HTTP_SSL_VALUE"`
	Timeout time.Duration `flag:",deprecated,hidden"`
}

type Config struct {
	HTTP   HTTPConfig
	Stats  StatsConfig
}
```

And you can use your favorite flag or cli library!

## Supported flags and cli libraries:

 - [x] [flag](https://golang.org/pkg/flag/) - [example](https://github.com/octago/sflags/blob/master/examples/flag/main.go)
 - [x] [spf13/pflag](https://github.com/spf13/pflag) - [example](https://github.com/octago/sflags/blob/master/examples/pflag/main.go)
 - [x] [spf13/cobra](https://github.com/spf13/cobra) - [example](https://github.com/octago/sflags/blob/master/examples/cobra/main.go)
 - [ ] [spf13/viper](https://github.com/spf13/viper)
 - [x] [urfave/cli](https://github.com/urfave/cli) [example](https://github.com/octago/sflags/blob/master/examples/urfave_cli/main.go)
 - [x] [kingpin](https://github.com/alecthomas/kingpin) [example](https://github.com/octago/sflags/blob/master/examples/kingpin/main.go)

## Features:

 - [x] Set environment name
 - [x] Set usage
 - [x] Long and short forms
 - [x] Skip field
 - [ ] Required
 - [ ] Placeholders (by `name`)
 - [x] Deprecated and hidden options
 - [ ] Multiple ENV names
 - [x] Interface for user types.
 - [x] [Validation](https://godoc.org/github.com/octago/sflags/validator/govalidator#New) (using [govalidator](https://github.com/asaskevich/govalidator) package)
 - [x] Anonymous nested structure support (anonymous structures flatten by default)

## Supported types in structures:

 - [x] `int`, `int8`, `int16`, `int32`, `int64`
 - [x] `uint`, `uint8`, `uint16`, `uint32`, `uint64`
 - [x] `float32`, `float64`
 - [x] slices for all previous numeric types (e.g. `[]int`, `[]float64`)
 - [x] `bool`
 - [x] `[]bool`
 - [x] `string`
 - [x] `[]string`
 - [x] nested structures
 - [x] net.TCPAddr
 - [x] net.IP
 - [x] time.Duration
 - [x] regexp.Regexp
 - [ ] map[string]string
 - [ ] map[string]int

## Custom types:
 - [x] HexBytes

 - [x] count
 - [ ] ipmask
 - [ ] enum values
 - [ ] enum list values
 - [ ] file
 - [ ] file list
 - [ ] url
 - [ ] url list
 - [ ] units (bytes 1kb = 1024b, speed, etc)

## Supported features matrix:

| Name | Hidden | Deprecated | Short | Env |
| --- | --- | --- | --- | --- |
| flag | - | - | - | - |
| pflag | [x] | [x] | [x] | - |
| kingpin | [x] | [ ] | [x] | [x] |
| urfave | [x] | - | [x] | [x] |
| cobra | [x] | [x] | [x] | - |
| viper | [ ] | [ ] | [ ] | [ ] |

  \[x] - feature is supported and implemented
  
  `-` - feature can't be implemented for this cli library

Simple example for flag library:

```golang
package main


import (
	"log"
	"time"
	"flag"

	"github.com/octago/sflags/gen/gflag"
)

type httpConfig struct {
	Host    string `desc:"HTTP host"`
	Port    int
	SSL     bool
	Timeout time.Duration
}

type config struct {
	HTTP   httpConfig
}

func main() {
	cfg := &config{
		HTTP: httpConfig{
			Host:    "127.0.0.1",
			Port:    6000,
			SSL:     false,
			Timeout: 15 * time.Second,
		},
	}
	err := gflag.ParseToDef(cfg)
	if err != nil {
		log.Fatalf("err: %v", err)
	}
	flag.Parse()
}
```

That code generates next output:
```
go run ./main.go --help
Usage of _obj/exe/main:
  -http-host value
    	HTTP host (default 127.0.0.1)
  -http-port value
    	 (default 6000)
  -http-ssl

  -http-timeout value
    	 (default 15s)
exit status 2
```

Look at the other [examples](https://github.com/octago/sflags/blob/master/examples) for different flag libraries.

## Options for flag tag

The flag default key string is the struct field name but can be specified in the struct field's tag value.
The "flag" key in the struct field's tag value is the key name, followed by an optional comma and options. Examples:
```
// Field is ignored by this package.
Field int `flag:"-"`

// Field appears in flags as "myName".
Field int `flag:"myName"`

// If this field is from nested struct, prefix from parent struct will be ingored.
Field int `flag:"~myName"`

// You can set short name for flags by providing it's value after a space
// Prefixes will not be applied for short names.
Field int `flag:"myName a"`

// this field will be removed from generated help text.
Field int `flag:",hidden"`

// this field will be marked as deprecated in generated help text
Field int `flag:",deprecated"`
```

## Options for desc tag
If you specify description in description tag (`desc` by default) it will be used in USAGE section.

```
Addr string `desc:"HTTP address"`
```
this description produces something like:
```
  -addr value
    	HTTP host (default 127.0.0.1)
```

## Options for env tag


## Options for Parse function:

```
// DescTag sets custom description tag. It is "desc" by default.
func DescTag(val string)

// FlagTag sets custom flag tag. It is "flag" be default.
func FlagTag(val string)

// Prefix sets prefix that will be applied for all flags (if they are not marked as ~).
func Prefix(val string)

// EnvPrefix sets prefix that will be applied for all environment variables (if they are not marked as ~).
func EnvPrefix(val string)

// FlagDivider sets custom divider for flags. It is dash by default. e.g. "flag-name".
func FlagDivider(val string)

// EnvDivider sets custom divider for environment variables.
// It is underscore by default. e.g. "ENV_NAME".
func EnvDivider(val string)

// Validator sets validator function for flags.
// Check existed validators in sflags/validator package.
func Validator(val ValidateFunc)

// Set to false if you don't want anonymous structure fields to be flatten.
func Flatten(val bool)
```


## Known issues

 - kingpin doesn't pass value for boolean arguments. Counter can't get initial value from arguments.
