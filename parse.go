package iac

import (
	"bytes"
	"fmt"
	"github.com/nullstone-io/iac/config"
	yaml2 "github.com/nullstone-io/iac/yaml"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"path"
	"path/filepath"
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
	Overrides map[string]*config.EnvConfiguration
}

func ParseMap(parseContext string, files map[string]string) (ParseMapResult, error) {
	result := ParseMapResult{
		Config:    nil,
		Overrides: map[string]*config.EnvConfiguration{},
	}

	for filepath, raw := range files {
		desc := getConfigFileDescription(filepath)
		if desc == "config" {
			parsed, err := ParseConfig(parseContext, filepath, false, bytes.NewBufferString(raw))
			if err != nil {
				return result, err
			}
			result.Config = parsed
		} else {
			eo, err := ParseConfig(parseContext, filepath, true, bytes.NewBufferString(raw))
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

func ParseConfig(parseContext, filename string, isOverrides bool, r io.Reader) (*config.EnvConfiguration, error) {
	decoder := yaml.NewDecoder(r)
	var obj yaml2.EnvConfiguration
	if err := decoder.Decode(&obj); err != nil {
		return nil, InvalidYamlError{ParseContext: parseContext, FileName: filename, Err: err}
	}
	return config.ConvertConfiguration(parseContext, filename, isOverrides, obj), nil
}

func ParseConfigFile(parseContext, filename string, isOverrides bool) (*config.EnvConfiguration, error) {
	raw, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return ParseConfig(parseContext, filename, isOverrides, bytes.NewReader(raw))
}

func ParseConfigDir(dir string) (*ParseMapResult, error) {
	pmr := &ParseMapResult{
		Overrides: map[string]*config.EnvConfiguration{},
	}
	entries, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		return pmr, nil
	}
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".yml") {
			continue
		}
		isOverrides := entry.Name() != "config.yml"
		ec, err := ParseConfigFile("TestApplyChanges", filepath.Join(dir, entry.Name()), isOverrides)
		if err != nil {
			return nil, fmt.Errorf("cannot parse config file: %w", err)
		}
		if !isOverrides {
			pmr.Config = ec
		} else {
			pmr.Overrides[strings.TrimSuffix(entry.Name(), ".yml")] = ec
		}
	}
	return pmr, nil
}
