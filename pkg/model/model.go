package model

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
	"github.com/sandrolain/gommander/pkg/fs"
	"github.com/sandrolain/gommander/pkg/rows"
)

var (
	footLSty = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColViolet))

	footVSty = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColOrange))

	dialogBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(ColViolet)).
			Padding(1, 2).
			Margin(1, 1).
			BorderTop(true).
			BorderLeft(true).
			BorderRight(true).
			BorderBottom(true)

	buttonStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColLightYellow)).
			Background(lipgloss.Color(ColDarkYellow)).
			Padding(0, 3).
			MarginLeft(1)

	activeButtonStyle = buttonStyle.
				Foreground(lipgloss.Color(ColLightYellow)).
				Background(lipgloss.Color(ColPink)).
				Padding(0, 3)

	tableNormalStyle = lipgloss.NewStyle().
				BorderForeground(lipgloss.Color(ColDarkGray)).Faint(true).
				Align(lipgloss.Right)

	tableActiveStyle = lipgloss.NewStyle().
				BorderForeground(lipgloss.Color(ColWhite)).
				Align(lipgloss.Right)
)

type model struct {
	leftPanelDir       string
	rightPanelDir      string
	leftTable          table.Model
	rightTable         table.Model
	active             string
	windowWidth        int
	windowHeight       int
	panelWidth         int
	selectedRows       map[string]bool // Track selected rows
	leftFilesInfo      rows.FilesInfo
	rightFilesInfo     rows.FilesInfo
	errorMessage       string
	confirmMessage     string
	confirmCallback    ConfirmCallback
	inputMessage       string
	inputCallback      InputCallback
	inputValue         string // Nome della nuova directory
	view               string
	key                string
	log                string
	updateLeftWatcher  UpdateWatcherFn
	updateRightWatcher UpdateWatcherFn
}

type UpdateWatcherFn func(string, func()) error

func InitialModel(ul UpdateWatcherFn, ur UpdateWatcherFn) model {
	currentDir, _ := os.Getwd()

	leftFilesInfo, leftTable := createTable(currentDir)
	rightFilesInfo, rightTable := createTable(currentDir)

	leftTable = leftTable.Focused(true)

	m := model{
		leftPanelDir:       currentDir,
		rightPanelDir:      currentDir,
		leftTable:          leftTable,
		rightTable:         rightTable,
		active:             "left",
		windowWidth:        0,
		windowHeight:       0,
		panelWidth:         0,
		selectedRows:       make(map[string]bool),
		leftFilesInfo:      leftFilesInfo,
		rightFilesInfo:     rightFilesInfo,
		updateLeftWatcher:  ul,
		updateRightWatcher: ur,
	}

	var err error

	err = m.updateLeftWatcher(currentDir, func() {
		m.refreshLeftTableRows()
	})
	if err != nil {
		m.log = fmt.Sprintf("Error creating left watcher: %s", err)
	}

	err = m.updateRightWatcher(currentDir, func() {
		m.refreshRightTableRows()
	})
	if err != nil {
		m.log = fmt.Sprintf("Error creating right watcher: %s", err)
	}

	return m
}

