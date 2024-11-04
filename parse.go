package iac

import (
	"bytes"
	"fmt"
	"github.com/nullstone-io/iac/config"
	yaml2 "github.com/nullstone-io/iac/yaml"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
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
	// Config contains `.nullstone/config.yml` validated and normalized as Nullstone objects
	// This was added to the state since TemporalIacSync.IacSync Config is intentionally redacted from json
	Config *config.EnvConfiguration `json:"config"`

	// Overrides contains `.nullstone/<env|previews>.yml` validated and normalized as Nullstone objects
	// This was added to the state since TemporalIacSync.IacSync Overrides is intentionally redacted from json
	Overrides map[string]*config.EnvConfiguration `json:"overrides"`
}

func (r ParseMapResult) BlockNames(env types.Environment) map[string]bool {
	blockNames := map[string]bool{}
	if r.Config != nil {
		blockNames = r.Config.BlockNames()
	}
	envName := env.Name
	if env.Type == types.EnvTypePreview {
		envName = "previews"
	}
	if cur, _ := r.Overrides[envName]; cur != nil {
		for k := range cur.BlockNames() {
			blockNames[k] = true
		}
	}
	return blockNames
}

func ParseMap(repoUrl, repoName string, files map[string]string) (ParseMapResult, error) {
	result := ParseMapResult{
		Config:    nil,
		Overrides: map[string]*config.EnvConfiguration{},
	}

	for filename, raw := range files {
		desc, isOverrides := getConfigFileDescription(filename)
		parsed, err := ParseConfig(repoUrl, repoName, filename, isOverrides, bytes.NewBufferString(raw))
		if err != nil {
			return result, err
		}
		if isOverrides {
			result.Overrides[desc] = parsed
		} else {
			result.Config = parsed
		}
	}
	return result, nil
}

func ParseConfig(repoUrl, repoName, filename string, isOverrides bool, r io.Reader) (*config.EnvConfiguration, error) {
	decoder := yaml.NewDecoder(r)
	var obj yaml2.EnvConfiguration
	if err := decoder.Decode(&obj); err != nil {
		return nil, InvalidYamlError{ParseContext: repoName, FileName: filename, Err: err}
	}
	return config.ConvertConfiguration(repoUrl, repoName, filename, isOverrides, obj), nil
}

func ParseConfigFile(repoUrl, repoName, filename string, isOverrides bool) (*config.EnvConfiguration, error) {
	raw, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return ParseConfig(repoUrl, repoName, filename, isOverrides, bytes.NewReader(raw))
}

func ParseConfigDir(repoUrl, repoName, dir string) (*ParseMapResult, error) {
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
		filename := entry.Name()
		if entry.IsDir() || !isYmlFile(filename) {
			continue
		}
		desc, isOverrides := getConfigFileDescription(filename)
		ec, err := ParseConfigFile(repoUrl, repoName, filepath.Join(dir, filename), isOverrides)
		if err != nil {
			return nil, fmt.Errorf("cannot parse config file: %w", err)
		}
		if isOverrides {
			pmr.Overrides[desc] = ec
		} else {
			pmr.Config = ec
		}
	}
	return pmr, nil
}

func getConfigFileDescription(filepath string) (string, bool) {
	_, filename := path.Split(filepath)
	woExt := strings.TrimSuffix(filename, path.Ext(filename))
	return woExt, woExt != "config"
}

func isYmlFile(filename string) bool {
	switch filepath.Ext(filename) {
	case ".yml":
		return true
	case ".yaml":
		return true
	}
	return false
}
