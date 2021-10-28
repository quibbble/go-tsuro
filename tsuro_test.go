package go_tsuro

import (
	bg "github.com/quibbble/go-boardgame"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const (
	TeamA = "TeamA"
	TeamB = "TeamB"
)

func Test_Tsuro(t *testing.T) {
	tsuro, err := NewTsuro(bg.BoardGameOptions{
		Teams: []string{TeamA, TeamB},
	}, time.Now().UnixNano())
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	tile := "ABCDEFG"
	rotated := "CDEFGHA"
	tsuro.state.hands[TeamA].hand[0] = newTile(tile)

	// rotate first tile in TeamA hand
	err = tsuro.Do(bg.BoardGameAction{
		Team:       TeamA,
		ActionType: ActionRotateTileRight,
		MoreDetails: RotateTileActionDetails{
			Tile: tile,
		},
	})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	assert.Equal(t, rotated, tsuro.state.hands[TeamA].hand[0].Edges)

	tsuro.state.turn = TeamB
	tsuro.state.tokens[TeamB].Row = 0
	tsuro.state.tokens[TeamB].Col = 0
	tsuro.state.tokens[TeamB].Notch = "A"

	// place the first tile in TeamB hand at 0,0
	err = tsuro.Do(bg.BoardGameAction{
		Team:       TeamB,
		ActionType: ActionPlaceTile,
		MoreDetails: PlaceTileActionDetails{
			Row:    0,
			Column: 0,
			Tile:   tsuro.state.hands[TeamB].hand[0].Edges,
		},
	})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}

func Test_Tsuro_Undo(t *testing.T) {
	tsuro, err := NewTsuro(bg.BoardGameOptions{
		Teams: []string{TeamA, TeamB},
	}, 123)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// place ADBGCHEF at 0,3 moving TeamB notch from 0,3,B to 1,3,G
	err = tsuro.Do(bg.BoardGameAction{
		Team:       TeamB,
		ActionType: ActionPlaceTile,
		MoreDetails: PlaceTileActionDetails{
			Row:    0,
			Column: 3,
			Tile:   tsuro.state.hands[TeamB].hand[0].Edges,
		},
	})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	assert.Equal(t, "ADBGCHEF", tsuro.state.board.board[0][3].Edges)
	assert.Equal(t, "G", tsuro.state.tokens[TeamB].Notch)

	// undo placement
	err = tsuro.Do(bg.BoardGameAction{
		ActionType: bg.ActionUndo,
	})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	assert.Equal(t, "B", tsuro.state.tokens[TeamB].Notch)
}
