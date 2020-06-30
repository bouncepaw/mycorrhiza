var isOpen = false
var sidebar = document.getElementById('sidebar')
var btn = document.getElementById('shroomburger')
btn.addEventListener('click', function() {
	if (isOpen) {
		sidebar.classList.add('hidden_mobile')
	} else {
		sidebar.classList.remove('hidden_mobile')
	}
	isOpen = !isOpen
})
