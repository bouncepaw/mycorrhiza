rrh.l10n('List of shortcuts', { ru: 'Горячие клавиши' })
rrh.l10n('Close this dialog', { ru: 'Закрыть диалог' })

rrh.l10n('Common', { ru: 'Общее' })
rrh.l10n('Home', { ru: 'Главная' })

rrh.l10n('Hypha', { ru: 'Гифа' })

rrh.l10n('Editor', { ru: 'Редактор' })

rrh.l10n('Format', { ru: 'Форматирование' })

function keyEventToShortcut(event) {
    let elideShift = event.key.toUpperCase() === event.key && event.shiftKey
    return (event.ctrlKey ? 'Ctrl+' : '') +
        (event.altKey ? 'Alt+' : '') +
        (event.metaKey ? 'Meta+' : '') +
        (!elideShift && event.shiftKey ? 'Shift+' : '') +
        (event.key === ',' ? 'Comma' : event.key === ' ' ? 'Space' : event.key)
}

function prettifyShortcut(shortcut) {
    let keys = shortcut.split('+')

    if (isMac) {
        let cmdIdx = keys.indexOf('Meta')
        if (cmdIdx !== -1 && keys.length - cmdIdx > 2) {
            let tmp = keys[cmdIdx + 1]
            keys[cmdIdx + 1] = 'Meta'
            keys[cmdIdx] = tmp
        }
    }

    let lastKey = keys[keys.length - 1]
    if (!keys.includes('Shift') && lastKey.toUpperCase() === lastKey && lastKey.toLowerCase() !== lastKey) {
        keys.splice(keys.length - 1, 0, 'Shift')
    }

    for (let i = 0; i < keys.length; i++) {
        if (isMac) {
            if (keys[i] === 'Ctrl') keys[i] = '⌃'
            if (keys[i] === 'Alt') keys[i] = '⌥'
            if (keys[i] === 'Shift') keys[i] = '⇧'
            if (keys[i] === 'Meta') keys[i] = '⌘'
        } else {
            if (keys[i] === 'Meta') keys[i] = 'Win'
        }

        if (i === keys.length - 1 && i > 0 && keys[i].length === 1) {
            keys[i] = keys[i].toUpperCase()
        }

        if (keys[i] === 'ArrowLeft') keys[i] = '←'
        if (keys[i] === 'ArrowRight') keys[i] = '→'
        if (keys[i] === 'ArrowUp') keys[i] = '↑'
        if (keys[i] === 'ArrowDown') keys[i] = '↓'
        if (keys[i] === 'Comma') keys[i] = ','
        if (keys[i] === 'Enter') keys[i] = '↩'
        if (keys[i] === ' ') keys[i] = 'Space'
        keys[i] = `<kbd>${keys[i]}</kbd>`
    }

    return keys.join(isMac ? '' : ' + ')
}

function isTextField(element) {
    let name = element.nodeName.toLowerCase()
    return name === 'textarea' ||
        name === 'select' ||
        (name === 'input' && !['submit', 'reset', 'checkbox', 'radio'].includes(element.type)) ||
        element.isContentEditable
}

class ShortcutHandler {
    constructor(element) {
        // This shit is built like this: map is the whole shortcut graph. active
        // points to the node of the shortcut graph that's currently "selected".
        // When the user presses a key, the code finds the key node in the
        // active subgraph and sets active to the found node. On each key press
        // we get a narrower view until we reach a leaf and execute the action.
        this.map = {}
        this.active = this.map
        element.addEventListener('keydown', event => this.handleKeyDown(event))
    }

    register(shortcuts, action, other = {}) {
        // shortcuts looks like this: ['g r', 'Ctrl+r']
        // Every item of shortcuts is a sequential key chord.
        for (let chord of strToArr(shortcuts)) {
            let leaf = this.map
            let keys = chord.trim().split(' ')
            for (let key of keys) {
                // If there's no existing edge, create one
                if (!leaf[key]) leaf[key] = {}
                leaf = leaf[key]
                if (leaf.action) throw new Error(`Shortcut ${chord} already exists`)
            }
            // Now we've traversed to the leaf. Bind the shortcut
            leaf.action = action
            Object.assign(leaf, other)
        }
    }

