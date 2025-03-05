# SberMortgageCalculator

SberMortgageCalculator is a Go-based project for mortgage calculation, leveraging Docker for development, linting, testing, and deployment. This repository includes a simple `Makefile` to facilitate common tasks like building images, running tests, and managing development environments.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Getting Started](#getting-started)
  - [Build and Run](#build-and-run)
  - [Development](#development)
  - [Testing](#testing)
  - [Linting](#linting)
- [Makefile Targets](#makefile-targets)

---

## Prerequisites

To work with this project, ensure you have the following installed:

- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)
- [GNU Make](https://www.gnu.org/software/make/)
- Go programming language (optional, only needed for local non-Dockerized development)

---

## Getting Started

This project uses a `Makefile` to streamline common workflows.

### Build and Run

To create and deploy an environment for testing the application image and launch the project, run:

```bash
make dev
```

This will:

1. Clean up any existing Docker artifacts.
2. Build the release image.
3. Start the application using Docker Compose in detached mode.

To view live logs of the application, use:

```bash
make logs
```

To stop the development environment, run:

```bash
make stop_dev
```

### Development

You can work on the project in a Dockerized environment. Use the following commands during development:

- Rebuild the project: `make image`
- Launching the app and swagger containers: `make dev`
- Stopping the app container and swagger: `make stop_dev`

### Testing

To execute tests in a Dockerized environment:

```bash
make test
```

This will:

1. Build the testing image (if not already built).
2. Run tests with `go test`, using `tparse` for enhanced output.
3. Display test results, including code coverage.

### Linting

Linting is performed using `golangci-lint`. To lint the project, run:

```bash
make lint
```

This will:

1. Build the linter image (if not already built).
2. Run `golangci-lint` over the codebase.

---

## Makefile Targets

| Target            | Description                                                                |
| ----------------- | -------------------------------------------------------------------------- |
| `image_lint`    | Build the Docker image for `golangci-lint` if it does not exist.         |
| `lint`          | Run the linter inside a Docker container using `golangci-lint`.          |
| `image`         | Build the release Docker image.                                            |
| `dev`           | Clean, build, and start the application in a Dockerized development setup. |
| `logs`          | Stream logs from the `calculator` container.                             |
| `stop_dev`      | Stop the Docker Compose development setup.                                 |
| `run`           | Run the image in standalone mode.                                          |
| `deps`          | Add project dependencies to the local `vendor` folder using Go modules.  |
| `image_testing` | Build the testing image if it does not exist.                              |
| `test`          | Run all tests with coverage using `tparse`.                              |
| `clean`         | Remove all tagged images and dangling Docker artifacts.                    |

## Notes

- **Dockerized Development**: The use of Docker ensures a consistent development environment across different machines.
- **Date-Based Tagging**: The release images are tagged with the current date (`YYYYMMDD`) for versioning purposes.
- **Clean Command**: The `make clean` command will attempt to remove all dangling Docker images to keep your system tidy, but unused images must be removed manually in some cases.

```

```
