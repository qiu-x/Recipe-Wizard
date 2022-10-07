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

// POST request wrapper
function post(url, ms, data) {
	const controller = new AbortController()
	const signal = controller.signal;
	const promise = fetch(url, { signal: signal, method: "POST", redirect: "follow", body: JSON.stringify(data)});
	if (signal) signal.addEventListener("abort", () => controller.abort());
	const timeout = setTimeout(() => controller.abort(), ms);
	return promise.finally(() => clearTimeout(timeout));
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
	let resp;
	try {
		resp = await post(window.location.href + "/req", 20000, {
			Ingredients: textarea_content,
		});
	} catch (err) {
		button.classList.remove("loading");
		if (err.name === "AbortError") {
			alert("Timeout: It took more than 20 seconds to get the result!");
		} else if (err.name === "TypeError") {
			alert("AbortSignal.timeout() method is not supported");
		} else {
			// A network error, or some other problem.
			console.log(err);
			alert("Error: Server does not respond. Please check your internet connection.");
		}
		return;
	}
	if (!resp.ok) {
		button.classList.remove("loading");
		if (resp.status == 429) {
			alert("Request limit exceeded. Please wait for a moment.");
		} else {
			alert(`Error: ${resp.status}`);
		}
		return;
	}
	let json;
	json = await resp.json();

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
