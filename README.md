# Semblance

A CLI tool to create projects with set templates. Add new templates by creating a new folder in `templates/` and adding the files and folders you want to have in your template.

## Features

-   Fast and lightweight project creator with templates
-   Easy to use

## Installation

### Prerequisites

-   `make` or `go`

To install using `go install`:

```bash
go install github.com/ZShadow01/semblance@latest
```

_Or clone and build manually_

```bash
git clone https://github.com/ZShadow01/semblance.git
cd semblance

# either use make
make

# or go's built-in go build command
go build -o ./bin/sem ./cmd/sem/main.go
```

## Usage

To build a project, specify the project path and the template:

```bash
# command to create a project with the default nodejs template
sem --template=nodejs SampleNodeJS
```

Project names are automatically set to lowercase, so `SampleNodeJS` becomes `samplenodejs`. It is possible to explicitly specify a project name separate from the project path

```bash
sem --template=nodejs SampleNodeJS --name="sample-node-js"
```

This will create a `nodejs` project at `./SampleNodeJS` with the project name set to `sample-node-js`.

Currently available project templates:

-   `c`
-   `nodejs`
