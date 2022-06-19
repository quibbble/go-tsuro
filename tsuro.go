package go_tsuro

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/mitchellh/mapstructure"
	bg "github.com/quibbble/go-boardgame"
	"github.com/quibbble/go-boardgame/pkg/bgerr"
	"github.com/quibbble/go-boardgame/pkg/bgn"
)

const (
	minTeams = 2
	maxTeams = 8
)

type Tsuro struct {
	state   *state
	actions []*bg.BoardGameAction
	options *TsuroMoreOptions
}

func NewTsuro(options *bg.BoardGameOptions) (*Tsuro, error) {
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
	var details TsuroMoreOptions
	if err := mapstructure.Decode(options.MoreOptions, &details); err != nil {
		return nil, &bgerr.Error{
			Err:    err,
			Status: bgerr.StatusInvalidOption,
		}
	}
	if details.Variant == "" {
		details.Variant = VariantClassic
	} else if !contains(Variants, details.Variant) {
		return nil, &bgerr.Error{
			Err:    fmt.Errorf("invalid Tsuro variant"),
			Status: bgerr.StatusInvalidOption,
		}
	}
	state, err := newState(options.Teams, rand.New(rand.NewSource(details.Seed)), details.Variant)
	if err != nil {
		return nil, &bgerr.Error{
			Err:    err,
			Status: bgerr.StatusInvalidOption,
		}
	}
	return &Tsuro{
		state:   state,
		actions: make([]*bg.BoardGameAction, 0),
		options: &details,
	}, nil
}

func (t *Tsuro) Do(action *bg.BoardGameAction) error {
	if len(t.state.winners) > 0 {
		return &bgerr.Error{
			Err:    fmt.Errorf("game already over"),
			Status: bgerr.StatusGameOver,
		}
	}
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
		t.actions = append(t.actions, action)
	case bg.ActionSetWinners:
		var details bg.SetWinnersActionDetails
		if err := mapstructure.Decode(action.MoreDetails, &details); err != nil {
			return &bgerr.Error{
				Err:    err,
				Status: bgerr.StatusInvalidActionDetails,
			}
		}
		if err := t.state.SetWinners(details.Winners); err != nil {
			return err
		}
		t.actions = append(t.actions, action)
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
	details := TsuroSnapshotData{
		Board:          t.state.board.board,
		TilesRemaining: len(t.state.deck.deck),
		Hands:          hands,
		Tokens:         t.state.tokens,
		Dragon:         t.state.dragon,
		Variant:        t.state.variant,
		Points:         t.state.points,
	}
	var targets []*bg.BoardGameAction
	if len(t.state.winners) == 0 {
		targets = t.state.targets(team...)
	}
	return &bg.BoardGameSnapshot{
		Turn:     t.state.turn,
		Teams:    t.state.teams,
		Winners:  t.state.winners,
		MoreData: details,
		Targets:  targets,
		Actions:  t.actions,
		Message:  t.state.message(),
	}, nil
}

func (t *Tsuro) GetBGN() *bgn.Game {
	tags := map[string]string{
		"Game":  key,
		"Teams": strings.Join(t.state.teams, ", "),
		"Seed":  fmt.Sprintf("%d", t.options.Seed),
	}
	actions := make([]bgn.Action, 0)
	for _, action := range t.actions {
		bgnAction := bgn.Action{
			TeamIndex: indexOf(t.state.teams, action.Team),
			ActionKey: rune(actionToNotation[action.ActionType][0]),
		}
		switch action.ActionType {
		case ActionPlaceTile:
			var details PlaceTileActionDetails
			_ = mapstructure.Decode(action.MoreDetails, &details)
			bgnAction.Details = details.encodeBGN()
		case bg.ActionSetWinners:
			var details bg.SetWinnersActionDetails
			_ = mapstructure.Decode(action.MoreDetails, &details)
			bgnAction.Details, _ = details.EncodeBGN(t.state.teams)
		}
		actions = append(actions, bgnAction)
	}
	return &bgn.Game{
		Tags:    tags,
		Actions: actions,
	}
}
