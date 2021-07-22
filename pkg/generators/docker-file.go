package generators

import (
	"fmt"
	"os"

	"github.com/aaronellington/projectl/pkg/projector"
)

// Dockerfile generator
type Dockerfile struct {
	Port   int
	Target string
	Custom bool
}

// Generate the config file
func (dockerfile *Dockerfile) Generate(service *projector.Service) error {
	if dockerfile.Custom {
		return nil
	}

	file, err := os.Create("Dockerfile")
	if err != nil {
		return err
	}

	if service.Npm.Enabled {
		_, _ = file.WriteString(`FROM node:16-buster as nodeBuilder
WORKDIR /build-staging
COPY . .
RUN make clean-full
RUN make lint-npm test-npm build-npm

`)
	}

	if service.Go.Enabled {
		dockerfile.modeGo(service, file)
	} else if service.PHP.Enabled {
		dockerfile.modePHP(service, file)
	} else if service.Npm.Enabled {
		dockerfile.modeNPM(service, file)
	}

	if dockerfile.Port != 0 {
		_, _ = file.WriteString(fmt.Sprintf("EXPOSE %d\n", dockerfile.Port))
	}

	return nil
}

func (dockerfile *Dockerfile) modeGo(service *projector.Service, file *os.File) {
	targetBin := dockerfile.Target

	if targetBin == "" {
		for _, bin := range service.Go.Targets {
			targetBin = bin
			break
		}
	}

	_, _ = file.WriteString(`FROM golang:` + service.Go.Version() + `-buster as goBuilder
WORKDIR /build-staging
COPY . .
RUN make clean-full
`)

	if service.Npm.Enabled && service.Npm.HasScript("build") {
		_, _ = file.WriteString("COPY --from=nodeBuilder /build-staging/resources/dist/ /build-staging/resources/dist/\n")
	}
	_, _ = file.WriteString(`RUN make lint-go test-go build-go

FROM debian:buster
RUN apt-get update
RUN apt-get install -y ca-certificates
WORKDIR /app
`)

	_, _ = file.WriteString(`COPY --from=goBuilder /build-staging/var/` + targetBin + ` ./` + targetBin + `
CMD ["./` + targetBin + `"]
`)
}

func (dockerfile *Dockerfile) modePHP(service *projector.Service, file *os.File) {
	_, _ = file.WriteString(`FROM aaronellington/php-fpm-webserver:latest
COPY . .
RUN make clean-full
`)

	if service.Npm.Enabled {
		_, _ = file.WriteString("COPY --from=nodeBuilder /build-staging/public/build/ ./public/build/\n")
	}

	_, _ = file.WriteString(`RUN make lint-php test-php build-php-prod
RUN mkdir var
RUN chown www-data:www-data var
`)
}

func (dockerfile *Dockerfile) modeNPM(service *projector.Service, file *os.File) {
	_, _ = file.WriteString(`CMD ["npm", "run", "start"]
`)
}