func createTable(dir string) (rows.FilesInfo, table.Model) {
	filesInfo, rows := rows.GetTableRows(dir)

	nameCol := table.NewFlexColumn("name", "Name", 10)
	nameCol = nameCol.WithStyle(nameCol.Style().Align(lipgloss.Left))

	km := table.DefaultKeyMap()
	km.RowSelectToggle.SetKeys(" ")

	t := table.New([]table.Column{
		nameCol,
		table.NewColumn("size", "Size", 8),
		table.NewColumn("mode", "Mode", 10).WithStyle(lipgloss.NewStyle().Align(lipgloss.Center)),
		table.NewColumn("modified", "Modified", 19).WithStyle(lipgloss.NewStyle().Align(lipgloss.Center)),
	}).WithRows(rows).
		BorderRounded().
		SelectableRows(true).
		WithRowStyleFunc(func(rsfi table.RowStyleFuncInput) lipgloss.Style {
			if rsfi.IsHighlighted {
				return lipgloss.NewStyle().Bold(true).Background(lipgloss.Color(ColViolet)).Foreground(lipgloss.Color(ColWhite))
			}

			if rsfi.Row.Data["name"] == ".." {
				return lipgloss.NewStyle().Foreground(lipgloss.Color(ColDarkYellow))
			}

			if rsfi.Row.Data["dir"] == true {
				return lipgloss.NewStyle().Foreground(lipgloss.Color(ColOrange))
			}

			return lipgloss.NewStyle().Foreground(lipgloss.Color(ColLightBlue))
		}).WithKeyMap(km).HeaderStyle(lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(ColPink)))

	return filesInfo, t
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.updateTablesWidth(msg)
	case tea.FocusMsg:
		m.refreshTablesRows(true, true)
	case tea.MouseMsg:
		change := 0
		if msg.Button == tea.MouseButtonWheelDown {
			change = -1
		} else if msg.Button == tea.MouseButtonWheelUp {
			change = 1
		}

		if change != 0 {
			if m.active == "left" {
				m.leftTable = m.leftTable.WithHighlightedRow(m.leftTable.GetHighlightedRowIndex() + change)
			} else {
				m.rightTable = m.rightTable.WithHighlightedRow(m.rightTable.GetHighlightedRowIndex() + change)
			}
		}
	case tea.KeyMsg:
		key := msg.String()

		m.key = key

		if m.inputMessage != "" {

			if key == "enter" {
				if m.inputCallback != nil {
					err := m.inputCallback(m.inputValue, &m)
					if err != nil {
						m.showError(err.Error())
					}
					m.refreshTablesRows(true, false)
				}
				m.inputCallback = nil
				m.inputMessage = ""
				m.inputValue = ""
				return m, nil
			}

			if key == "esc" {
				m.inputCallback = nil
				m.inputMessage = ""
				m.inputValue = ""
				return m, nil
			}

			if key == "backspace" {
				if len(m.inputValue) > 0 {
					m.inputValue = m.inputValue[:len(m.inputValue)-1]
				}
				return m, nil
			}

			m.inputValue += key
			return m, nil
		}

		if m.errorMessage != "" || m.confirmMessage != "" {
			switch key {
			case "enter":
				if m.errorMessage != "" {
					m.errorMessage = "" // Clear error message
					return m, nil
				}

				if m.confirmCallback != nil {
					err := m.confirmCallback(&m)
					if err != nil {
						m.showError(err.Error())
					}
					m.confirmCallback = nil
					m.confirmMessage = ""

					return m, nil
				}

			case "esc":

				if m.errorMessage != "" {
					m.errorMessage = "" // Clear error message
					return m, nil
				}

				if m.confirmCallback != nil {
					m.confirmCallback = nil
					m.confirmMessage = ""
					return m, nil
				}
			}

			return m, nil
		}

		switch key {
		case KeyQuit:
			return m, tea.Quit
		case KeySwitch:
			if m.active == "left" {
				m.active = "right"
				m.leftTable = m.leftTable.Focused(false)
				m.rightTable = m.rightTable.Focused(true)
			} else {
				m.active = "left"
				m.leftTable = m.leftTable.Focused(true)
				m.rightTable = m.rightTable.Focused(false)
			}
		case "space":
			var currentTable *table.Model
			if m.active == "left" {
				currentTable = &m.leftTable
			} else {
				currentTable = &m.rightTable
			}

			selectedRows := currentTable.SelectedRows()
			for _, selectedRow := range selectedRows {
				id := selectedRow.Data["path"].(string)
				if m.selectedRows[id] {
					delete(m.selectedRows, id) // Deselect row
				} else {
					m.selectedRows[id] = true // Select row
				}
			}

		case KeyEnter:

			newPath, err := m.getHighlightedRowPath(false)
			if err != nil {
				m.showError("Error getting path")
				return m, nil
			}

			err = m.enterFile(newPath)
			if err != nil {
				m.showError(err.Error())
			}

		case KeyExit, KeyBack:

			var currentPath string
			if m.active == "left" {
				currentPath = m.leftPanelDir
			} else {
				currentPath = m.rightPanelDir
			}

			newPath := filepath.Join(currentPath, "..")

			err := m.enterFile(newPath)
			if err != nil {
				m.showError(err.Error())
			}

		case KeyCopy:

			paths, err := m.getCurrentRowsPaths()
			if err != nil {
				m.showError("Error getting paths to copy")
				return m, nil
			}

			m.confirmDialog(fmt.Sprintf("Are you sure you want to copy\n%s?", strings.Join(paths, "\n")), func(m *model) error {
				err := m.copyFiles(false)
				if err != nil {
					return fmt.Errorf("Error copying file: %v", err)
				}
				m.refreshTablesRows(false, true)
				return nil
			})

		case KeyCopyO:

			paths, err := m.getCurrentRowsPaths()
			if err != nil {
				m.showError("Error getting paths to copy")
				return m, nil
			}

			m.confirmDialog(fmt.Sprintf("Are you sure you want to copy with overwrite\n%s?", strings.Join(paths, "\n")), func(m *model) error {
				err := m.copyFiles(true)
				if err != nil {
					return fmt.Errorf("Error copying file: %v", err)
				}
				m.refreshTablesRows(false, true)
				return nil
			})

		case KeyMove:

			err := m.moveFiles(false)
			if err != nil {
				m.showError(err.Error())
			}

			m.refreshTablesRows(true, true)

		case "ctrl+d":

			paths, err := m.getCurrentRowsPaths()
			if err != nil {
				m.showError("Error getting path")
				return m, nil
			}

			m.confirmMessage = fmt.Sprintf("Are you sure you want to delete\n%s?", strings.Join(paths, "\n"))
			m.confirmCallback = func(m *model) error {
				err := fs.DeleteFiles(paths)
				if err != nil {
					return fmt.Errorf("Error deleting file: %v", err)
				}
				m.refreshTablesRows(true, false)
				return nil
			}

		case KeyTrash:

			paths, err := m.getCurrentRowsPaths()
			if err != nil {
				m.showError(fmt.Sprintf("Error getting paths: %v", err))
				return m, nil
			}

			m.confirmDialog(fmt.Sprintf("Are you sure you want move to trash\n%s?", strings.Join(paths, "\n")), func(m *model) error {
				err := fs.TrashFiles(paths)
				if err != nil {
					return fmt.Errorf("Error moving to trash file: %v", err)
				}
				m.refreshTablesRows(true, false)
				return nil
			})

		case KeyMkdir:

			m.inputDialog("Enter the name of the new directory:", func(value string, m *model) error {
				if value == "" {
					return fmt.Errorf("Directory name cannot be empty")
				}
				err := m.createDirectory(value)
				if err != nil {
					return fmt.Errorf("Error creating directory: %v", err)
				}
				m.refreshTablesRows(true, false)
				return nil
			})

		case KeyMkfile:

			m.inputDialog("Enter the name of the new file:", func(value string, m *model) error {
				if value == "" {
					return fmt.Errorf("File name cannot be empty")
				}
				err := m.createFile(value)
				if err != nil {
					return fmt.Errorf("Error creating file: %v", err)
				}
				m.refreshTablesRows(true, false)
				return nil
			})

		case KeyVscode:

			err := m.openVsCode()
			if err != nil {
				m.showError(err.Error())
			}

		}

		if m.active == "left" {
			var cmd tea.Cmd
			m.leftTable, cmd = m.leftTable.Update(msg)
			return m, cmd
		}

		var cmd tea.Cmd
		m.rightTable, cmd = m.rightTable.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m *model) getTable() *table.Model {
	if m.active == "left" {
		return &m.leftTable
	}
	return &m.rightTable
}

