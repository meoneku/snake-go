package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/font/opentype"
	"image/color"
	"log"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	cols     = 25
	rows     = 25
	cell     = 24
	margin   = 16
	hudH     = 96
	winW     = cols*cell + margin*2
	winH     = hudH + rows*cell + margin*2
	startGap = 9
	minGap   = 4
)

var (
	cBGTop    = color.RGBA{0x1b, 0x23, 0x38, 0xff}
	cBGBot    = color.RGBA{0x0b, 0x0f, 0x1a, 0xff}
	cDot      = color.RGBA{0xff, 0xff, 0xff, 0x10}
	cBoardA   = color.RGBA{0x20, 0x2b, 0x45, 0xff}
	cBoardB   = color.RGBA{0x24, 0x31, 0x4e, 0xff}
	cBorder   = color.RGBA{0x38, 0xbd, 0xf8, 0xff}
	cBorderHi = color.RGBA{0xa5, 0xf3, 0xfc, 0xff}
	cPill     = color.RGBA{0x2b, 0x38, 0x58, 0xff}
	cPillEdge = color.RGBA{0x3b, 0x4c, 0x74, 0xff}
	cAccent   = color.RGBA{0x38, 0xbd, 0xf8, 0xff}
	cGold     = color.RGBA{0xfb, 0xbf, 0x24, 0xff}
	cText     = color.RGBA{0xf1, 0xf5, 0xf9, 0xff}
	cDim      = color.RGBA{0x94, 0xa3, 0xb8, 0xff}
	cHead     = color.RGBA{0x4a, 0xde, 0x80, 0xff}
	cTail     = color.RGBA{0x14, 0xb8, 0xa6, 0xff}
	cOutline  = color.RGBA{0x05, 0x2e, 0x1b, 0xff}
	cApple    = color.RGBA{0xef, 0x44, 0x44, 0xff}
	cAppleDk  = color.RGBA{0x99, 0x1b, 0x1b, 0xff}
	cStem     = color.RGBA{0x92, 0x4d, 0x0e, 0xff}
	cLeaf     = color.RGBA{0x4a, 0xde, 0x80, 0xff}
	cTongue   = color.RGBA{0xf4, 0x72, 0xb6, 0xff}
	cWhite    = color.RGBA{0xff, 0xff, 0xff, 0xff}
	cPanel    = color.RGBA{0x15, 0x1d, 0x33, 0xf2}
	cFlash    = color.RGBA{0xef, 0x44, 0x44, 0xff}
)
var (
	fTitle font.Face
	fHUD   font.Face
	fSmall font.Face
	fBig   font.Face
)

func loadFace(paths []string, size float64) font.Face {
	for _, p := range paths {
		b, err := os.ReadFile(p)
		if err != nil {
			continue
		}
		f, err := opentype.Parse(b)
		if err != nil {
			continue
		}
		face, err := opentype.NewFace(f, &opentype.FaceOptions{
			Size:    size,
			DPI:     72,
			Hinting: font.HintingFull,
		})
		if err != nil {
			continue
		}
		return face
	}
	return basicfont.Face7x13
}
func loadFonts() {
	bold := []string{
		"/usr/share/fonts/TTF/DejaVuSans-Bold.ttf",
		"/usr/share/fonts/truetype/dejavu/DejaVuSans-Bold.ttf",
		"C:/Windows/Fonts/arialbd.ttf",
		"/System/Library/Fonts/Helvetica.ttc",
	}
	regular := []string{
		"/usr/share/fonts/TTF/DejaVuSans.ttf",
		"/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf",
		"C:/Windows/Fonts/arial.ttf",
		"/System/Library/Fonts/Helvetica.ttc",
	}
	fTitle = loadFace(bold, 20)
	fHUD = loadFace(bold, 15)
	fSmall = loadFace(regular, 12)
	fBig = loadFace(bold, 40)
}

type Point struct{ X, Y int }
type Dir struct{ DX, DY int }

var (
	Up    = Dir{0, -1}
	Down  = Dir{0, 1}
	Left  = Dir{-1, 0}
	Right = Dir{1, 0}
)

