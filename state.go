package go_tsuro

import (
	"fmt"
	"math/rand"
	"strings"

	bg "github.com/quibbble/go-boardgame"
	"github.com/quibbble/go-boardgame/pkg/bgerr"
)

type state struct {
	turn    string
	teams   []string
	winners []string
	board   *board
	deck    *deck
	tokens  map[string]*token
	hands   map[string]*hand
	dragon  string
	active  map[string]bool // teams that have placed and still alive
	alive   map[string]bool // teams that are alive
	variant string
	points  map[string]int
}

func newState(teams []string, random *rand.Rand, variant string) (*state, error) {
	hands := make(map[string]*hand)
	tokens := make(map[string]*token)
	alive := make(map[string]bool)
	deck := newDeck(random)
	var points map[string]int

	switch variant {
	case VariantClassic, VariantSolo:
		for _, team := range teams {
			hand := newHand()
			for i := 0; i < 3; i++ {
				tile, _ := deck.Draw()
				hand.Add(tile)
			}
			hands[team] = hand
			token := uniqueRandomToken(tokens, random)
			tokens[team] = token
			alive[team] = true
		}
	case VariantLongestPath, VariantMostLoops:
		points = make(map[string]int)
		for _, team := range teams {
			hand := newHand()
			for i := 0; i < 3; i++ {
				tile, _ := deck.Draw()
				hand.Add(tile)
			}
			hands[team] = hand
			token := uniqueRandomToken(tokens, random)
			tokens[team] = token
			alive[team] = true
			points[team] = 0
		}
	case VariantOpenTiles:
		hand := newHand()
		for i := 0; i < 3; i++ {
			tile, _ := deck.Draw()
			hand.Add(tile)
		}
		for _, team := range teams {
			hands[team] = hand
			token := uniqueRandomToken(tokens, random)
			tokens[team] = token
			alive[team] = true
		}
	}
	return &state{
		turn:    teams[0],
		teams:   teams,
		winners: make([]string, 0),
		board:   newBoard(),
		deck:    deck,
		tokens:  tokens,
		hands:   hands,
		dragon:  "",
		active:  make(map[string]bool),
		alive:   alive,
		variant: variant,
		points:  points,
	}, nil
}

func (s *state) RotateTileRight(team, tile string) error {
	if s.variant == VariantOpenTiles && team != s.turn {
		return &bgerr.Error{
			Err:    fmt.Errorf("%s cannot rotate tile on %s turn", team, s.turn),
			Status: bgerr.StatusWrongTurn,
		}
	}
	if !contains(s.teams, team) {
		return &bgerr.Error{
			Err:    fmt.Errorf("%s not a valid team", team),
			Status: bgerr.StatusUnknownTeam,
		}
	}
	t := newTile(tile)
	if !t.in(s.hands[team].hand) {
		return &bgerr.Error{
			Err:    fmt.Errorf("%s's hand does not contain %s", team, tile),
			Status: bgerr.StatusInvalidActionDetails,
		}
	}
	s.hands[team].hand[s.hands[team].IndexOf(t)].RotateRight()
	return nil
}

func (s *state) RotateTileLeft(team, tile string) error {
	if s.variant == VariantOpenTiles && team != s.turn {
		return &bgerr.Error{
			Err:    fmt.Errorf("%s cannot rotate tile on %s turn", team, s.turn),
			Status: bgerr.StatusWrongTurn,
		}
	}
	if !contains(s.teams, team) {
		return &bgerr.Error{
			Err:    fmt.Errorf("%s not a valid team", team),
			Status: bgerr.StatusUnknownTeam,
		}
	}
	t := newTile(tile)
	if !t.in(s.hands[team].hand) {
		return &bgerr.Error{
			Err:    fmt.Errorf("%s's hand does not contain %s", team, tile),
			Status: bgerr.StatusInvalidActionDetails,
		}
	}
	s.hands[team].hand[s.hands[team].IndexOf(t)].RotateLeft()
	return nil
}

