package go_tsuro

// Action types
const (
	ActionPlaceTile       = "PlaceTile"
	ActionRotateTileRight = "RotateTileRight"
	ActionRotateTileLeft  = "RotateRileLeft"
)

// TsuroMoreOptions are the additional options for creating a game of Tsuro
type TsuroMoreOptions struct {
	Seed int64
}

// RotateTileActionDetails is the action details for rotating a tile in hand
type RotateTileActionDetails struct {
	Tile string
}

// PlaceTileActionDetails is the action details for placing a tile in the desired location on the board
type PlaceTileActionDetails struct {
	Row, Column int
	Tile        string
}

// TsuroSnapshotData is the game data unique to Tsuro
type TsuroSnapshotData struct {
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
