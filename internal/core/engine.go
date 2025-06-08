package core

import (
	"embed"
	"errors"
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
	err := CopyTemplateContents(project, template)
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


func CopyTemplateContents(project Project, template string) error {
	// Virtual path to templates/template
	templatePath := path.Join("templates", template)

	// Check whether templatePath exists
	info, err := fs.Stat(templates, templatePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("the template '%s' does not exist: %w", template, err)
		}
		return fmt.Errorf("unexpected error occurred when reading '%s': %v", templatePath, err)
	}

	// Check whether templatePath is a directory
	if !info.IsDir() {
		return errors.New("unexpectedly read a file instead of a directory " + templatePath)
	}

	// Walk through all of the template directory's contents
	// Real paths use filepath, virtual path use path
	return fs.WalkDir(templates, templatePath, func(virtualPath string, d fs.DirEntry, _ error) error {
		relPath, err := filepath.Rel(templatePath, virtualPath)
		if err != nil {
			return err
		}
		targetPath := filepath.Join(project.Path, relPath)

		// All files starting with _ are "hidden" files
		fileName := filepath.Base(targetPath)
		if strings.HasPrefix(fileName, "_") {
			if fileName == "_keep" {
				return nil
			}

			targetPath = path.Join(filepath.Dir(targetPath), strings.Replace(fileName, "_", ".", 1))
		}

		if d.IsDir() {
			err := filesystem.CreateDirectory(targetPath)
			if err != nil {
				return err
			}
		} else {
			// Render the template
			err := RenderTemplate(templates, virtualPath, targetPath, project)
			if err != nil {
				return err
			}
		}

		return nil
	})
}


func RenderTemplate(fileSystem fs.FS, virtualSrc, destination string, data any) error {
	// Create new file
	file, err := filesystem.CreateFile(destination)
	if err != nil {
		return err
	}
	defer file.Close()

	// Setup template
	tmpl, err := template.ParseFS(fileSystem, virtualSrc)
	if err != nil {
		return fmt.Errorf("failed to setup template file %s: %w", virtualSrc, err)
	}

	// Insert additional information from the project object into the template file
	err = tmpl.Execute(file, data)
	if err != nil {
		return fmt.Errorf("failed to execute template modifier on %s: %w", virtualSrc, err)
	}

	return nil
}
