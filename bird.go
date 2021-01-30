package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type bird struct {
	time         int
	y, speed     float64
	birdTextures []*sdl.Texture
}

const (
	gravity   = 0.25
	jumpSpeed = 5
)

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
	return &bird{birdTextures: birdTextures, y: 300}, nil

}
func (b *bird) paint(r *sdl.Renderer) error {
	b.time++
	b.y -= b.speed
	if b.y < 0 {
		b.y = 0
		b.speed = -b.speed
	}
	b.speed += gravity

	rect := &sdl.Rect{X: 10, Y: (600 - int32(b.y)) - 43/2, W: 50, H: 43}
	i := b.time / 10 % len(b.birdTextures)
	if err := r.Copy(b.birdTextures[i], nil, rect); err != nil {
		return fmt.Errorf("Could not copy bird image texture: %v", err)

	}
	return nil

}
func (b *bird) jump() {
	b.speed = -jumpSpeed
}
func (b *bird) destroy() {
	for _, birdTexture := range b.birdTextures {
		birdTexture.Destroy()
	}
}
