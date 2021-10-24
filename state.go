package go_tsuro

import (
	"fmt"
	"github.com/quibbble/go-boardgame/pkg/bgerr"
	"math/rand"
	"strings"
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
	active  map[string]bool // players that have placed and still alive
	alive   map[string]bool // players that are alive
}

func newState(teams []string) *state {
	hands := make(map[string]*hand)
	tokens := make(map[string]*token)
	alive := make(map[string]bool)
	deck := newDeck()
	for _, team := range teams {
		hand := newHand()
		for i := 0; i < 3; i++ {
			tile, _ := deck.Draw()
			hand.Add(tile)
		}
		hands[team] = hand
		token := uniqueRandomToken(tokens)
		tokens[team] = token
		alive[team] = true
	}
	return &state{
		turn:    teams[rand.Intn(len(teams))],
		teams:   teams,
		winners: make([]string, 0),
		board:   newBoard(),
		deck:    deck,
		tokens:  tokens,
		hands:   hands,
		dragon:  "",
		active:  make(map[string]bool),
		alive:   alive,
	}
}

func (s *state) RotateTileRight(team, tile string) error {
	if len(s.winners) > 0 {
		return &bgerr.Error{
			Err:    fmt.Errorf("%s game already completed", key),
			Status: bgerr.StatusGameOver,
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
	if len(s.winners) > 0 {
		return &bgerr.Error{
			Err:    fmt.Errorf("%s game already completed", key),
			Status: bgerr.StatusGameOver,
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
	if len(s.winners) > 0 {
		return &bgerr.Error{
			Err:    fmt.Errorf("%s game already completed", key),
			Status: bgerr.StatusGameOver,
		}
	}
	if team != s.turn {
		return &bgerr.Error{
			Err:    fmt.Errorf("currently %s's turn", s.turn),
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
	s.updateAlive()
	s.handleDraws()
	s.nextTurn()
	return nil
}

func (s *state) moveTokens() {
	moved := 0
	move := map[string]string{"A": "F", "B": "E", "C": "H", "D": "G", "E": "B", "F": "A", "G": "D", "H": "C"}
	for player, token := range s.tokens {
		if s.active[player] {
			t := s.board.board[token.Row][token.Col]
			if !mapContainsVal(t.Paths, player) {
				// first placement so move through the just placed tile
				destination := t.GetDestination(token.Notch)
				t.Paths[token.Notch+destination] = player
				token.Notch = destination
				// token was moved
				moved++
			} else if s.collided(s.tokens, player, token) {
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
				nextTile.Paths[startNotch+endNotch] = player
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

func (s *state) collided(tokens map[string]*token, player string, token *token) bool {
	for player2, token2 := range tokens {
		if player != player2 && (token.collided(token2) || token.equals(token2)) {
			return true
		}
	}
	return false
}

func (s *state) updateAlive() {
	if len(s.winners) > 0 {
		return
	}
	// alive before checking
	initialAlive := make([]string, 0)
	for _, player := range s.teams {
		if s.alive[player] {
			initialAlive = append(initialAlive, player)
		}
	}
	// update who is still alive
	for player, token := range s.tokens {
		if s.active[player] {
			if (token.Row == 0 && strings.Contains("AB", token.Notch)) ||
				(token.Row == rows-1 && strings.Contains("EF", token.Notch)) ||
				(token.Col == 0 && strings.Contains("GH", token.Notch)) ||
				(token.Col == columns-1 && strings.Contains("CD", token.Notch)) {
				// check on board edge
				s.setLost(player)
			} else if s.collided(s.tokens, player, token) {
				// check if collided with another token
				s.setLost(player)
			}
		}
	}
	// who is still alive
	stillAlive := make([]string, 0)
	for _, player := range s.teams {
		if s.alive[player] {
			stillAlive = append(stillAlive, player)
		}
	}
	if len(stillAlive) == 0 {
		// no more alive so initial alive all win
		s.winners = initialAlive
	} else if len(stillAlive) == 1 {
		// one alive so they win
		s.winners = stillAlive
	} else if s.board.getTileCount() == len(Tiles) {
		// all tiles have been placed remaining alive are winners
		s.winners = stillAlive
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
	for idx, player := range s.teams {
		if player == turn {
			nextTurn = s.teams[(idx+1)%len(s.teams)]
			if !s.alive[nextTurn] {
				return s.getNextTurn(nextTurn)
			}
			return nextTurn
		}
	}
	return nextTurn
}

func (s *state) setLost(player string) {
	s.alive[player] = false
	s.active[player] = false
	s.deck.Add(s.hands[player].hand...)
	s.hands[player].Clear()
	if s.aliveCount() <= 0 {
		return
	}
	next := s.getNextTurn(s.turn)
	if s.dragon == player && len(s.hands[next].hand) < 3 {
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

func uniqueRandomToken(tokens map[string]*token) *token {
	token := RandomToken()
	for _, tok := range tokens {
		if token.Row == tok.Row && token.Col == tok.Col {
			return uniqueRandomToken(tokens)
		}
	}
	return token
}
