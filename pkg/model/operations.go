package model

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sandrolain/gommander/pkg/fs"
)

func (m *model) createDirectory(name string) error {
	var currentPath string
	if m.active == "left" {
		currentPath = m.leftPanelDir
	} else {
		currentPath = m.rightPanelDir
	}
	newDirPath := filepath.Join(currentPath, name)
	err := os.Mkdir(newDirPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("Error creating directory: %v", err)
	}
	if m.active == "left" {
		m.leftPanelDir = newDirPath
	} else {
		m.rightPanelDir = newDirPath
	}
	return nil
}

func (m *model) createFile(name string) error {
	var currentPath string
	if m.active == "left" {
		currentPath = m.leftPanelDir
	} else {
		currentPath = m.rightPanelDir
	}
	newFilePath := filepath.Join(currentPath, name)
	err := os.WriteFile(newFilePath, []byte{}, os.ModePerm)
	if err != nil {
		return fmt.Errorf("Error creating file: %v", err)
	}
	return nil
}

func (m *model) copyFiles(overwrite bool) error {
	destPath, err := m.getDestinationDirPath()
	if err != nil {
		return err
	}

	paths, err := m.getCurrentRowsPaths()
	if err != nil {
		return err
	}

	for _, filePath := range paths {
		err := fs.CopyFile(filePath, destPath, overwrite)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *model) moveFiles(overwrite bool) error {
	destPath, err := m.getDestinationDirPath()
	if err != nil {
		return err
	}

	paths, err := m.getCurrentRowsPaths()
	if err != nil {
		return err
	}

	for _, filePath := range paths {
		err := fs.MoveFile(filePath, destPath, overwrite)
		if err != nil {
			return err
		}
	}

	return nil
}
