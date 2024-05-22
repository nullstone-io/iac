package iac

import (
	"bytes"
	"fmt"
	"github.com/nullstone-io/iac/config"
	"github.com/nullstone-io/iac/overrides"
	yaml2 "github.com/nullstone-io/iac/yaml"
	"gopkg.in/yaml.v3"
	"io"
	"path"
	"strings"
)

const (
	ConfigFileTypeConfig   = "config"
	ConfigFileTypeOverride = "previews"
)

type InvalidYamlError struct {
	ParseContext string
	FileName     string
	Err          error
}

func (e InvalidYamlError) Error() string {
	return fmt.Sprintf("Error parsing YAML file (%s) from (%s): %v", e.FileName, e.ParseContext, e.Err)
}

func (e InvalidYamlError) Unwrap() error {
	return e.Err
}

func ParseMap(parseContext string, files map[string][]byte) (config.EnvConfiguration, map[string]overrides.EnvOverrides, error) {
	ec := config.EnvConfiguration{}
	eos := map[string]overrides.EnvOverrides{}

	for filepath, raw := range files {
		desc := getConfigFileDescription(filepath)
		if desc == "config" {
			var err error
			if ec, err = ParseConfig(parseContext, filepath, bytes.NewReader(raw)); err != nil {
				return ec, eos, err
			}
		} else {
			if eo, err := ParseOverrides(parseContext, filepath, bytes.NewReader(raw)); err != nil {
				return ec, eos, err
			} else {
				eos[desc] = eo
			}
		}
	}
	return ec, eos, nil
}

func getConfigFileDescription(filepath string) string {
	_, filename := path.Split(filepath)
	woExt := strings.TrimSuffix(filename, path.Ext(filename))
	return woExt
}

func ParseConfig(parseContext, filename string, r io.Reader) (config.EnvConfiguration, error) {
	decoder := yaml.NewDecoder(r)
	var obj yaml2.EnvConfiguration
	if err := decoder.Decode(&obj); err != nil {
		return config.EnvConfiguration{}, InvalidYamlError{ParseContext: parseContext, FileName: filename, Err: err}
	}
	return config.ConvertConfiguration(parseContext, filename, obj), nil
}

func ParseOverrides(parseContext, filename string, r io.Reader) (overrides.EnvOverrides, error) {
	decoder := yaml.NewDecoder(r)
	var obj yaml2.EnvOverrides
	if err := decoder.Decode(&obj); err != nil {
		return overrides.EnvOverrides{}, InvalidYamlError{ParseContext: parseContext, FileName: filename, Err: err}
	}
	return overrides.ConvertOverrides(parseContext, filename, obj), nil
}
