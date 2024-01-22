package main

import (
	"fmt"
	_ "image/png"
	"os"
	"time"

	"ui_sprite_viewer/lib/jsreader"
	"ui_sprite_viewer/lib/spritereader"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

var (
	// App config, command line & env var configuration
	app = cli.App{
		Version:     "0.0.1",
		Name:        "plane.watch UI Sprite Viewer",
		Usage:       "Viewer of plane-watch/pw-ui sprites",
		Description: `An easy way to see the aircraft sprites`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "aircraft_sprite_js",
				Usage:    "Path to plane-watch/pw-ui/app/javascript/aircraft_sprite.js",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "sprites_png",
				Usage:    "Path to plane-watch/pw-ui/app/javascript/images/sprites.png",
				Required: true,
			},
			&cli.IntFlag{
				Name:        "sprite_w",
				Usage:       "Width of sprites in pixels",
				Value:       72,
				DefaultText: "72",
			},
			&cli.IntFlag{
				Name:        "sprite_h",
				Usage:       "Height of sprites in pixels",
				Value:       72,
				DefaultText: "72",
			},
		},
	}

	id int
)

type Game struct {
	Sprites *spritereader.Sprites
}

func (g *Game) Update() error {
	_, dy := ebiten.Wheel()
	if dy > 0 && dy < 1 {
		id += 1
	} else if dy < 0 && dy > -1 {
		id -= 1
	} else {
		id += int(dy / 5.0)
	}
	if id < 0 {
		id = 0
	}
	if id > len(*g.Sprites) {
		id = len(*g.Sprites) - 1
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := ebiten.DrawImageOptions{}
	x := (screen.Bounds().Dx() / 2) - ((*g.Sprites)[id].Bounds().Dx() / 2)
	y := (screen.Bounds().Dx() / 2) - ((*g.Sprites)[id].Bounds().Dx() / 2)
	op.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage((*g.Sprites)[id], &op)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("viewing sprite id: %d", id), 1, 1)
	ebitenutil.DebugPrintAt(screen, "mousewheel changes sprite", 1, 15)

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {

	// set up logging
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.UnixDate})
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	// set action when run
	app.Action = runApp

	// run & final exit
	err := app.Run(os.Args)
	if err != nil {
		log.Err(err).Msg("finished with error")
		os.Exit(1)
	} else {
		log.Info().Msg("finished without error")
		os.Exit(0)
	}
}

func runApp(cliContext *cli.Context) error {

	// open aircraft_sprite.js
	_, err := jsreader.LoadSpriteDefinitions(cliContext.String("aircraft_sprite_js"))
	if err != nil {
		return err
	}

	// open sprites.png
	sprites, err := spritereader.LoadSprites(cliContext.String("sprites_png"), cliContext.Int("sprite_w"), cliContext.Int("sprite_h"))
	if err != nil {
		return err
	}

	// prep view
	g := &Game{
		Sprites: sprites,
	}
	ebiten.SetWindowSize(300, 300)
	ebiten.SetWindowTitle("plane.watch sprite viewer")

	log.Info().Msg("starting \"game\" window")
	err = ebiten.RunGame(g)

	return nil
}
