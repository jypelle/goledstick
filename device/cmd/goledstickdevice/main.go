package main

import (
	"image/color"
	"machine"
	"time"
	"tinygo.org/x/drivers/ws2812"
)

type Theme int

const (
	SwimmingPoolTheme Theme = iota
	FireplaceTheme
	FiestaTheme
	WhiteTheme
)

const ThemeCount = 4

const LedCount = 16

func main() {

	var ledList [LedCount]color.RGBA

	var sinIdx = [LedCount]uint32{3, 182, 47, 32, 129, 210, 76, 27, 253, 38, 164, 117, 211, 0, 39, 145}
	var speedFactor = [LedCount]uint32{180, 138, 128, 135, 128, 135, 210, 135, 128, 135, 128, 135, 128, 163, 128, 195}

	var swimmingPoolThemeColorId = [LedCount]uint8{0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1}

	var fireplaceThemeColorId = [LedCount]uint8{0, 2, 1, 1, 0, 2, 0, 1, 2, 1, 0, 1, 0, 2, 0, 1}

	var fiestaThemeIdx = 0

	// Setup led strip device
	var neo = machine.D8
	neo.Configure(machine.PinConfig{Mode: machine.PinOutput})
	ledStripDevice := ws2812.New(neo)

	// Setup button device
	buttonDevice := machine.D3
	buttonDevice.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	var previousButtonState = false
	var buttonState = false

	// Init sinusoidal index for each led
	for i := range sinIdx {
		sinIdx[i] = sinIdx[i] % speedFactor[i]
	}

	// Setup starting theme
	var currentTheme = SwimmingPoolTheme

	var i int

	for {
		// Check new button state
		buttonState = !buttonDevice.Get()
		if buttonState && !previousButtonState {
			currentTheme = (currentTheme + 1) % ThemeCount
		}
		previousButtonState = buttonState

		// Update leds
		switch currentTheme {

		case SwimmingPoolTheme:
			for i = 0; i < LedCount; i++ {

				if swimmingPoolThemeColorId[i] == 0 {
					ledList[i] = color.RGBA{R: 0, G: uint8(uint32(SinTable[256*sinIdx[i]/speedFactor[i]]) * 20 / 256), B: uint8(uint32(SinTable[256*sinIdx[i]/speedFactor[i]]))}
				} else {
					ledList[i] = color.RGBA{R: 0, G: uint8(uint32(SinTable[256*sinIdx[i]/speedFactor[i]]) * 136 / 256), B: uint8(uint32(SinTable[256*sinIdx[i]/speedFactor[i]]) * 140 / 256)}
				}

				sinIdx[i] = (sinIdx[i] + 1) % speedFactor[i]
				if sinIdx[i] == 0 {
					swimmingPoolThemeColorId[i] = (swimmingPoolThemeColorId[i] + 1) % 2
				}
			}

		case FireplaceTheme:
			for i = 0; i < LedCount; i++ {

				switch fireplaceThemeColorId[i] {
				case 0:
					ledList[i] = color.RGBA{R: uint8(uint32(SinTable[256*sinIdx[i]/speedFactor[i]])), G: uint8(uint32(SinTable[256*sinIdx[i]/speedFactor[i]]) * 20 / 256), B: 0}
				case 1:
					ledList[i] = color.RGBA{R: uint8(uint32(SinTable[256*sinIdx[i]/speedFactor[i]]) * 185 / 256), G: uint8(uint32(SinTable[256*sinIdx[i]/speedFactor[i]]) * 91 / 256), B: 0}
				case 2:
					ledList[i] = color.RGBA{R: uint8(uint32(SinTable[256*sinIdx[i]/speedFactor[i]]) * 220 / 256), G: uint8(uint32(SinTable[256*sinIdx[i]/speedFactor[i]]) * 56 / 256), B: 0}
				}

				sinIdx[i] = (sinIdx[i] + 1) % speedFactor[i]
				if sinIdx[i] == 0 {
					fireplaceThemeColorId[i] = (fireplaceThemeColorId[i] + 1) % 3
				}
			}

		case FiestaTheme:
			ledList[(fiestaThemeIdx/10+0)%(LedCount)] = color.RGBA{R: 10, G: 255, B: 10}
			ledList[(fiestaThemeIdx/10+1)%(LedCount)] = color.RGBA{R: 0, G: 0, B: 0}
			ledList[(fiestaThemeIdx/10+2)%(LedCount)] = color.RGBA{R: 255, G: 10, B: 10}
			ledList[(fiestaThemeIdx/10+3)%(LedCount)] = color.RGBA{R: 0, G: 0, B: 0}
			ledList[(fiestaThemeIdx/10+4)%(LedCount)] = color.RGBA{R: 10, G: 10, B: 255}
			ledList[(fiestaThemeIdx/10+5)%(LedCount)] = color.RGBA{R: 0, G: 0, B: 0}
			ledList[(fiestaThemeIdx/10+6)%(LedCount)] = color.RGBA{R: 133, G: 10, B: 132}
			ledList[(fiestaThemeIdx/10+7)%(LedCount)] = color.RGBA{R: 0, G: 0, B: 0}
			ledList[(fiestaThemeIdx/10+8)%(LedCount)] = color.RGBA{R: 133, G: 132, B: 10}
			ledList[(fiestaThemeIdx/10+9)%(LedCount)] = color.RGBA{R: 0, G: 0, B: 0}
			ledList[(fiestaThemeIdx/10+10)%(LedCount)] = color.RGBA{R: 10, G: 153, B: 112}
			ledList[(fiestaThemeIdx/10+11)%(LedCount)] = color.RGBA{R: 0, G: 0, B: 0}
			ledList[(fiestaThemeIdx/10+12)%(LedCount)] = color.RGBA{R: 235, G: 40, B: 0}
			ledList[(fiestaThemeIdx/10+13)%(LedCount)] = color.RGBA{R: 0, G: 0, B: 0}
			ledList[(fiestaThemeIdx/10+14)%(LedCount)] = color.RGBA{R: 60, G: 60, B: 155}
			ledList[(fiestaThemeIdx/10+15)%(LedCount)] = color.RGBA{R: 0, G: 0, B: 0}
			fiestaThemeIdx = (fiestaThemeIdx + 1) % (10 * LedCount)

		case WhiteTheme:
			for i = 0; i < LedCount; i++ {
				ledList[i] = color.RGBA{R: 255, G: 230, B: 190}
			}
		}

		ledStripDevice.WriteColors(ledList[:])
		time.Sleep(18 * time.Millisecond)
	}
}
