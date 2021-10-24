# Go-tsuro

Go-tsuro is a [Go](https://golang.org) implementation of the board game [Tsuro](https://boardgamegeek.com/boardgame/16992/tsuro). Please note that this repo only includes game logic and a basic API to interact with the game but does NOT include any form of GUI.

Check out [quibbble.com](https://quibbble.com/paths) if you wish to view and play a live version of this game which utilizes this project along with a separate custom UI.
![Quibbble Paths](https://i.imgur.com/xdtMvHf.png)

## Usage

To play a game create a new Tsuro instance:
```go
game, err := NewTsuro(bg.BoardGameOptions{
    Teams: []string{"TeamA", "TeamB"} // must contain at least 2 and at most 8 teams
})
```

To rotate a tile in your hand do the following action:
```go
err := game.Do(bg.BoardGameAction{
    Team: "TeamA",
    ActionType: "RotateTileRight", // can also be "RotateTileLeft"
    MoreDetails: RotateTileActionDetails{
        Tile: "ABCDEFGH"
    },
})
```

To place a tile on the board do the following action:
```go
err := game.Do(bg.BoardGameAction{
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
