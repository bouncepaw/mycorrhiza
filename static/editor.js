(function () {
	let changed = false;
	let textarea = document.querySelector('.edit-form__textarea');
	let form = document.querySelector('.edit-form');

	let warnBeforeClosing = function (ev) {
		if (!changed) return;
		ev.preventDefault();
		return ev.returnValue = 'Are you sure you want to exit? You have unsaved changes.';
	};

	textarea.addEventListener('input', function () {
		changed = true;
	});

	form.addEventListener('submit', function () {
		changed = false;
	});

	window.addEventListener('beforeunload', warnBeforeClosing);
})();
