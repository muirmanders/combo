// Copyright 2015 Muir Manders.  All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

var combo = {
  moves: {},
  splits: {}
};

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
              to: {x: +$(this).attr("data-x"), y: +$(this).attr("data-y")},
              split: combo.splitting
            }
          }));
        }
      });
    }
    row.appendTo("#board");
  }

  this.cells = cells;
};

combo.display_message = function(msg) {
  $("#messages").text(msg);
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
    case "game_over":
      combo.display_message(cmd.args.message);
      $(".piece").draggable({disabled: true});
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
  this.moves = {};
  this.splits = {};

  for (var i = 0; i < args.moves.length; i++) {
    var m = args.moves[i];

    var from = m.from.x + "-" + m.from.y;
    var to = m.to.x + "-" + m.to.y;

    var type = m.split ? this.splits : this.moves;
    type[from] = type[from] || [];
    type[from].push(to);
  }

  for (var x = 0; x < this.width; x++) {
    for (var y = 0; y < this.height; y++) {
      var square = args.board[x][y];

      var cell = combo.cells[y][x];

      cell.find(".piece").remove();

      if (square.piece_count == 0) {
        continue;
      }

      var piece = $("<div class=piece>").appendTo(cell);
      piece.attr("data-x", x);
      piece.attr("data-y", y);

      piece.text(square.piece_count);
      piece.addClass("piece");
      piece.addClass(square.piece_color);
    }
  }

  this.set_move_type();
};

combo.set_move_type = function() {
  var old_type, new_type;
  if (this.splitting) {
    old_type = this.moves;
    new_type = this.splits;
  } else {
    old_type = this.splits;
    new_type = this.moves;
  }

  for (var x = 0; x < combo.width; x++) {
    for (var y = 0; y < combo.height; y++) {
      var piece = $(".piece[data-x="+x+"][data-y="+y+"]");
      if (piece.length == 0) {
        continue;
      }

      var tos = old_type[x+"-"+y];
      if (tos) {
        for (var i = 0; i < tos.length; i++) {
          piece.removeClass("to-"+tos[i]);
        }
      }

      tos = new_type[x+"-"+y];
      if (tos) {
        piece.draggable({
          revert: "invalid",
          opacity: 0.75,
          disabled: false
        });
        for (var i = 0; i < tos.length; i++) {
          piece.addClass("to-"+tos[i]);
        }
      } else {
        piece.draggable({disabled: true});
      }
    }
  }
};

combo.init_key_handlers = function() {
  $(document).focus();

  $(document).keydown(function(event) {
    if (event.keyCode == 16) {
      combo.splitting = true;
      combo.set_move_type();
    }
  });

  $(document).keyup(function(event) {
    if (event.keyCode == 16) {
      combo.splitting = false;
      combo.set_move_type();
    }
  });
};

$(combo.init_key_handlers.bind(combo));
$(combo.open_ws.bind(combo));
