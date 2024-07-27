
const login_username = (function() {
	return document.getElementById("login-username").innerText.split(" ")[2];
})();
console.log("login_user: ", login_username);

const recipient = document.getElementById("username").innerText
console.log("recipient: ", recipient)

chat_socket = new WebSocket(
	'ws://' +
	window.location.host +
	'/ws/private/' +
	login_username
);

chat_socket.onopen = function(e) {
	console.log("Websocket is now open")
}

chat_socket.onmessage = function(e) {
	console.log("Message received: ", e.data);
	const data = JSON.parse(e.data);
	console.log(data);
	const message = data.message;
	const sender = data.sender;
	// const recipient = data.recipient;

	document.getElementById("chat-room").innerHTML += `<div class="recipient-bubble-container">` + sender + `: ` + message + `</div>`;
}

function handleSendMessage(e) {
	if (e.type === "click" || (e.type === "keypress" && e.key === "Enter")) {
		const message = document.getElementById("chat-input").value;
		if (message === "") {
			return;
		}
		// console.log(message);
		document.getElementById("chat-input").value = "";
		chat_socket.send(JSON.stringify({
			"message": message,
			"sender": login_username,
			"recipient": recipient
		}))
	}
}

document.getElementById("submit").addEventListener("click", handleSendMessage);

document.getElementById("chat-input").addEventListener("keypress", handleSendMessage);
