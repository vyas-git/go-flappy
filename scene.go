package main

import (
	"fmt"
	"log"
	"time"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type scene struct {
	bg   *sdl.Texture
	bird *bird
}

func newScene(r *sdl.Renderer) (*scene, error) {
	// Load background image texture
	bgTexture, err := img.LoadTexture(r, "res/imgs/background.png")
	if err != nil {
		return nil, fmt.Errorf("Could not load image texture: %v", err)
	}
	// load birds
	bird, err := newBird(r)
	if err != nil {
		return nil, err
	}

	return &scene{bg: bgTexture, bird: bird}, nil
}

// Run scene for every 0.01 seconds
func (s *scene) run(events <-chan sdl.Event, r *sdl.Renderer) <-chan error {
	errc := make(chan error)
	go func() {
		defer close(errc)
		tick := time.Tick(10 * time.Millisecond)
		done := false
		for !done {
			//	fmt.Println(<-ctx.Done())
			select {
			case e := <-events:
				if done = s.handleEvent(e); done {
					return
				}
				//time.Sleep(5 * time.Millisecond)
			case <-tick:

				if err := s.paint(r); err != nil {
					errc <- err

				}
			}

		}

	}()
	return errc
}

// Paint : Will copy each texture to renderer
func (s *scene) paint(r *sdl.Renderer) error {
	r.Clear()
	if err := r.Copy(s.bg, nil, nil); err != nil {
		return fmt.Errorf("Could not copy background image texture: %v", err)
	}
	if err := s.bird.paint(r); err != nil {
		return err
	}
	r.Present()
	return nil
}
func (s *scene) handleEvent(event sdl.Event) bool {
	switch e := event.(type) {
	case *sdl.QuitEvent:
		return true
	case *sdl.KeyboardEvent:
		s.bird.jump()
	case *sdl.WindowEvent:
	case *sdl.MouseMotionEvent:
	default:
		log.Printf("unkown event %T", e)
	}
	return false

}

func (s *scene) destroy() {
	s.bg.Destroy()
	s.bird.destroy()
}
