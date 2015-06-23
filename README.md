# combo

Combo is an 8x8 board game. The top two rows start off with white pieces, and the bottom two rows black pieces. The players take turns moving a piece until the game is over. A player wins by capturing all his oponent's pieces.

A piece has two properties: color (black/white), and count (number of pieces that have merged to make this piece). Pieces may combine with friendly pieces by moving into them. Two friendly pieces of any count are able to merge (e.g. a count=3 piece can merge with a count=4 to produce a count=7). A piece is allowed to move up to count squares horizontally, vertically, or diagonally.  When moving, pieces are optionally able to split off any number of pieces less than their count. This split off piece can then be moved like a normal piece of that count.

You capture enemy pieces by moving onto their square. Your pieces cannot jump over other pieces. Any count piece can destroy any enemy piece, except count=1 pieces are not able to capture at all.

# setup

1. Make sure you have the [go tool installed](https://golang.org/dl/).
1. Make sure you have your [GOPATH configured](https://golang.org/doc/code.html) ("The GOPATH environment variable").
1. `git clone git@github.com:muirmanders/combo.git`
1. `cd combo; rake`
1. `$GOPATH/bin/combo http -l localhost:8080`
1. Open browser to [http://localhost:8080](http://localhost:8080)

# development notes

The HTTP resources are compiled into the binary for easy distribution, so when you change a resource you need to re-rake and restart combo. Another option is to run `rake dev_install`, which will cause combo to re-read the source files from disk each request, so you can edit resource files without recompiling/restarting combo. You still need to restart combo once after running `rake dev_install`. Don't forget to run `rake install` before checking in to re-generate the non-debug resources file.

# game design issues

- Currently it is possible to get into a stalemate/draw-like situation (e.g. both players are left with a count=2 piece at the end). It's not technically a stalemate, but unless somene makes a stupid mistake no one will ever win. I'm not sure yet how to deal with this. One option is to decide based on the game clock (i.e. first player to use up all their time loses).