func (m *model) getHighlightedRowPath(noSup bool) (string, error) {
	currentTable := m.getTable()
	return rows.GetRowPath(currentTable.HighlightedRow(), noSup)
}

func (m *model) getSelectedRowsPaths() ([]string, error) {
	currentTable := m.getTable()
	return rows.GetRowsPaths(currentTable.SelectedRows())
}

func (m *model) getCurrentRowsPaths() ([]string, error) {
	paths, err := m.getSelectedRowsPaths()
	if err != nil {
		return nil, err
	}
	if len(paths) > 0 {
		return paths, nil
	}
	path, err := m.getHighlightedRowPath(true)
	if err != nil {
		return nil, err
	}

	if path != "" && !strings.HasSuffix(path, string(filepath.Separator)+"..") {
		return []string{path}, nil
	}

	return nil, fmt.Errorf("no path selected")
}

func (m *model) getDestinationDirPath() (string, error) {
	var destPath string
	if m.active == "left" {
		destPath = m.rightPanelDir
	} else {
		destPath = m.leftPanelDir
	}
	return destPath, nil
}

func (m *model) enterFile(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	if info.IsDir() {
		filesInfo, newRows := rows.GetTableRows(path)
		// Enter directory
		if m.active == "left" {
			err := m.updateLeftWatcher(path, func() {
				m.refreshLeftTableRows()
			})
			if err != nil {
				return fmt.Errorf("error creating watcher: %v", err)
			}
			m.leftPanelDir = path
			m.leftTable = m.leftTable.WithRows(newRows).WithHighlightedRow(0).WithAllRowsDeselected()
			m.leftFilesInfo = filesInfo
		} else {
			err := m.updateRightWatcher(path, func() {
				m.refreshRightTableRows()
			})
			if err != nil {
				return fmt.Errorf("error creating watcher: %v", err)
			}
			m.rightPanelDir = path
			m.rightTable = m.rightTable.WithRows(newRows).WithHighlightedRow(0).WithAllRowsDeselected()
			m.rightFilesInfo = filesInfo
		}

		return nil
	}

	// If it's a file, open it with the default application
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", path)
	case "darwin":
		cmd = exec.Command("open", path)
	default: // Assume Linux or other Unix-like OS
		cmd = exec.Command("xdg-open", path)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}

	return nil
}

