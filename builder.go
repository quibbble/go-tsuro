package go_tsuro

import (
	bg "github.com/quibbble/go-boardgame"
)

const key = "tsuro"

type Builder struct{}

func (b *Builder) Create(options bg.BoardGameOptions) (bg.BoardGame, error) {
	return NewTsuro(options)
}

func (b *Builder) Key() string {
	return key
}
