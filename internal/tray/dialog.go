package tray

import "fmt"

type Dialog struct{}

func (d *Dialog) ShowSettings() {
	fmt.Println("Settings dialog - runtime-driven capability metadata available")
}
