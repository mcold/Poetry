package main

import "fyne.io/fyne/v2"

func (app *Config) setKeys(win fyne.Window) {
	win.Canvas().SetOnTypedKey(func(k *fyne.KeyEvent) {
		switch k.Name {
		case "F2":
			{

			}
		}
	})
}
