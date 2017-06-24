package main

import (
	"strconv"
	"time"

	"github.com/nsf/termbox-go"

	"github.com/AaronO/go2048/board"
)

var (
	COLORS = []termbox.Attribute{
		termbox.ColorWhite,
		termbox.ColorGreen,
		termbox.ColorBlue,
		termbox.ColorCyan,
		termbox.ColorMagenta,
		termbox.ColorRed,
	}
)

func pad(str string, pad string, length int) string {
	for len(str) < length {
		str = pad + str
	}
	return str[0:length]
}

func cell_color(power int) termbox.Attribute {
	return COLORS[power%len(COLORS)]
}

func iPow(x int, n int) int {
	accu := 1
	for i := 0; i < n; i++ {
		accu *= x
	}

	return accu
}

func cell_str(power int) string {
	if power == 0 {
		return "."
	}
	// Convert to power of two and then string
	return strconv.Itoa(iPow(2, power))
}

func draw_cell(x int, y int, power int) {
	var tx = x * 4
	var color = cell_color(power)
	var str = pad(cell_str(power), " ", 4)
	//fmt.Printf("out str = %#v\nsf", str)

	for i, c := range str {
		//fmt.Println(i, c)
		//if(c == ' ' && tx != 999 && color != termbox.ColorYellow && i != 99) {
		if c == ' ' {
			continue
		}
		termbox.SetCell(4+tx+i, 4+y, c, color, termbox.ColorDefault)
	}
}

func draw_score(b board.Board) {
	// Compute score
	score := 0
	for _, v := range b.Values() {
		for n := v; n > 1; n-- {
			score += iPow(2, v)
		}
	}

	// Build score string
	str := "Score: " + strconv.Itoa(score)

	// Draw Score string
	draw_text(str, 24, 4, termbox.ColorDefault, termbox.ColorDefault)
}

func draw_cells(b board.Board) {
	// Draw the cells
	for y := 0; y < board.Y; y++ {
		for x := 0; x < board.X; x++ {
			draw_cell(x, y, b.Cells[y][x])
		}
	}
}

func draw_board(b board.Board) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	draw_cells(b)
	draw_score(b)

	termbox.Flush()
}

func draw_loser(b board.Board) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	draw_cells(b)
	draw_text("WOW SUCH LOSER", 24, 4, termbox.ColorRed, termbox.ColorDefault)

	termbox.Flush()
}

func draw_text(str string, x, y int, fg, bg termbox.Attribute) {
	for i, c := range str {
		termbox.SetCell(x+i, y, c, fg, bg)
	}
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}

	// Our Game board
	var _board = board.New()

	// Cleanup on exit
	defer termbox.Close()

	// Keyboard only
	termbox.SetInputMode(termbox.InputEsc)

	// Clear empty
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	termbox.Flush()

	draw_board(_board)

	// Event queue
	event_queue := make(chan termbox.Event)
	go func() {
		for {
			event_queue <- termbox.PollEvent()
		}
	}()

	// Event loop
loop:
	for {
		select {
		case ev := <-event_queue:
			switch ev.Type {
			case termbox.EventKey:
				// Exit
				if ev.Key == termbox.KeyCtrlC || ev.Key == termbox.KeyEsc {
					break loop
				}

				// Can no longer play
				if !_board.Playable() {
					draw_loser(_board)
					continue
				}

				switch ev.Key {
				case termbox.KeyArrowLeft:
					_board.Move(board.LEFT)
				case termbox.KeyArrowRight:
					_board.Move(board.RIGHT)
				case termbox.KeyArrowUp:
					_board.Move(board.UP)
				case termbox.KeyArrowDown:
					_board.Move(board.DOWN)
				}

				draw_board(_board)
			case termbox.EventResize:
				draw_board(_board)
			}

		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
}