func (m *model) openVsCode() error {
	path, err := m.getHighlightedRowPath(false)
	if err != nil {
		return fmt.Errorf("error getting highlighted row path: %v", err)
	}
	cmd := exec.Command("code", path)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("error opening VS Code: %v", err)
	}
	return nil
}

func (m *model) refreshTablesRows(source bool, dest bool) {
	same := m.leftPanelDir == m.rightPanelDir
	left := (m.active == "left" && source) || (m.active == "right" && dest)
	right := (m.active == "right" && source) || (m.active == "left" && dest)

	if left || (right && same) {
		m.refreshLeftTableRows()
	}
	if right || (left && same) {
		m.refreshRightTableRows()
	}
}

func (m *model) refreshLeftTableRows() {
	path := m.leftPanelDir
	filesInfo, newRows := rows.GetTableRows(path)
	m.leftTable = m.leftTable.WithRows(newRows).WithHighlightedRow(0)
	m.leftFilesInfo = filesInfo
}

func (m *model) refreshRightTableRows() {
	path := m.rightPanelDir
	filesInfo, newRows := rows.GetTableRows(path)
	m.rightTable = m.rightTable.WithRows(newRows).WithHighlightedRow(0)
	m.rightFilesInfo = filesInfo
}

func (m *model) showError(msg string) {
	m.errorMessage = msg
}

func (m *model) updateTablesWidth(msg tea.WindowSizeMsg) {
	m.windowWidth = msg.Width
	m.windowHeight = msg.Height
	m.panelWidth = m.windowWidth / 2
	m.leftTable = m.leftTable.WithTargetWidth(m.panelWidth).WithMinimumHeight(m.windowHeight).WithPageSize(m.windowHeight - extraRows)
	m.rightTable = m.rightTable.WithTargetWidth(m.panelWidth).WithMinimumHeight(m.windowHeight).WithPageSize(m.windowHeight - extraRows)
}

func fL(faint bool, s string) string {
	return footLSty.Faint(faint).Render(s)
}

func fV(faint bool, s string) string {
	return footVSty.Faint(faint).Render(s)
}

func (m model) View() string {

	leftFaint := false
	rightFaint := false
	if m.active == "left" {
		rightFaint = true
	} else {
		leftFaint = true
	}

	leftFooter := lipgloss.JoinVertical(
		lipgloss.Left,
		fL(leftFaint, "Path: ")+fV(leftFaint, m.leftPanelDir),
		fL(leftFaint, "Total: ")+
			fV(leftFaint, fmt.Sprintf("%d", m.leftFilesInfo.Total))+fL(leftFaint, " | Dirs: ")+
			fV(leftFaint, fmt.Sprintf("%d", m.leftFilesInfo.Dirs))+fL(leftFaint, " | Files: ")+
			fV(leftFaint, fmt.Sprintf("%d", m.leftFilesInfo.Files))+fL(leftFaint, " - ")+fL(leftFaint, fmt.Sprintf("%d/%d", m.leftTable.CurrentPage(), m.leftTable.MaxPages())),
	)

	rightFooter := lipgloss.JoinVertical(
		lipgloss.Left,
		fL(rightFaint, "Path: ")+fV(rightFaint, m.rightPanelDir),
		fL(rightFaint, "Total: ")+
			fV(rightFaint, fmt.Sprintf("%d", m.rightFilesInfo.Total))+fL(rightFaint, " | Dirs: ")+
			fV(rightFaint, fmt.Sprintf("%d", m.rightFilesInfo.Dirs))+fL(rightFaint, " | Files: ")+
			fV(rightFaint, fmt.Sprintf("%d", m.rightFilesInfo.Files))+fL(rightFaint, " - ")+fL(rightFaint, fmt.Sprintf("%d/%d", m.rightTable.CurrentPage(), m.rightTable.MaxPages())),
	)

	leftTable := m.leftTable.WithStaticFooter(fL(leftFaint, m.log+" - "+m.key+" - ") + leftFooter)
	rightTable := m.rightTable.WithStaticFooter(rightFooter)

	if m.active == "left" {
		leftTable = leftTable.WithBaseStyle(tableActiveStyle)
		rightTable = rightTable.WithBaseStyle(tableNormalStyle)
	} else {
		leftTable = leftTable.WithBaseStyle(tableNormalStyle)
		rightTable = rightTable.WithBaseStyle(tableActiveStyle)
	}

	leftContent := leftTable.View()
	rightContent := rightTable.View()

	m.view = lipgloss.JoinHorizontal(lipgloss.Top, leftContent, rightContent)

	if m.errorMessage != "" {
		return m.renderAlertDialog(m.errorMessage)
	}

	if m.confirmMessage != "" {
		return m.renderConfirmDialog(m.confirmMessage)
	}

	return m.view
}
