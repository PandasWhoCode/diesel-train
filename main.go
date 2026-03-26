package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"golang.org/x/term"
)

const (
	hideCursor = "\033[?25l"
	showCursor = "\033[?25h"
	clearAll   = "\033[2J\033[H"
)

// Roof line — the @ beacon flashes on a 30-frame cycle (20 on, 10 off).
// When off the @ is replaced with * so it dims rather than disappears.
// The two exhaust stacks (/___\) are fixed on the roofline.
const (
	roofOn  = `                __@______________________/___\__/___\____  `
	roofOff = `                __*______________________/___\__/___\____  `
)

func roofLine(frame int) string {
	if frame%30 < 20 {
		return roofOn
	}
	return roofOff
}

// Locomotive body — drawn below the animated roof line.
// Nose faces left; train moves right-to-left.
var locoBody = []string{
	`          _____/ ___   ___   ___   ___   ___   ___   ___ |        `,
	`         / [==] |\o/| |   | |   | |   | |   | |   | |   ||       `,
	`        |  [**] |===| | @ | | @ | | @ | | @ | | @ | | @ ||       `,
	`        |  [==] |___| |___| |___| |___| |___| |___| |___||       `,
	`         \_______________________________________________/         `,
}

// Exhaust trail — 2 rows, drifts right of the stacks using o/* chars.
// Each frame the smoke shifts rightward, giving a trailing plume effect.
const exhaustOffset = 41

var exhaustFrames = [4][2]string{
	{`  o  o              `, `                    `},
	{`  *  *  o  o        `, `  o  o              `},
	{`        *  *  o  o  `, `  *  *              `},
	{`              *  *  `, `        *  *        `},
}

// Wheel animation frames.
var wheelFrames = [4]string{
	`          (oo)(oo)         (oo)(oo)         (oo)(oo)              `,
	`          (OO)(OO)         (OO)(OO)         (OO)(OO)              `,
	`          (oo)(oo)         (oo)(oo)         (oo)(oo)              `,
	`          (O-)(O-)         (O-)(O-)         (O-)(O-)              `,
}


func csi(row, col int) string {
	return fmt.Sprintf("\033[%d;%dH", row, col)
}

func termSize() (w, h int) {
	w, h, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return 80, 24
	}
	return
}

func locoWidth() int {
	w := len(roofOn)
	for _, l := range locoBody {
		if len(l) > w {
			w = len(l)
		}
	}
	for _, wf := range wheelFrames {
		if len(wf) > w {
			w = len(wf)
		}
	}
	return w
}

// locoHeight is roof + body + wheels.
func locoHeight() int { return 1 + len(locoBody) + 1 }

// printAt writes s clipped to terminal bounds at 0-indexed (row, col).
func printAt(row, col int, s string, termW int) {
	if row < 1 || col >= termW || len(s) == 0 {
		return
	}
	start := 0
	if col < 0 {
		start = -col
		col = 0
	}
	if start >= len(s) {
		return
	}
	end := len(s)
	if col+(end-start) > termW {
		end = start + (termW - col)
	}
	if end <= start {
		return
	}
	fmt.Printf("%s%s", csi(row, col+1), s[start:end])
}

func render(x, termW, termH, frame int) {
	lh := locoHeight()
	startRow := (termH-lh)/2 + 1
	if startRow < 5 {
		startRow = 5
	}

	// Erase 2 smoke rows + all loco rows.
	for i := -2; i < lh; i++ {
		r := startRow + i
		if r >= 1 && r <= termH {
			fmt.Printf("%s%s", csi(r, 1), strings.Repeat(" ", termW))
		}
	}

	// Draw exhaust trail above the roof.
	ef := exhaustFrames[(frame/2)%4]
	for i, s := range ef {
		printAt(startRow-2+i, x+exhaustOffset, s, termW)
	}

	// Draw animated roof (flashing beacon).
	printAt(startRow, x, roofLine(frame), termW)

	// Draw body rows.
	for i, l := range locoBody {
		printAt(startRow+1+i, x, l, termW)
	}

	// Draw animated wheels.
	printAt(startRow+1+len(locoBody), x, wheelFrames[frame%4], termW)
}

func main() {
	slow := flag.Bool("s", false, "slow mode (half speed)")
	fast := flag.Bool("f", false, "fast mode (double speed)")
	flag.Parse()

	delay := 30 * time.Millisecond
	switch {
	case *slow:
		delay = 60 * time.Millisecond
	case *fast:
		delay = 15 * time.Millisecond
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	done := make(chan struct{})

	cleanup := func() {
		fmt.Print(showCursor)
		fmt.Print(clearAll)
	}

	go func() {
		select {
		case <-sigs:
			cleanup()
			os.Exit(0)
		case <-done:
		}
	}()

	fmt.Print(hideCursor)
	fmt.Print(clearAll)

	termW, termH := termSize()
	lw := locoWidth()

	for x, frame := termW, 0; x > -lw; x, frame = x-1, frame+1 {
		render(x, termW, termH, frame)
		time.Sleep(delay)
	}

	close(done)
	cleanup()
}
