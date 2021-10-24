package go_tsuro

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	bg "github.com/quibbble/go-boardgame"
	"github.com/quibbble/go-boardgame/pkg/bgerr"
)

const (
	minTeams = 2
	maxTeams = 8
)

type Tsuro struct {
	state   *state
	actions []*bg.BoardGameAction
}

func NewTsuro(options bg.BoardGameOptions) (*Tsuro, error) {
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
		state:   newState(options.Teams),
		actions: make([]*bg.BoardGameAction, 0),
	}, nil
}

func (t *Tsuro) Do(action bg.BoardGameAction) error {
	switch action.ActionType {
	case RotateTileRight:
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
	case RotateTileLeft:
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
	case PlaceTile:
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
	case Reset:
		t.state = newState(t.state.teams)
		t.actions = make([]*bg.BoardGameAction, 0)
	default:
		return &bgerr.Error{
			Err:    fmt.Errorf("cannot process action type %s", action.ActionType),
			Status: bgerr.StatusUnknownActionType,
		}
	}
	return nil
}

func (t *Tsuro) GetSnapshot(team string) (bg.BoardGameSnapshot, error) {
	if !contains(t.state.teams, team) {
		return bg.BoardGameSnapshot{}, &bgerr.Error{
			Err:    fmt.Errorf("%s not a valid team", team),
			Status: bgerr.StatusUnknownTeam,
		}
	}
	return bg.BoardGameSnapshot{
		Turn:    t.state.turn,
		Teams:   t.state.teams,
		Winners: t.state.winners,
		MoreData: TsuroSnapshotDetails{
			Board:          t.state.board.board,
			TilesRemaining: len(t.state.deck.deck),
			Hand:           t.state.hands[team].hand,
			Tokens:         t.state.tokens,
			Dragon:         t.state.dragon,
		},
		Actions: t.actions,
	}, nil
}