type Particle struct {
	x, y, vx, vy float64
	life, max    int
	size         float32
	clr          color.Color
}
type FloatText struct {
	x, y float64
	txt  string
	life int
	max  int
	clr  color.Color
	face font.Face
}
type Bonus struct {
	pos  Point
	life int
	max  int
}
type Game struct {
	snake     []Point
	prev      []Point
	dir       Dir
	nextDir   Dir
	food      Point
	bonus     *Bonus
	obstacles []Point
	score     int
	best      int
	eaten     int
	over      bool
	paused    bool
	started   bool
	wrap      bool
	gap       int
	tick      int
	frame     int
	flash     float64
	shake     float64
	parts     []Particle
	floats    []FloatText
	bg        *ebiten.Image
	board     *ebiten.Image
	boardBG   *ebiten.Image
	px, py    []float64
	audioCtx  *audio.Context
}

func NewGame() *Game {
	g := &Game{gap: startGap}
	g.bg = makeBackground()
	g.board = ebiten.NewImage(cols*cell, rows*cell)
	g.boardBG = ebiten.NewImage(cols*cell, rows*cell)
	g.boardBG.Fill(cBoardA)
	g.audioCtx = audio.NewContext(44100)
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			if (x+y)%2 == 1 {
				vector.FillRect(g.boardBG, float32(x*cell), float32(y*cell), cell, cell, cBoardB, false)
			}
		}
	}
	g.Reset()
	g.started = false
	g.best = loadBest()
	return g
}
func (g *Game) Reset() {
	cx, cy := cols/2, rows/2
	g.snake = []Point{{cx, cy}, {cx - 1, cy}, {cx - 2, cy}}
	g.prev = append([]Point(nil), g.snake...)
	g.dir, g.nextDir = Right, Right
	g.score, g.eaten = 0, 0
	g.over, g.paused = false, false
	g.started = true
	g.gap, g.tick = startGap, 0
	g.flash, g.shake = 0, 0
	g.bonus = nil
	g.parts = nil
	g.floats = nil
	g.food = Point{-1, -1}
	g.obstacles = nil
	nObs := cols * rows / 40
	if nObs < 8 {
		nObs = 8
	}
	for i := 0; i < nObs; i++ {
		for tries := 0; tries < 50; tries++ {
			p := Point{rand.Intn(cols), rand.Intn(rows)}
			if g.free(p) {
				g.obstacles = append(g.obstacles, p)
				break
			}
		}
	}
	g.placeFood()
}
func (g *Game) placeFood() {
	for {
		p := Point{rand.Intn(cols), rand.Intn(rows)}
		if g.free(p) {
			g.food = p
			return
		}
	}
}
func (g *Game) free(p Point) bool {
	for _, s := range g.snake {
		if s == p {
			return false
		}
	}
	if g.bonus != nil && *&g.bonus.pos == p {
		return false
	}
	for _, o := range g.obstacles {
		if o == p {
			return false
		}
	}
	return p != g.food
}
func (g *Game) obsOccupied(p Point) bool {
	for _, o := range g.obstacles {
		if o == p {
			return true
		}
	}
	return false
}
func (g *Game) setDir(d Dir) {
	if d.DX == -g.dir.DX && d.DY == -g.dir.DY {
		return
	}
	if d.DX == -g.nextDir.DX && d.DY == -g.nextDir.DY {
		return
	}
	g.nextDir = d
}
func (g *Game) step() {
	g.prev = append([]Point(nil), g.snake...)
	g.dir = g.nextDir
	head := g.snake[0]
	nh := Point{head.X + g.dir.DX, head.Y + g.dir.DY}
	if g.wrap {
		nh.X = (nh.X%cols + cols) % cols
		nh.Y = (nh.Y%rows + rows) % rows
	} else if nh.X < 0 || nh.Y < 0 || nh.X >= cols || nh.Y >= rows {
		g.finish()
		g.playBeep(180, 260, 0.35)
		return
	}
	if g.obsOccupied(nh) {
		g.finish()
		g.playBeep(180, 260, 0.35)
		return
	}
	eatApple := nh == g.food
	eatBonus := g.bonus != nil && nh == g.bonus.pos
	eating := eatApple || eatBonus
	body := g.snake
	if !eating {
		body = body[:len(body)-1]
	}
	for _, s := range body {
		if s == nh {
			g.finish()
			return
		}
	}
	g.snake = append([]Point{nh}, g.snake...)
	if len(g.snake) > len(g.prev) {
		g.prev = append(g.prev, g.prev[len(g.prev)-1])
	}
	if eatApple {
		g.eaten++
		g.score += 10
		if g.score > g.best {
			g.best = g.score
			saveBest(g.best)
		}
		g.playBeep(880, 70, 0.25)
		if g.gap > minGap {
			g.gap--
		}
		fx, fy := g.cellCenter(g.food)
		g.burst(fx, fy, cApple, 10)
		g.burst(fx, fy, cHead, 5)
		g.addFloat(fx, fy-cell/2, "+10", cHead, fHUD)
		g.placeFood()
		if g.bonus == nil && g.eaten%5 == 0 {
			for i := 0; i < 40; i++ {
				p := Point{rand.Intn(cols), rand.Intn(rows)}
				if g.free(p) {
					g.bonus = &Bonus{pos: p, life: 480, max: 480}
					break
				}
			}
		}
	}
	if eatBonus {
		g.score += 50
		if g.score > g.best {
			g.best = g.score
			saveBest(g.best)
		}
		g.playBeep(1320, 90, 0.28)
		fx, fy := g.cellCenter(g.bonus.pos)
		g.burst(fx, fy, cGold, 16)
		g.addFloat(fx, fy-cell/2, "+50", cGold, fTitle)
		g.bonus = nil
	}
	if !eating {
		g.snake = g.snake[:len(g.snake)-1]
	}
}
func (g *Game) finish() {
	g.over = true
	g.flash = 0.55
	g.shake = 9
	if g.score > g.best {
		g.best = g.score
		saveBest(g.best)
	}
	g.playBeep(220, 320, 0.35)
	hx, hy := g.cellCenter(g.snake[0])
	g.burst(hx, hy, cApple, 18)
}
func (g *Game) burst(x, y float64, clr color.Color, n int) {
	for i := 0; i < n; i++ {
		a := rand.Float64() * 2 * math.Pi
		sp := 1 + rand.Float64()*3
		g.parts = append(g.parts, Particle{
			x: x, y: y,
			vx: math.Cos(a) * sp, vy: math.Sin(a)*sp - 1,
			life: 25 + rand.Intn(20), max: 45,
			size: 2 + float32(rand.Float64()*3),
			clr:  clr,
		})
	}
}
func (g *Game) addFloat(x, y float64, txt string, clr color.Color, face font.Face) {
	g.floats = append(g.floats, FloatText{x: x, y: y, txt: txt, life: 50, max: 50, clr: clr, face: face})
}
func (g *Game) cellCenter(p Point) (float64, float64) {
	return float64(margin + p.X*cell + cell/2), float64(hudH + margin + p.Y*cell + cell/2)
}
func (g *Game) Update() error {
	g.frame++
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		g.setDir(Up)
		g.started = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		g.setDir(Down)
		g.started = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		g.setDir(Left)
		g.started = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		g.setDir(Right)
		g.started = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyP) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		if g.started && !g.over {
			g.paused = !g.paused
		} else {
			g.started = true
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.Reset()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		if g.over {
			g.Reset()
		} else {
			g.started = true
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) || inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		return ebiten.Termination
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyM) {
		g.wrap = !g.wrap
	}
	if g.started && !g.over && !g.paused {
		g.tick++
		if g.tick >= g.gap {
			g.tick = 0
			g.step()
		}
		if g.bonus != nil {
			g.bonus.life--
			if g.bonus.life <= 0 {
				g.bonus = nil
			}
		}
	}
	for i := range g.parts {
		p := &g.parts[i]
		p.x += p.vx
		p.y += p.vy
		p.vx *= 0.96
		p.vy = p.vy*0.96 + 0.15
		p.life--
	}
	keep := g.parts[:0]
	for _, p := range g.parts {
		if p.life > 0 {
			keep = append(keep, p)
		}
	}
	g.parts = keep
	for i := range g.floats {
		f := &g.floats[i]
		f.y -= 0.8
		f.life--
	}
	kf := g.floats[:0]
	for _, f := range g.floats {
		if f.life > 0 {
			kf = append(kf, f)
		}
	}
	g.floats = kf
	if g.flash > 0 {
		g.flash -= 0.03
	}
	if g.shake > 0 {
		g.shake *= 0.85
		if g.shake < 0.3 {
			g.shake = 0
		}
	}
	return nil
}
func makeBackground() *ebiten.Image {
	img := ebiten.NewImage(winW, winH)
	for y := 0; y < winH; y++ {
		t := float64(y) / float64(winH)
		r := uint8(float64(cBGTop.R)*(1-t) + float64(cBGBot.R)*t)
		gg := uint8(float64(cBGTop.G)*(1-t) + float64(cBGBot.G)*t)
		b := uint8(float64(cBGTop.B)*(1-t) + float64(cBGBot.B)*t)
		for x := 0; x < winW; x++ {
			img.Set(x, y, color.RGBA{r, gg, b, 0xff})
		}
	}
	for y := 12; y < winH; y += 26 {
		for x := 12; x < winW; x += 26 {
			img.Set(x, y, cDot)
		}
	}
	return img
}
func roundRect(dst *ebiten.Image, x, y, w, h, r float32, clr color.Color) {
	vector.FillRect(dst, x+r, y, w-2*r, h, clr, true)
	vector.FillRect(dst, x, y+r, w, h-2*r, clr, true)
	vector.DrawFilledCircle(dst, x+r, y+r, r, clr, true)
	vector.DrawFilledCircle(dst, x+w-r, y+r, r, clr, true)
	vector.DrawFilledCircle(dst, x+r, y+h-r, r, clr, true)
	vector.DrawFilledCircle(dst, x+w-r, y+h-r, r, clr, true)
}
func lerp(a, b, t float64) float64 { return a + (b-a)*t }
func mixColor(a, b color.RGBA, t float64) color.Color {
	return color.RGBA{
		uint8(lerp(float64(a.R), float64(b.R), t)),
		uint8(lerp(float64(a.G), float64(b.G), t)),
		uint8(lerp(float64(a.B), float64(b.B), t)),
		0xff,
	}
}
func drawCentered(dst *ebiten.Image, s string, face font.Face, cx, y int, clr color.Color) {
	b := text.BoundString(face, s)
	text.Draw(dst, s, face, cx-b.Dx()/2, y, clr)
}
func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(g.bg, nil)
	g.drawHUD(screen)
	g.board.Clear()
	bw, bh := float32(cols*cell), float32(rows*cell)
	g.drawObstacles(g.board)
	g.drawFood(g.board)
	g.drawSnake(g.board)
	op := &ebiten.DrawImageOptions{}
	sx, sy := 0.0, 0.0
	if g.shake > 0 {
		sx = (rand.Float64()*2 - 1) * g.shake
		sy = (rand.Float64()*2 - 1) * g.shake
	}
	op.GeoM.Translate(float64(margin)+sx, float64(hudH+margin)+sy)
	screen.DrawImage(g.boardBG, op)
	screen.DrawImage(g.board, op)
	bx, by := float32(margin)+float32(sx), float32(hudH+margin)+float32(sy)
	vector.StrokeRect(screen, bx-2, by-2, bw+4, bh+4, 2, cBorder, true)
	vector.StrokeRect(screen, bx-4, by-4, bw+8, bh+8, 1, cBorderHi, true)
	g.drawParticles(screen)
	g.drawFloats(screen)
	if g.flash > 0 {
		a := uint8(math.Min(g.flash, 0.6) * 255)
		vector.FillRect(screen, 0, 0, winW, winH, color.RGBA{0xef, 0x44, 0x44, a}, false)
	}
	if !g.started {
		g.drawStart(screen)
	} else if g.over {
		g.drawGameOver(screen)
	} else if g.paused {
		g.drawPaused(screen)
	}
}
func (g *Game) drawHUD(screen *ebiten.Image) {
	text.Draw(screen, "Ular-ularan", fTitle, margin, 30, cText)
	pills := []struct {
		label string
		val   string
		clr   color.Color
	}{
		{"SKOR", fmt.Sprint(g.score), cAccent},
		{"TERBAIK", fmt.Sprint(g.best), cGold},
		{"PANJANG", fmt.Sprint(len(g.snake)), cHead},
		{"LEVEL", fmt.Sprint(startGap - g.gap + 1), cText},
	}
	x := float32(winW - margin)
	for i := len(pills) - 1; i >= 0; i-- {
		p := pills[i]
		lb := text.BoundString(fSmall, p.label)
		vb := text.BoundString(fHUD, p.val)
		pw := float32(max(lb.Dx(), vb.Dx()) + 28)
		x -= pw
		roundRect(screen, x, 12, pw, 56, 12, cPill)
		vector.StrokeRect(screen, x+1, 13, pw-2, 54, 1, cPillEdge, true)
		drawCentered(screen, p.label, fSmall, int(x+pw/2), 30, cDim)
		drawCentered(screen, p.val, fHUD, int(x+pw/2), 56, p.clr)
		x -= 8
	}
	text.Draw(screen, "WASD / Panah gerak   •   P pause   •   R restart   •   M mode   •   Esc keluar", fSmall, margin, 100, cDim)
}
func (g *Game) drawSnake(dst *ebiten.Image) {
	n := len(g.snake)
	if n == 0 {
		return
	}
	t := 0.0
	if g.gap > 0 {
		t = float64(g.tick) / float64(g.gap)
	}
	if cap(g.px) < n {
		g.px = make([]float64, n)
	} else {
		g.px = g.px[:n]
	}
	if cap(g.py) < n {
		g.py = make([]float64, n)
	} else {
		g.py = g.py[:n]
	}
	px, py := g.px, g.py
	for i := 0; i < n; i++ {
		c, p := g.snake[i], g.prev[min(i, len(g.prev)-1)]
		px[i] = (lerp(float64(p.X), float64(c.X), t) + 0.5) * cell
		py[i] = (lerp(float64(p.Y), float64(c.Y), t) + 0.5) * cell
	}
	headW, tailW := float32(cell*0.86), float32(cell*0.34)
	for pass := 0; pass < 2; pass++ {
		for i := n - 1; i >= 1; i-- {
			f := float64(i) / float64(max(n-1, 1))
			w := float32(lerp(float64(headW), float64(tailW), f))
			var clr color.Color
			if pass == 0 {
				clr = cOutline
				w += 3
			} else {
				clr = mixColor(cHead, cTail, f)
			}
			vector.StrokeLine(dst, float32(px[i]), float32(py[i]), float32(px[i-1]), float32(py[i-1]), w, clr, true)
			vector.DrawFilledCircle(dst, float32(px[i]), float32(py[i]), w/2, clr, true)
		}
	}
	hx, hy := float32(px[0]), float32(py[0])
	vector.DrawFilledCircle(dst, hx, hy, headW/2+1.5, cOutline, true)
	vector.DrawFilledCircle(dst, hx, hy, headW/2, cHead, true)
	dx, dy := float32(g.dir.DX), float32(g.dir.DY)
	pxx, pyy := -dy, dx
	er, pr := float32(cell)*0.17, float32(cell)*0.085
	for _, s := range []float32{-1, 1} {
		ex := hx + dx*cell*0.16 + pxx*s*cell*0.20
		ey := hy + dy*cell*0.16 + pyy*s*cell*0.20
		vector.DrawFilledCircle(dst, ex, ey, er, cWhite, true)
		vector.DrawFilledCircle(dst, ex+dx*cell*0.06, ey+dy*cell*0.06, pr, cOutline, true)
	}
	if g.frame/18%2 == 0 && !g.over {
		tx, ty := hx+dx*headW/2, hy+dy*headW/2
		ex, ey := tx+dx*cell*0.42, ty+dy*cell*0.42
		vector.StrokeLine(dst, tx, ty, ex, ey, 2.5, cTongue, true)
		vector.StrokeLine(dst, ex, ey, ex+dx*3+pxx*4, ey+dy*3+pyy*4, 2, cTongue, true)
		vector.StrokeLine(dst, ex, ey, ex+dx*3-pxx*4, ey+dy*3-pyy*4, 2, cTongue, true)
	}
}
func (g *Game) drawObstacles(dst *ebiten.Image) {
	for _, o := range g.obstacles {
		x := float32(o.X * cell)
		y := float32(o.Y * cell)
		vector.FillRect(dst, x+2, y+2, cell-4, cell-4, cPill, true)
		vector.StrokeRect(dst, x+1, y+1, cell-2, cell-2, 1, cPillEdge, true)
	}
}
func (g *Game) drawFood(dst *ebiten.Image) {
	cx := float32(g.food.X*cell + cell/2)
	cy := float32(g.food.Y*cell + cell/2 + 1)
	r := float32(cell) * 0.34
	pulse := 1 + 0.06*float32(math.Sin(float64(g.frame)*0.08))
	vector.DrawFilledCircle(dst, cx, cy, r*pulse+1.5, cAppleDk, true)
	vector.DrawFilledCircle(dst, cx, cy, r*pulse, cApple, true)
	vector.DrawFilledCircle(dst, cx-r*0.35, cy-r*0.35, r*0.28, color.RGBA{0xff, 0xff, 0xff, 0x99}, true)
	vector.FillRect(dst, cx-1.5, cy-r-5, 3, 6, cStem, true)
	vector.DrawFilledCircle(dst, cx+6, cy-r-3, 4, cLeaf, true)
	if g.bonus != nil {
		bx := float32(g.bonus.pos.X*cell + cell/2)
		by := float32(g.bonus.pos.Y*cell + cell/2)
		br := float32(cell) * 0.30 * (1 + 0.12*float32(math.Sin(float64(g.frame)*0.2)))
		vector.DrawFilledCircle(dst, bx, by, br+2, color.RGBA{0x78, 0x35, 0x06, 0xff}, true)
		vector.DrawFilledCircle(dst, bx, by, br, cGold, true)
		vector.DrawFilledCircle(dst, bx-br*0.3, by-br*0.3, br*0.3, color.RGBA{0xff, 0xff, 0xff, 0xcc}, true)
		vector.FillRect(dst, bx-1, by-br-5, 2, (br+5)*2, cWhite, true)
		vector.FillRect(dst, bx-br-5, by-1, (br+5)*2, 2, cWhite, true)
		frac := float32(g.bonus.life) / float32(g.bonus.max)
		vector.FillRect(dst, bx-cell/2, by-cell/2-7, cell, 4, cPillEdge, false)
		vector.FillRect(dst, bx-cell/2, by-cell/2-7, cell*frac, 4, cGold, false)
	}
}
func (g *Game) drawParticles(screen *ebiten.Image) {
	for _, p := range g.parts {
		a := uint8(255 * float64(p.life) / float64(p.max))
		var clr color.Color = p.clr
		if c, ok := p.clr.(color.RGBA); ok {
			c.A = a
			clr = c
		}
		vector.DrawFilledCircle(screen,
			float32(p.x)+float32(margin), float32(p.y)+float32(hudH+margin),
			p.size*float32(p.life)/float32(p.max)+0.5, clr, true)
	}
}
func (g *Game) drawFloats(screen *ebiten.Image) {
	for _, f := range g.floats {
		a := uint8(255 * float64(f.life) / float64(f.max))
		var clr color.Color = f.clr
		if c, ok := f.clr.(color.RGBA); ok {
			c.A = a
			clr = c
		}
		b := text.BoundString(f.face, f.txt)
		text.Draw(screen, f.txt, f.face,
			int(f.x)+margin-b.Dx()/2, int(f.y)+hudH+margin, clr)
	}
}
func (g *Game) dim(screen *ebiten.Image) {
	vector.FillRect(screen, 0, 0, winW, winH, color.RGBA{0x02, 0x06, 0x12, 0xb8}, false)
}
func (g *Game) drawStart(screen *ebiten.Image) {
	g.dim(screen)
	pw, ph := float32(440), float32(250)
	px, py := float32(winW)/2-pw/2, float32(winH)/2-ph/2
	roundRect(screen, px, py, pw, ph, 20, cPanel)
	vector.StrokeRect(screen, px+1, py+1, pw-2, ph-2, 1.5, cAccent, true)
	drawCentered(screen, "Ular-ularan", fBig, winW/2, int(py)+72, cHead)
	drawCentered(screen, "papan 25 × 25", fSmall, winW/2, int(py)+100, cDim)
	drawCentered(screen, "Makan apel +10  •  Bintang emas +50", fHUD, winW/2, int(py)+140, cGold)
	blink := g.frame/30%2 == 0
	if blink {
		drawCentered(screen, "Tekan WASD / Panah / Enter untuk mulai", fHUD, winW/2, int(py)+185, cText)
	}
	mode := "Tembok"
	if g.wrap {
		mode = "Wrap"
	}
	drawCentered(screen, "Mode: "+mode+" (M)", fSmall, winW/2, int(py)+205, cDim)
	drawCentered(screen, "P pause   •   R restart   •   Esc keluar", fSmall, winW/2, int(py)+225, cDim)
}
func (g *Game) drawGameOver(screen *ebiten.Image) {
	g.dim(screen)
	pw, ph := float32(440), float32(260)
	px, py := float32(winW)/2-pw/2, float32(winH)/2-ph/2
	roundRect(screen, px, py, pw, ph, 20, cPanel)
	vector.StrokeRect(screen, px+1, py+1, pw-2, ph-2, 1.5, cApple, true)
	drawCentered(screen, "GAME OVER", fBig, winW/2, int(py)+74, cApple)
	drawCentered(screen, fmt.Sprintf("Skor %d   •   Terbaik %d   •   Panjang %d", g.score, g.best, len(g.snake)),
		fHUD, winW/2, int(py)+116, cText)
	if g.score >= g.best && g.score > 0 {
		drawCentered(screen, "★ Rekor baru! ★", fHUD, winW/2, int(py)+148, cGold)
	}
	if g.frame/30%2 == 0 {
		drawCentered(screen, "R / Enter = main lagi     Esc = keluar", fHUD, winW/2, int(py)+200, cText)
	}
}
func (g *Game) drawPaused(screen *ebiten.Image) {
	g.dim(screen)
	pw, ph := float32(300), float32(110)
	px, py := float32(winW)/2-pw/2, float32(winH)/2-ph/2
	roundRect(screen, px, py, pw, ph, 18, cPanel)
	vector.StrokeRect(screen, px+1, py+1, pw-2, ph-2, 1.5, cAccent, true)
	drawCentered(screen, "PAUSED", fTitle, winW/2, int(py)+52, cText)
	drawCentered(screen, "tekan P untuk lanjut", fSmall, winW/2, int(py)+78, cDim)
}
func (g *Game) Layout(outsideW, outsideH int) (int, int) { return winW, winH }
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func (g *Game) playBeep(freq float64, durMs int, vol float64) {
	if g.audioCtx == nil {
		return
	}
	const sampleRate = 44100
	n := sampleRate * durMs / 1000
	if n <= 0 {
		return
	}
	buf := make([]byte, n*4)
	for i := 0; i < n; i++ {
		t := float64(i) / sampleRate
		env := 1 - float64(i)/float64(n)
		val := math.Sin(2*math.Pi*freq*t) * 0.45 * env
		s := int16(val * 32767)
		buf[i*4] = byte(s)
		buf[i*4+1] = byte(s >> 8)
		buf[i*4+2] = byte(s)
		buf[i*4+3] = byte(s >> 8)
	}
	p := g.audioCtx.NewPlayerFromBytes(buf)
	p.SetVolume(vol)
	p.Play()
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
	loadFonts()
	ebiten.SetWindowSize(winW, winH)
	ebiten.SetWindowTitle("Ular-ularan")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	if err := ebiten.RunGame(NewGame()); err != nil && err != ebiten.Termination {
		log.Fatal(err)
	}
}
