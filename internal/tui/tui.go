package tui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/shaunmolloy/bugbox/internal/logging"
	"github.com/shaunmolloy/bugbox/internal/storage/config"
)

// Start initializes and runs the TUI
func Start() error {
	app := tview.NewApplication()

	handleKeyboardShortcuts(app)

	rootFlex := layout()
	if err := app.SetRoot(rootFlex, true).Run(); err != nil {
		logging.Error(fmt.Sprintf("Error: Failed to run TUI: %v", err))
		return err
	}

	return nil
}

func layout() tview.Primitive {
	rootFlex := tview.NewFlex().SetDirection(tview.FlexRow)

	innerFlex := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(issuesView(), 0, 1, false).
		AddItem(orgsView(), 30, 0, true)

	rootFlex.
		AddItem(innerFlex, 0, 1, true).
		AddItem(shortcutsView(), 1, 0, false)

	return rootFlex
}

func handleKeyboardShortcuts(app *tview.Application) {
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// If "q" is pressed, stop the application
		if event.Key() == tcell.KeyRune && event.Rune() == 'q' {
			app.Stop()
		}
		return event
	})
}

func orgsView() tview.Primitive {
	conf, _ := config.LoadConfig()

	searchField := tview.NewInputField().
		SetPlaceholder("Search for org...").
		SetPlaceholderTextColor(tcell.ColorWhite).
		SetFieldBackgroundColor(tcell.ColorBlue).
		SetFieldTextColor(tcell.ColorWhite)

	table := tview.NewTable().SetFixed(1, 0)
	for row, org := range conf.Orgs {
		table.SetCell(row+1, 0, tview.NewTableCell(org))
	}

	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	flex.SetTitle("Orgs").SetBorder(true)
	flex.AddItem(searchField, 1, 1, true)
	flex.AddItem(table, 0, 1, false)
	return flex
}

func issuesView() tview.Primitive {
	issues, _ := config.LoadIssues()

	table := tview.NewTable().SetFixed(1, 0) // lock header row

	// Header row
	headers := []string{"Title", "Org", "Created"}
	for i, h := range headers {
		cell := tview.NewTableCell(fmt.Sprintf("[::b]%s", h)).SetTextColor(tcell.ColorLightGrey)
		table.SetCell(0, i, cell)
	}

	// Data rows
	for row, issue := range issues {
		// Truncate title to 50 characters if it's too long
		title := issue.Title
		if len(title) > 72 {
			title = title[:72] + "..."
		}

		table.SetCell(row+1, 0, tview.NewTableCell(title))
		table.SetCell(row+1, 1, tview.NewTableCell(issue.Org))
		table.SetCell(row+1, 2, tview.NewTableCell(issue.CreatedAt.Format("2006-01-02")))
	}

	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	flex.SetTitle("Issues").SetBorder(true)
	flex.AddItem(table, 0, 1, true)
	return flex
}

func shortcutsView() tview.Primitive {
	return tview.NewTextView().
		SetText("Keyboard Shortcuts: Q - Quit").
		SetTextAlign(tview.AlignCenter)
}