func (s *state) PlaceTile(team, tile string, row, column int) error {
	if team != s.turn {
		return &bgerr.Error{
			Err:    fmt.Errorf("%s cannot play on %s turn", team, s.turn),
			Status: bgerr.StatusWrongTurn,
		}
	}
	if !s.active[s.turn] && (s.tokens[team].Row != row || s.tokens[team].Col != column) {
		return &bgerr.Error{
			Err:    fmt.Errorf("%s cannot place in row %d column %d", team, row, column),
			Status: bgerr.StatusInvalidAction,
		}
	} else if s.active[s.turn] {
		adj, err := s.tokens[team].getAdjacent()
		if err != nil {
			return err
		}
		if row != adj.Row || column != adj.Col {
			return &bgerr.Error{
				Err:    fmt.Errorf("%s cannot place in row %d column %d", team, row, column),
				Status: bgerr.StatusInvalidAction,
			}
		}
	}
	t := newTile(tile)
	if !t.in(s.hands[team].hand) {
		return &bgerr.Error{
			Err:    fmt.Errorf("%s's hand does not contain %s", team, tile),
			Status: bgerr.StatusInvalidAction,
		}
	}
	if err := s.hands[team].Remove(t); err != nil {
		return &bgerr.Error{
			Err:    err,
			Status: bgerr.StatusInvalidAction,
		}
	}
	if err := s.board.Place(t, row, column); err != nil {
		return &bgerr.Error{
			Err:    err,
			Status: bgerr.StatusInvalidAction,
		}
	}
	if !s.active[s.turn] {
		s.active[s.turn] = true
	}
	s.moveTokens()
	s.score()
	s.updateAlive()
	s.handleDraws()
	s.nextTurn()
	return nil
}

func (s *state) SetWinners(winners []string) error {
	for _, winner := range winners {
		if !contains(s.teams, winner) {
			return &bgerr.Error{
				Err:    fmt.Errorf("winner not in teams"),
				Status: bgerr.StatusInvalidActionDetails,
			}
		}
	}
	s.winners = winners
	return nil
}

func (s *state) moveTokens() {
	moved := 0
	move := map[string]string{"A": "F", "B": "E", "C": "H", "D": "G", "E": "B", "F": "A", "G": "D", "H": "C"}
	for team, token := range s.tokens {
		if s.active[team] {
			t := s.board.board[token.Row][token.Col]
			if !mapContainsVal(t.Paths, team) {
				// first placement so move through the just placed tile
				destination := t.GetDestination(token.Notch)
				t.Paths[token.Notch+destination] = team
				token.Notch = destination
				// token was moved
				moved++
			} else if s.collided(s.tokens, team, token) {
				// token collided with other token
				continue
			} else {
				// normal case
				var nextTile *tile
				if strings.Contains("AB", token.Notch) && token.Row-1 >= 0 && s.board.board[token.Row-1][token.Col] != nil {
					nextTile = s.board.board[token.Row-1][token.Col]
					token.Row -= 1
				} else if strings.Contains("CD", token.Notch) && token.Col+1 < columns && s.board.board[token.Row][token.Col+1] != nil {
					nextTile = s.board.board[token.Row][token.Col+1]
					token.Col += 1
				} else if strings.Contains("EF", token.Notch) && token.Row+1 < rows && s.board.board[token.Row+1][token.Col] != nil {
					nextTile = s.board.board[token.Row+1][token.Col]
					token.Row += 1
				} else if strings.Contains("GH", token.Notch) && token.Col-1 >= 0 && s.board.board[token.Row][token.Col-1] != nil {
					nextTile = s.board.board[token.Row][token.Col-1]
					token.Col -= 1
				} else {
					continue
				}
				// move the token to the notch on the next tile
				startNotch := move[token.Notch]
				// where the token ends up on the next tile
				endNotch := nextTile.GetDestination(startNotch)
				// update token location
				nextTile.Paths[startNotch+endNotch] = team
				token.Notch = endNotch
				// token was moved
				moved++
			}
		}
	}
	if moved > 0 {
		s.moveTokens()
	}
}

func (s *state) collided(tokens map[string]*token, team string, token *token) bool {
	for team2, token2 := range tokens {
		if team != team2 && (token.collided(token2) || token.equals(token2)) {
			return true
		}
	}
	return false
}

