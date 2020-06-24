var menu = document.getElementById('shroomburgerMenu');
document.getElementById('shroomBtn').addEventListener('click', function() {
    menu.classList.add('active');
});
document.getElementById('mushroomBtn').addEventListener('click', function() {
    menu.classList.remove('active');
});
