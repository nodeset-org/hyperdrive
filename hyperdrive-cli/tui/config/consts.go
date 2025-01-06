package config

import (
	"github.com/gdamore/tcell/v2"
)

const (
	// Background for non-interactive elements
	NonInteractiveBackgroundColor tcell.Color = tcell.ColorBlack

	// Background for all UI elements
	BackgroundColor tcell.Color = tcell.ColorDarkSlateGray

	// Border
	BorderColor tcell.Color = tcell.ColorGold

	// Unfocused buttons
	ButtonUnfocusedBackgroundColor tcell.Color = tcell.ColorBlack
	ButtonUnfocusedTextColor       tcell.Color = tcell.ColorLightGray

	// Focused buttons
	ButtonFocusedBackgroundColor tcell.Color = tcell.Color46 // A lovely bright green
	ButtonFocusedTextColor       tcell.Color = tcell.ColorBlack

	// Unfocused home buttons
	HomeButtonUnfocusedBackgroundColor tcell.Color = tcell.ColorDarkSlateGray
	HomeButtonUnfocusedTextColor       tcell.Color = tcell.ColorLightGray
)
