package go_tsuro

import (
	"fmt"
	"github.com/quibbble/go-boardgame/pkg/bgerr"
	"strconv"
	"strings"
)

// Notation - "'number of teams':'seed':'MoreOptions':'team index','action type number','details','details';..."

var (
	notationActionToInt = map[string]int{ActionPlaceTile: 0}
	notationIntToAction = map[string]string{"0": ActionPlaceTile}
)

func (p *PlaceTileActionDetails) encode() string {
	return fmt.Sprintf("%d,%d,%d", p.Column, p.Row, indexOf(tiles, p.Tile))
}

func decodeNotationPlaceTileActionDetails(notation string) (*PlaceTileActionDetails, error) {
	split := strings.Split(notation, ",")
	row, err := strconv.Atoi(split[0])
	if err != nil {
		return nil, loadFailure(err)
	}
	column, err := strconv.Atoi(split[1])
	if err != nil {
		return nil, loadFailure(err)
	}
	tileIndex, err := strconv.Atoi(split[2])
	if err != nil {
		return nil, loadFailure(err)
	}
	if tileIndex < 0 || tileIndex >= len(tiles) {
		return nil, loadFailure(fmt.Errorf("got %d but want index between %d and %d when decoding %s action", tileIndex, 0, len(tiles), ActionPlaceTile))
	}
	return &PlaceTileActionDetails{
		Row:    row,
		Column: column,
		Tile:   tiles[tileIndex],
	}, nil
}

func loadFailure(err error) error {
	return &bgerr.Error{
		Err:    err,
		Status: bgerr.StatusGameLoadFailure,
	}
}
