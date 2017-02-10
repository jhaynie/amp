package cli

import (
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

func parseSuite(dirs []string) ([]*SuiteSpec, error) {
	for _, dir := range dirs {
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			return nil, err
		}
		suites := []*SuiteSpec{}
		for _, file := range files {
			suite, err := createSuite(path.Join(dir, file.Name()))
			if err != nil {
				return nil, err
			}
			if suite != nil {
				suites = append(suites, suite)
			}
		}
		return suites, nil
	}
	return nil, nil
}

func createSuite(name string) (*SuiteSpec, error) {

	if filepath.Ext(name) != ".yml" {
		if filepath.Ext(name) == "" {
			// ignore file but don't fail
			return nil, nil
		}
		return nil, fmt.Errorf("Cannot parse non-yaml file: %s", name)
	}
	contents, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf("Unable to read yaml suite spec: %s. Error: %v", name, err)
	}
	suite := &SuiteSpec{}
	if err = yaml.Unmarshal(contents, suite); err != nil {
		return nil, fmt.Errorf("Unable to unmarshal yaml suite spec: %s. Error: %v", name, err)
	}

	regexMap, err = parseLookup(suite.LookupDir)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse yaml lookup spec: %s. Error: %v", suite.LookupDir, err)
	}

	setup, err := parseSpec(suite.SetupDir)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse yaml setup spec: %s. Error: %v", suite.SetupDir, err)
	}
	suite.Setup = append(suite.Setup, *setup...)

	for _, dir := range suite.TestDirs {
		test, err := parseSpec(dir)
		if err != nil {
			return nil, fmt.Errorf("Unable to parse yaml test spec: %s. Error: %v", test, err)
		}
		suite.Tests = append(suite.Tests, *test...)
	}

	tearDown, err := parseSpec(suite.TearDownDir)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse yaml teardown spec: %s. Error: %v", suite.TearDownDir, err)
	}

	suite.TearDown = append(suite.TearDown, *tearDown...)
	return suite, nil
}

// read lookup directory by parsing its contents
func parseLookup(dir string) (map[string]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	rgxMap := make(map[string]string)
	for _, file := range files {
		lookup, err := createLookup(path.Join(dir, file.Name()))
		if err != nil {
			return nil, err
		}
		for expect, rgx := range lookup {
			rgxMap[expect] = rgx
		}
	}
	return rgxMap, nil
}

// parse lookup directory and unmarshal its contents
func createLookup(name string) (map[string]string, error) {
	if filepath.Ext(name) != ".yml" {
		return nil, fmt.Errorf("Cannot parse non-yaml file: %s", name)
	}
	contents, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf("Unable to read yaml regex lookup: %s. Error: %v", name, err)
	}
	lookup := make(map[string]string)
	if err := yaml.Unmarshal(contents, &lookup); err != nil {
		return nil, fmt.Errorf("Unable to unmarshal yaml lookup: %s. Error: %v", name, err)
	}
	return lookup, nil
}

// read specs from directory by parsing its contents
func parseSpec(dir string) (*[]TestSpec, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	tests := []TestSpec{}
	for _, file := range files {
		test, err := createSpec(path.Join(dir, file.Name()))
		if err != nil {
			return nil, err
		}
		if test != nil {
			tests = append(tests, *test)
		}
	}
	return &tests, nil
}

// parse samples directory and unmarshal its contents
func createSpec(name string) (*TestSpec, error) {
	if filepath.Ext(name) != ".yml" {
		return nil, fmt.Errorf("Cannot parse non-yaml file: %s", name)
	}
	contents, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf("Unable to read yaml test spec: %s. Error: %v", name, err)
	}
	test := &TestSpec{}
	if err = yaml.Unmarshal(contents, &test); err != nil {
		return nil, fmt.Errorf("Unable to unmarshal yaml test spec: %s. Error: %v", name, err)
	}
	for i, _ := range test.Commands {
		if test.Commands[i].Timeout == "" {
			// default command spec timeout
			test.Commands[i].Timeout = "30s"
		}
		if test.Commands[i].Skip == true {
			test.Commands = append(test.Commands[:i], test.Commands[i+1:]...)
		}
	}
	return test, nil
}
