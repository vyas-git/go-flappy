package main

import (
	"fmt"
	"log"
	"time"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type scene struct {
	bg    *sdl.Texture
	bird  *bird
	pipes *pipes
}

// NewScene : every scene creates from here
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

	// load pipes
	pipes, err := newPipes(r)
	if err != nil {
		return nil, err
	}
	return &scene{bg: bgTexture, bird: bird, pipes: pipes}, nil
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
				s.update()
				if s.bird.isDead() {
					drawTitle(r, "Game Over")
					time.Sleep(time.Second)
					s.restart()
				}
				if err := s.paint(r); err != nil {
					errc <- err

				}

			}

		}

	}()
	return errc
}

// Update : Will update bird animation , pipe touches bird
func (s *scene) update() {
	s.bird.update()
	s.pipes.update()
	s.pipes.touch(s.bird)
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
	if err := s.pipes.paint(r); err != nil {
		return err
	}
	r.Present()
	return nil
}

// Restart: will restat  when game over or bird dies
func (s *scene) restart() {
	s.bird.restart()
	s.pipes.restart()
}

// HandleEvent: will handle quit click event ,jump keyboard space event
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

// Destroy : will destroy background , bird , pipes textures
func (s *scene) destroy() {
	s.bg.Destroy()
	s.bird.destroy()
	s.pipes.destroy()
}
