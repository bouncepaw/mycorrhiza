package bxjson

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type JSONFile[T any] struct {
	path string
}

func New[T any](path string) *JSONFile[T] {
	return &JSONFile[T]{
		path: path,
	}
}

func (f *JSONFile[T]) Read() (T, error) {
	var value T
	data, err := os.ReadFile(f.path)
	if errors.Is(err, os.ErrNotExist) {
		return value, nil
	}
	if err != nil {
		return value, fmt.Errorf("failed to read %s: %w", f.path, err)
	}

	if err = json.Unmarshal(data, &value); err != nil {
		var zero T
		return zero, fmt.Errorf("failed to decode %s: %w", f.path, err)
	}

	return value, nil
}

func (f *JSONFile[T]) Write(value T) error {
	data, err := json.MarshalIndent(value, "", "\t")
	if err != nil {
		return fmt.Errorf("failed to encode %s: %w", f.path, err)
	}

	mode := os.FileMode(0600)
	if info, err := os.Stat(f.path); err == nil {
		mode = info.Mode().Perm()
	} else if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("failed to stat %s: %w", f.path, err)
	}

	temp, err := os.CreateTemp(filepath.Dir(f.path), "."+filepath.Base(f.path)+"-*")
	if err != nil {
		return fmt.Errorf("failed to create temporary file for %s: %w", f.path, err)
	}
	defer os.Remove(temp.Name())
	defer temp.Close()

	if _, err = temp.Write(data); err != nil {
		return fmt.Errorf("failed to write %s: %w", f.path, err)
	}

	if err = temp.Chmod(mode); err != nil {
		return fmt.Errorf("failed to chmod %s: %w", f.path, err)
	}

	if err = temp.Sync(); err != nil {
		return fmt.Errorf("failed to sync %s: %w", f.path, err)
	}

	if err = temp.Close(); err != nil {
		return fmt.Errorf("failed to close %s: %w", f.path, err)
	}

	if err = os.Rename(temp.Name(), f.path); err != nil {
		return fmt.Errorf("failed to replace %s: %w", f.path, err)
	}

	return nil
}
