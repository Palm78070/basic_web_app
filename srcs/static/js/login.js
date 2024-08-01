document.getElementById("loginForm").addEventListener("submit", function(event){
	// the default action is to send the form data to the server and then reload the page
	event.preventDefault();

	//Create FormData object from the form that event handler is attached to => this refer to DOM element "loginForm"
	var formData = new FormData(this);
	document.getElementById("errorMessage").style.display = "none";

	fetch("/login", {
		method: "POST",
		body: formData,
		headers: {
			'Accept': 'application/json',
		}
	}).then(response=>{
		return response.json(); //ensures that the next .then block only runs after the JSON has been successfully parsed
	}).then(data=>{
		if (!data.success) {
			console.log("response failed");
			console.log(data.message);
			document.getElementById("errorMessage").innerText = "Login failed:" + data.message + "!";
			document.getElementById("errorMessage").style.display = "block";
		}
	}).catch(error=>{
		console.error('Error:', error);
		alert('An error occurred')
	})
})
