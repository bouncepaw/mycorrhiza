(() => {
    const branches = document.querySelectorAll('#hypha-list details')
    if (!branches.length || !('open' in document.createElement('details'))) return

    // Keep every hypha visible without JavaScript; otherwise start collapsed.
    for (const branch of branches) branch.open = false

    const controls = document.getElementById('hypha-list-controls')
    for (const button of controls.querySelectorAll('button')) {
        button.addEventListener('click', () => {
            for (const branch of branches) branch.open = button.dataset.open === 'true'
        })
    }
    controls.hidden = false
})()