func (s *state) score() {
	switch s.variant {
	case VariantLongestPath:
		points := make(map[string]int)
		for _, team := range s.teams {
			points[team] = 0
		}
		for _, row := range s.board.board {
			for _, tile := range row {
				if tile != nil {
					for _, team := range tile.Paths {
						points[team]++
					}
				}
			}
		}
		s.points = points
	case VariantMostLoops:
		points := make(map[string]int)
		for _, team := range s.teams {
			points[team] = 0
		}
		for _, row := range s.board.board {
			for _, tile := range row {
				if tile != nil {
					// calculates the number of times path of the same team cross on a tile
					// this is equivalent to the number of loops
					pathScores := make(map[string]int)
					for path1, team1 := range tile.Paths {
						for path2, team2 := range tile.Paths {
							if path1 == path2 {
								continue
							}
							// check if same team and paths cross
							if team1 == team2 && contains(crossing[path1], path2) {
								pathScores[path1]++
							}
						}
					}
					// find the max number of crossings for each team
					tilePoints := make(map[string]int)
					for path, score := range pathScores {
						team := tile.Paths[path]
						if tilePoints[team] < score {
							tilePoints[team] = score
						}
					}
					for team, score := range tilePoints {
						points[team] += score
					}
				}
			}
		}
		s.points = points
	}
}

func (s *state) updateAlive() {
	if len(s.winners) > 0 {
		return
	}
	// alive before checking
	initialAlive := make([]string, 0)
	for _, team := range s.teams {
		if s.alive[team] {
			initialAlive = append(initialAlive, team)
		}
	}
	// update who is still alive
	for team, token := range s.tokens {
		if s.active[team] {
			if (token.Row == 0 && strings.Contains("AB", token.Notch)) ||
				(token.Row == rows-1 && strings.Contains("EF", token.Notch)) ||
				(token.Col == 0 && strings.Contains("GH", token.Notch)) ||
				(token.Col == columns-1 && strings.Contains("CD", token.Notch)) {
				// check on board edge
				s.setLost(team)
			} else if s.collided(s.tokens, team, token) {
				// check if collided with another token
				s.setLost(team)
			}
		}
	}
	// who is still alive
	stillAlive := make([]string, 0)
	for _, team := range s.teams {
		if s.alive[team] {
			stillAlive = append(stillAlive, team)
		}
	}
	switch s.variant {
	case VariantClassic, VariantOpenTiles:
		if len(stillAlive) == 0 { // no more alive so initial alive all win
			s.winners = initialAlive
		} else if len(stillAlive) == 1 { // one alive so they win
			s.winners = stillAlive
		} else if s.board.getTileCount() == len(tiles) { // all tiles have been placed remaining alive are winners
			s.winners = stillAlive
		}
	case VariantLongestPath, VariantMostLoops:
		if len(stillAlive) == 0 { // no more alive
			s.winners = max(s.points)
		} else if s.board.getTileCount() == len(tiles) { // all tiles have been placed
			s.winners = max(s.points)
		}
	case VariantSolo:
		if s.board.getTileCount() == len(tiles) { // win if all tokens are still on board and all tiles have been placed
			if len(stillAlive) == len(s.teams) {
				s.winners = []string{"true"}
			} else {
				s.winners = []string{"false"}
			}
		}
	}
}

func (s *state) handleDraws() {
	if len(s.winners) > 0 {
		return
	}
	current := s.turn
	if s.dragon != "" {
		current = s.dragon
	}
	for s.alive[current] && len(s.deck.deck) > 0 && len(s.hands[current].hand) < 3 {
		tile, err := s.deck.Draw()
		if err != nil {
			return
		}
		s.hands[current].Add(tile)
		current = s.getNextTurn(current)
	}
	if len(s.deck.deck) == 0 && len(s.hands[current].hand) < 3 {
		s.dragon = current
	} else {
		s.dragon = ""
	}
}

func (s *state) nextTurn() {
	if len(s.winners) > 0 {
		return
	}
	s.turn = s.getNextTurn(s.turn)
}

