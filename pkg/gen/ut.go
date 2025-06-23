package gen

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

const (
	UtTmpDir = "/tmp/gen-ut"
)

func GetGenUtPrompt(goCode []byte) string {
	return fmt.Sprintf(`
	You are an expert Go developer and a meticulous test engineer.
	I will provide you with a Go source code file.
	Your task is to write comprehensive unit test cases for this Go file.

	Here are the requirements for the test cases:
	- Use the standard Go testing package ("testing").
	- For each public function/method, write at least one test case.
	- Consider edge cases, error handling, and typical usage scenarios.
	- Ensure test names follow Go conventions (e.g., "TestFunctionName").
	- Provide clear test descriptions within the tests.
	- Use "t.Run" for subtests where appropriate.
	- Do NOT include the original source code in your response. Only provide the test code.
	- Do include the package declaration or "import" statements at the very beginning;
	- If the code contains structs, write tests for their methods.
	- Aim for good code coverage.
	- If there is any function called from the interface object, using the "github.com/stretchr/testify/mock" package to mock the interface object.

	Here is the Go source code:

	%s
	`, string(goCode))
}

func GetCoverage(goCode []byte, testCode []byte, id uuid.UUID) (float64, error) {
	// create a tmp dir
	if err := os.MkdirAll(path.Join(UtTmpDir, id.String()), 0755); err != nil {
		return 0, err
	}
	// write the go code to a file
	if err := os.WriteFile(path.Join(UtTmpDir, id.String(), "main.go"), goCode, 0644); err != nil {
		return 0, err
	}
	// write the test code to a file
	if err := os.WriteFile(path.Join(UtTmpDir, id.String(), "main_test.go"), testCode, 0644); err != nil {
		return 0, err
	}
	// run go mod init genut
	cmd := exec.Command("go", "mod", "init", "genut")
	// set the working directory
	cmd.Dir = path.Join(UtTmpDir, id.String())
	// run the command
	if err := cmd.Run(); err != nil {
		return 0, err
	}
	// run go mod tidy
	cmd = exec.Command("go", "mod", "tidy")
	// set the working directory
	cmd.Dir = path.Join(UtTmpDir, id.String())
	// run the command
	if err := cmd.Run(); err != nil {
		return 0, err
	}
	// run go test -cover -coverprofile=coverage.out ./...
	cmd = exec.Command("go", "test", "-cover", "-coverprofile=coverage.out", "./...")
	// set the working directory
	cmd.Dir = path.Join(UtTmpDir, id.String())
	out, err := cmd.CombinedOutput()
	if err != nil {
		return 0, err
	}
	outStr := string(out)
	// find the coverage rate
	idx := strings.Index(outStr, "coverage: ")
	endIdx := strings.Index(outStr[idx:], `% of statements`)
	cover := outStr[idx+len("coverage: ") : idx+endIdx]
	// parse to float64
	coverStr := strings.TrimSpace(cover)
	rate, err := strconv.ParseFloat(coverStr, 64)
	if err != nil {
		return 0, err
	}
	return rate, nil
}
