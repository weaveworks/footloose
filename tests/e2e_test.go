package tests

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func find(dir string) ([]string, error) {
	var files []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		switch {
		case err != nil:
			return err
		case info.IsDir():
			return nil
		case strings.HasSuffix(path, "~"):
			return nil
		}
		files = append(files, strings.TrimPrefix(path, dir))
		return nil
	})

	return files, err
}

// test is a end to end test, corresponding to one test-$testname.cmd file.
type test struct {
	file string // name of the test file (test-*.cmd), without the extentension.
}

func newTest(testFile string) *test {
	ext := filepath.Ext(testFile)
	file := testFile[:len(testFile)-len(ext)]
	return &test{
		file: file,
	}
}

func (t *test) name() string {
	return t.file
}

func exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func (t *test) shouldErrorOut() bool {
	return exists(t.file + ".error")
}

func (t *test) shouldSkip() bool {
	return exists(t.file + ".skip")
}

func (t *test) outputDir() string {
	return t.file + ".got"
}

type cmd struct {
	name string
	args []string
	// should we capture the command output to be tested against the golden
	// output?
	captureOutput bool
}

func (t *test) parseCmd(line string) cmd {
	parts := strings.Split(line, " ")
	replacer := strings.NewReplacer(
		"%d", t.outputDir(),
		"%t", t.name(),
	)
	// Replace special strings
	for i := range parts {
		parts[i] = replacer.Replace(parts[i])
	}

	cmd := cmd{}
	switch parts[0] {
	case "%out":
		cmd.captureOutput = true
		parts = parts[1:]
	}

	cmd.name = parts[0]
	cmd.args = parts[1:]
	return cmd

}

func (t *test) run() (string, error) {
	f, err := os.Open(t.file + ".cmd")
	if err != nil {
		return "", err
	}
	defer f.Close()

	var capturedOutput strings.Builder

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		testCmd := t.parseCmd(scanner.Text())
		cmd := exec.Command(testCmd.name, testCmd.args...)
		if testCmd.captureOutput {
			output, err := cmd.CombinedOutput()
			if err != nil {
				// Display the captured output in case of failure.
				fmt.Print(string(output))
				return "", err
			}
			capturedOutput.Write(output)
		} else {
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				return "", err
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return capturedOutput.String(), nil
}

func runTest(t *testing.T, test *test) {
	base := test.file
	goldenDir := base + ".golden"
	gotDir := base + ".got"

	if test.shouldSkip() {
		return
	}

	output, err := test.run()

	// 0. Check process exit code.
	if test.shouldErrorOut() {
		_, ok := err.(*exec.ExitError)
		assert.True(t, ok, err.Error())
	} else {
		if err != nil {
			fmt.Print(string(output))
		}
		assert.NoError(t, err)
	}

	// 1. Compare stdout/err.
	golden, _ := ioutil.ReadFile(test.file + ".golden.output")
	assert.Equal(t, string(golden), string(output))

	// 2. Compare produced files.
	goldenFiles, _ := find(goldenDir)
	gotFiles, _ := find(gotDir)

	// 2. a) Compare the list of files.
	if !assert.Equal(t, goldenFiles, gotFiles) {
		assert.FailNow(t, "generated files not equivalent; bail")
	}

	// 2. b) Compare file content.
	for i := range goldenFiles {
		golden, err := ioutil.ReadFile(goldenDir + goldenFiles[i])
		assert.NoError(t, err)
		got, err := ioutil.ReadFile(gotDir + gotFiles[i])
		assert.NoError(t, err)

		assert.Equal(t, string(golden), string(got))
	}
}

func listTestFiles(t *testing.T) []string {
	files, err := filepath.Glob("test-*.cmd")
	assert.NoError(t, err)

	sort.Strings(files)
	return files
}

func TestEndToEnd(t *testing.T) {
	files := listTestFiles(t)

	for _, file := range files {
		test := newTest(file)
		t.Run(test.name(), func(t *testing.T) {
			runTest(t, test)
		})
	}
}
