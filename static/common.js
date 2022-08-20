const $ = document.querySelector.bind(document)
const $$ = (...args) => Array.prototype.slice.call(document.querySelectorAll(...args))
const isMac = /Macintosh/.test(window.navigator.userAgent)
const arrToStr = a => Array.isArray(a) ? a.join('') : a

const rrh = {
    html(s, ...parts) {
        s = s.reduce((acc, cur, i) => (`${acc}${cur}${parts[i] ? arrToStr(parts[i]) : ''}`), '')
        const wrapper = document.createElement('div')
        wrapper.innerHTML = s
        return wrapper.children[0]
    },
    escape(text) {
        return text
            .replace('&', '&amp;')
            .replace('<', '&lt;')
            .replace('>', '&gt;')
    }
}