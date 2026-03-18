package ui

import (
	"fmt"
	"strings"
	"time"
)

func (p *Printer) Progress(label string, seconds int) {
	if seconds <= 0 {
		return
	}

	const width = 20
	totalDuration := time.Duration(seconds) * time.Second
	stepDuration := totalDuration / time.Duration(width)

	for i := 0; i <= width; i++ {
		filled := strings.Repeat("█", i)
		empty := strings.Repeat("░", width-i)

		color := ansiRed
		if i > width/3 {
			color = ansiYellow
		}
		if i > (2*width)/3 {
			color = ansiGreen
		}

		fmt.Fprintf(
			p.out,
			"\r%s %s%s%s",
			label,
			color,
			filled+empty,
			ansiReset,
		)

		if i < width {
			time.Sleep(stepDuration)
		}
	}

	fmt.Fprintln(p.out)
}
