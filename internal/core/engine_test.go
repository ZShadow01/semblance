package core

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"testing"
	"text/template"
)

type TestTemplate struct {
	Name 	string 	`json:"name"`
	Exists 	bool 	`json:"exists"`
}

//go:embed testdata/sample-templates.json
var SampleTemplatesData []byte

func TestGenerateByDefaultTemplate(t *testing.T) {
	var sampleTemplates []TestTemplate
	err := json.Unmarshal(SampleTemplatesData, &sampleTemplates)
	if err != nil {
		t.Fatal(err)
	}
	
	for _, sample := range sampleTemplates {
		fmt.Println("New template")
		tempDir := t.TempDir()

		project := Project{
			Name: "test",
			Path: filepath.Join(tempDir, "test"),
			Description: "Test",
		}

		err := GenerateByDefaultTemplate(project, sample.Name)
		_, err2 := isValidDefaultTemplate(sample.Name)

		if err == err2 && err != nil {
			continue
		}

		if err2 != nil && fmt.Sprintf("%T", err) != fmt.Sprintf("%T", err2) {
			t.Fatalf("GenerateByDefaultTemplate(%v, %s) validated an invalid template", project, sample.Name)
		}

		if err2 == nil && err != nil {
			t.Fatalf("GenerateByDefaultTemplate(%v, %s) = %v", project, sample.Name, err)
		}

		// Get the amount of files in the template
		nTemplateFiles := 0
		fs.WalkDir(templates, path.Join("templates", sample.Name), func(p string, _ fs.DirEntry, _ error) error {
			if path.Base(p) == "_keep" {
				return nil
			}

			nTemplateFiles++
			return nil
		})

		nGeneratedFiles := 0
		fs.WalkDir(os.DirFS(project.Path), ".", func(_ string, _ fs.DirEntry, _ error) error {
			nGeneratedFiles++
			return nil
		})

		if nTemplateFiles != nGeneratedFiles {
			t.Fatalf("GenerateByDefaultTemplate(%v, %s) does not generate the same amount of files as the template", project, sample.Name)
		}
	}
}


func TestIsValidDefaultTemplate(t *testing.T) {
	var sampleTemplates []TestTemplate
	err := json.Unmarshal(SampleTemplatesData, &sampleTemplates)
	if err != nil {
		t.Fatal (err)
	}
	
	for _, sample := range sampleTemplates {
		_, err := isValidDefaultTemplate(sample.Name)

		// Check whether isValidDefaultTemplate correctly finds the templates
		if (err != nil) && sample.Exists {
			t.Fatalf("valid template flagged as invalid template '%s'\n", sample.Name)
		} else if (err == nil) && !sample.Exists {
			t.Fatalf("invalid template flagged as valid template '%s'\n", sample.Name)
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

	// Render the source file
	tmpl, err := template.ParseFiles(sourcePath)
	if err != nil {
		t.Fatal(err)
	}
	
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, project); err != nil {
		t.Fatal(err)
	}

	rendered := buf.String()

	// Read the rendered target file
	fileContents, err := os.ReadFile(targetPath)
	if err != nil {
		t.Fatal(err)
	}

	// Compare the source file with the target file with the templates rendered
	if rendered != string(fileContents) {
		t.Fatal("rendered template file has not been rendered correctly")
	}
}
