package model

const (
	extraRows         = 7
	headFootExtraRows = 2 // Rows for header and footer bars
)

const (
	ColWhite       = "#FFFFFF"
	ColDarkGray    = "#333333"
	ColViolet      = "#874BFD"
	ColLightBlue   = "#87afff"
	ColPink        = "#F25D94"
	ColLightYellow = "#FFF7DB"
	ColDarkYellow  = "#888B7E"
	ColYellow      = "#ffd703"
	ColOrange      = "#ffaf00"
)

const (
	KeyQuit   = "q"
	KeyEnter  = "enter"
	KeyCancel = "esc"
	KeyBack   = "backspace"
	KeySwitch = "tab"
	KeyHelp   = "ctrl+h"
	KeyCopy   = "ctrl+c"
	KeyCopyO  = "ctrl+r"
	KeyMove   = "ctrl+x"
	KeyMoveO  = "ctrl+t"
	KeyDelete = "ctrl+d"
	KeyTrash  = "delete"
	KeyMkdir  = "ctrl+f"
	KeyMkfile = "ctrl+n"
	KeyVscode = "ctrl+k"
	KeySelect = "space"
)

var helpArray = [][2]string{
	{KeyHelp, "Show help"},
	{KeyQuit, "Quit program"},
	{KeyEnter, "Enter directory / Open file / Confirm"},
	{KeyCancel, "Cancel"},
	{KeyBack, "Upper directory"},
	{KeySwitch, "Switch panel"},
	{KeySelect, "Select file"},
	{KeyCopy, "Copy files"},
	{KeyCopyO, "Copy files with overwrite"},
	{KeyMove, "Move files"},
	{KeyMoveO, "Move files with overwrite"},
	{KeyDelete, "Delete files"},
	{KeyTrash, "Move files to trash"},
	{KeyMkdir, "Create new directory"},
	{KeyMkfile, "Create new file"},
	{KeyVscode, "Open in VSCode"},
}
