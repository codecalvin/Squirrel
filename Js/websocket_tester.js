// https://127.0.0.1:10443/ws

socket = new WebSocket("wss:/127.0.0.1:10443/ws/api?uname=wangha")
socket.onmessage = function(evt) { console.log("insocket1", evt.data); }
socket.send("from client1")

socket2 = new WebSocket("wss:/127.0.0.1:10443/ws/api?uname=wangha")
socket2.onmessage = function(evt) { console.log("insocket2", evt.data); }
socket2.send("from client2")


