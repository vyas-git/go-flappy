package main

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/veandco/go-sdl2/ttf"

	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(2)
	}
}
func run() error {

	// Init SDL
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		return fmt.Errorf("could not initialize SDL: %v", err)
	}
	defer sdl.Quit()

	// Init TTF fonts
	if err := ttf.Init(); err != nil {
		return fmt.Errorf("could not initialize TTF: %v", err)
	}
	defer ttf.Quit()

	// Create window and renderer
	w, r, err := sdl.CreateWindowAndRenderer(800, 600, sdl.WINDOW_SHOWN)
	if err != nil {
		return fmt.Errorf("could not create window: %v", err)
	}
	defer w.Destroy()

	sdl.PumpEvents()
	// Draw Title Flappy Gopher
	if err := drawTitle(r); err != nil {
		return fmt.Errorf("Could not draw title: %v", err)
	}
	time.Sleep(1 * time.Second)

	// Draw scene
	scene, err := newScene(r)
	if err != nil {
		return fmt.Errorf("Could not draw image: %v", err)
	}

	defer scene.destroy()
	events := make(chan sdl.Event)
	errc := scene.run(events, r)
	runtime.LockOSThread()
	for {
		select {
		case events <- sdl.WaitEvent():

		case err := <-errc:
			return err
		}
	}

}

func drawTitle(r *sdl.Renderer) error {
	r.Clear()

	// Open font file
	fontFile, err := ttf.OpenFont("res/fonts/Flappy.ttf", 50)
	if err != nil {
		return fmt.Errorf("Could not open font: %v", err)
	}
	defer fontFile.Close()

	// Render font
	color := sdl.Color{R: 255, G: 100, B: 0, A: 255}
	surface, err := fontFile.RenderUTF8Solid("Flappy Gopher", color)
	if err != nil {
		return fmt.Errorf("Could not render font: %v", err)
	}
	defer surface.Free()

	// create texture from surface
	texture, err := r.CreateTextureFromSurface(surface)
	if err != nil {
		return fmt.Errorf("Could not create texture: %v", err)
	}
	defer texture.Destroy()

	// copy texture to render
	if err := r.Copy(texture, nil, nil); err != nil {
		return fmt.Errorf("Could not copy texture: %v", err)
	}

	r.Present()
	return nil
}
