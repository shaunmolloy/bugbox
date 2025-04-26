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
	conf, _            = config.LoadConfig()
	showSearch         = false
	searchQuery        = ""
	orgFilter          = ""
	useVerticalLayout  = false
	currentScreenWidth = 0
	// Colors
	primaryColor   = tcell.ColorLimeGreen
	secondaryColor = tcell.ColorDarkOliveGreen
	grayColor      = tcell.ColorDarkGray
)

const (
	breakpointSmall  = 72
	breakpointXSmall = 100
	breakpointMedium = 130
	breakpointLarge  = 180
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
			logging.Info(fmt.Sprintf("Refreshing TUI. Width: %d", currentScreenWidth))
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
	useVerticalLayout = currentScreenWidth < breakpointLarge
	if useVerticalLayout {
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

		// If "Tab" is pressed, filter by next org
		if event.Key() == tcell.KeyTab && !showSearch {
			cycleOrgFilter()
			logging.Info(fmt.Sprintf("Filtering by org: %s", fallback(orgFilter, "(none)")))
			RefreshChan <- struct{}{}
			return nil // Consume the event
		}

		// If "Esc" is pressed, clear the org filter
		if event.Key() == tcell.KeyEscape && orgFilter != "" {
			orgFilter = ""
			logging.Info("Clearing org filter")
			RefreshChan <- struct{}{}
			return nil // Consume the event
		}

		return event
	})
}

func cycleOrgFilter() {
	if len(conf.Orgs) == 0 {
		return
	}

	// If orgFilter is empty, set it to the first org
	if orgFilter == "" {
		orgFilter = conf.Orgs[0]
		return
	}

	// Find the index of the current orgFilter
	index := indexOf(conf.Orgs, orgFilter)

	// Can we switch to the next org?
	if index+1 < len(conf.Orgs) {
		orgFilter = conf.Orgs[index+1]
		return
	}

	// If orgFilter was last org, reset to empty
	orgFilter = ""
}

func orgsView() tview.Primitive {
	table := tview.NewTable().SetFixed(1, 0)
	for row, org := range conf.Orgs {
		cell := tview.NewTableCell(org)
		if org == orgFilter {
			cell.SetBackgroundColor(tcell.ColorWhite).
				SetTextColor(tcell.ColorBlack).
				SetExpansion(1)
		}
		table.SetCell(row+0, 0, cell)
	}

	// Set title to indicate filtering
	title := fmt.Sprintf("Orgs (%d)", len(conf.Orgs))

	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	flex.SetTitle(title).SetTitleColor(primaryColor).SetBorder(true)

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
	colExpansions := []int{8, 1, 1} // Title, Org, Created

	// Header row
	headers := []string{"Title", "Org", "Created"}

	// Add "Repo" column if the screen is wide enough
	if currentScreenWidth > breakpointMedium {
		headers = append(headers, "")            // extend slice by 1
		index := 2                               // insert after "Org"
		copy(headers[index+1:], headers[index:]) // shift elements right
		headers[index] = "Repo"                  // insert new value
		colExpansions = []int{7, 1, 1, 1}        // Title, Repo, Org, Created
	}

	for i, h := range headers {
		cell := tview.NewTableCell(fmt.Sprintf("[::b]%s", h)).
			SetTextColor(tcell.ColorLightGrey).
			SetExpansion(colExpansions[i])
		table.SetCell(0, i, cell)
	}

	// Filter issues based on searchQuery and orgFilter
	filteredIssues := issues

	// Apply filters
	if searchQuery != "" || orgFilter != "" {
		filteredIssues = nil
		query := strings.ToLower(searchQuery)
		for _, issue := range issues {
			// Check if issue matches all active filters
			matchesSearch := true
			matchesOrg := true

			// Apply search filter if active
			if searchQuery != "" {
				matchesSearch = strings.Contains(strings.ToLower(issue.Title), query) ||
					strings.Contains(strings.ToLower(issue.Org), query)
			}

			// Apply org filter if active
			if orgFilter != "" {
				matchesOrg = issue.Org == orgFilter
			}

			// Add issue if it matches all active filters
			if matchesSearch && matchesOrg {
				filteredIssues = append(filteredIssues, issue)
			}
		}
	}

	// Data rows
	for row, issue := range filteredIssues {
		// Truncate title based on screen width
		title := issue.Title
		switch {
		case currentScreenWidth < breakpointSmall && len(title) > 30:
			title = title[:30]
			break
		case currentScreenWidth < breakpointXSmall && len(title) > 50:
			title = title[:50]
			break
		case currentScreenWidth < breakpointMedium && len(title) > 72:
			title = title[:72]
			break
		case currentScreenWidth < breakpointLarge && len(title) > 100:
			title = title[:100]
			break
		case len(title) > 120:
			title = title[:120]
			break
		}

		cells := []*tview.TableCell{
			tview.NewTableCell(title),
			tview.NewTableCell(issue.Org),
		}

		if currentScreenWidth > breakpointMedium {
			cells = append(cells, tview.NewTableCell(issue.Repo))
		}

		cells = append(cells, tview.NewTableCell(issue.CreatedAt.Format("2006-01-02")))

		for col, cell := range cells {
			table.SetCell(row+1, col, cell)
		}
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
	title := fmt.Sprintf("Issues (%d)", len(filteredIssues))
	if len(filteredIssues) != len(issues) {
		title = fmt.Sprintf("Issues (%d/%d)", len(filteredIssues), len(issues))
	}

	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	flex.SetTitle(title).SetTitleColor(primaryColor).SetBorder(true)
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
			RefreshChan <- struct{}{} // Trigger a refresh
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
		"Tab - Next Org",
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
		SetText(strings.Join(shortcuts, "    |    ")).
		SetTextAlign(tview.AlignCenter).
		SetTextColor(grayColor)
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
