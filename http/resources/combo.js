// Copyright 2015 Muir Manders.  All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

var combo = {
  moves: {},
  split_piece_count: null
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
        drop: function(event, ui) {
          ui.draggable.addClass("dropped");

          // avoid piece flicker when dropping a piece and the helper disappears
          ui.helper.clone().appendTo(ui.helper.parent());

          combo.ws.send(JSON.stringify({
            command: "move",
            args: {
              from: {x: +ui.draggable.attr("data-x"), y: +ui.draggable.attr("data-y")},
              to: {x: +$(this).attr("data-x"), y: +$(this).attr("data-y")},
              piece_count: combo.split_piece_count || +ui.draggable.attr("data-count")
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

  for (var i = 0; i < args.moves.length; i++) {
    var m = args.moves[i];

    var from = m.from.x + "-" + m.from.y + "-" + m.piece_count;

    this.moves[from] = this.moves[from] || [];
    this.moves[from].push(m.to);
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
      piece.attr("data-count", square.piece_count);
      piece.attr("data-color", square.piece_color);

      piece.text(square.piece_count);
      piece.addClass("piece");
      piece.addClass(square.piece_color);
    }
  }

  this.set_move_type();
};

combo.set_move_type = function() {
  for (var x = 0; x < combo.width; x++) {
    for (var y = 0; y < combo.height; y++) {
      var piece = $(".piece[data-x="+x+"][data-y="+y+"]");
      if (piece.length == 0) {
        continue;
      }

      var split_piece_count = this.split_piece_count || +piece.attr("data-count");

      // clear out existing tos
      var classes = piece.attr("class").split(/\s+/);
      for (var i = 0; i < classes.length; i++) {
        if (classes[i].indexOf("to-") == 0) {
          piece.removeClass(classes[i]);
        }
      }

      var is_dragging = piece.hasClass("ui-draggable-dragging");

      if (is_dragging) {
        $(".valid-move").removeClass("valid-move");

        if (split_piece_count < +piece.attr("data-count")) {
          piece.show().css("opacity", 0.75).text(+piece.attr("data-count")-split_piece_count);
          $(".ui-draggable-helper").text(split_piece_count);
        } else {
         piece.hide();
         $(".ui-draggable-helper").show().text(split_piece_count);
        }
      } else {
        piece.text(piece.attr("data-count"));
      }

      var tos = this.moves[x+"-"+y+"-"+split_piece_count];
      if (tos) {
        piece.draggable({
          revert: "invalid",
          opacity: 0.75,
          disabled: false,
          helper: "clone",
          start: (function(tos, split_piece_count) {
            return function(event, ui) {
              var me = $(this);

              me.removeClass("dropped");

              // ack
              ui.helper.width(me.width()).height(me.height());
              ui.helper.addClass("ui-draggable-helper");

              if (split_piece_count < +me.attr("data-count")) {
                me.text(+me.attr("data-count") - split_piece_count).css("opacity", 0.75);
                ui.helper.text(split_piece_count);
              } else {
                me.hide();
              }

              for (var i = 0; i < tos.length; i++) {
                $(".cell[data-x="+tos[i].x+"][data-y="+tos[i].y+"]").addClass("valid-move");
              }
            };
          })(tos, split_piece_count),

          stop: function(event, ui) {
            var me = $(this);
            if (!me.hasClass("dropped")) {
              me.show().css("opacity", 1);
              me.text(me.attr("data-count"));
            }
            $(".valid-move").removeClass("valid-move");
          }
        });

        for (var i = 0; i < tos.length; i++) {
          piece.addClass("to-"+tos[i].x + "-" + tos[i].y);
          if (is_dragging) {
            $(".cell[data-x="+tos[i].x+"][data-y="+tos[i].y+"]").addClass("valid-move");
          }
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
    if (event.keyCode > 48 && event.keyCode <= 57) {
      if (combo.split_piece_count == event.keyCode - 48) {
        return;
      }
      combo.split_piece_count = event.keyCode - 48;
      combo.display_message("move piece count: " + combo.split_piece_count);
      combo.set_move_type();
    }
  });

  $(document).keyup(function(event) {
    if (event.keyCode > 48 && event.keyCode <= 57) {
      combo.split_piece_count = null;
      combo.display_message("");
      combo.set_move_type();
    }
  });
};

$(combo.init_key_handlers.bind(combo));
$(combo.open_ws.bind(combo));
