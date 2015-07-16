# combo

Combo is an 8x8 board game. The top two rows start off with white pieces, and the bottom two rows black pieces. The players take turns moving a piece until the game is over. A player loses when all his pieces are captured or he has no available moves. Black always goes first.

A piece has two properties: color (black/white), and count (number of pieces that have merged to make this piece). Pieces may combine with friendly pieces by moving into them. Two friendly pieces of any count are able to merge (e.g. a count=3 piece can merge with a count=4 to produce a count=7). A piece is allowed to move up to count squares horizontally, vertically, or diagonally.  When moving, pieces are optionally able to split off any number of pieces less than their count. This split off piece can then be moved like a normal piece of that count.

You capture enemy pieces by moving onto their square. Your pieces cannot jump over other pieces. Any count piece can destroy any enemy piece, except count=1 pieces are not able to capture at all.

# setup

1. Make sure you have the [go tool installed](https://golang.org/dl/).
1. Make sure you have your [GOPATH configured](https://golang.org/doc/code.html) ("The GOPATH environment variable").
1. `git clone git@github.com:muirmanders/combo.git`
1. `cd combo; rake`
1. `$GOPATH/bin/combo http`
1. Open browser to [http://localhost:8080](http://localhost:8080)

# external player interface

This is the API contest submissions must implement. The combo game server will start your program once and then communicate with it over STDIN/STDOUT for the duration of the game, sending and receiving newline terminated JSON messages. There is only one type of message sent to your program, and only one type expected as a response. If your player makes an illegal move, you lose immediately.

From combo server to player application on each turn (whitespace added for clarity, squares omitted for brevity):

    {
      // the color you are playing as
      "color": "white",

      // the board state
      "board": {
        "width": 8,
        "height": 8,

        // two dimensional array of board squares
        "squares": [
          [
            {"x": 0, "y": 0, "piece_color": "white", "piece_count": 1},
            {"x": 0, "y": 1, "piece_color": "white", "piece_count": 1},

            ...
          ],

          ...

          [
            {"x": 0, "y": 4, "piece_count": 0},
            {"x": 0, "y": 5, "piece_count": 0},

            ...
          ],

          ...

          [
            ...

            {"x": 7, "y": 6, "piece_color": "black", "piece_count": 1},
            {"x": 7, "y": 7, "piece_color": "black", "piece_count": 1}
          ]
        ]
      }
    }

Response from player application indicating move choice (whitespace added for clarity):

    {
      "from": {
        "x": 3,
        "y": 2
      },
      "to": {
        "x": 5,
        "y": 4,
      },
      "piece_count": 3
    }


There is an example player in external/example.rb. You can invoke it (or your player) via `$GOPATH/bin/combo http --cpu=external/example.rb` (assuming you are in the combo project root directory).

# development notes

The HTTP resources are compiled into the binary for easy distribution, so when you change a resource you need to re-rake and restart combo. Another option is to run `rake dev_install`, which will cause combo to re-read the source files from disk each request, so you can edit resource files without recompiling/restarting combo. You still need to restart combo once after running `rake dev_install`. Don't forget to run `rake install` before checking in to re-generate the non-debug resources file.

# game design issues

- Currently it is possible to get into a stalemate/draw-like situation (e.g. both players are left with a count=2 piece at the end). It's not technically a stalemate, but unless somene makes a stupid mistake no one will ever win. I'm not sure yet how to deal with this. One option is to decide based on the game clock (i.e. first player to use up all their time loses).
