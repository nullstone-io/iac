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

type ParseMapResult struct {
	Config    *config.EnvConfiguration
	Overrides map[string]overrides.EnvOverrides
}

func ParseMap(parseContext string, files map[string]string) (ParseMapResult, error) {
	result := ParseMapResult{
		Config:    nil,
		Overrides: map[string]overrides.EnvOverrides{},
	}

	for filepath, raw := range files {
		desc := getConfigFileDescription(filepath)
		if desc == "config" {
			parsed, err := ParseConfig(parseContext, filepath, bytes.NewBufferString(raw))
			if err != nil {
				return result, err
			}
			result.Config = &parsed
		} else {
			eo, err := ParseOverrides(parseContext, filepath, bytes.NewBufferString(raw))
			if err != nil {
				return result, err
			}
			result.Overrides[desc] = eo
		}
	}
	return result, nil
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
