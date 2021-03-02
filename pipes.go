package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	img "github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type pipes struct {
	mu sync.RWMutex

	texture *sdl.Texture
	speed   int32

	pipes []*pipe
}

// NewPipes : Will load  pipe images and appends to new pipes list

func newPipes(r *sdl.Renderer) (*pipes, error) {
	// Load Pipe and append to list
	pipeTexture, err := img.LoadTexture(r, "res/imgs/pipe.png")
	if err != nil {
		return nil, fmt.Errorf("Could not load pipe texture: %v", err)
	}
	ps := &pipes{
		texture: pipeTexture,
		speed:   1,
	}
	go func() {
		for {
			ps.mu.Lock()
			ps.pipes = append(ps.pipes, newPipe())
			ps.mu.Unlock()
			time.Sleep(time.Second)

		}
	}()
	return ps, nil
}

// Touch : When pipe touch bird
func (ps *pipes) touch(b *bird) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	for _, p := range ps.pipes {
		p.touch(b)
	}
}

// Paint : Will call each pipe paint
func (ps *pipes) paint(r *sdl.Renderer) error {
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	for _, p := range ps.pipes {
		p.paint(r, ps.texture)
	}
	return nil
}

// Restart : will do all pipes nill
func (ps *pipes) restart() {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	ps.pipes = nil
}

// Update : Will update pipe speed
func (ps *pipes) update() {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	var rem []*pipe // remove pass over pipes from list
	for _, p := range ps.pipes {
		p.mu.Lock()
		p.x -= ps.speed
		p.mu.Unlock()
		if p.x+p.w > 0 {
			rem = append(rem, p)
		}

	}
	ps.pipes = rem
}

// Destroy : Will destroy all pipes textures
func (ps *pipes) destroy() {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	ps.texture.Destroy()

}

type pipe struct {
	mu       sync.RWMutex
	x        int32
	h        int32
	w        int32
	inverted bool
}

// NewPipe : Will return new pipe object
func newPipe() *pipe {

	return &pipe{
		x:        800,
		h:        100 + int32(rand.Intn(50)), // random height of pipe
		w:        50,
		inverted: rand.Float32() > 0.5, // random flip position of pipe
	}
}

// Paint : Will copy each pipe texture to render
func (p *pipe) paint(r *sdl.Renderer, texture *sdl.Texture) error {
	p.mu.RLock()
	defer p.mu.RUnlock()
	rect := &sdl.Rect{X: p.x, Y: 600 - p.h, W: p.w, H: p.h}
	flip := sdl.FLIP_NONE
	if p.inverted {
		rect.Y = 0
		flip = sdl.FLIP_VERTICAL
	}
	if err := r.CopyEx(texture, nil, rect, 0, nil, flip); err != nil {
		return fmt.Errorf("Could not copy bird image texture: %v", err)

	}
	return nil
}

// Touch : When bird touch pipe
func (p *pipe) touch(b *bird) {
	p.mu.Lock()
	defer p.mu.Unlock()
	b.touch(p)

}
