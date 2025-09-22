import $ from "/static/framework.js";

(() => {
	document.addEventListener("DOMContentLoaded", () => {
		$.$registerRoot($.$byId("root"));
		$.$create("p").$textContent("Hello World!");
	});
})();
