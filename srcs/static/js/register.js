document.getElementById("registerForm").addEventListener("submit", function(event) {
	// the default action is to send the form data to the server and then reload the page
	event.preventDefault();

	//Create FormData object from the form that event handler is attached to
	var formData = new FormData(this);
	document.getElementById("successMessage").style.display = "none";
	document.getElementById("errorMessage").style.display = "none";

	for (var pair of formData.entries()) {
		console.log(pair[0]+ ': ' + pair[1]);
	}

	fetch("/register", {
		method: "POST",
		body: formData,
		headers: {
			'Accept': 'application/json',
		}
	}).then(response => {
		return response.json();
	}).then(data => {
		if (data.success) {
			//changes the CSS display property of the HTML element with the id successMessage to block
			console.log("response ok");
			document.getElementById("successMessage").style.display = "block";
		} else {
			console.log("response failed");
			console.log(data.message);
			document.getElementById("errorMessage").innerText = "Register failed: " + data.message + "!";
			document.getElementById("errorMessage").style.display = "block";
		}
	}).catch(error => {
		console.error('Error:', error);
        alert('An error occurred');
	})
});
