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
	cancelStyle := activeButtonStyle
	confirmStyle := buttonStyle

	if m.confirmBtn == 1 {
		cancelStyle = buttonStyle
		confirmStyle = activeButtonStyle
	}

	cancelButton := cancelStyle.Render("Cancel")
	confirmButton := confirmStyle.Render("Confirm")

	width := min(lipgloss.Width(text), m.windowWidth)

	question := lipgloss.NewStyle().Width(width).Align(lipgloss.Center).MarginBottom(1).Render(text)
	if m.inputValue != "" {
		question += "\n" + lipgloss.NewStyle().Align(lipgloss.Center).Render(m.inputValue) + "\n"
	}
	buttons := lipgloss.JoinHorizontal(lipgloss.Top, cancelButton, confirmButton)
	ui := lipgloss.JoinVertical(lipgloss.Center, question, buttons)

	modal := dialogBoxStyle.Render(ui)

	return m.renderOverlayViews(modal)
}

func (m *model) renderAlertDialog(text string) string {
	okButton := activeButtonStyle.Render("Ok")

	width := min(lipgloss.Width(text), m.windowWidth-20)

	question := lipgloss.NewStyle().Width(width).Align(lipgloss.Center).MarginBottom(1).Render(text)
	buttons := lipgloss.JoinHorizontal(lipgloss.Top, okButton)
	ui := lipgloss.JoinVertical(lipgloss.Center, question, buttons)

	modal := dialogBoxStyle.Render(ui)

	return m.renderOverlayViews(modal)
}

func (m *model) renderHelpDialog() string {
	okButton := activeButtonStyle.Render("Ok")

	text := ""

	maxLength := 0
	for _, item := range helpArray {
		if len(item[0]) > maxLength {
			maxLength = len(item[0])
		}
	}

	for _, item := range helpArray {
		paddedKey := lipgloss.NewStyle().Bold(true).Render(item[0]) + strings.Repeat(" ", maxLength-len(item[0]))
		text += paddedKey + " : " + item[1] + "\n"
	}

	width := min(lipgloss.Width(text), m.windowWidth-20)

	question := lipgloss.NewStyle().Width(width).Align(lipgloss.Left).MarginBottom(1).Render(text)
	buttons := lipgloss.JoinHorizontal(lipgloss.Top, okButton)
	ui := lipgloss.JoinVertical(lipgloss.Center, question, buttons)

	modal := dialogBoxStyle.Render(ui)

	return m.renderOverlayViews(modal)
}
