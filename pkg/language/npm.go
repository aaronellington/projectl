package language

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// NewNpm generates a ready-to-use StateNpm
func NewNpm() (*Npm, error) {
	languageNpm := &Npm{}
	fileBytes, err := ioutil.ReadFile("package.json")
	if err != nil {
		// package.json file not able to be opened,
		// this probably means it's not a npm project
		return languageNpm, nil
	}

	languageNpm.Enabled = true

	err = json.Unmarshal(fileBytes, &languageNpm.packageJSON)
	if err != nil {
		return nil, fmt.Errorf("%w while parsing package.json", err)
	}

	lockFileBytes, err := ioutil.ReadFile("package-lock.json")
	if err == nil {
		err = json.Unmarshal(lockFileBytes, &languageNpm.packageLockJSON)
		if err != nil {
			return nil, fmt.Errorf("%w while parsing package-lock.json", err)
		}
	}

	return languageNpm, nil
}

// Npm is the state of the npm project
type Npm struct {
	Enabled         bool
	packageJSON     PackageDotJSON
	packageLockJSON PackageLockDotJSON
}

// HasScript checks if a script is defined
func (languageNpm Npm) HasScript(targetName string) bool {
	for scriptName := range languageNpm.packageJSON.Scripts {
		if scriptName == targetName {
			return true
		}
	}

	return false
}

// HasDependency checks if the project has a dependency
func (languageNpm Npm) HasDependency(targetName string) bool {
	_, found := languageNpm.packageLockJSON.Dependencies[targetName]

	return found
}

// PackageDotJSON is the structure of the package.json file
type PackageDotJSON struct {
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
	Scripts         map[string]string `json:"scripts"`
}

// PackageLockDotJSON is the structure of the package-lock.json file
type PackageLockDotJSON struct {
	Dependencies map[string]interface{} `json:"dependencies"`
}
