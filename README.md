# combo

Combo is an 8x8 board game. The top two rows start off with white pieces, and the bottom two rows black pieces. The players take turns moving a piece until the game is over. A player loses when he has no moves left.

Pieces may combine with friendly pieces by moving into them. A piece has two properties: color (black/white), and count (number of pieces that have merged to make this piece). A piece is allowed to move up to count-1 (count minus one) squares horizontally or vertically. The exception is a count=1 piece, which is allowed to merge with directly adjacent friendly pieces, but otherwise unable to move. Two friendly pieces of any count are able to merge (e.g. a count=3 piece can merge with a count=4 to produce a count=7).

You can destroy enemy pieces by moving onto their square. Your pieces cannot jump over friendly or enemy pieces. Note that count=1 pieces are not able to destroy enemy pieces (they are only able to merge with adjacent friendlies). As in chess where a pawn can capture a queen with no problem, any of your pieces can capture any enemy piece so long as your piece is able to move onto the enemy piece's square.

Merged pieces are able to split off one count=1 piece and place it on an adjacent empty square, or merge it with an adjacent friendly piece.