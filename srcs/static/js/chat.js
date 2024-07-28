
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

	const bubble_type = sender === login_username ? "sender-bubble-container": "recipient-bubble-container";

	const msg_div = document.createElement("div");
	msg_div.className = bubble_type;
	msg_div.textContent = sender + ": " + message;

	document.getElementById("chat-room").appendChild(msg_div);

	// document.getElementById("chat-room").innerHTML += `<div class=${bubble_type}>` + sender + `: ` + message + `</div>`;
	scrollToBottom();

	// Force reflow
	// chat_room.offsetHeight;
	// document.getElementById("chat-room").offsetHeight;

	// Use a timeout to ensure the DOM updates before scrolling
	// setTimeout(scrollToBottom, 0);
}

function scrollToBottom() {
	const chat_room = document.getElementById("chat-room");
	console.log("Scrolling to bottom, chat_room.scrollHeight: ", chat_room.scrollHeight);
	chat_room.scrollTop = chat_room.scrollHeight;
	console.log("chat_room.scrollTop after scroll: ", chat_room.scrollTop);
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
