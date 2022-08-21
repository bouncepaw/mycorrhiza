const $ = document.querySelector.bind(document)
const $$ = (...args) => Array.prototype.slice.call(document.querySelectorAll(...args))
const isMac = /Macintosh/.test(window.navigator.userAgent)
const arrToStr = a => Array.isArray(a) ? a.join('') : a
const strToArr = a => Array.isArray(a) ? a : [a]

const rrh = {
    html(s, ...parts) {
        s = s.reduce((acc, cur, i) => (`${acc}${cur}${parts[i] ? arrToStr(parts[i]) : ''}`), '')
        const wrapper = document.createElement('div')
        wrapper.innerHTML = s
        return wrapper.children[0]
    },

    l10nMap: {},
    l10n(text, translations) {
        // Choose the translation on load to be consistent with the
        // server-rendered interface.
        if (translations) {
            translations.en = text
            this.l10nMap[text] = translations[navigator.languages
                .map(lang => lang.split('-')[0])
                .find(lang => translations[lang])] || text
        }
        return this.l10nMap[text] || text
    },
    l10nify(object) {
        return new Proxy(object, {
            get(target, prop, receiver) {
                const value = target[prop]
                if (value instanceof Function) {
                    return function (...args) {
                        args = args.map(arg => typeof arg === 'string' ? rrh.l10n(arg) : arg)
                        let result = value.apply(this === receiver ? target : this, args)
                        if (typeof result === 'object' && result !== null) {
                            result = rrh.l10nify(result)
                        }
                        return result
                    }
                }
                return value
            },
        })
    },
}