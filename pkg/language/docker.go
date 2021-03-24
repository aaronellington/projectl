package language

import (
	"io/ioutil"
)

// NewDocker generates a ready-to-use StateGo
func NewDocker() (*Docker, error) {
	languageDocker := &Docker{}

	_, err := ioutil.ReadFile("Dockerfile")
	if err != nil {
		// Dockerfile not able to be opened,
		// this probably means it's not a Docker project
		return languageDocker, nil
	}

	languageDocker.Enabled = true

	return languageDocker, nil
}

// Docker is the state of the Docker project
type Docker struct {
	Enabled bool
}
