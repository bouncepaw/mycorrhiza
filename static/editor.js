(function () {
	window.hyphaChanged = false;
	let textarea = document.querySelector('.edit-form__textarea');
	let form = document.querySelector('.edit-form');

	let warnBeforeClosing = function (ev) {
		if (!window.hyphaChanged) return;
		ev.preventDefault();
		return ev.returnValue = 'Are you sure you want to exit? You have unsaved changes.';
	};

	textarea.addEventListener('input', function () {
		window.hyphaChanged = true;
	});

	form.addEventListener('submit', function () {
		window.hyphaChanged = false;
	});

	window.addEventListener('beforeunload', warnBeforeClosing);
})();