    handleKeyDown(event) {
        if (event.defaultPrevented) return
        if (['Control', 'Alt', 'Shift', 'Meta'].includes(event.key)) return
        if ((!event.ctrlKey && !event.metaKey && !event.shiftKey && !event.altKey) &&
            event.target instanceof Node && isTextField(event.target)) return

        let shortcut = keyEventToShortcut(event)

        if (!this.active[shortcut]) {
            this.resetActive()
            return
        }

        this.active = this.active[shortcut]
        if (this.active.action && (!this.active.element || event.target === this.active.element)) {
            event.stopPropagation()
            if (this.active.force) event.preventDefault()
            this.active.action(event)
            this.resetActive()
            return
        }

        if (this.timeout) clearTimeout(this.timeout)
        this.timeout = window.setTimeout(() => this.resetActive(), 1500)
    }

    resetActive() {
        this.active = this.map
        if (this.timeout) {
            clearTimeout(this.timeout)
            this.timeout = null
        }
    }
}

class ShortcutGroup {
    constructor(handler, element = null) {
        this.shortcuts = []
        this.handler = handler
        this.element = element
    }

    bind(shortcuts, target, description = null, other = {}) {
        shortcuts = strToArr(shortcuts)
        if (!other.element) other.element = this.element
        if (target instanceof Function) {
            this.shortcuts.push({ shortcuts, description })
            this.handler.register(shortcuts, target, other)
        } else if (target instanceof Node) {
            this.bind(shortcuts, () => {
                if (isTextField(target)) target.focus()
                target.click()
            }, description)
        } else if (Array.isArray(target) && (target.length === 0 || target[0] instanceof Node)) {
            this.shortcuts.push({
                shortcuts: shortcuts.map(s => `${s} 1 — 9`),
                description,
            })
            for (let i = 0; i < target.length && i < 9; i++) {
                let element = target[i]
                this.handler.register(shortcuts.map(s => `${s} ${i + 1}`), () => {
                    if (isTextField(element)) element.focus()
                    else element.click()
                }, other)
            }
        } else if (typeof target === 'string') {
            this.bind(shortcuts, () => window.location.href = target, description)
        } else if (target !== undefined && target !== null) {
            throw new Error('Invalid target type')
        }
    }

    apply(fn) {
        fn(this) // Kotlin moment
    }
}

rrh.shortcuts = {
    handler: new ShortcutHandler(window),
    groups: [],

    // Creates or finds an existing group by name.
    group(name, element = null) {
        let group = this.groups.find(group => group.name === name)
        if (!group) {
            group = { name, group: new ShortcutGroup(this.handler, element) }
            this.groups.push(group)
        }
        return group.group
    },
}

function openHelp() {
    if ($('.shortcuts-help')) return

    document.body.overflow = 'hidden'
    let prevActiveElement = document.activeElement

    let backdrop = rrh.html`<div class="dialog-backdrop"></div>`
    backdrop.onclick = close
    document.body.appendChild(backdrop)

    let dialog = rrh.html`
        <div class="dialog shortcuts-help" tabindex="0">
            <div class="dialog__header">
                <h1 class="dialog__title">${rrh.l10n('List of shortcuts')}</h1>
                <button class="dialog__close-button" aria-label="${rrh.l10n('Close this dialog')}"></button>
            </div>
            <div class="dialog__content"></div>
        </div>
    `
    dialog.querySelector('.dialog__close-button').onclick = () => this.close()
    dialog.onkeydown = event => event.key === 'Escape' && close()
    document.body.appendChild(dialog)
    dialog.focus()

    function close() {
        document.body.overflow = ''
        document.body.removeChild(backdrop)
        document.body.removeChild(dialog)
        if (prevActiveElement) prevActiveElement.focus()
    }

    for (let { name, group } of rrh.shortcuts.groups) {
        if (group.shortcuts.length === 0) continue
        dialog.querySelector('.dialog__content').appendChild(rrh.html`
            <div class="shortcuts-group">
                <h2 class="shortcuts-group-heading">${name}</h2>
                <ul class="shortcuts-list">
                    ${group.shortcuts.map(({ description, shortcuts }) => `
                        <li class="shortcut-row">
                            <div class="shortcut-row__description">${description}</div>
                            <div class="shortcut-row__keys">
                                ${shortcuts.map(s => s.trim()
            .split(' ')
            .map(prettifyShortcut)
            .join(' '))
            .join(' <span class="kbd-or">or</span> ')}
                            </div>
                        </li>
                    `)}
                </ul>
            </div>
        `)
    }
}

