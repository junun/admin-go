let style = {};


// let bindTerminalResize = (term, websocket) => {
//     let onTermResize = size => {
//         websocket.send(
//             JSON.stringify({
//                 type: "resize",
//                 rows: size.rows,
//                 cols: size.cols,
//             })
//         );
//     };
//     // register resize event.
//     term.onResize("resize", onTermResize);
//     // unregister resize event when WebSocket closed.
//     websocket.addEventListener("close", function () {
//         term.off("resize", onTermResize);
//     });
// };

let bindTerminal = (term, websocket, bidirectional, bufferedTime) => {
    term.write('\r\nWelcome to Web SSH Terminal\n\r\n');
    term.socket = websocket;
    
    let messageBuffer = null;
    let handleWebSocketMessage = function (ev) {
        if (bufferedTime && bufferedTime > 0) {
            if (messageBuffer) {
                messageBuffer += ev.data;
            } else {
                messageBuffer = ev.data;
                setTimeout(function () {
                    term.write(messageBuffer);
                }, bufferedTime);
            }
        } else {
            term.write(ev.data);
        }
    };

    websocket.onmessage = handleWebSocketMessage;
    if (bidirectional) {
         term.onData(function (data) {
            websocket.send(JSON.stringify({'data': data, "type": "cmd"}));
          });
    }

    // send heartbeat package to avoid closing webSocket connection in some proxy environmental such as nginx.
    let heartBeatTimer = setInterval(function () {
        websocket.send(JSON.stringify({type: "heartbeat", data: ""}));
    }, 20 * 1000);

    websocket.addEventListener("close", function () {
        websocket.removeEventListener("message", handleWebSocketMessage);
        setTimeout(() => term.write('\r\nConnection is closed.\r\n'), 200)
        // term.dispose(function (data) {
        //     websocket.send(JSON.stringify({'data': data, "type": "cmd"}));
        //   });
        delete term.socket;
        clearInterval(heartBeatTimer);
    });
};


// function get_cell_size(term) {
//   style.width = term._core._renderService._renderer.dimensions.actualCellWidth;
//   style.height = term._core._renderService._renderer.dimensions.actualCellHeight;
// }


// function current_geometry(term) {
//   if (!style.width || !style.height) {
//     get_cell_size(term);
//   }
//   console.log(window.innerHeight);

//   const cols = parseInt(window.innerWidth / style.width, 10) - 1;
//   const rows = parseInt(window.innerHeight / style.height, 10);
//   console.log(rows);
//   return {'cols': cols, 'rows': rows};
// }


// function resize_terminal(term) {
//   console.log(geometry)
//   const geometry = current_geometry(term);
//   term.on_resize(geometry.cols, geometry.rows);
// }


// function read_as_text_with_decoder(file, callback, decoder) {
//   console.log(file)

//   let reader = new window.FileReader();

//   if (decoder === undefined) {
//     decoder = new window.TextDecoder('utf-8', {'fatal': true});
//   }

//   reader.onload = function () {
//     let text;
//     try {
//       // text = decoder.decode(reader.result);
//       text = reader.result;
//     } catch (TypeError) {
//       console.log('Decoding error happened.');
//     } finally {
//       if (callback) {
//         callback(text);
//       }
//     }
//   };

//   reader.onerror = function (e) {
//     console.error(e);
//   };

//   reader.readAsArrayBuffer(file);
// }


// function read_as_text_with_encoding(file, callback, encoding) {
//   let reader = new window.FileReader();

//   if (encoding === undefined) {
//     encoding = 'utf-8';
//   }

//   reader.onload = function () {
//     if (callback) {
//       callback(reader.result);
//     }
//   };

//   reader.onerror = function (e) {
//     console.error(e);
//   };

//   reader.readAsText(file, encoding);
// }


// function read_file_as_text(file, callback, decoder) {
//   console.log(!window.TextDecoder)
//   if (!window.TextDecoder) {
//     read_as_text_with_encoding(file, callback, decoder);
//   } else {
//     read_as_text_with_decoder(file, callback, decoder);
//   }
// }


function run(id, token) {
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
  const hostname = window.location.hostname;
  const port = window.location.port;
  const sock = new window.WebSocket(`${protocol}//${hostname}:${port}/admin/ws/${id}/ssh/${token}`),
  // const  sock = new window.WebSocket(`${protocol}//127.0.0.1:8080/admin/exec/ws/${id}/ssh/${token}/`),
    encoding = 'utf-8',
    decoder = window.TextDecoder ? new window.TextDecoder(encoding) : encoding,
    terminal = document.getElementById('terminal'),
    term = new window.Terminal({
      fontSize: 18,
      cursorBlink: true,
      cursorStyle: 'bar',
      bellStyle: "sound",
      theme: {
        background: 'black',
      }
    });

  term.fitAddon = new window.FitAddon.FitAddon();
  term.loadAddon(term.fitAddon);

  // function term_write(text) {
  //   if (term) {
  //     term.write(text);
  //     if (!term.resized) {
  //       resize_terminal(term);
  //       term.resized = true;
  //     }
  //   }
  // }

  // term.on_resize = function (cols, rows) {
  //   if (cols !== this.cols || rows !== this.rows) {
  //     sock.send(JSON.stringify({
  //                               "type": "resize",
  //                               "rows": rows,
  //                               "cols": cols,
  //                             }));
  //   }
  // };

  // term.onData(function (data) {
  //   sock.send(JSON.stringify({'data': data, "type": "cmd"}));
  // });

  sock.onopen = function () {
    term.open(terminal);
    term.fitAddon.fit();
    term.focus();
  };

  bindTerminal(term, sock, true, -1);
  // bindTerminalResize(term, sock);

  // sock.onmessage = function (msg) {
  //   read_file_as_text(msg.data, term_write, decoder);
  // };

  // sock.onerror = function (e) {
  //   console.error(e);
  // };

  // sock.onclose = function (e) {
  //   term.setOption("cursorBlink", false);
  //   if (e.code === 1005) {
  //     window.location.href = "about:blank";
  //     window.close()
  //   } 
  //   setTimeout(() => term.write('\r\nConnection is closed.\r\n'), 20000)
  //   if (sock) {
  //       sock.close()
  //   }
  //   // if (term) {
  //   //     term.dispose()
  //   // }
  // };

  // window.onresize = function () {
  //   if (term) {
  //     resize_terminal(term);
  //   }
  // };
  // window.addEventListener("resize", this.onWindowResize);
}
