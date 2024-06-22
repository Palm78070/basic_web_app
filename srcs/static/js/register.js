document.getElementById("register").addEventListener("submit", function(event) {
	// the default action is to send the form data to the server and then reload the page
	event.preventDefault();

	//Create FormData object from the form that event handler is attached to
	var formData = new FormData(this);

	fetch("/register", {
		method: "POST",
		body: formData
	}).then(response => {
		if (response.ok) {
			//changes the CSS display property of the HTML element with the id successMessage to block
			document.getElementById("successMessage").style.display = "block";
		} else {
			alert("Failed to register. Please try again.");
		}
	}).catch(error => {
		console.error('Error:', error);
        alert('An error occurred');
	})
});
