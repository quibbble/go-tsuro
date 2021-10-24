package go_tsuro

/*
   tile representation

            A  B
           ——  ——
       H |        | C
       G |        | D
           ——  ——
            F  E
*/
type tile struct {
	Edges string            // defines the tile
	Paths map[string]string // map from path to team defining section team path
}

func newTile(edges string) *tile {
	return &tile{
		Edges: edges,
		Paths: make(map[string]string),
	}
}

func (t *tile) GetDestination(start string) string {
	for idx, char := range t.Edges {
		if string(char) == start && idx%2 == 0 {
			return string(t.Edges[idx+1])
		} else if string(char) == start && idx%2 == 1 {
			return string(t.Edges[idx-1])
		}
	}
	return ""
}

func (t *tile) RotateRight() {
	transform := map[string]string{"A": "C", "B": "D", "C": "E", "D": "F", "E": "G", "F": "H", "G": "A", "H": "B"}
	transformed := ""
	for _, char := range t.Edges {
		transformed += transform[string(char)]
	}
	t.Edges = transformed
}

func (t *tile) RotateLeft() {
	transform := map[string]string{"A": "G", "B": "H", "C": "A", "D": "B", "E": "C", "F": "D", "G": "E", "H": "F"}
	transformed := ""
	for _, char := range t.Edges {
		transformed += transform[string(char)]
	}
	t.Edges = transformed
}

func (t *tile) equals(t2 *tile) bool {
	if t.Edges == t2.Edges {
		return true
	}
	return false
}

func (t *tile) in(list []*tile) bool {
	for _, t2 := range list {
		if t.equals(t2) {
			return true
		}
	}
	return false
}
