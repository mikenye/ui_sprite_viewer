package spritereader

import (
	"image"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/rs/zerolog/log"
)

type Sprites map[int]*ebiten.Image

func LoadSprites(filename string, spriteWidth, spriteHeight int) (*Sprites, error) {

	sprites := make(Sprites)

	f, err := os.Open(filename)
	if err != nil {
		return &sprites, err
	}

	log.Debug().Msg("reading spritesheet")
	img, _, err := image.Decode(f)
	spriteSheet := ebiten.NewImageFromImage(img)
	id := 0
	for sy := 0; sy <= spriteSheet.Bounds().Dy()-spriteHeight; sy += spriteHeight {
		for sx := 0; sx <= spriteSheet.Bounds().Dx()-spriteWidth; sx += spriteWidth {
			subImg := spriteSheet.SubImage(image.Rect(sx, sy, sx+spriteWidth, sy+spriteHeight))
			sprites[id] = ebiten.NewImageFromImage(subImg)
			log.Debug().Int("id", id).Ints("bounds", []int{sx, sy, sx + spriteWidth, sy + spriteHeight}).Msg("read sprite")
			id++
		}
	}

	return &sprites, nil
}
