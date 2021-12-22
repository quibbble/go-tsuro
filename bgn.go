package go_tsuro

import (
	"fmt"
	"strconv"

	bg "github.com/quibbble/go-boardgame"
	"github.com/quibbble/go-boardgame/pkg/bgerr"
)

var (
	actionToNotation = map[string]string{ActionPlaceTile: "p", bg.ActionSetWinners: "w"}
	notationToAction = reverseMap(actionToNotation)
)

func (p *PlaceTileActionDetails) encodeBGN() []string {
	return []string{strconv.Itoa(p.Row), strconv.Itoa(p.Column), strconv.Itoa(indexOf(tiles, p.Tile))}
}

func decodePlaceTileActionDetailsBGN(notation []string) (*PlaceTileActionDetails, error) {
	if len(notation) != 3 {
		return nil, loadFailure(fmt.Errorf("invalid place tile notation"))
	}
	row, err := strconv.Atoi(notation[0])
	if err != nil {
		return nil, loadFailure(err)
	}
	column, err := strconv.Atoi(notation[1])
	if err != nil {
		return nil, loadFailure(err)
	}
	tileIndex, err := strconv.Atoi(notation[2])
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
		Status: bgerr.StatusBGNDecodingFailure,
	}
}
