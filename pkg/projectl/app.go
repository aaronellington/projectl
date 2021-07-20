package projectl

import (
	"github.com/aaronellington/projectl/pkg/configuration"
	"github.com/aaronellington/projectl/pkg/generators"
	"github.com/aaronellington/projectl/pkg/projector"
)

// App is the projectl app
type App struct{}

// Execute the app
func (app *App) Execute() error {
	config, err := configuration.NewConfig(".projectl/config.json")
	if err != nil {
		return err
	}

	service, err := projector.NewService()
	if err != nil {
		return err
	}

	config.DistedFiles = append(config.DistedFiles, "/app/config/parameters.yml")
	config.DistedFiles = append(config.DistedFiles, "/.env")

	service.DistedFiles = config.DistedFiles

	service.Generators = append(service.Generators, []projector.Generator{
		generators.NewGitignore(service, config),
		generators.NewMakefile(service, config),
		&generators.GithubWorkflow{},
		&generators.EslintGenerator{},
		&generators.PHPConfig{},
	}...)

	if config.DockerName != "" {
		service.Generators = append(service.Generators, &generators.Dockerfile{
			Port:   config.DockerPort,
			Target: config.DockerTarget,
			Custom: config.CustomDockerFile,
		})
	}

	return service.Generate()
}
