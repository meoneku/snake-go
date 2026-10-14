package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type Point struct{ X, Y int }
type Direction struct{ DX, DY int }

var (
	Up    = Direction{0, -1}
	Down  = Direction{0, 1}
	Left  = Direction{-1, 0}
	Right = Direction{1, 0}
)

type Game struct {
	W, H      int
	Snake     []Point
	Dir       Direction
	NextDir   Direction
	Food      Point
	Obstacles []Point
	Score     int
	Best      int
	Over      bool
	Paused    bool
	Wrap      bool
	Tick      time.Duration
	NeedDraw  bool
}

func NewGame(w, h int) *Game {
	g := &Game{W: w, H: h, Tick: 150 * time.Millisecond}
	g.Reset()
	g.Best = loadBest()
	if g.Best < 0 {
		g.Best = 0
	}
	return g
}
func (g *Game) Reset() {
	cx, cy := g.W/2, g.H/2
	g.Snake = []Point{{cx, cy}, {cx - 1, cy}, {cx - 2, cy}}
	g.Dir = Right
	g.NextDir = Right
	g.Score = 0
	g.Over = false
	g.Paused = false
	g.Tick = 150 * time.Millisecond
	g.Obstacles = nil
	nObs := g.W * g.H / 40
	if nObs < 5 {
		nObs = 5
	}
	for i := 0; i < nObs; i++ {
		for tries := 0; tries < 50; tries++ {
			p := Point{1 + rand.Intn(g.W-2), 1 + rand.Intn(g.H-2)}
			if g.obsOccupied(p) {
				continue
			}
			snakeHit := false
			for _, s := range g.Snake {
				if s == p {
					snakeHit = true
					break
				}
			}
			if snakeHit {
				continue
			}
			g.Obstacles = append(g.Obstacles, p)
			break
		}
	}
	g.placeFood()
}
func (g *Game) placeFood() {
	for {
		p := Point{1 + rand.Intn(g.W-2), 1 + rand.Intn(g.H-2)}
		collide := false
		for _, s := range g.Snake {
			if s == p {
				collide = true
				break
			}
		}
		if !collide && !g.obsOccupied(p) {
			g.Food = p
			return
		}
	}
}
func (g *Game) obsOccupied(p Point) bool {
	for _, o := range g.Obstacles {
		if o == p {
			return true
		}
	}
	return false
}
func (g *Game) setDir(d Direction) {
	if d.DX == -g.Dir.DX && d.DY == -g.Dir.DY {
		return
	}
	if d.DX == -g.NextDir.DX && d.DY == -g.NextDir.DY {
		return
	}
	g.NextDir = d
}
func (g *Game) Update() {
	if g.Over || g.Paused {
		return
	}
	g.Dir = g.NextDir
	head := g.Snake[0]
	nh := Point{head.X + g.Dir.DX, head.Y + g.Dir.DY}
	if g.Wrap {
		nh.X = (nh.X + g.W) % g.W
		nh.Y = (nh.Y + g.H) % g.H
	} else {
		if nh.X <= 0 || nh.Y <= 0 || nh.X >= g.W-1 || nh.Y >= g.H-1 {
			g.gameOver()
			return
		}
	}
	if g.obsOccupied(nh) {
		g.gameOver()
		return
	}
	eating := nh == g.Food
	body := g.Snake
	if !eating {
		body = body[:len(body)-1]
	}
	for _, s := range body {
		if s == nh {
			g.gameOver()
			return
		}
	}
	g.Snake = append([]Point{nh}, g.Snake...)
	if eating {
		g.Score += 10
		if g.Score > g.Best {
			g.Best = g.Score
			saveBest(g.Best)
		}
		if g.Tick > 60*time.Millisecond {
			g.Tick -= 5 * time.Millisecond
		}
		g.placeFood()
	} else {
		g.Snake = g.Snake[:len(g.Snake)-1]
	}
}
func (g *Game) gameOver() {
	g.Over = true
	if g.Score > g.Best {
		g.Best = g.Score
		saveBest(g.Best)
	}
}
func sttyCmd(extra ...string) *exec.Cmd {
	flag := "-F"
	if runtime.GOOS == "darwin" {
		flag = "-f"
	}
	return exec.Command("stty", append([]string{flag, "/dev/tty"}, extra...)...)
}
func enableRaw() func() {
	old, err := sttyCmd("-g").Output()
	restore := func() {
		fmt.Print("\x1b[?25h")
		if err == nil {
			s := strings.TrimSpace(string(old))
			sttyCmd(s).Run()
		} else {
			sttyCmd("sane").Run()
		}
	}
	sttyCmd("cbreak", "min", "1", "-echo").Run()
	fmt.Print("\x1b[?25l")
	return restore
}
func readInput(ch chan<- string) {
	buf := make([]byte, 3)
	for {
		n, err := os.Stdin.Read(buf[:1])
		if err != nil || n == 0 {
			continue
		}
		b := buf[0]
		if b == 0x1b {
			n2 := readWithTimeout(buf[1:2], 30*time.Millisecond)
			if n2 == 0 {
				ch <- "esc"
				continue
			}
			if buf[1] != '[' {
				continue
			}
			n3 := readWithTimeout(buf[2:3], 30*time.Millisecond)
			if n3 == 0 {
				continue
			}
			switch buf[2] {
			case 'A':
				ch <- "up"
			case 'B':
				ch <- "down"
			case 'C':
				ch <- "right"
			case 'D':
				ch <- "left"
			}
			continue
		}
		switch b {
		case 'w', 'W':
			ch <- "up"
		case 's', 'S':
			ch <- "down"
		case 'a', 'A':
			ch <- "left"
		case 'd', 'D':
			ch <- "right"
		case 'p', 'P', ' ':
			ch <- "pause"
		case 'r', 'R':
			ch <- "restart"
		case 'm', 'M':
			ch <- "mode"
		case 'q', 'Q', 0x03:
			ch <- "quit"
		case '\r', '\n':
			ch <- "enter"
		}
	}
}
func readWithTimeout(p []byte, d time.Duration) int {
	type res struct {
		n   int
		err error
	}
	c := make(chan res, 1)
	go func() {
		n, err := os.Stdin.Read(p)
		c <- res{n, err}
	}()
	select {
	case r := <-c:
		if r.err != nil {
			return 0
		}
		return r.n
	case <-time.After(d):
		return 0
	}
}
func (g *Game) Draw() string {
	var sb strings.Builder
	sb.WriteString("\x1b[H\x1b[2J")
	mode := "Tembok"
	if g.Wrap {
		mode = "Wrap"
	}
	sb.WriteString(fmt.Sprintf("  SNAKE  |  Skor: %d  |  Terbaik: %d  |  Panjang: %d  |  Mode: %s (M)\n", g.Score, g.Best, len(g.Snake), mode))
	inSnake := make(map[Point]bool, len(g.Snake))
	for _, s := range g.Snake {
		inSnake[s] = true
	}
	head := g.Snake[0]
	for y := 0; y < g.H; y++ {
		for x := 0; x < g.W; x++ {
			p := Point{x, y}
			if x == 0 || y == 0 || x == g.W-1 || y == g.H-1 {
				sb.WriteString("\x1b[36m█\x1b[0m")
				continue
			}
			if g.obsOccupied(p) {
				sb.WriteString("\x1b[90m█\x1b[0m")
				continue
			}
			if p == g.Food {
				sb.WriteString("\x1b[31m★\x1b[0m")
				continue
			}
			if p == head {
				sb.WriteString("\x1b[32m●\x1b[0m")
				continue
			}
			if inSnake[p] {
				sb.WriteString("\x1b[32m■\x1b[0m")
				continue
			}
			sb.WriteString(" ")
		}
		sb.WriteString("\n")
	}
	if g.Over {
		sb.WriteString("\n  \x1b[31mGAME OVER!\x1b[0m  Tekan R / Enter untuk main lagi, Q untuk keluar.\n")
	} else if g.Paused {
		sb.WriteString("\n  PAUSE — tekan P / Spasi untuk lanjut.\n")
	} else {
		sb.WriteString("\n  Kontrol: WASD / Panah = gerak, P = pause, Q = keluar.\n")
	}
	return sb.String()
}
func bestFile() string {
	if home, err := os.UserHomeDir(); err == nil && home != "" {
		return filepath.Join(home, ".ular_best")
	}
	return ".ular_best"
}
func loadBest() int {
	b, err := os.ReadFile(bestFile())
	if err != nil {
		return 0
	}
	n, _ := strconv.Atoi(strings.TrimSpace(string(b)))
	if n < 0 {
		n = 0
	}
	return n
}
func saveBest(n int) {
	_ = os.WriteFile(bestFile(), []byte(strconv.Itoa(n)+"\n"), 0644)
}
func main() {
	rand.Seed(time.Now().UnixNano())
	g := NewGame(50, 20)
	restore := enableRaw()
	defer restore()
	input := make(chan string, 16)
	go readInput(input)
	ticker := time.NewTicker(g.Tick)
	defer ticker.Stop()
	lastTick := g.Tick
	fmt.Print(g.Draw())
	for {
		select {
		case key := <-input:
			switch key {
			case "up":
				g.setDir(Up)
			case "down":
				g.setDir(Down)
			case "left":
				g.setDir(Left)
			case "right":
				g.setDir(Right)
			case "pause":
				if !g.Over {
					g.Paused = !g.Paused
				}
			case "restart", "enter":
				if g.Over {
					g.Reset()
					g.Best = loadBest()
				}
			case "mode":
				g.Wrap = !g.Wrap
			case "quit", "esc":
				fmt.Print("\x1b[H\x1b[2J")
				fmt.Printf("Keluar. Skor terakhir: %d, Terbaik: %d. Bye!\n", g.Score, g.Best)
				return
			}
			fmt.Print(g.Draw())
		case <-ticker.C:
			if g.Tick != lastTick {
				ticker.Reset(g.Tick)
				lastTick = g.Tick
			}
			g.Update()
			fmt.Print(g.Draw())
		}
	}
}
