package main

import (
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/micaelmcarvalho/coolify-tui/internal/config"
	"github.com/micaelmcarvalho/coolify-tui/internal/coolify"
	"github.com/micaelmcarvalho/coolify-tui/internal/ui"
)

func main() {
	log.SetFlags(0)

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "configure":
			if err := config.Configure(); err != nil {
				log.Fatalf("configuration error: %v", err)
			}
			return

		case "version", "--version", "-v":
			printVersion()
			return

		case "demo", "--demo", "-d":
			runProgram(
				ui.NewModel(
					coolify.NewDemoClient(),
				),
			)
			return
		}
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("configuration error: %v", err)
	}

	client := coolify.NewClient(
		cfg.CoolifyURL,
		cfg.CoolifyToken,
	)

	program := tea.NewProgram(
		ui.NewModel(client),
		tea.WithAltScreen(),
	)

	if _, err := program.Run(); err != nil {
		log.Fatalf("TUI error: %v", err)
	}
}

func runProgram(model ui.Model) {
	program := tea.NewProgram(
		model,
		tea.WithAltScreen(),
	)

	if _, err := program.Run(); err != nil {
		log.Fatalf("TUI error: %v", err)
	}
}