rrh.l10nify(rrh.shortcuts).group('Common').apply(common => {
    common.bind('g', $$('.top-bar__highlight-link'), 'First 9 header links')
    common.bind('g h', '/', 'Home')
    common.bind('g l', '/list/', 'List of hyphae')
    common.bind('g r', '/recent-changes/', 'Recent changes')
    common.bind('g u', $('.auth-links__user-link'), 'Your profile′s hypha')
    common.bind(['?', isMac ? 'Meta+/' : 'Ctrl+/'], openHelp, 'Shortcut help')
})

if (document.body.dataset.rrhAddr.startsWith('/hypha')) {
    rrh.l10nify(rrh.shortcuts).group('Hypha').apply(hypha => {
        hypha.bind('', $$('article .wikilink'), 'First 9 hypha′s links')
        hypha.bind(['p', 'Alt+ArrowLeft', 'Ctrl+Alt+ArrowLeft'], $('.prevnext__prev'), 'Next hypha')
        hypha.bind(['n', 'Alt+ArrowRight', 'Ctrl+Alt+ArrowRight'], $('.prevnext__next'), 'Previous hypha')
        hypha.bind(['s', 'Alt+ArrowUp', 'Ctrl+Alt+ArrowUp'], $$('.navi-title a').slice(1, -1).slice(-1)[0], 'Parent hypha')
        hypha.bind(['c', 'Alt+ArrowDown', 'Ctrl+Alt+ArrowDown'], $('.subhyphae__link'), 'First child hypha')
        hypha.bind(['e', isMac ? 'Meta+Enter' : 'Ctrl+Enter'], $('.btn__link_navititle[href^="/edit/"]'), 'Edit this hypha')
        hypha.bind('v', $('.hypha-info__link[href^="/hypha/"]'), 'Go to hypha′s page')
        hypha.bind('a', $('.hypha-info__link[href^="/media/"]'), 'Go to media management')
        hypha.bind('h', $('.hypha-info__link[href^="/history/"]'), 'Go to history')
        hypha.bind('r', $('.hypha-info__link[href^="/rename/"]'), 'Rename this hypha')
        hypha.bind('b', $('.hypha-info__link[href^="/backlinks/"]'), 'Backlinks')
    })
}

if (document.body.dataset.rrhAddr.startsWith('/edit')) {
    rrh.l10nify(rrh.shortcuts).group('Editor').apply(editor => {
        editor.bind(isMac ? 'Meta+Enter' : 'Ctrl+Enter', $('.edit-form__save'), 'Save changes')
        editor.bind(isMac ? 'Meta+Shift+Enter' : 'Ctrl+Shift+Enter', $('.edit-form__preview'), 'Preview changes')
    })

    if (editTextarea) {
        rrh.l10nify(rrh.shortcuts).group('Format', editTextarea).apply(format => {
            format.bind(isMac ? 'Meta+b' : 'Ctrl+b', wrapBold, 'Bold', { force: true })
            format.bind(isMac ? 'Meta+i' : 'Ctrl+i', wrapItalic, 'Italic', { force: true })
            format.bind(isMac ? 'Meta+Shift+m' : 'Ctrl+M', wrapMonospace, 'Monospaced', { force: true })
            format.bind(isMac ? 'Meta+Shift+i' : 'Ctrl+I', wrapHighlighted, 'Highlight', { force: true })
            format.bind(isMac ? 'Meta+.' : 'Ctrl+.', wrapLifted, 'Superscript', { force: true })
            format.bind(isMac ? 'Meta+Comma' : 'Ctrl+Comma', wrapLowered, 'Subscript', { force: true })
            format.bind(isMac ? 'Meta+Shift+x' : 'Ctrl+X', wrapStrikethrough, 'Strikethrough', { force: true })
            format.bind(isMac ? 'Meta+k' : 'Ctrl+k', wrapLink, 'Inline link', { force: true })
            // Apparently, ⌘; conflicts with a Safari's hotkey. Whatever.
            format.bind(isMac ? 'Meta+;' : 'Ctrl+;', insertDate, 'Insert date UTC', { force: true })
        })
    }
}
