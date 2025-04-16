package rows

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	humanize "github.com/dustin/go-humanize"
	"github.com/evertras/bubble-table/table"
)

func getTableRowForPath(dirPath string, name string) (table.Row, error) {
	path := filepath.Join(dirPath, name)
	path, err := filepath.Abs(path)
	if err != nil {
		return table.Row{}, err
	}

	info, err := os.Stat(path)
	if err != nil {
		return table.Row{}, err
	}

	isDir := info.IsDir()
	permissions := info.Mode().String()
	usize := uint64(info.Size())

	formattedSize := ""
	if !isDir {
		formattedSize = humanize.Bytes(usize)
	}

	row := table.NewRow(table.RowData(map[string]interface{}{
		"path":        path,
		"dir":         isDir,
		"name":        name,
		"size":        formattedSize,
		"usize":       usize,
		"permissions": permissions,
	}))

	return row, nil
}

type FilesInfo struct {
	Total int
	Dirs  int
	Files int
}

func GetTableRows(dir string) (FilesInfo, []table.Row) {
	files, _ := os.ReadDir(dir)
	dirs := []table.Row{}
	regularFiles := []table.Row{}

	row, err := getTableRowForPath(dir, "..")
	if err == nil {
		dirs = append(dirs, row)
	}

	for _, file := range files {
		if file.Name() == "." || file.Name() == ".." {
			continue
		}
		row, err := getTableRowForPath(dir, file.Name())
		if err != nil {
			continue
		}

		if file.IsDir() {
			dirs = append(dirs, row)
		} else {
			regularFiles = append(regularFiles, row)
		}
	}

	// Sort directories and files alphabetically
	sort.Slice(dirs[1:], func(i, j int) bool {
		return dirs[i+1].Data["name"].(string) < dirs[j+1].Data["name"].(string)
	})
	sort.Slice(regularFiles, func(i, j int) bool {
		return regularFiles[i].Data["name"].(string) < regularFiles[j].Data["name"].(string)
	})

	totalDirs := len(dirs) - 1 // Exclude ".."
	totalFiles := len(regularFiles)

	info := FilesInfo{
		Total: totalDirs + totalFiles,
		Dirs:  totalDirs,
		Files: totalFiles,
	}

	// Combine directories and files
	return info, append(dirs, regularFiles...)
}

func CountTableRows(rows []table.Row) (int, int, int) {
	total := len(rows) - 1 // Exclude ".."
	dirs := 0
	for _, row := range rows[1:] {
		if row.Data["dir"] == true {
			dirs++
		}
	}
	files := total - dirs
	return total, dirs, files
}

func GetRowPath(row table.Row, noSup bool) (string, error) {
	if noSup {
		n, ok := row.Data["name"]
		if !ok {
			return "", fmt.Errorf("name not found")
		}
		name, ok := n.(string)
		if !ok {
			return "", fmt.Errorf("name is not a string")
		}

		if name == ".." {
			return "", nil
		}
	}

	p, ok := row.Data["path"]
	if !ok {
		return "", fmt.Errorf("path not found")
	}

	path, ok := p.(string)
	if !ok {
		return "", fmt.Errorf("path is not a string")
	}

	return path, nil
}

func GetRowsPaths(rows []table.Row) ([]string, error) {
	paths := make([]string, len(rows))
	for i, row := range rows {
		path, err := GetRowPath(row, true)
		if err != nil {
			return nil, err
		}
		paths[i] = path
	}
	return paths, nil
}
