package model

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

type ConfirmCallback func(*model) error
type InputCallback func(string, *model) error

func (m *model) confirmDialog(text string, cb ConfirmCallback) {
	m.confirmMessage = text
	m.confirmCallback = cb
}

func (m *model) inputDialog(text string, callback InputCallback) {
	m.inputValue = ""
	m.inputMessage = text
	m.inputCallback = callback
}

func (m *model) renderOverlayViews(modal string) string {
	modalWidth := lipgloss.Width(modal)
	modalHeight := lipgloss.Height(modal)
	left := (m.windowWidth - modalWidth) / 2
	top := (m.windowHeight - modalHeight) / 2

	mainViewLines := strings.Split(m.view, "\n")
	modalLines := strings.Split(modal, "\n")

	for i, overlayLine := range modalLines {
		bgLine := mainViewLines[i+top] // TODO: index handling
		if len(bgLine) < left {
			bgLine += strings.Repeat(" ", left-len(bgLine)) // add padding
		}

		bgLeft := ansi.Truncate(bgLine, left, "")
		bgRight := truncateLeft(bgLine, left+ansi.StringWidth(overlayLine))

		mainViewLines[i+top] = bgLeft + overlayLine + bgRight
	}

	return strings.Join(mainViewLines, "\n")
}

func (m *model) renderConfirmDialog(text string) string {
	okButton := activeButtonStyle.Render("Ok [Enter]")
	cancelButton := buttonStyle.Render("Cancel [Esc]")

	width := min(lipgloss.Width(text), m.windowWidth)

	question := lipgloss.NewStyle().Width(width).Align(lipgloss.Center).MarginBottom(1).Render(text)
	if m.inputValue != "" {
		question += "\n" + lipgloss.NewStyle().Align(lipgloss.Center).Render(m.inputValue) + "\n"
	}
	buttons := lipgloss.JoinHorizontal(lipgloss.Top, okButton, cancelButton)
	ui := lipgloss.JoinVertical(lipgloss.Center, question, buttons)

	modal := dialogBoxStyle.Render(ui)

	return m.renderOverlayViews(modal)
}

func (m *model) renderAlertDialog(text string) string {
	okButton := activeButtonStyle.Render("Ok [Enter / Esc]")

	width := min(lipgloss.Width(text), m.windowWidth-20)

	question := lipgloss.NewStyle().Width(width).Align(lipgloss.Center).MarginBottom(1).Render(text)
	buttons := lipgloss.JoinHorizontal(lipgloss.Top, okButton)
	ui := lipgloss.JoinVertical(lipgloss.Center, question, buttons)

	modal := dialogBoxStyle.Render(ui)

	return m.renderOverlayViews(modal)
}
