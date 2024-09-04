package views

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type Cursor struct {
	visible bool
	symbol  string
}

func NewCursor(symbol string) Cursor {
	return Cursor{
		visible: true,
		symbol:  symbol,
	}
}

func (c *Cursor) Blink() tea.Cmd {
	return tea.Tick(time.Second/2, func(t time.Time) tea.Msg {
		return toggleCursorMsg{}
	})
}

func (c *Cursor) Update(msg tea.Msg) tea.Cmd {
	if _, ok := msg.(toggleCursorMsg); ok {
		c.visible = !c.visible
		return c.Blink()
	}
	return nil
}

func (c *Cursor) Render() string {
	if c.visible {
		return c.symbol
	}
	return ""
}

type toggleCursorMsg struct{}
