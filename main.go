package main

import (
    "log"
    "AP/tui"
    "fmt"
    "os"
		tea "github.com/charmbracelet/bubbletea"
)

func main() {
        f, err := tea.LogToFile("assets/debug.log", "debug")
        if err != nil {
          log.Fatal(err)
        }

        defer f.Close()

				p := tea.NewProgram(tui.InitialModel(), tea.WithAltScreen())
				if _, err := p.Run(); err != nil {
								fmt.Println(err)
								os.Exit(1)
				}
}


