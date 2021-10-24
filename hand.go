package go_tsuro

import "errors"

type hand struct {
	hand []*tile
}

func newHand() *hand {
	return &hand{hand: make([]*tile, 0)}
}

func (h *hand) Add(tiles ...*tile) {
	for _, tile := range tiles {
		h.hand = append(h.hand, tile)
	}
}

func (h *hand) Remove(tile *tile) error {
	for idx, t := range h.hand {
		if tile.equals(t) {
			h.hand = append(h.hand[:idx], h.hand[idx+1:]...)
			return nil
		}
	}
	return errors.New("tile not found")
}

func (h *hand) Clear() {
	h.hand = make([]*tile, 0)
}

func (h *hand) IndexOf(tile *tile) int {
	for idx, t := range h.hand {
		if tile.equals(t) {
			return idx
		}
	}
	return -1
}
