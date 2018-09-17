package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path"

	"github.com/UnnoTed/fileb0x/utils"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v2"
	"fmt"
)

// File holds config file info
type File struct {
	FilePath string
	Data     []byte
	Mode     string // "json" || "yaml" || "yml" || "toml"
}

// FromArg gets the json/yaml/toml file from args
func (f *File) FromArg(read bool) error {
	// (length - 1)
	arg := os.Args[len(os.Args)-1:][0]

	// get extension
	ext := path.Ext(arg)
	if len(ext) > 1 {
		ext = ext[1:] // remove dot
	}

	// when json/yaml/toml file isn't found on last arg
	// it searches for a ".json", ".yaml", ".yml" or ".toml" string in all args
	if ext != "json" && ext != "yaml" && ext != "yml" && ext != "toml" {
		// loop through args
		for _, a := range os.Args {
			// get extension
			ext := path.Ext(a)

			// check for valid extensions
			if ext == ".json" || ext == ".yaml" || ext == ".yml" || ext == ".toml" {
				f.Mode = ext[1:] // remove dot
				ext = f.Mode
				arg = a
				break
			}
		}
	} else {
		f.Mode = ext
	}

	// check if extension is json, yaml or toml
	// then get it's absolute path
	if ext == "json" || ext == "yaml" || ext == "yml" || ext == "toml" {
		f.FilePath = arg

		// so we can test without reading a file
		if read {
			if !utils.Exists(f.FilePath) {
				return errors.New("Error: I Can't find the config file at [" + f.FilePath + "]")
			}
		}
	} else {
		return errors.New("Error: You must specify a json, yaml or toml file")
	}

	return nil
}

// Parse gets the config file's content from File.Data
func (f *File) Parse() (*Config, error) {
	// remove comments
	f.RemoveJSONComments()

	to := &Config{}
	switch f.Mode {
	case "json":
		return to, json.Unmarshal(f.Data, to)
	case "yaml", "yml":
		return to, yaml.Unmarshal(f.Data, to)
	case "toml":
		return to, toml.Unmarshal(f.Data, to)
	default:
		return nil, fmt.Errorf("unknown mode '%s'", f.Mode)
	}
}

// Load the json/yaml file that was specified from args
// and transform it into a config struct
func (f *File) Load() (*Config, error) {
	var err error
	if !utils.Exists(f.FilePath) {
		return nil, errors.New("Error: I Can't find the config file at [" + f.FilePath + "]")
	}

	// read file
	f.Data, err = ioutil.ReadFile(f.FilePath)
	if err != nil {
		return nil, err
	}

	// parse file
	return f.Parse()
}

// RemoveJSONComments from the file
func (f *File) RemoveJSONComments() {
	if f.Mode == "json" {
		// remove inline comments
		f.Data = []byte(regexComments.ReplaceAllString(string(f.Data), ""))
	}
}
