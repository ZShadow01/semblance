package filesystem

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func CreateDirectory(path string) error {
	err := os.Mkdir(path, os.ModePerm)
	if err == nil {
		return nil
	}
	
	var pathErr *os.PathError
	if errors.As(err, &pathErr) {
		switch {
		case errors.Is(pathErr.Err, os.ErrExist):
			return fmt.Errorf("directory '%s' already exists: %w", path, err)
		case errors.Is(pathErr.Err, os.ErrNotExist):
			return fmt.Errorf("could not find the path '%s': %w", filepath.Dir(path), err)
		case errors.Is(pathErr.Err, os.ErrPermission):
			return fmt.Errorf("permission denied creating directory %q: %w", path, err)
		default:
			return fmt.Errorf("failed to create directory '%s' (op: %s): %w", path, pathErr.Op, err)
		}
	}

	return fmt.Errorf("unexpected error creating directory %q: %w", path, err)
}


func CreateFile(path string) (*os.File, error) {
	file, err := os.Create(path)
	if err == nil {
		return file, nil
	}

	var pathErr *os.PathError
	if errors.As(err, &pathErr) {
		switch {
		case errors.Is(pathErr.Err, os.ErrNotExist):
			return nil, fmt.Errorf("could not find the path '%s': %w", filepath.Dir(path), err)
		case errors.Is(pathErr.Err, os.ErrPermission):
			return nil, fmt.Errorf("permission denied creating file %q: %w", path, err)
		default:
			return nil, fmt.Errorf("failed to create file %q (op: %s): %w", path, pathErr.Op, err)
		}
	}
	return nil, fmt.Errorf("unexpected error creating file %q: %w", path, err)
}
