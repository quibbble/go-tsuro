# Go-tsuro

Go-tsuro is a [Go](https://golang.org) implementation of the board game [Tsuro](https://en.wikipedia.org/wiki/Tsuro).

Check out [tsuro.quibbble.com](https://tsuro.quibbble.com) to play a live version of this game. This website utilizes [tsuro](https://github.com/quibbble/tsuro) frontend code, [go-tsuro](https://github.com/quibbble/go-tsuro) game logic, and [go-quibbble](https://github.com/quibbble/go-quibbble) server logic.

[![Quibbble Tsuro](https://raw.githubusercontent.com/quibbble/tsuro/main/screenshot.png)](https://tsuro.quibbble.com)

## Usage

To play a game create a new Tsuro instance:
```go
builder := Builder{}
game, err := builder.Create(&bg.BoardGameOptions{
    Teams: []string{"TeamA", "TeamB"}, // must contain at least 2 and at most 8 teams
    MoreOptions: TsuroMoreOptions{
        Seed: 123, // OPTIONAL - seed used to generate deterministic randomness which defaults to 0
        Variant: "Classic" // OPTIONAL - variants that change the game rules i.e. Classic (default), LongestPath, MostCrossings, OpenTiles, or Solo
    }
})
```

To rotate a tile in your hand do the following action:
```go
err := game.Do(&bg.BoardGameAction{
    Team: "TeamA",
    ActionType: "RotateTileRight", // can also be "RotateTileLeft"
    MoreDetails: RotateTileActionDetails{
        Tile: "ABCDEFGH"
    },
})
```

To place a tile on the board do the following action:
```go
err := game.Do(&bg.BoardGameAction{
    Team: "TeamA",
    ActionType: "PlaceTile",
    MoreDetails: PlaceTileActionDetails{
        Row: 0,
        Column: 1,
        Tile: "ABCDEFGH"
    },
})
```

To get the current state of the game call the following:
```go
snapshot, err := game.GetSnapshot("TeamA")
```
