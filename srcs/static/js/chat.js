
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

chat_socket.addEventListener('open', function(e){
	console.log("Websocket is now open")
})

document.getElementById("submit").addEventListener("click", function() {
	const message = document.getElementById("chat-input").value;
	if (message === "") {
		return;
	}
	console.log(message);
	document.getElementById("chat-input").value = "";
})
