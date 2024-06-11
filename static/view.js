let wrapper = document.getElementsByClassName("top-bar__wrapper")[0],
    auth = document.getElementsByClassName("top-bar__section_auth")[0],
    highlights = document.getElementsByClassName("top-bar__section_highlights")[0]

const toggleElement = el => el.classList.toggle("top-bar__section_hidden-on-mobile")
toggleElement(auth)
toggleElement(highlights)

let hamburger = document.createElement("button")
hamburger.classList.add("top-bar__hamburger")
hamburger.onclick = _ => {
    toggleElement(auth)
    toggleElement(highlights)
}
hamburger.innerText = "Menu"

let hamburgerWrapper = document.createElement("div")
hamburgerWrapper.classList.add("top-bar__hamburger-wrapper")

let hamburgerSection = document.createElement("li")
hamburgerSection.classList.add("top-bar__section", "top-bar__section_hamburger")

hamburgerWrapper.appendChild(hamburger)
hamburgerSection.appendChild(hamburgerWrapper)
wrapper.appendChild(hamburgerSection);

(async () => {
    const input = document.querySelector('.js-add-cat-name'),
        datalist = document.querySelector('.js-add-cat-list')
    if (!input || !datalist) return;

    const categories = await fetch('/category')
        .then(resp => resp.text())
        .then(html => {
            return Array
                .from(new DOMParser()
                    .parseFromString(html, 'text/html')
                    .querySelectorAll('.mv-tags .p-name'))
                .map(a => a.innerText);
        });

    for (let cat of categories) {
        let optionElement = document.createElement('option')
        optionElement.value = cat
        datalist.appendChild(optionElement)
    }
    input.setAttribute('list', 'cat-name-options')
})();

