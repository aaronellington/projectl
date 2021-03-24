package language

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"golang.org/x/mod/modfile"
)

// NewGo generates a ready-to-use StateGo
func NewGo() (*Go, error) {
	languageGo := &Go{
		Targets: make(map[string]string),
	}

	fileBytes, err := ioutil.ReadFile("go.mod")
	if err != nil {
		// go.mod file not able to be opened,
		// this probably means it's not a go project
		return languageGo, nil
	}

	languageGo.Enabled = true

	modfile, err := modfile.Parse("go.mod", fileBytes, nil)
	if err != nil {
		return nil, fmt.Errorf("%w while parsing go.mod", err)
	}

	languageGo.modfile = modfile

	if _, err := os.Open("main.go"); err == nil {
		languageGo.Targets["."] = path.Base(modfile.Module.Mod.Path)
	}

	if commandDirectories, err := ioutil.ReadDir("./cmd/"); err == nil {
		for _, commandDirectory := range commandDirectories {
			languageGo.Targets["./cmd/"+commandDirectory.Name()] = commandDirectory.Name()
		}
	}

	return languageGo, nil
}

// Go is the state of the go project
type Go struct {
	Enabled bool
	Targets map[string]string
	modfile *modfile.File
}

// Version gets the go version
func (languageGo *Go) Version() string {
	return languageGo.modfile.Go.Version
}

// ModRequired checks if a module is required
func (languageGo *Go) ModRequired(modPath string) bool {
	if languageGo.modfile == nil {
		return false
	}

	for _, require := range languageGo.modfile.Require {
		if require.Mod.Path == modPath {
			return true
		}
	}

	return false
}
