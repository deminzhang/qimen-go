package ui

import "github.com/hajimehoshi/ebiten/v2"

type Container struct {
	X, Y float32
	uis  []IUIPanel
}

func (c *Container) Add(u ...IUIPanel) {
	c.uis = append(c.uis, u...)
}

func (c *Container) Remove(u IUIPanel) {
	for i, p := range c.uis {
		if u == p {
			if i == 0 {
				c.uis = c.uis[1:]
			} else {
				if i+1 == len(c.uis) {
					c.uis = c.uis[:i]
				} else {
					c.uis = append(c.uis[:i], c.uis[i+1:]...)
				}
			}
		}
	}
}

func (c *Container) Update() {
	for _, u := range c.uis {
		if !u.IsDisabled() && u.IsVisible() {
			u.Update()
		}
	}
}

func (c *Container) Draw(dst *ebiten.Image) {
	for _, u := range c.uis {
		if u.IsVisible() {
			u.Draw(dst)
		}
	}
}
