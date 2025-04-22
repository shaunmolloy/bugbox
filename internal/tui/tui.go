package tui

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/shaunmolloy/bugbox/internal/logging"
	"github.com/shaunmolloy/bugbox/internal/storage/config"
)

// RefreshChan is a channel that receives signals to refresh the TUI
var RefreshChan = make(chan struct{}, 1)

// Global state for controlling UI elements
var (
	showSearch = false
)

// Init function to set up the global styles
func init() {
	// Set global styles to prevent focus highlighting
	tview.Styles.PrimitiveBackgroundColor = tcell.ColorBlack
	tview.Styles.BorderColor = tcell.ColorGrey
	tview.Styles.TitleColor = tcell.ColorGrey
}

// Start initializes and runs the TUI
func Start() error {
	app := tview.NewApplication()
	
	// Disable mouse input
	app.EnableMouse(false)

	handleKeyboardShortcuts(app)

	// Create the initial layout
	rootFlex := layout()

	// Set up refresh handler
	go func() {
		for range RefreshChan {
			logging.Info("Refreshing TUI with updated issues")
			app.QueueUpdateDraw(func() {
				// Replace the layout with a refreshed one
				newLayout := layout()
				app.SetRoot(newLayout, true)
			})
		}
	}()

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
		AddItem(orgsView(), 30, 0, false) // Only focus on orgs view if search is not showing

	rootFlex.AddItem(innerFlex, 0, 1, !showSearch)
	
	if showSearch {
		rootFlex.AddItem(searchView(), 1, 0, true) // Reduced height from 3 to 1
	}
	
	rootFlex.AddItem(shortcutsView(), 1, 0, false)

	return rootFlex
}

func handleKeyboardShortcuts(app *tview.Application) {
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// If "q" is pressed, stop the application
		if event.Key() == tcell.KeyRune && event.Rune() == 'q' {
			app.Stop()
			return nil
		}
		
		// If "/" is pressed, toggle search view
		if event.Key() == tcell.KeyRune && event.Rune() == '/' {
			showSearch = !showSearch
			// Use the existing refresh mechanism
			RefreshChan <- struct{}{}
			return nil // Consume the event
		}
		
		return event
	})
}

func orgsView() tview.Primitive {
	conf, _ := config.LoadConfig()

	table := tview.NewTable().SetFixed(1, 0)
	for row, org := range conf.Orgs {
		table.SetCell(row+0, 0, tview.NewTableCell(org))
	}

	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	flex.SetTitle("Orgs").SetBorder(true)

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

func searchView() tview.Primitive {
	searchField := tview.NewInputField().
		SetPlaceholder("Search").
		SetPlaceholderTextColor(tcell.ColorWhite).
		SetFieldBackgroundColor(tcell.ColorBlue).
		SetFieldTextColor(tcell.ColorWhite).
		SetDoneFunc(func(key tcell.Key) {
			// When user finishes input (hits Enter/Esc), hide the search view
			if key == tcell.KeyEnter || key == tcell.KeyEscape {
				showSearch = false
				RefreshChan <- struct{}{} // Trigger a refresh
			}
		})

	flex := tview.NewFlex()

	flex.
		SetTitle("Search").
		SetTitleColor(tcell.ColorWhite)
	
	flex.AddItem(searchField, 0, 1, true)
	return flex
}

func shortcutsView() tview.Primitive {
	shortcuts := []string{
		"/ - Search",
		"Q - Quit",
	}

	if showSearch {
		shortcuts = []string{
			"Enter - Search",
			"Esc - Cancel",
			"Q - Quit",
		}
	}

	return tview.NewTextView().
		SetText(strings.Join(shortcuts, ", ")).
		SetTextAlign(tview.AlignCenter)
}
