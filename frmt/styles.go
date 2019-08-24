// +build !windows

package frmt

var (
	iconInitial = icon{fGBlue, "?"}
	iconGood    = icon{fGGreen, "✔"}
	iconWarn    = icon{fGYellow, "⚠"}
	iconBad     = icon{fGRed, "✗"}
	iconSelect  = icon{fGBold, "▸"}
)
