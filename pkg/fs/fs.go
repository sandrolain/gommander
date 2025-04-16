package fs

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	gotrash "github.com/laurent22/go-trash"
)

func CopyFile(srcFilePath, dstDirPath string, permitOverwrite bool) error {
	fileName := filepath.Base(srcFilePath)

	srcFile, err := os.Open(srcFilePath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	srcInfo, err := os.Stat(srcFilePath)
	if err != nil {
		return fmt.Errorf("error stating file %v: %v", srcFilePath, err)
	}

	dstFilePath := filepath.Join(dstDirPath, fileName)

	if !permitOverwrite {
		if _, err := os.Stat(dstFilePath); err == nil {
			return fmt.Errorf("file %v already exists", dstFilePath)
		}
	}

	dstFile, err := os.Create(dstFilePath)
	if err != nil {
		return fmt.Errorf("error creating file %v: %v", dstFilePath, err)
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return fmt.Errorf("error copying file %v to %v: %v", srcFilePath, dstFilePath, err)
	}

	err = os.Chmod(dstFilePath, srcInfo.Mode())
	if err != nil {
		return fmt.Errorf("error setting permissions for file %v: %v", dstFilePath, err)
	}

	return nil
}

func DeleteFile(filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		return err
	}
	return nil
}

func DeleteFiles(filePaths []string) error {
	for _, filePath := range filePaths {
		err := DeleteFile(filePath)
		if err != nil {
			return err
		}
	}
	return nil
}

func TrashFile(filePath string) error {
	_, err := gotrash.MoveToTrash(filePath)
	if err != nil {
		return err
	}
	return nil
}

func TrashFiles(filePaths []string) error {
	for _, filePath := range filePaths {
		err := TrashFile(filePath)
		if err != nil {
			return err
		}
	}
	return nil
}

func MoveFile(srcFilePath, dstDirPath string, permitOverwrite bool) error {
	fileName := filepath.Base(srcFilePath)

	srcFile, err := os.Open(srcFilePath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	srcInfo, err := os.Stat(srcFilePath)
	if err != nil {
		return err
	}

	dstFilePath := filepath.Join(dstDirPath, fileName)

	if !permitOverwrite {
		if _, err := os.Stat(dstFilePath); err == nil {
			return os.ErrExist
		}
	}

	err = os.Rename(srcFilePath, dstFilePath)
	if err != nil {
		return err
	}

	err = os.Chmod(dstFilePath, srcInfo.Mode())
	if err != nil {
		return err
	}

	return nil
}
