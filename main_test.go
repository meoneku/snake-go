package main

import "testing"

func TestMove(t *testing.T) {
	g := NewGame(30, 20)
	g.Obstacles = nil
	head := g.Snake[0]
	g.NextDir = Right
	g.Dir = Right
	g.Food = Point{1, 1}
	n := len(g.Snake)
	g.Update()
	if g.Snake[0] != (Point{head.X + 1, head.Y}) {
		t.Fatalf("kepala tidak maju: got %+v", g.Snake[0])
	}
	if len(g.Snake) != n {
		t.Fatalf("panjang berubah tanpa makan: got %d want %d", len(g.Snake), n)
	}
}
func TestEat(t *testing.T) {
	g := NewGame(30, 20)
	g.Obstacles = nil
	g.Dir = Right
	g.NextDir = Right
	head := g.Snake[0]
	g.Food = Point{head.X + 1, head.Y}
	n := len(g.Snake)
	g.Update()
	if len(g.Snake) != n+1 {
		t.Fatalf("tidak tumbuh saat makan: got %d want %d", len(g.Snake), n+1)
	}
	if g.Score != 10 {
		t.Fatalf("skor salah: got %d", g.Score)
	}
}
func TestWallCollision(t *testing.T) {
	g := NewGame(10, 10)
	g.Snake = []Point{{1, 1}}
	g.Dir = Left
	g.NextDir = Left
	g.Food = Point{5, 5}
	g.Update()
	if !g.Over {
		t.Fatal("seharusnya game over saat menabrak dinding")
	}
}
func TestNoReverse(t *testing.T) {
	g := NewGame(30, 20)
	g.Dir = Right
	g.NextDir = Right
	g.setDir(Left)
	if g.NextDir != Right {
		t.Fatal("putar balik 180 derajat seharusnya ditolak")
	}
}
