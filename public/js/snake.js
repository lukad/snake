(function () {

	var id = window.location.pathname.slice(1);
	var players = [];
	var connected = false;
	var state = {
		Others: [],
		You: {
			Points: []
		}
	};

	var connect = function() {
		var wsAddress = "ws://" + window.location.hostname + ":" + window.location.port + "/connect";
		console.log("connecting to:", wsAddress);
		var socket = new WebSocket(wsAddress);

		socket.onopen = function(event) {
			if (id == "") {
				socket.send("new");
			} else {
				socket.send(id);
			}
		};

		socket.onmessage = function(event) {
			if (!connected) {
				if (id == "new") {
					id = event.data;
					console.log("Connected to game", id);
					connected = true;
					return;
				} else {
					if (event.data == "notfound") {
						console.log("Game not found");
					} else {
						id = event.data;
						console.log("Connected to game", id);
						connected = true;
					}
					return
				}
			}

			state = JSON.parse(event.data);
		};

		socket.onclose = function (event) {
			console.log("Disconnected");
		};

		return socket;
	};

	var canvas = document.getElementById('canvas');
	var ctx = canvas.getContext('2d');

	var x = 0.0;
	var y = 0.0;

	var GRID_SIZE = 16;

	var update = function(dt) {};

	var drawCell = function(x, y) {
		ctx.fillRect(Math.floor(x) * GRID_SIZE,
				Math.floor(y) * GRID_SIZE,
				GRID_SIZE,
				GRID_SIZE);
	};

	var drawPlayer = function(player, color) {
		if (typeof(color) === "undefined") ctx.fillStyle = "#333";
		ctx.fillStyle = color;
		for (var i = player.Points.length - 1; i >= 0; i--) {
			drawCell(player.Points[i].X, player.Points[i].Y);
		}
	}

	var draw = function() {
		ctx.clearRect(0, 0, canvas.width, canvas.height);
	
		drawPlayer(state.You, "#1e1");
		if (state.Others) {
			for (var i = state.Others.length - 1; i >= 0; i--) {
				drawPlayer(state.Others[i]);
			}
		}

		ctx.beginPath();
		for (var i = Math.floor(GRID_SIZE); i < canvas.width; i += Math.floor(GRID_SIZE)) {
			ctx.moveTo(i, 0);
			ctx.lineTo(i, canvas.height);
			ctx.moveTo(0, i);
			ctx.lineTo(canvas.width, i);
		}

		ctx.strokeStyle = "#bbb";
		ctx.stroke();
	};

	var lastUpdate = Date.now();

	var run = function() {
		var now = Date.now();
		var dt = (now - lastUpdate) / 1000.0;
		lastUpdate = now;

		update(dt);
		draw();
	};

	var animate = function() {
		run();
		requestAnimationFrame(animate);
	}

	animate();

	var socket = connect();

	document.onkeydown = function() {
		switch (window.event.keyCode) {
			case 37:
				socket.send("left");
				break;
			case 38:
				socket.send("up");
				break;
			case 39:
				socket.send("right");
				break;
			case 40:
				socket.send("down");
				break;
		}
	};
})();
