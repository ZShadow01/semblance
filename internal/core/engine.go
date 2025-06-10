package core

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/ZShadow01/semblance/internal/filesystem"
	"github.com/go-git/go-git/v5"
)

type Project struct {
	Name 		string
	Path 		string
	Description string
}

//go:embed all:templates/**
var templates embed.FS

func CreateProject(project Project, template string) error {
	// Copy all contents from the template
	err := GenerateByDefaultTemplate(project, template)
	if err != nil {
		return err
	}

	// Initialize a git repo
	_, err = git.PlainInit(project.Path, false)
	if err != nil {
		return fmt.Errorf("failed to initialize the git repo: %w", err)
	}

	return nil
}


func GenerateByDefaultTemplate(project Project, template string) error {
	templatePath, err := isValidDefaultTemplate(template)
	if err != nil {
		return err
	}

	return fs.WalkDir(templates, templatePath, func(currentPath string, entry fs.DirEntry, _ error) error {
		rel, err := filepath.Rel(templatePath, currentPath)
		if err != nil {
			return fmt.Errorf("could not compute relative path: %v", err)
		}

		// Rename target path if filename starts with '_'
		target := filepath.Join(project.Path, rel)

		filename := filepath.Base(target)
		if strings.HasPrefix(filename, "_") {
			// Don't add files that are called '_keep'
			if filename == "_keep" {
				return nil
			}

			filename = strings.Replace(filename, "_", ".", 1)
			target = filepath.Join(filepath.Dir(target), filename)
		}

		// Add directory or file
		if entry.IsDir() {
			return filesystem.CreateDirectory(target)
		}

		return RenderTemplate(templates, currentPath, target, project)
	})
}


func isValidDefaultTemplate(template string) (string, error) {
	templatePath := path.Join("templates", template)

	info, err := fs.Stat(templates, templatePath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("the template '%s' does not exist", template)
		}
		return "", fmt.Errorf("unexpected error occurred trying to read the template '%s': %v", template, err)
	}

	if !info.IsDir() {
		return "", fmt.Errorf("file type of '%s' is not a directory", template)
	}

	return templatePath, nil
}


func RenderTemplate(fileSystem fs.FS, source, destination string, data any) error {
	// Create new file
	file, err := filesystem.CreateFile(destination)
	if err != nil {
		return err
	}
	defer file.Close()

	// Setup template
	tmpl, err := template.ParseFS(fileSystem, source)
	if err != nil {
		return fmt.Errorf("failed to setup template file %s: %w", source, err)
	}

	// Insert additional information from the project object into the template file
	err = tmpl.Execute(file, data)
	if err != nil {
		return fmt.Errorf("failed to execute template modifier on %s: %w", source, err)
	}

	return nil
}
