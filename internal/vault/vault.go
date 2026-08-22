package vault

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

var ErrAlredyExists = errors.New("note alredy exits")

type Note struct {
	Name    string
	ModTime time.Time
}

func DefaultDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("Get home directory: %w", err)
	}

	return filepath.Join(homeDir, ".totion"), nil
}

func EnsureDir(dir string) error {
	if err := os.MkdirAll(dir, 0750); err != nil {
		return fmt.Errorf("Creat vault dir: %w", err)
	}

	return nil
}

func List(dir string) ([]Note, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("Read the directory: %w", err)
	}

	notes := make([]Note, 0, len(entries))

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		notes = append(notes, Note{
			Name:    entry.Name(),
			ModTime: info.ModTime(),
		})
	}

	return notes, nil
}

func Creat(dir, name string) (*os.File, error) {
	path := filepath.Join(dir, name+".md")

	if _, err := os.Stat(path); err == nil {
		return nil, ErrAlredyExists
	}

	f, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("Creating file: %w", err)
	}

	return f, nil
}

func Open(dir, name string) (*os.File, []byte, error) {
	path := filepath.Join(dir, name)

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, fmt.Errorf("read file %w", err)
	}

	f, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		return nil, nil, fmt.Errorf("Open file: %w", err)
	}

	return f, content, nil
}

func Write(f *os.File, content string) error {
	if err := f.Truncate(0); err != nil {
		return fmt.Errorf("truncated Notes: %w", err)
	}

	if _, err := f.Seek(0, 0); err != nil {
		return fmt.Errorf("Seek note: %w", err)
	}

	if _, err := f.WriteString(content); err != nil {
		return fmt.Errorf("Write note: %w", err)
	}

	return nil
}

func Save(f *os.File, content string) error {
	if err := Write(f, content); err != nil {
		return err
	}

	if err := f.Close(); err != nil {
		return fmt.Errorf("close note: %w", err)
	}

	return nil
}

func Delet(dir, name string) error {
	path := filepath.Join(dir, name)

	if err := os.Remove(path); err != nil {
		return fmt.Errorf("delete note: %w", err)
	}

	return nil
}

func Rename(dir, oldname, newname string) error {
	oldPath := filepath.Join(dir, oldname)
	newPath := filepath.Join(dir, newname)

	if _, err := os.Stat(newPath); err == nil {
		return ErrAlredyExists
	}

	if err := os.Rename(oldPath, newPath); err != nil {
		return fmt.Errorf("rename note : %w", err)
	}

	return nil
}
