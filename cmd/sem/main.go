package main

import (
	"log"
	"path/filepath"
	"strings"

	"github.com/ZShadow01/semblance/internal/core"
	"github.com/spf13/pflag"
)

func main() {
	var template, name, description string
	pflag.StringVar(&template, "template", "", "the project template")
	pflag.StringVar(&name, "name", "", "the project name")
	pflag.StringVar(&description, "description", "", "the project description (short)")

	pflag.Parse()

	args := pflag.Args()
	if len(args) == 0 {
		log.Fatal("Missing path argument")
	}

	if template == "" {
		log.Fatal("Missing required project template")
	}

	if name == "" {
		name = filepath.Base(args[0])
	}

	name = strings.ToLower(name)

	project := core.Project{
		Name: name,
		Path: args[0],
		Description: description,
	}

	err := core.CreateProject(project, template)
	if err != nil {
		log.Fatalf("failed to create a new project: %v\n", err)
	}

	log.Printf("Successfully created a new project with the '%s' template\n", template)
}
