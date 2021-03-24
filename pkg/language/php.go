package language

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"
)

// NewPHP generates a ready-to-use StatePHP
func NewPHP() (*PHP, error) {
	languagePHP := &PHP{}
	fileBytes, err := ioutil.ReadFile("composer.json")
	if err != nil {
		// composer.json file not able to be opened,
		// this probably means it's not a npm project
		return languagePHP, nil
	}

	languagePHP.Enabled = true

	err = json.Unmarshal(fileBytes, &languagePHP.composerJSON)
	if err != nil {
		return nil, fmt.Errorf("%w while parsing composer.json", err)
	}

	return languagePHP, nil
}

// PHP is the state of the php project
type PHP struct {
	Enabled      bool
	composerJSON ComposerDotJSON
}

// IsSymfony3 checks if the project is a symfony 3.4 project or not
func (languagePHP PHP) IsSymfony3() bool {
	var versionMatcher = regexp.MustCompile(`(?m)^(\W)?3\.4`)
	for packageName, version := range languagePHP.composerJSON.Require {
		if packageName == "symfony/symfony" && versionMatcher.MatchString(version) {
			return true
		}
	}

	return false
}

// ComposerDotJSON is the structure of the composer.json file
type ComposerDotJSON struct {
	Autoload   ComposerDotJSONAutoload `json:"autoload"`
	Require    map[string]string       `json:"require"`
	RequireDev map[string]string       `json:"require-dev"`
}

// ComposerDotJSONConfig is the config struct for the composer.json file
type ComposerDotJSONConfig struct {
	BinDir string `json:"bin-dir"`
}

// ComposerDotJSONAutoload is the autoload struct for the composer.json file
type ComposerDotJSONAutoload struct {
	PSR0 map[string]string `json:"psr-0"`
	PSR4 map[string]string `json:"psr-4"`
}
