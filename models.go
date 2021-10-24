package go_tsuro

// Action types
const (
	PlaceTile       = "PlaceTile"
	RotateTileRight = "RotateTileRight"
	RotateTileLeft  = "RotateRileLeft"
	Reset           = "Reset"
)

// RotateTileActionDetails is the action details for rotating a tile in hand
type RotateTileActionDetails struct {
	Tile string
}

// PlaceTileActionDetails is the action details for placing a tile in the desired location on the board
type PlaceTileActionDetails struct {
	Row, Column int
	Tile        string
}

// TsuroSnapshotDetails are the details unique to tsuro
type TsuroSnapshotDetails struct {
	Board          [][]*tile
	TilesRemaining int
	Hands          map[string][]*tile
	Tokens         map[string]*token
	Dragon         string
}

var tiles = []string{
	"ABCDEFGH", "AHBGCDEF", "AHBCDGEF", "AHBCDEFG", "AGBHCDEF",
	"ABCHDGEF", "ABCGDHEF", "AGBCDHEF", "ABCGDEFH", "AGBCDEFH",
	"ACBGDEFH", "ACBGDHEF", "ACBHDGEF", "ADBHCGEF", "ADBGCHEF",
	"ADBCEHFG", "ADBCEGFH", "AEBCDGFH", "AEBCDHFG", "AFBHCDEG",
	"AFBGCHDE", "AFBCDHEG", "AFBDCHEG", "AFBDCGEH", "AEBDCGFH",
	"ACBDEGFH", "AFBECHDG", "AFBECGDH", "AEBFCGDH", "ADBFCGEH",
	"ADBFCHEG", "ACBFDHEG", "ADBGCEFH", "AGBDCEFH", "ADBGCFEH",
}
