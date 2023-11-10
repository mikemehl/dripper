package main

import "github.com/charmbracelet/bubbles/key"

var keyBindings = struct {
	quit     key.Binding
	menuNext key.Binding
	menuPrev key.Binding
	addSub   key.Binding
	confirm  key.Binding
	up       key.Binding
	down     key.Binding
	left     key.Binding
	right    key.Binding
}{
	quit:     key.NewBinding(key.WithKeys("q", "esc", "ctrl+c")),
	menuNext: key.NewBinding(key.WithKeys("tab")),
	menuPrev: key.NewBinding(key.WithKeys("shift+tab")),
	addSub:   key.NewBinding(key.WithKeys("a")),
	confirm:  key.NewBinding(key.WithKeys("enter")),
	up:       key.NewBinding(key.WithKeys("up", "k")),
	down:     key.NewBinding(key.WithKeys("down", "j")),
	left:     key.NewBinding(key.WithKeys("left", "h")),
	right:    key.NewBinding(key.WithKeys("right", "l")),
}
