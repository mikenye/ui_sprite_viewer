package jsreader

import (
	"bufio"
	"os"
	"regexp"
	"strings"

	"github.com/rs/zerolog/log"
)

type SpriteDefinitions map[string]SpriteDefinition
type SpriteDefinition struct {
	id          string
	w           string
	h           string
	strokeScale string
	noRotate    string
	noAspect    string
	viewBox     string
	transform   string
	accentMult  string
	size        string
}

var (
	reComment    = regexp.MustCompile(`^\s*\/\/.*$`)
	reClosure    = regexp.MustCompile(`^\s*\},?\s*$`)
	reWhitespace = regexp.MustCompile(`^\s*$`)

	reConstSpriteDefinitions = regexp.MustCompile(`^\s*const\s+spriteDefinitions\s+=\s+{\s*$`)

	reAircraftName = regexp.MustCompile(`^\s*"[\w-]+":\s+\{\s*$`)
	reAircraftId   = regexp.MustCompile(`^\s*id:\s+\d+,?\s*$`)
	reWidth        = regexp.MustCompile(`^\s*w:\s+[\d\.]+,?\s*$`)
	reHeight       = regexp.MustCompile(`^\s*h:\s+[\d\.]+,?\s*$`)
	reStrokeScale  = regexp.MustCompile(`^\s*strokeScale:\s+[\d\.]+,?\s*$`)
	reNoRotate     = regexp.MustCompile(`^\s*noRotate:\s+(true|false),?\s*$`)
	reViewBox      = regexp.MustCompile(`^\s*viewBox:\s+"([-0-9\.]+\s+){3}[-0-9\.]+"\s*,?\s*$`)
	reTransform    = regexp.MustCompile(`^\s*transform:\s+".*?"\s*,?\s*$`)
	reNoAspect     = regexp.MustCompile(`^\s*noAspect:\s+(true|false),?\s*$`)
	reAccentMult   = regexp.MustCompile(`^\s*accentMult:\s+[\d\.]+,?\s*$`)
	reSize         = regexp.MustCompile(`^\s*size:\s+\[[-0-9\.]+,\s*[-0-9\.]+\]\s*,?\s*$`)
)

func LoadSpriteDefinitions(filename string) (SpriteDefinitions, error) {
	f, err := os.Open(filename)
	if err != nil {
		return SpriteDefinitions{}, err
	}
	defer f.Close()

	// scan through file
	var (
		state        int
		spriteDefs   SpriteDefinitions
		aircraftName string
		ae           SpriteDefinition
	)
	const (
		stateLookingForConst = iota
		stateGetAircraftName
		stateGetAircraftInfo
	)
	spriteDefs = make(SpriteDefinitions)
	s := bufio.NewScanner(f) //scan the contents of a file and print line by line
	state = stateLookingForConst
	for s.Scan() {

		// ignore comments
		if reComment.Match(s.Bytes()) {
			continue
		}

		// else interpret
		switch state {

		// look for `const spriteDefinitions = {`
		case stateLookingForConst:
			if reConstSpriteDefinitions.Match(s.Bytes()) {
				log.Debug().Msg(`found: "const spriteDefinitions = {"`)
				state = stateGetAircraftName
				continue
			}

		// get aircraft sprite name
		case stateGetAircraftName:
			if reAircraftName.Match(s.Bytes()) {
				aircraftName = strings.Split(s.Text(), "\"")[1]
				state = stateGetAircraftInfo
				continue
			}
			if reClosure.Match(s.Bytes()) {
				state = stateLookingForConst
			}

		// get aircraft id
		case stateGetAircraftInfo:
			switch {
			case reAircraftId.Match(s.Bytes()):
				ae.id = strings.Fields(strings.Split(s.Text(), ",")[0])[1]

				continue
			case reWidth.Match(s.Bytes()):
				ae.w = strings.Fields(strings.Split(s.Text(), ",")[0])[1]
				continue

			case reHeight.Match(s.Bytes()):
				ae.h = strings.Fields(strings.Split(s.Text(), ",")[0])[1]
				continue

			case reStrokeScale.Match(s.Bytes()):
				ae.strokeScale = strings.Fields(strings.Split(s.Text(), ",")[0])[1]
				continue

			case reNoRotate.Match(s.Bytes()):
				ae.noRotate = strings.Fields(strings.Split(s.Text(), ",")[0])[1]
				continue

			case reNoAspect.Match(s.Bytes()):
				ae.noAspect = strings.Fields(strings.Split(s.Text(), ",")[0])[1]
				continue

			case reViewBox.Match(s.Bytes()):
				ae.viewBox = strings.Split(strings.Split(s.Text(), ",")[0], "\"")[1]
				continue

			case reTransform.Match(s.Bytes()):
				ae.transform = strings.Split(s.Text(), "\"")[1]
				continue

			case reAccentMult.Match(s.Bytes()):
				ae.accentMult = strings.Fields(strings.Split(s.Text(), ",")[0])[1]
				continue

			case reSize.Match(s.Bytes()):
				x := strings.Split(s.Text(), ",")
				ae.size = strings.TrimSpace(strings.Split(strings.Join(x[0:], ","), ":")[1])
				continue

			case reClosure.Match(s.Bytes()):
				spriteDefs[aircraftName] = ae
				log.Debug().
					Str("name", aircraftName).
					Str("id", ae.id).
					Str("w", ae.w).
					Str("h", ae.h).
					Str("strokeScale", ae.strokeScale).
					Str("noRotate", ae.noRotate).
					Str("viewBox", ae.viewBox).
					Str("transform", ae.transform).
					Str("accentMult", ae.accentMult).
					Str("size", ae.size).
					Msg("aircraft found")
				state = stateGetAircraftName
				ae = SpriteDefinition{}
				continue

			case reWhitespace.Match(s.Bytes()):
				continue

			default:
				log.Panic().Str("unknown", s.Text()).Msg("unknown")
			}

		}
	}
	return spriteDefs, nil
}
