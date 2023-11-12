package pages

import "github.com/charmbracelet/bubbles/key"

var KeyBindings = struct {
	Back     key.Binding
	Quit     key.Binding
	MenuNext key.Binding
	MenuPrev key.Binding
	AddSub   key.Binding
	Confirm  key.Binding
	Up       key.Binding
	Down     key.Binding
	Left     key.Binding
	Right    key.Binding
}{
	Back:     key.NewBinding(key.WithKeys("esc")),
	Quit:     key.NewBinding(key.WithKeys("q", "ctrl+c")),
	MenuNext: key.NewBinding(key.WithKeys("tab")),
	MenuPrev: key.NewBinding(key.WithKeys("shift+tab")),
	AddSub:   key.NewBinding(key.WithKeys("a")),
	Confirm:  key.NewBinding(key.WithKeys("enter")),
	Up:       key.NewBinding(key.WithKeys("up", "k")),
	Down:     key.NewBinding(key.WithKeys("down", "j")),
	Left:     key.NewBinding(key.WithKeys("left", "h")),
	Right:    key.NewBinding(key.WithKeys("right", "l")),
}
