// Pixel Life is a graphical interface for an underlying CA rule-set
package main

import (
	"fmt"
	"math/rand"
	"time"

	"pixel-life/life"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

// alive and dead used to clarify when a cell's state is on or off
const (
	alive = true
	dead  = false
)

// Settings for session variables
var (
	screenColor  = colornames.Black
	lifeWidth    = 150
	lifeHeight   = 80
	cellWidth    = 8.0
	cellBorder   = 1.0
	block        = cellWidth + cellBorder
	screenWidth  = block*float64(lifeWidth) + cellBorder
	screenHeight = block*float64(lifeHeight) + cellBorder
)

type (
	Cell pixel.Rect

	Cells [][]pixel.Rect

	position struct {
		x, y int
	}

	// A PixelLife is the combination of a traditionally defined Life and Cells
	PixelLife struct {
		life  *life.Life
		cells *Cells
	}
)

// Main calls pixelgl.Run and is then unused
func main() {
	pixelgl.Run(run)
}

// run is functionally the start of the game logic
func run() {
	rand.Seed(time.Now().UnixNano())

	cfg := pixelgl.WindowConfig{
		Title:  "Pixel Life",
		Bounds: pixel.R(0, 0, screenWidth, screenHeight),
		VSync:  false,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	win.Clear(screenColor)

	imd := imdraw.New(nil)

	p := PixelLife{
		life:  life.NewLife(lifeWidth, lifeHeight),
		cells: NewCells(lifeWidth, lifeHeight),
	}

	var (
		frames      = 0
		generations = 1
		auto        = false
		second      = time.Tick(time.Second)
		rate        = time.Tick(time.Second / 60)
	)

	last := time.Now()
	for !win.Closed() {
		// calculate frame time in ms
		dt := time.Since(last).Milliseconds()
		last = time.Now()
		// Clear content for next frame
		imd.Clear()
		win.Clear(screenColor)

		// Controls
		if win.JustReleased(pixelgl.KeyQ) {
			break
		}
		if win.JustPressed(pixelgl.KeyC) {
			auto = false
			p.life.Clear()
			generations = 1
		}
		if win.JustPressed(pixelgl.KeyF) {
			auto = false
			p.life.Fill()
			generations = 1
		}
		if win.JustPressed(pixelgl.KeyN) && !auto {
			p.life.Next()
			generations++
		}
		if win.JustPressed(pixelgl.KeyG) {
			p.life.Rand()
			generations = 1
		}
		if win.Pressed(pixelgl.MouseButtonLeft) {
			mp := win.MousePosition()
			if yes, v := p.cells.Contains(win, mp); yes {
				(*p.life)[v.x][v.y] = alive
			}
		}
		if win.Pressed(pixelgl.MouseButtonRight) {
			mp := win.MousePosition()
			if yes, v := p.cells.Contains(win, mp); yes {
				(*p.life)[v.x][v.y] = dead
			}
		}
		if win.JustPressed(pixelgl.KeySpace) {
			auto = !auto
		}

		// auto determines whether state automatically updates each frame
		if auto {
			select {
			case <-rate:
				p.life.Next()
				generations++
			default:
			}
		}

		// If the state of a cell is alive, draw the cell.
		for i := range *p.life {
			for j := range (*p.life)[i] {
				if (*p.life)[i][j] {
					Cell((*p.cells)[i][j]).Draw(imd)
				}
			}
		}

		imd.Draw(win)
		win.Update()

		frames++
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d | FT: %d | Generations: %d", cfg.Title, frames, dt, generations))
			frames = 0
		default:
		}
	}
}

// Contains checks if the mouse position is inside a cell rectangle
func (c *Cells) Contains(win *pixelgl.Window, mp pixel.Vec) (bool, position) {
	for i := range *c {
		for j, v := range (*c)[i] {
			if v.Contains(mp) {
				return true, position{i, j}
			}
		}
	}
	return false, position{0, 0}
}

// Draw draws a cell to the specified IMDraw object
func (c Cell) Draw(imd *imdraw.IMDraw) {
	imd.Color = colornames.White
	imd.Push(c.Min)
	imd.Push(c.Max)
	imd.Rectangle(0)
}

// NewCells creates a new field of Cells with the specified dimensions
func NewCells(w, h int) *Cells {
	c := make(Cells, w)
	for i := range c {
		x := cellBorder + (block * float64(i))
		dx := block + (block * float64(i))
		c[i] = make([]pixel.Rect, h)
		for j := range c[i] {
			y := (screenHeight - block) - (block * float64(j))
			dy := (screenHeight - cellBorder) - (block * float64(j))
			c[i][j] = pixel.R(x, y, dx, dy)
		}
	}
	return &c
}