func (s *state) getNextTurn(turn string) string {
	nextTurn := ""
	if len(s.winners) > 0 {
		return nextTurn
	}
	for idx, team := range s.teams {
		if team == turn {
			nextTurn = s.teams[(idx+1)%len(s.teams)]
			if !s.alive[nextTurn] {
				return s.getNextTurn(nextTurn)
			}
			return nextTurn
		}
	}
	return nextTurn
}

func (s *state) setLost(team string) {
	s.alive[team] = false
	s.active[team] = false
	s.deck.Add(s.hands[team].hand...)
	s.hands[team].Clear()
	if s.aliveCount() <= 0 {
		return
	}
	next := s.getNextTurn(s.turn)
	if s.dragon == team && len(s.hands[next].hand) < 3 {
		s.dragon = next
	}
}

func (s *state) aliveCount() int {
	count := 0
	for _, alive := range s.alive {
		if alive {
			count++
		}
	}
	return count
}

func (s *state) targets(team ...string) []*bg.BoardGameAction {
	targets := make([]*bg.BoardGameAction, 0)
	// rotate tile actions
	if len(team) == 0 {
		for _, t := range s.teams {
			for _, tile := range s.hands[t].hand {
				targets = append(targets, &bg.BoardGameAction{
					Team:       s.turn,
					ActionType: ActionRotateTileLeft,
					MoreDetails: RotateTileActionDetails{
						Tile: tile.Edges,
					},
				}, &bg.BoardGameAction{
					Team:       s.turn,
					ActionType: ActionRotateTileRight,
					MoreDetails: RotateTileActionDetails{
						Tile: tile.Edges,
					},
				})
			}
		}
	} else if len(team) == 1 {
		for _, tile := range s.hands[team[0]].hand {
			targets = append(targets, &bg.BoardGameAction{
				Team:       s.turn,
				ActionType: ActionRotateTileLeft,
				MoreDetails: RotateTileActionDetails{
					Tile: tile.Edges,
				},
			}, &bg.BoardGameAction{
				Team:       s.turn,
				ActionType: ActionRotateTileRight,
				MoreDetails: RotateTileActionDetails{
					Tile: tile.Edges,
				},
			})
		}
	}
	// place tile actions
	if len(team) == 0 || (len(team) == 1 && team[0] == s.turn) {
		row := s.tokens[s.turn].Row
		col := s.tokens[s.turn].Col
		switch s.tokens[s.turn].Notch {
		case "A", "B":
			row++
		case "C", "D":
			col++
		case "E", "F":
			row--
		case "G", "H":
			col--
		default:
		}
		for _, tile := range s.hands[s.turn].hand {
			targets = append(targets, &bg.BoardGameAction{
				Team:       s.turn,
				ActionType: ActionPlaceTile,
				MoreDetails: PlaceTileActionDetails{
					Row:    row,
					Column: col,
					Tile:   tile.Edges,
				},
			})
		}
	}
	return targets
}

func (s *state) message() string {
	message := fmt.Sprintf("%s must place a tile", s.turn)
	switch s.variant {
	case VariantClassic, VariantOpenTiles, VariantLongestPath, VariantMostLoops:
		if len(s.winners) > 0 {
			message = fmt.Sprintf("%s tie", strings.Join(s.winners, ", "))
			if len(s.winners) == 1 {
				message = fmt.Sprintf("%s wins", s.winners[0])
			}
		}
	case VariantSolo:
		if len(s.winners) > 0 {
			if s.winners[0] == "true" {
				message = "you saved all the tokens"
			} else {
				saved := 0
				for _, team := range s.teams {
					if s.alive[team] {
						saved++
					}
				}
				message = fmt.Sprintf("you saved %d tokens", saved)
			}
		}
	}
	return message
}

func uniqueRandomToken(tokens map[string]*token, random *rand.Rand) *token {
	token := randomToken(random)
	for _, tok := range tokens {
		if token.Row == tok.Row && token.Col == tok.Col {
			return uniqueRandomToken(tokens, random)
		}
	}
	return token
}
