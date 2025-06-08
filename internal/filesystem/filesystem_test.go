package filesystem

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestCreateDirectory(t *testing.T) {
	tempDir := t.TempDir()
	testDir := "test-project-dir"
	
	outputPath := filepath.Join(tempDir, testDir)

	err := CreateDirectory(outputPath)

	if err != nil {
		t.Errorf(`CreateDirectory("%s") = %v`, outputPath, err)
	}
}


func TestCreateDirectory_Exists(t *testing.T) {
	tempDir := t.TempDir()
	testDir := "test-project-dir"
	
	outputPath := filepath.Join(tempDir, testDir)

	_ = CreateDirectory(outputPath)
	err := CreateDirectory(outputPath)
	if err == nil || !errors.Is(err, os.ErrExist) {
		t.Errorf(`CreateDirectory("%s") = %v`, outputPath, err)
	}
}


func TestCreateDirectory_NotExists(t *testing.T) {
	tempDir := t.TempDir()
	testDir := "test-project-dir"
	
	outputPath := filepath.Join(tempDir, "temp", testDir)

	err := CreateDirectory(outputPath)
	if err == nil || !errors.Is(err, os.ErrNotExist) {
		t.Errorf(`CreateDirectory("%s") = %v`, outputPath, err)
	}
}


func TestCreateFile(t *testing.T) {
	tempDir := t.TempDir()
	testFile := "test.txt"

	outputPath := filepath.Join(tempDir, testFile)

	file, err := CreateFile(outputPath)
	if err != nil {
		t.Errorf(`CreateFile("%s") = %v`, outputPath, err)
	}
	defer file.Close()
}


func TestCreateFile_NotExist(t *testing.T) {
	tempDir := t.TempDir()
	testFile := "test.txt"

	outputPath := filepath.Join(tempDir, "temp", testFile)

	file, err := CreateFile(outputPath)
	if err == nil || !errors.Is(err, os.ErrNotExist) {
		t.Errorf(`CreateFile("%s") = %v`, outputPath, err)
	}
	defer file.Close()
}
