var navLinks = document.querySelectorAll("nav a");
for (var i = 0; i < navLinks.length; i++) {
	var link = navLinks[i]
	if (link.getAttribute('href') == window.location.pathname) {
		link.classList.add("live");
		break;
	}
}

// Code to add and delete rules from the scoresheet edit page
if (!window.scoresheetInitialized) {
	window.scoresheetInitialized = true;

	document.addEventListener("DOMContentLoaded", function() {
		const container = document.getElementById("rules-container");
		const addButton = document.getElementById('add-rule-button');
		
		if (!addButton || !container) return;

		function updateRuleIndices() {
			const groups = document.querySelectorAll('.rule-group');
			groups.forEach((group, index) => {
				group.setAttribute('data-index', index);

				const titleEl = group.querySelector(".rule-title");
				if (titleEl) titleEl.textContent = 'Rule ' + (index + 1) ;

				const nameInput = group.querySelector('input[name$="_name"]');
				const multInput = group.querySelector('input[name$="_multiplier"]');
				const descInput = group.querySelector('input[name$="_description"]');

				if (nameInput) nameInput.name = 'rule' + index + '_name';
				if (multInput) multInput.name = 'rule' + index + '_multiplier';
				if (descInput) descInput.name = 'rule' + index + '_description';
			});
		}
		
		addButton.addEventListener("click", function() {
			const existingGroups = document.querySelectorAll(".rule-group")
			const nextIndex = existingGroups.length;

			const ruleDiv = document.createElement("div");
			ruleDiv.className = "rule-group";
			ruleDiv.setAttribute("data-index", nextIndex);

			ruleDiv.innerHTML = `
				<div class="rule-title">Rule ${nextIndex + 1}</div>
				<div class="label-wrap">
					<label class="rule rule-label">Name: </label>
					<input
					class="rule rule-text"
					type="text"
					name="rule${nextIndex}_name"
					value="">
				</div>
				<div class="label-wrap">
					<label class="rule rule-label">Multiplier: </label>
					<input
					class="rule rule-number"
					type="number"
					name="rule${nextIndex}_multiplier"
					value="0"
					min="0"
					max="20"
					step="0.5">
				</div>
				<div class="label-wrap">
					<label class="rule rule-label">Short Description: </label>
					<input
					class="rule rule-text"
					type="text"
					name="rule${nextIndex}_description"
					value="">
				</div>
				<div class="rule-button">
                    <button type="button" class="delete-rule-button">Delete Rule</button>
                </div>
			`;
			
			container.appendChild(ruleDiv);
		});

		container.addEventListener('click', function(e) {
			if (e.target && e.target.classList.contains('delete-rule-button')) {
				const group = e.target.closest('.rule-group');
				if (group) {
					group.remove();
					updateRuleIndices();
				}
			}
		});
	});
};


// Code to add and delete players from the session create page
if (!window.sessionInitialized) {
	window.sessionInitialized = true;

	document.addEventListener("DOMContentLoaded", function() {
		const container = document.getElementById("players-container");
		const addButton = document.getElementById('add-player-button');
		
		if (!addButton || !container) return;

		function updateRuleIndices() {
			const groups = document.querySelectorAll('.rule-group');
			groups.forEach((group, index) => {
				group.setAttribute('data-index', index);

				const titleEl = group.querySelector(".rule-title");
				if (titleEl) titleEl.textContent = 'Rule ' + (index + 1) ;

				const nameInput = group.querySelector('input[name$="_name"]');
				const multInput = group.querySelector('input[name$="_multiplier"]');
				const descInput = group.querySelector('input[name$="_description"]');

				if (nameInput) nameInput.name = 'rule' + index + '_name';
				if (multInput) multInput.name = 'rule' + index + '_multiplier';
				if (descInput) descInput.name = 'rule' + index + '_description';
			});
		}
		
		addButton.addEventListener("click", function() {
			const existingGroups = document.querySelectorAll(".rule-group")
			const nextIndex = existingGroups.length;

			const ruleDiv = document.createElement("div");
			ruleDiv.className = "rule-group";
			ruleDiv.setAttribute("data-index", nextIndex);

			ruleDiv.innerHTML = `
				<div class="rule-title">Rule ${nextIndex + 1}</div>
				<div class="label-wrap">
					<label class="rule rule-label">Name: </label>
					<input
					class="rule rule-text"
					type="text"
					name="rule${nextIndex}_name"
					value="">
				</div>
				<div class="label-wrap">
					<label class="rule rule-label">Multiplier: </label>
					<input
					class="rule rule-number"
					type="number"
					name="rule${nextIndex}_multiplier"
					value="0"
					min="0"
					max="20"
					step="0.5">
				</div>
				<div class="label-wrap">
					<label class="rule rule-label">Short Description: </label>
					<input
					class="rule rule-text"
					type="text"
					name="rule${nextIndex}_description"
					value="">
				</div>
				<div class="rule-button">
                    <button type="button" class="delete-rule-button">Delete Rule</button>
                </div>
			`;
			
			container.appendChild(ruleDiv);
		});

		container.addEventListener('click', function(e) {
			if (e.target && e.target.classList.contains('delete-rule-button')) {
				const group = e.target.closest('.rule-group');
				if (group) {
					group.remove();
					updateRuleIndices();
				}
			}
		});
	});
};