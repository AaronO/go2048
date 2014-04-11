package main

import (
    "fmt"
    "math"
    "time"
    "strconv"

    "github.com/nsf/termbox-go"

    "./board"
)

// Board size
const (
    X = 4
    Y = 4
)

var (
    COLORS = []termbox.Attribute {
        termbox.ColorWhite,
        termbox.ColorGreen,
        termbox.ColorBlue,
        termbox.ColorCyan,
        termbox.ColorMagenta,
        termbox.ColorRed,
    }
)

func pad(str string, pad string, length int) string {
    for  {
        str = pad + str
        if len(str) > length-1 {
            return str[0:length]
        }
    }
 }

func cell_color(power int)termbox.Attribute {
    return COLORS[power % len(COLORS)]
}

func iPow(x int, n int) int {
    return int(math.Pow(float64(x), float64(n)))
}

func cell_str(power int)string {
    if(power == 0) {
        return ".";
    }
    // Convert to power of two and then string
    return strconv.Itoa(iPow(2, power))
}

func draw_cell(x int, y int, power int) {
    var tx = x*4;
    var color = cell_color(power)
    var str = pad(cell_str(power), " ", 4)
    //fmt.Printf("out str = %#v\nsf", str)

    for i, c := range str {
        //fmt.Println(i, c)
        //if(c == ' ' && tx != 999 && color != termbox.ColorYellow && i != 99) {
        if(c == ' ') {
            continue
        }
        termbox.SetCell(4+tx+i, 4+y, c, color, termbox.ColorDefault)
    }
}

func draw_score(b board.Board) {
    score := 0
    x := 24
    y := 4

    for _, v := range b.Values() {
        for i := 1; i < v+1; i++ {
            score = score + iPow(2, i)
        }
    }

    str := "Score: "+ strconv.Itoa(score)

    for i, c := range str {
        termbox.SetCell(x+i, y, c, termbox.ColorDefault, termbox.ColorDefault)
    }
}

func draw_board(b board.Board) {
    termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

    // Draw the cells
    for y := 0; y < Y; y++ {
    for x := 0; x < X; x++ {
        draw_cell(x, y, b.Cells[y][x])
    }}

    draw_score(b)

    termbox.Flush()
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

                // Can no longer play
                if !_board.Playable() {
                    fmt.Println("\n\n\n\n\nLOST !")
                    break loop
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

