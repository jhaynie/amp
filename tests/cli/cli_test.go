package cli

import (
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
	lookupDir   = "./lookup"
	setupDir    = "./setup"
	sampleDir   = "./samples"
	tearDownDir = "./tearDown"
	wg          sync.WaitGroup
	regexMap    map[string]string
)

// read, parse and execute test commands
func TestCliCmds(t *testing.T) {

	// test suite timeout
	suiteTimeout := "10m"
	duration, err := time.ParseDuration(suiteTimeout)
	if err != nil {
		t.Errorf("Unable to create duration for timeout: Suite. Error:", err)
		return
	}
	// create test suite context
	cancel := createTimeout(t, duration, "Suite")
	defer cancel()

	// parse regexes
	regexMap, err = parseLookup(lookupDir)
	if err != nil {
		t.Errorf("Unable to load lookup specs, reason:", err)
		return
	}

	// create setup timeout and parse setup specs
	setupTimeout := "8m"
	setup, err := createTestSpecs(setupDir, setupTimeout)
	if err != nil {
		t.Errorf("Unable to create setup specs, reason:", err)
		return
	}

	// create samples timeout and parse sample specs
	sampleTimeout := "30s"
	samples, err := createTestSpecs(sampleDir, sampleTimeout)
	if err != nil {
		t.Errorf("Unable to create sample specs, reason:", err)
		return
	}

	// create teardown timeout and parse tearDown specs
	tearDownTimeout := "1.5m"
	tearDown, err := createTestSpecs(tearDownDir, tearDownTimeout)
	if err != nil {
		t.Errorf("Unable to create tearDown specs, reason:", err)
		return
	}

	noOfSpecs := len(samples)

	runFramework(t, setup)
	wg.Add(noOfSpecs)
	runTests(t, samples)
	wg.Wait()
	runFramework(t, tearDown)
}

func createTestSpecs(directory string, timeout string) ([]*TestSpec, error) {
	// test spec timeout
	duration, err := time.ParseDuration(timeout)
	if err != nil {
		return nil, fmt.Errorf("Unable to create duration for timeout: %s. Error: %v", directory, err)
	}
	// parse tests
	testSpecs, err := parseSpec(directory, duration)
	if err != nil {
		return nil, fmt.Errorf("Unable to load test specs: %s. reason: %v", directory, err)
	}
	return testSpecs, nil
}

// runs a framework (setup/tearDown)
func runFramework(t *testing.T, commands []*TestSpec) {
	for _, command := range commands {
		runFrameworkSpec(t, command)
	}
}

// runs test commands
func runTests(t *testing.T, samples []*TestSpec) {
	for _, sample := range samples {
		go runSampleSpec(t, sample)
	}
}

// execute framework commands and check for timeout, delay and retry
func runFrameworkSpec(t *testing.T, test *TestSpec) {

	var cache = map[string]string{}

	// create test spec context
	cancel := createTimeout(t, test.Timeout, test.Name)
	defer cancel()

	// iterate through all the testSpec
	for _, command := range test.Commands {
		runCmdSpec(t, command, cache)
	}
}

// execute sample commands and decrement waitgroup counter
func runSampleSpec(t *testing.T, test *TestSpec) {

	var cache = map[string]string{}

	// decrements wg counter
	defer wg.Done()

	// create test spec context
	cancel := createTimeout(t, test.Timeout, test.Name)
	defer cancel()

	// iterate through all the testSpec
	for _, command := range test.Commands {
		runCmdSpec(t, command, cache)
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
