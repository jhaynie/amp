package cli

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"

	gexpect "github.com/Thomasrooney/gexpect"
)

// TestSpec contains all the CommandSpec objects
type TestSpec struct {
	Name     string
	Timeout  time.Duration
	Commands []CommandSpec
}

// CommandSpec defines the commands with arguments and options
type CommandSpec struct {
	Cmd         string   `yaml:"cmd"`
	Args        []string `yaml:"args"`
	Options     []string `yaml:"options"`
	Input 			[]string `yaml:"input"`
	Expectation string   `yaml:"expectation"`
	Skip        bool     `yaml:"skip"`
	Retry       int      `yaml:"retry"`
	Timeout     string   `yaml:"timeout"`
	Delay       string   `yaml:"delay"`
}

var (
	suiteDir = []string{"./suite"}
	regexMap map[string]string
)

func TestAllCmds(t *testing.T) {

	// setup = parse, run platform
	suites, err := parseSuite(suiteDir)
	if err != nil {
		t.Errorf("Unable to parse suites, reason: %v", err)
	}
	for _, suite := range suites {
		for _, setup := range suite.Setup {
			t.Run(setup.Name, func(t *testing.T) {
				runTestSpec(t, setup)
			})
		}
	}

	// Main tests
	for _, suite := range suites {
		for _, test := range suite.Tests {
			test := test
			t.Run(test.Name, func(t *testing.T) {
				t.Parallel()
				runTestSpec(t, test)
			})
		}
	}

	// teardown = run platform, clean goroutines
	for _, suite := range suites {
		for _, teardown := range suite.TearDown {
			t.Run(teardown.Name, func(t *testing.T) {
				runTestSpec(t, teardown)
			})
		}
	}
}

// execute commands and check for timeout, delay and retry
func runTestSpec(t *testing.T, test TestSpec) {

	var cache = map[string]string{}

	// create test spec context
	cancel := createTimeout(t, test.Timeout, test.Name)
	defer cancel()

	// iterate through all the testSpec
	for _, cmd := range test.Commands {
		runCmdSpec(t, cmd, cache)
	}
}

// TODO: add description for method, remove unnecessary comments
func runCmdSpec(t *testing.T, spec CommandSpec, cache map[string]string) {

	i := 0

	duration, err := time.ParseDuration(spec.Timeout)
	if err != nil {
		t.Fatal("Unable to create duration for timeout:", spec.Cmd, "Error:", err)
	}
	cancel := createTimeout(t, duration, spec.Cmd)
	defer cancel()

	cmd := generateCmdString(spec)

	// TODO: move into unmarshalling of yaml files in parse.go
	cmd, err = templating(cmd, cache)
	if err != nil {
		t.Fatal("Executing templating failed:", cmd, "Error:", err)
	}

	// TODO: move into unmarshalling of yaml files in parse.go
	regex, err := templating(regexMap[spec.Expectation], cache)
	if err != nil {
		t.Fatal("Executing templating failed:", spec.Expectation, "Error:", err)
	}

	// TODO: move into unmarshalling of yaml files in parse.go
	if spec.Expectation != "" && regexMap[spec.Expectation] == "" {
		t.Fatal("Unable to fetch regex for command:", cmd, "reason: no regex for given expectation:", spec.Expectation)
	}

	for i = 0; i <= spec.Retry; i++ {
		// err is set to nil, ensure no carryover
		err = nil

		output := runCmd(cmd, spec)
		fmt.Println(output)

		expectedOutput := regexp.MustCompile(regex)
		if !expectedOutput.MatchString(string(output)) {
			err = errors.New("Error: mismatched return:" + cmd + ", " + regex + ", " + output)
		}

		if spec.Delay != "" {
			del, err := time.ParseDuration(spec.Delay)
			if err != nil {
				t.Fatal("Invalid delay specified: ", spec.Delay, "Error:", err)
			}
			time.Sleep(del)
		}

		if err == nil {
			break
		}
	}
	if i > 1 {
		t.Log("The command :", cmd, "has re-run", i, "times.")
	}
	if err != nil {
		t.Fatal(err)
	}
}

func runCmd(cmd string, spec CommandSpec) string {
	child, err := gexpect.Spawn(cmd)
	if err != nil {
		panic(err)
	}
	for _, input := range spec.Input {
		child.Send(input + "\r")
		child.Expect(input)
	}
	output, _ := child.ReadUntil('$')
	fmt.Println(string(output))
	return string(output)
}

// concatenates struct fields into string
func generateCmdString(cmdSpec CommandSpec) string {
	cmdSplit := strings.Fields(cmdSpec.Cmd)
	optionsSplit := []string{}
	// Possible to have multiple options
	for _, val := range cmdSpec.Options {
		optionsSplit = append(optionsSplit, strings.Fields(val)...)
	}
	cmdString := append(cmdSplit, cmdSpec.Args...)
	cmdString = append(cmdString, optionsSplit...)
	return strings.Join(cmdString, " ")
}
