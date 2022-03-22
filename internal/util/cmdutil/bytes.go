package cmdutil

import (
	"fmt"
	"math"
)

func FormatBytes(bytes uint64) string {
	// See also BytesLabel.
	step := 1024
	if bytes < uint64(step) {
		return fmt.Sprintf("%d B", bytes)
	}
	units := []string{"KiB", "MiB", "GiB", "TiB", "PiB", "EiB"}
	var i int
	var unit string
	var quantity float64
	for i, unit = range units {
		quantity = float64(bytes) / math.Pow(float64(step), float64(i+1))
		if quantity < 1024 {
			break
		}
	}
	return fmt.Sprintf("%.3g %s", float64(float64(quantity)/100)*100, unit)
}
