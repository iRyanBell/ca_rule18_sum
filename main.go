package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten"
)

const (
	WIDTH  = 600
	HEIGHT = 300
	SCALE  = 2
)

var kp bool
var ruleNumber int
var grid [HEIGHT][WIDTH]bool

func ruleNumToBits(r int) [5]bool {
	return [5]bool{
		r&16 == 16,
		r&8 == 8,
		r&4 == 4,
		r&2 == 2,
		r&1 == 1,
	}
}

// Sum of Left + Right (1 point) + Center (2 points)
func applyRule(row [WIDTH]bool) [WIDTH]bool {
	rule := ruleNumToBits(ruleNumber)

	rowNext := [WIDTH]bool{}

	for i := 0; i < len(row); i++ {
		var l, r bool
		// Wrap
		if i == 0 {
			l = row[WIDTH-1]
		} else {
			l = row[i-1]
		}

		if i == WIDTH-1 {
			r = row[0]
		} else {
			r = row[i+1]
		}

		c := row[i]
		sum := 0
		if l {
			sum++
		}
		if c {
			sum += 2
		}
		if r {
			sum++
		}

		switch sum {
		case 0:
			rowNext[i] = rule[0]
		case 1:
			rowNext[i] = rule[1]
		case 2:
			rowNext[i] = rule[2]
		case 3:
			rowNext[i] = rule[3]
		case 4:
			rowNext[i] = rule[4]
		}
	}

	return rowNext
}

func incRule() {
	ruleNumber += 1
	initialize()
}

func decRule() {
	if ruleNumber > 0 {
		ruleNumber -= 1
		initialize()
	}
}

func update(screen *ebiten.Image) error {
	if ebiten.IsDrawingSkipped() {
		return nil
	}

	// Key input for rule change
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		if !kp {
			// Increment rule once.
			kp = true
			incRule()
		}
	} else if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		if !kp {
			// Decrement rule once.
			kp = true
			decRule()
		}
	} else {
		// No keys are pressed.
		kp = false
	}

	for y := 0; y < HEIGHT; y++ {
		for x := 0; x < WIDTH; x++ {
			if grid[y][x] {
				screen.Set(x, y, color.RGBA{0xFF, 0xFF, 0xFF, 0xFF})
			}
		}

		if y < HEIGHT-1 {
			// Scroll
			grid[y] = grid[y+1]
		}
	}
	grid[HEIGHT-1] = applyRule(grid[HEIGHT-1])

	return nil
}

func initialize() {
	// Flip center bit on top row
	grid = [HEIGHT][WIDTH]bool{}
	grid[HEIGHT-1][int(WIDTH/2)] = true

	fmt.Println("RULE:", ruleNumber) // Print rule number to console
}

func main() {
	ruleNumber = 18
	initialize()
	err := ebiten.Run(update, WIDTH, HEIGHT, SCALE, "CA")
	if err != nil {
		fmt.Println(err)
	}
}
