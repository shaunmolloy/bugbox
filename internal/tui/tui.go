package tui

import (
	"fmt"
	"os/exec"
	"runtime"
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
	showSearch         = false
	searchQuery        = ""
	currentScreenWidth = 0
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
	app.EnableMouse(false) // Disable mouse input

	handleKeyboardShortcuts(app)

	// Go routine for refresh handler and resize handler
	go func() {
		for range RefreshChan {
			logging.Info("Refreshing TUI")
			app.QueueUpdateDraw(func() {
				// Replace the layout with a refreshed one
				newLayout := layout()
				app.SetRoot(newLayout, true)
			})
		}
	}()

	// Create a handler for terminal resize events
	var lastWidth int
	app.SetBeforeDrawFunc(func(screen tcell.Screen) bool {
		width, _ := screen.Size()
		currentScreenWidth = width

		// If width crosses our breakpoint of 130, refresh the layout
		if width != lastWidth {
			RefreshChan <- struct{}{} // Trigger a refresh
		}

		lastWidth = width
		return false // Allow normal drawing
	})

	// Create the initial layout
	rootFlex := layout()

	if err := app.SetRoot(rootFlex, true).Run(); err != nil {
		logging.Error(fmt.Sprintf("Error: Failed to run TUI: %v", err))
		return err
	}

	return nil
}

func layout() tview.Primitive {
	rootFlex := tview.NewFlex().SetDirection(tview.FlexRow)

	// Use vertical layout for small screens
	useVerticalLayout := false
	if currentScreenWidth < 130 {
		useVerticalLayout = true
	}

	if useVerticalLayout {
		// Vertical layout for small screens
		rootFlex.AddItem(issuesView(), 0, 3, !showSearch) // Issues take 3/4 of height
		rootFlex.AddItem(orgsView(), 0, 1, false)         // Orgs take 1/4 of height
	} else {
		// Default horizontal layout for wider screens
		innerFlex := tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(issuesView(), 0, 1, !showSearch). // Focus on issues when search is hidden
			AddItem(orgsView(), 30, 0, false)         // Orgs take fixed 30 columns

		rootFlex.AddItem(innerFlex, 0, 1, !showSearch)
	}

	if showSearch {
		rootFlex.AddItem(searchView(), 1, 0, true) // Focus on search when visible
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

	// Create a selectable table
	table := tview.NewTable().
		SetFixed(1, 0).            // Lock header row
		SetSelectable(true, false) // Enable row selection only

	// Define column widths - use expansion to control the width ratios
	colExpansions := []int{8, 1, 1} // Title takes 80%, Org and Created take 10% each

	// Header row
	headers := []string{"Title", "Org", "Created"}
	for i, h := range headers {
		cell := tview.NewTableCell(fmt.Sprintf("[::b]%s", h)).
			SetTextColor(tcell.ColorLightGrey).
			SetExpansion(colExpansions[i])
		table.SetCell(0, i, cell)
	}

	// Filter issues based on searchQuery if it's not empty
	filteredIssues := issues
	if searchQuery != "" {
		filteredIssues = nil
		query := strings.ToLower(searchQuery)
		for _, issue := range issues {
			// Search in title and org
			if strings.Contains(strings.ToLower(issue.Title), query) ||
				strings.Contains(strings.ToLower(issue.Org), query) {
				filteredIssues = append(filteredIssues, issue)
			}
		}
	}

	// Data rows
	for row, issue := range filteredIssues {
		// Truncate title to 72 characters if it's too long
		title := issue.Title
		if len(title) > 72 {
			title = title[:72] + "..."
		}

		table.SetCell(row+1, 0, tview.NewTableCell(title))
		table.SetCell(row+1, 1, tview.NewTableCell(issue.Org))
		table.SetCell(row+1, 2, tview.NewTableCell(issue.CreatedAt.Format("2006-01-02")))
	}

	// Handle selection - only allow selecting data rows, not the header
	if len(filteredIssues) > 0 {
		table.Select(1, 0) // Select first data row by default
	}

	// Custom selection handler to prevent selecting header
	table.SetSelectionChangedFunc(func(row, column int) {
		// If header row is selected, move to first data row if available
		if row == 0 && len(filteredIssues) > 0 {
			table.Select(1, 0)
		}
	})

	table.SetSelectedFunc(func(row, column int) {
		// Handle row selection (user pressed Enter on a row)
		if row > 0 && row <= len(filteredIssues) {
			issue := filteredIssues[row-1]
			logging.Info(fmt.Sprintf("Opening issue in browser: %s", issue.Title))

			// Open the issue URL in the default browser
			if err := openBrowser(issue.URL); err != nil {
				logging.Error(fmt.Sprintf("Failed to open browser: %v", err))
			}
		}
	})

	// Set title to indicate filtering
	title := "Issues"
	if len(filteredIssues) != len(issues) {
		title = fmt.Sprintf("Issues (%d/%d)", len(filteredIssues), len(issues))
	}

	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	flex.SetTitle(title).SetBorder(true)
	flex.AddItem(table, 0, 1, true)
	return flex
}

func searchView() tview.Primitive {
	searchField := tview.NewInputField().
		SetPlaceholder("Search").
		SetPlaceholderTextColor(tcell.ColorWhite).
		SetPlaceholderStyle(tcell.StyleDefault.Background(tcell.Color238)).
		SetFieldBackgroundColor(tcell.Color238).
		SetFieldTextColor(tcell.ColorWhite).
		SetText(searchQuery). // Initialize with current search query
		SetChangedFunc(func(text string) {
			// Update search query as user types
			searchQuery = text
		}).
		SetDoneFunc(func(key tcell.Key) {
			// When user finishes input (hits Enter/Esc), hide the search view
			if key == tcell.KeyEnter {
				// Keep the current search query
				showSearch = false
				RefreshChan <- struct{}{} // Trigger a refresh
			} else if key == tcell.KeyEscape {
				// Clear search query
				searchQuery = ""
				showSearch = false
				RefreshChan <- struct{}{} // Trigger a refresh
			}
		})

	// Create a flex with horizontal padding
	flex := tview.NewFlex().
		AddItem(nil, 1, 0, false). // Left padding
		AddItem(searchField, 0, 1, true).
		AddItem(nil, 1, 0, false) // Right padding

	return flex
}

func shortcutsView() tview.Primitive {
	shortcuts := []string{
		"↑↓ - Navigate",
		"Enter - Open",
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

func openBrowser(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		return fmt.Errorf("unsupported platform")
	}

	return cmd.Start()
}
