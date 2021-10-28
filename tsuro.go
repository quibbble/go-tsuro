package go_tsuro

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	bg "github.com/quibbble/go-boardgame"
	"github.com/quibbble/go-boardgame/pkg/bgerr"
	"math/rand"
	"time"
)

const (
	minTeams = 2
	maxTeams = 8
)

type Tsuro struct {
	state   *state
	actions []*bg.BoardGameAction
	seed    int64
}

func NewTsuro(options bg.BoardGameOptions, seed int64) (*Tsuro, error) {
	if len(options.Teams) < minTeams {
		return nil, &bgerr.Error{
			Err:    fmt.Errorf("at least %d teams required to create a game of %s", minTeams, key),
			Status: bgerr.StatusTooFewTeams,
		}
	} else if len(options.Teams) > maxTeams {
		return nil, &bgerr.Error{
			Err:    fmt.Errorf("at most %d teams allowed to create a game of %s", maxTeams, key),
			Status: bgerr.StatusTooManyTeams,
		}
	}
	return &Tsuro{
		state:   newState(options.Teams, rand.New(rand.NewSource(seed))),
		actions: make([]*bg.BoardGameAction, 0),
		seed:    seed,
	}, nil
}

func (t *Tsuro) Do(action bg.BoardGameAction) error {
	switch action.ActionType {
	case ActionRotateTileRight:
		var details RotateTileActionDetails
		if err := mapstructure.Decode(action.MoreDetails, &details); err != nil {
			return &bgerr.Error{
				Err:    err,
				Status: bgerr.StatusInvalidActionDetails,
			}
		}
		if err := t.state.RotateTileRight(action.Team, details.Tile); err != nil {
			return err
		}
	case ActionRotateTileLeft:
		var details RotateTileActionDetails
		if err := mapstructure.Decode(action.MoreDetails, &details); err != nil {
			return &bgerr.Error{
				Err:    err,
				Status: bgerr.StatusInvalidActionDetails,
			}
		}
		if err := t.state.RotateTileLeft(action.Team, details.Tile); err != nil {
			return err
		}
	case ActionPlaceTile:
		var details PlaceTileActionDetails
		if err := mapstructure.Decode(action.MoreDetails, &details); err != nil {
			return &bgerr.Error{
				Err:    err,
				Status: bgerr.StatusInvalidActionDetails,
			}
		}
		if err := t.state.PlaceTile(action.Team, details.Tile, details.Row, details.Column); err != nil {
			return err
		}
		t.actions = append(t.actions, &action)
	case bg.ActionReset:
		seed := time.Now().UnixNano()
		t.state = newState(t.state.teams, rand.New(rand.NewSource(seed)))
		t.actions = make([]*bg.BoardGameAction, 0)
		t.seed = seed
	case bg.ActionUndo:
		if len(t.actions) > 0 {
			undo, _ := NewTsuro(bg.BoardGameOptions{Teams: t.state.teams}, t.seed)
			for _, a := range t.actions[:len(t.actions)-1] {
				if err := undo.Do(*a); err != nil {
					return err
				}
			}
			t.state = undo.state
			t.actions = undo.actions
		} else {
			return &bgerr.Error{
				Err:    fmt.Errorf("no actions to undo"),
				Status: bgerr.StatusInvalidAction,
			}
		}
	default:
		return &bgerr.Error{
			Err:    fmt.Errorf("cannot process action type %s", action.ActionType),
			Status: bgerr.StatusUnknownActionType,
		}
	}
	return nil
}

func (t *Tsuro) GetSnapshot(team ...string) (*bg.BoardGameSnapshot, error) {
	if len(team) > 1 {
		return nil, &bgerr.Error{
			Err:    fmt.Errorf("get snapshot requires zero or one team"),
			Status: bgerr.StatusTooManyTeams,
		}
	}
	hands := make(map[string][]*tile)
	for t, hand := range t.state.hands {
		if len(team) == 0 {
			hands[t] = hand.hand
		} else {
			if team[0] == t {
				hands[t] = hand.hand
			}
		}
	}
	details := TsuroSnapshotDetails{
		Board:          t.state.board.board,
		TilesRemaining: len(t.state.deck.deck),
		Hands:          hands,
		Tokens:         t.state.tokens,
		Dragon:         t.state.dragon,
	}
	return &bg.BoardGameSnapshot{
		Turn:     t.state.turn,
		Teams:    t.state.teams,
		Winners:  t.state.winners,
		MoreData: details,
		Actions:  t.actions,
	}, nil
}

func (t *Tsuro) GetSeed() int64 {
	return t.seed
}
