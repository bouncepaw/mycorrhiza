(function () {
	let changed = false;
	let textarea = document.querySelector('.edit-form__textarea');

	let warnBeforeClosing = function (ev) {
		ev.preventDefault();
		return ev.returnValue = 'Are you sure you want to exit? You have unsaved changes.';
	};

	textarea.addEventListener('input', function () {
		if (!changed) window.addEventListener('beforeunload', warnBeforeClosing);
		changed = true;
	});
})();
