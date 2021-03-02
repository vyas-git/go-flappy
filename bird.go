package main

import (
	"fmt"
	"sync"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type bird struct {
	mu           sync.RWMutex
	time         int
	x, y         int32
	w, h         int32
	speed        float64
	birdTextures []*sdl.Texture
	dead         bool
}

const (
	gravity   = 0.1
	jumpSpeed = 5
)

// NewBird : Will Load Bird Images with positions and width,height
func newBird(r *sdl.Renderer) (*bird, error) {
	// Load bird image textures
	var birdTextures []*sdl.Texture
	for i := 1; i <= 4; i++ {
		path := fmt.Sprintf("res/imgs/bird_frame_%d.png", i)
		birdTexture, err := img.LoadTexture(r, path)
		if err != nil {
			return nil, fmt.Errorf("Could not load image texture: %v", err)
		}
		birdTextures = append(birdTextures, birdTexture)
	}
	return &bird{birdTextures: birdTextures, x: 10, y: 300, w: 50, h: 43}, nil

}

// Update : Will update bird  speed
func (b *bird) update() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.time++
	b.y -= int32(b.speed)
	if b.y < 0 {
		b.y = 0
		b.speed = -b.speed
		b.dead = true
	}
	b.speed += gravity
}

// Dead : Return Bird Dead boolean
func (b *bird) isDead() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.dead
}

// Restart : Will restart bird position y,dead, speed
func (b *bird) restart() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.y = 300
	b.speed = 0
	b.dead = false

}

// Paint : Will copy bird textures to renders
func (b *bird) paint(r *sdl.Renderer) error {
	b.mu.RLock()
	defer b.mu.RUnlock()
	rect := &sdl.Rect{X: b.x, Y: 600 - b.y - b.h/2, W: b.w, H: b.h}
	i := b.time / 10 % len(b.birdTextures)
	if err := r.Copy(b.birdTextures[i], nil, rect); err != nil {
		return fmt.Errorf("Could not copy bird image texture: %v", err)

	}
	return nil

}

// Jump : Reduces bird speed by constant jump speed
//So that Bird position will goes to top automatically
func (b *bird) jump() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.speed = -jumpSpeed
}

// Touch : When bird touch pipe
func (b *bird) touch(p *pipe) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if p.x > b.x+b.w { // too far right
		return
	}

	if p.x+p.w < b.x { // too far left
		return
	}

	if !p.inverted && p.h < b.y-b.h/2 { // pipe is too low
		return
	}
	if p.inverted && 600-p.h > b.y+b.h/2 { // pipe is too high
		return
	}
	b.dead = true

}

// Destory : Will destroy bird texture
func (b *bird) destroy() {
	b.mu.Lock()
	defer b.mu.Unlock()
	for _, birdTexture := range b.birdTextures {
		birdTexture.Destroy()
	}
}
