// Generates test data for log viewers containing colors and fonts.
package main

import (
	"fmt"

	"github.com/deref/rgbterm"
)

func main() {
	{
		fmt.Println("8 Colors")
		for i := 0; i < 8; i++ {
			fmt.Printf("\u001b[%dmX", 30+i)
		}
		fmt.Println("\u001b[0m")
	}

	fmt.Println()
	{
		fmt.Println("16 Colors")
		for i := 0; i < 16; i++ {
			fmt.Printf("\u001b[%d;1mX", 30+i)
		}
		fmt.Println("\u001b[0m")
	}

	fmt.Println()
	{
		fmt.Println("256 Colors")
		for i := 0; i < 16; i++ {
			for j := 0; j < 16; j++ {
				fmt.Printf("\u001b[38;5;%04dmX", i*16+j)
			}
			fmt.Println("\u001b[0m")
		}
	}

	fmt.Println()
	{
		fmt.Println("24-bit Color")
		i := 0
		for h := 0; h < 256; h++ {
			r, g, b := rgbterm.HSLtoRGB(float64(h)/256.0, 0.7, 0.5)
			fmt.Print(rgbterm.FgString("X", r, g, b))
			i++
			if i%32 == 0 {
				fmt.Println()
			}
		}
	}

	fmt.Println()
	{
		fmt.Println("decorations")
		fmt.Println("\u001b[1mBold\u001b[0m")
		fmt.Println("\u001b[2mFaint\u001b[0m")
		fmt.Println("\u001b[3mItalic\u001b[0m")
		fmt.Println("\u001b[4mUnderline\u001b[0m")
		fmt.Println("\u001b[5mSlow Blink\u001b[0m")
		fmt.Println("\u001b[6mFast Blink\u001b[0m") // Unsupported in Terminal.app on Mac.
		fmt.Println("\u001b[7mInvert\u001b[0m")
		fmt.Println("\u001b[8mConceal\u001b[0m")
		fmt.Println("\u001b[9mStrikethrough\u001b[0m") // Unsupported in Terminal.app on Mac.
	}
}
