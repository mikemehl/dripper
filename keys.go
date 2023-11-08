package main

import "github.com/charmbracelet/bubbles/key"

var keyBindings = struct {
	quit     key.Binding
	menuNext key.Binding
	menuPrev key.Binding
	addSub   key.Binding
	confirm  key.Binding
}{
	quit:     key.NewBinding(key.WithKeys("q", "esc", "ctrl+c")),
	menuNext: key.NewBinding(key.WithKeys("tab")),
	menuPrev: key.NewBinding(key.WithKeys("shift+tab")),
	addSub:   key.NewBinding(key.WithKeys("a")),
	confirm:  key.NewBinding(key.WithKeys("enter")),
}
