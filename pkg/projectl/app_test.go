package projectl_test

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"testing"

	"github.com/aaronellington/projectl/pkg/configuration"
	"github.com/aaronellington/projectl/pkg/projectl"
)

type TestCase struct {
	Path          string
	ExpectedError error
}

func TestExecute(t *testing.T) {
	testCases := []TestCase{
		{
			Path:          buildPath("full_go"),
			ExpectedError: nil,
		},
		{
			Path:          buildPath("full_php"),
			ExpectedError: nil,
		},
		{
			Path:          buildPath("missing_config_file"),
			ExpectedError: configuration.ErrMissingConfigFile,
		},
		{
			Path:          buildPath("invalid_config_file"),
			ExpectedError: configuration.ErrInvalidConfigFile,
		},
	}

	for _, testCase := range testCases {
		testName := path.Base(testCase.Path)
		t.Run(testName, func(tt *testing.T) {
			testProject(tt, testCase)
		})
	}
}

func testProject(t *testing.T, testCase TestCase) {
	_ = os.Chdir(testCase.Path)
	_, _ = os.Create("Dockerfile")

	app := projectl.App{}

	err := app.Execute()
	if !errors.Is(err, testCase.ExpectedError) {
		if testCase.ExpectedError == nil {
			t.Fatalf("Unexpected error: %v", err)
		} else {
			t.Fatalf("Incorrect error: %v Got: %v", testCase.ExpectedError.Error(), err)
		}
	}
	if testCase.ExpectedError != nil {
		return
	}

	filesToCompare := []string{
		".gitignore",
		"Makefile",
		".github/workflows/main.yml",
		"Dockerfile",
	}

	for _, fileToCompare := range filesToCompare {
		if err := compareTwoFiles(fileToCompare); err != nil {
			t.Fatal(err)
		}
	}
}

func buildPath(testName string) string {
	wd, _ := os.Getwd()
	return wd + "/test_projects/" + testName
}

func compareTwoFiles(sourcePath string) error {
	targetPath := sourcePath + "-target"

	targetFile, err := os.Open(targetPath)
	if err != nil {
		return err
	}
	defer targetFile.Close()

	targetBytes, err := ioutil.ReadAll(targetFile)
	if err != nil {
		return err
	}

	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	sourceBytes, err := ioutil.ReadAll(sourceFile)
	if err != nil {
		return err
	}

	if !reflect.DeepEqual(sourceBytes, targetBytes) {
		return errors.New("File contents do not match target: " + sourcePath)
	}

	return nil
}
