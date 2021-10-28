package go_tsuro

import (
	bg "github.com/quibbble/go-boardgame"
	"time"
)

const key = "tsuro"

type Builder struct{}

func (b *Builder) Create(options bg.BoardGameOptions, seed ...int64) (bg.BoardGame, error) {
	if len(seed) > 0 {
		return NewTsuro(options, seed[0])
	}
	return NewTsuro(options, time.Now().Unix())
}

func (b *Builder) Key() string {
	return key
}
