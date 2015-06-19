// Copyright 2015 Muir Manders.  All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

var combo = {};

combo.init_board = function(width, height) {
  var cells = [];

  this.width = width;
  this.height = height;

  for (var y = 0; y < height; y++) {
    cells[y] = [];
    var row = $("<div class=row>");
    for (var x = 0; x < width; x++) {
      cells[y][x] = $("<div class=cell>").attr("data-x", x).attr("data-y", y).appendTo(row).droppable({
        tolerance: "intersect",
        accept: ".to-"+x+"-"+y,
        activeClass: "valid-move",
        drop: function(event, ui) {
          combo.ws.send(JSON.stringify({
            command: "move",
            args: {
              from: {x: +ui.draggable.attr("data-x"), y: +ui.draggable.attr("data-y")},
              to: {x: +$(this).attr("data-x"), y: +$(this).attr("data-y")}
            }
          }));
        }
      });
    }
    row.appendTo("#board");
  }

  this.cells = cells;
};

combo.open_ws = function() {
  var width = 8, height = 8;

  var ws = this.ws = new WebSocket("ws://" + location.host + "/connect");

  ws.onmessage = function(event) {
    var cmd = JSON.parse(event.data);
    switch (cmd.command) {
      case "move":
      combo.move(cmd.args);
      break;
    }
  };

  ws.onopen = function() {
    ws.send(JSON.stringify({
      command: "new_game",
      args: {
        width: width,
        height: height
      }
    }));

    combo.init_board(width, height);
  };

  ws.onerror = function(e) {
    console.log("websocket error: " + e);
  };
};

combo.move = function(args) {
  var moves = {};
  var splits = {};

  for (var i = 0; i < args.moves.length; i++) {
    var m = args.moves[i];

    var from = m.from.x + "-" + m.from.y;
    var to = m.to.x + "-" + m.to.y;

    var type = m.split ? splits : moves;
    type[from] = type[from] || [];
    type[from].push(to);
  }

  for (var x = 0; x < this.width; x++) {
    for (var y = 0; y < this.height; y++) {
      var square = args.board[x][y];

      var cell = combo.cells[y][x];
      var piece = cell.find(".piece");

      if (square.piece_count == 0) {
        piece.remove();
      } else {
        if (piece.length == 0) {
          piece = $("<div class=piece>").appendTo(cell);
          piece.attr("data-x", x);
          piece.attr("data-y", y);
        }

        piece.text(square.piece_count);
      }

      piece.removeClass();
      piece.addClass("piece");

      if (square.piece_color == "white") {
        piece.addClass("white");
      } else if (square.piece_color == "black") {
        piece.addClass("black");
      }

      var tos = moves[x+"-"+y];
      if (tos) {
        piece.draggable({
          revert: "invalid",
          opacity: 0.75
        });

        for (var i = 0; i < tos.length; i++) {
          piece.addClass("to-" + tos[i]);
        }
      }
    }
  }
};

$(combo.open_ws.bind(combo));
