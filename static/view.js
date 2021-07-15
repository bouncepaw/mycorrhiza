let wrapper = document.getElementsByClassName("top-bar__wrapper")[0],
    auth = document.getElementsByClassName("top-bar__section_auth")[0],
    highlights = document.getElementsByClassName("top-bar__section_highlights")[0]

const toggleElement = el => el.classList.toggle("top-bar__section_hidden-on-mobile")

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
wrapper.appendChild(hamburgerSection)