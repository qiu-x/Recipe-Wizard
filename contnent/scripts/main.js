function openTab(evt, tabName) {
	var i, tabcontent, tablinks;
	tabcontent = document.getElementsByClassName("tabcontent");
	for (i = 0; i < tabcontent.length; i++) {
		tabcontent[i].style.display = "none";
	}
	tablinks = document.getElementsByClassName("tablinks");
	for (i = 0; i < tablinks.length; i++) {
		tablinks[i].className = tablinks[i].className.replace(" active", "");
	}
	document.getElementById(tabName).style.display = "block";
	if (evt != null) {
		evt.currentTarget.className += " active";
	}
}

// Open app tab on start
document.getElementById("app-button").click();

function post(url, data) {
	return fetch(url, {method: "POST", redirect: "follow", body: JSON.stringify(data)})
}

async function getRecipe() {
	let button = document.getElementById("generate-button");
	button.classList.add("loading");
	let textarea = document.getElementById("ingredients-input");
	let textarea_content = textarea.value;
	if (textarea_content.length < 3) {
		button.classList.remove("loading");
		alert("Please enter at least one ingredient");
		return;
	}
	const resp = await post("/req", {
		Ingredients: textarea_content,
	});
	const json = await resp.json();
	console.log(json);
	if (json.errorcode != "none") {
		button.classList.remove("loading");
		alert(json.errorcode);
		return;
	}
	button.classList.remove("loading");
	showCompletion(json.completion);
}

async function showCompletion(recipe_text) {
	let recipe_container = document.getElementById("recipe");
	let lines = recipe_text.split("\n");
	let firstline = "<b>" + lines.shift() + "</b>";
	let rest = lines.join("\n");

	recipe_container.innerHTML = firstline + rest;
	openTab(null, "app-result");
}
