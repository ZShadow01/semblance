package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"
	"text/template"
)

type TemplateTest struct {
	Name 	string 	`json:"name"`
	Exists 	bool 	`json:"exists"`
}

func TestCopyTemplateContents(t *testing.T) {
	data, err := os.ReadFile("testdata/sample-templates.json")
	if err != nil {
		t.Fatal(err)
	}

	var sampleTemplates []TemplateTest
	err = json.Unmarshal(data, &sampleTemplates)
	if err != nil {
		t.Fatal (err)
	}
	
	for _, sample := range sampleTemplates {
		tempDir := t.TempDir()

		project := Project{
			Name: "test",
			Path: filepath.Join(tempDir, "test"),
			Description: "Test",
		}

		err := CopyTemplateContents(project, sample.Name)
		if (sample.Exists && err != nil) || (!sample.Exists && err == nil) {
			t.Fatal(err)
		}

		if (!sample.Exists) {
			continue
		}

		templatePath := path.Join("templates", sample.Name)
		err = fs.WalkDir(templates, templatePath, func(p string, d fs.DirEntry, _ error) error {
			relPath, err := filepath.Rel(templatePath, p)
			if err != nil {
				return err
			}

			fileName := filepath.Base(relPath)
			if strings.HasPrefix(fileName, "_") {
				if fileName == "_keep" {
					return nil
				}
				relPath = path.Join(filepath.Dir(relPath), strings.Replace(fileName, "_", ".", 1))
			}

			filePath := filepath.Join(project.Path, relPath)
			info, err := os.Stat(filePath)
			if err != nil {
				return err
			}

			if info.IsDir() != d.IsDir() {
				return fmt.Errorf("mismatch in file types: %s | %s ", filePath, p)
			}
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}

		projectFS := os.DirFS(project.Path)
		err = fs.WalkDir(projectFS, ".", func(p string, d fs.DirEntry, _ error) error {
			virtualFileName := path.Base(p)
			if virtualFileName != "." && strings.HasPrefix(virtualFileName, ".") {
				virtualFileName = strings.Replace(virtualFileName, ".", "_", 1)
				p = path.Join(path.Dir(p), virtualFileName)
			}

			if virtualFileName == "_keep" {
				return nil
			}

			info, err := fs.Stat(templates, path.Join(templatePath, p))
			if err != nil {
				return err
			}

			filePath := filepath.Join(project.Path, p)
			if info.IsDir() != d.IsDir() {
				return fmt.Errorf("mismatch in file types: %s | %s ", filePath, p)
			}
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
	}
}


func TestRenderTemplate(t *testing.T) {
	tempDir := t.TempDir()

	sourcePath := filepath.Join("testdata", "sample-template-file.txt")
	targetPath := filepath.Join(tempDir, "sample-template-file.txt")

	project := Project{
		Name: "name",
		Path: tempDir,
		Description: "This is a sample description",
	}

	err := RenderTemplate(os.DirFS("."), filepath.ToSlash(sourcePath), targetPath, project)
	if err != nil {
		t.Fatal(err)
	}

	tmpl, err := template.ParseFiles(sourcePath)
	if err != nil {
		t.Fatal(err)
	}
	
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, project); err != nil {
		t.Fatal(err)
	}

	rendered := buf.String()

	fileContents, err := os.ReadFile(targetPath)
	if err != nil {
		t.Fatal(err)
	}

	if rendered != string(fileContents) {
		t.Fatal("rendered template file has not been rendered correctly")
	}
}
