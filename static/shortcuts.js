rrh.l10n('List of shortcuts', { ru: 'Горячие клавиши' })
rrh.l10n('Close this dialog', { ru: 'Закрыть диалог' })

rrh.l10n('Common', { ru: 'Общее' })
rrh.l10n('Home', { ru: 'Главная' })

rrh.l10n('Hypha', { ru: 'Гифа' })

rrh.l10n('Editor', { ru: 'Редактор' })

rrh.l10n('Format', { ru: 'Форматирование' })

rrh.shortcuts = {
    // map is the whole shortcut graph. active points to the node of the
    // shortcut graph that's currently "selected". When the user presses a
    // key, the code finds the key node in the active subgraph and sets
    // active to the found node. On each key press we get a narrower view
    // until we reach a leaf and execute the action. View this property in
    // the JavaScript console of your browser for easier understanding.
    map: {},
    groups: [],

    addGroup(group) {
        this.groups.push(group)
    },

    addBindingToGroup(name, binding) {
        let group = this.groups.find(group => group.name === name)
        if (!group) {
            console.warn('Shortcut group', name, 'not found')
            return
        }
        group.bind(binding)
    },

    _register(shortcuts, action, other = {}) {
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
    },

    _handleKeyDown(event) {
        if (event.defaultPrevented) return
        if (['Control', 'Alt', 'Shift', 'Meta'].includes(event.key)) return
        if ((!event.ctrlKey && !event.metaKey && !event.shiftKey && !event.altKey) &&
            event.target instanceof Node && isTextField(event.target)) return

        let shortcut = keyEventToShortcut(event)

        if (!this.active[shortcut]) {
            this._resetActive()
            return
        }

        this.active = this.active[shortcut]
        if (this.active.action && (!this.active.element || event.target === this.active.element)) {
            event.stopPropagation()
            if (this.active.force) event.preventDefault()
            this.active.action(event)
            this._resetActive()
            return
        }

        if (this.timeout) clearTimeout(this.timeout)
        this.timeout = window.setTimeout(() => this._resetActive(), 1500)
    },

    _resetActive() {
        this.active = this.map
        if (this.timeout) {
            clearTimeout(this.timeout)
            this.timeout = null
        }
    },
}
window.addEventListener('keydown', event => rrh.shortcuts._handleKeyDown(event))
rrh.shortcuts._resetActive()

// Convert a KeyboardEvent into a shortcut string for matching by
// ShortcutHandler.
function keyEventToShortcut(event) {
    let elideShift = event.key.toUpperCase() === event.key && event.shiftKey
    return (event.ctrlKey ? 'Ctrl+' : '') +
        (event.altKey ? 'Alt+' : '') +
        (event.metaKey ? 'Meta+' : '') +
        (!elideShift && event.shiftKey ? 'Shift+' : '') +
        (event.key === ',' ? 'Comma' : event.key === ' ' ? 'Space' : event.key)
}

// Prettify the shortcut string by replacing modifiers and arrow codes with
// Unicode symbol for presentation to the user.
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

    // Add Shift into shortcut strings like Ctrl+L
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
            // Make every key uppercase. This does not introduce any ambiguous
            // cases because we insert a Shift modifier for upper-case keys
            // earlier.
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

class Shortcut {
    constructor(shortcuts, target, description, other) {
        this.shortcuts = shortcuts
        this.target = target
        this.description = rrh.l10n(description)
        Object.assign(this, other)
    }
}

class ShortcutGroup {
    constructor(name, element = null, bindings = []) {
        this.name = rrh.l10n(name)
        this.shortcuts = []
        this.element = element
        bindings.forEach(binding => this.bind(binding))
    }

    bind({ shortcuts, target, description, ...other }) {
        shortcuts = strToArr(shortcuts).map(s => s.trim())
        if (!other.element) other.element = this.element
        if (target instanceof Function) {
            this.shortcuts.push({ shortcuts, description })
            rrh.shortcuts._register(shortcuts, target, other)
        } else if (target instanceof Node) {
            this.bind({
                shortcuts,
                target: () => {
                    if (isTextField(target)) target.focus()
                    target.click()
                },
                description, ...other,
            })
        } else if (Array.isArray(target) && (target.length === 0 || target[0] instanceof Node)) {
            this.shortcuts.push({
                shortcuts: shortcuts.map(s => `${s} 1 — 9`),
                description,
            })
            for (let i = 0; i < target.length && i < 9; i++) {
                let element = target[i]
                rrh.shortcuts._register(shortcuts.map(s => `${s} ${i + 1}`), () => {
                    if (isTextField(element)) element.focus()
                    else element.click()
                }, other)
            }
        } else if (typeof target === 'string') {
            this.bind({
                shortcuts,
                target: () => window.location.href = target,
                description,
                ...other,
            })
        } else if (target !== undefined && target !== null) {
            throw new Error('Invalid target type')
        }
    }
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
    dialog.querySelector('.dialog__close-button').onclick = close
    dialog.onkeydown = event => {
        if (event.key === 'Escape') close()
    }
    document.body.appendChild(dialog)
    dialog.focus()

    function close() {
        document.body.overflow = ''
        document.body.removeChild(backdrop)
        document.body.removeChild(dialog)
        if (prevActiveElement) prevActiveElement.focus()
    }

    function formatShortcuts(shortcuts) {
        return shortcuts.map(s => s.split(' ')
            .map(prettifyShortcut).join(' '))
            .join(' <span class="kbd-or">or</span> ')
    }

    for (let group of rrh.shortcuts.groups) {
        if (group.shortcuts.length === 0) continue
        dialog.querySelector('.dialog__content').appendChild(rrh.html`
            <div class="shortcuts-group">
                <h2 class="shortcuts-group-heading">${group.name}</h2>
                <ul class="shortcuts-list">
                    ${group.shortcuts.map(({ description, shortcuts }) => `
                        <li class="shortcut-row">
                            <div class="shortcut-row__description">${description}</div>
                            <div class="shortcut-row__keys">
                                ${formatShortcuts(shortcuts)}
                            </div>
                        </li>
                    `)}
                </ul>
            </div>
        `)
    }
}

rrh.shortcuts.addGroup(new ShortcutGroup('Common', null, [
    new Shortcut('g', $$('.top-bar__highlight-link'), 'First 9 header links'),
    new Shortcut('g h', '/', 'Home'),
    new Shortcut('g l', '/list/', 'List of hyphae'),
    new Shortcut('g r', '/recent-changes/', 'Recent changes'),
    new Shortcut('g u', $('.auth-links__user-link'), 'Your profile′s hypha'),
    new Shortcut(['?', isMac ? 'Meta+/' : 'Ctrl+/'], openHelp, 'Shortcut help'),
]))

if (document.body.dataset.rrhAddr.startsWith('/hypha')) {
    rrh.shortcuts.addGroup(new ShortcutGroup('Hypha', null, [
        new Shortcut('', $$('article .wikilink'), 'First 9 hypha′s links'),
        new Shortcut(['p', 'Alt+ArrowLeft', 'Ctrl+Alt+ArrowLeft'], $('.prevnext__prev'), 'Next hypha'),
        new Shortcut(['n', 'Alt+ArrowRight', 'Ctrl+Alt+ArrowRight'], $('.prevnext__next'), 'Previous hypha'),
        new Shortcut(['s', 'Alt+ArrowUp', 'Ctrl+Alt+ArrowUp'], $$('.navi-title a').slice(1, -1).slice(-1)[0], 'Parent hypha'),
        new Shortcut(['c', 'Alt+ArrowDown', 'Ctrl+Alt+ArrowDown'], $('.subhyphae__link'), 'First child hypha'),
        new Shortcut(['e', isMac ? 'Meta+Enter' : 'Ctrl+Enter'], $('.btn__link_navititle[href^="/edit/"]'), 'Edit this hypha'),
        new Shortcut('v', $('.hypha-info__link[href^="/hypha/"]'), 'Go to hypha′s page'),
        new Shortcut('a', $('.hypha-info__link[href^="/media/"]'), 'Go to media management'),
        new Shortcut('h', $('.hypha-info__link[href^="/history/"]'), 'Go to history'),
        new Shortcut('r', $('.hypha-info__link[href^="/rename/"]'), 'Rename this hypha'),
        new Shortcut('b', $('.hypha-info__link[href^="/backlinks/"]'), 'Backlinks'),
    ]))
}

if (document.body.dataset.rrhAddr.startsWith('/edit')) {
    rrh.shortcuts.addGroup(new ShortcutGroup('Editor', null, [
        new Shortcut(isMac ? 'Meta+Enter' : 'Ctrl+Enter', $('.edit-form__save'), 'Save changes'),
        new Shortcut(isMac ? 'Meta+Shift+Enter' : 'Ctrl+Shift+Enter', $('.edit-form__preview'), 'Preview changes'),
    ]))

    if (editTextarea) {
        rrh.shortcuts.addGroup(new ShortcutGroup('Format', null, [
            new Shortcut(isMac ? 'Meta+b' : 'Ctrl+b', wrapBold, 'Bold', { force: true }),
            new Shortcut(isMac ? 'Meta+i' : 'Ctrl+i', wrapItalic, 'Italic', { force: true }),
            new Shortcut(isMac ? 'Meta+Shift+m' : 'Ctrl+M', wrapMonospace, 'Monospaced', { force: true }),
            new Shortcut(isMac ? 'Meta+Shift+i' : 'Ctrl+I', wrapHighlighted, 'Highlight', { force: true }),
            new Shortcut(isMac ? 'Meta+.' : 'Ctrl+.', wrapLifted, 'Superscript', { force: true }),
            new Shortcut(isMac ? 'Meta+Comma' : 'Ctrl+Comma', wrapLowered, 'Subscript', { force: true }),
            new Shortcut(isMac ? 'Meta+Shift+x' : 'Ctrl+X', wrapStrikethrough, 'Strikethrough', { force: true }),
            new Shortcut(isMac ? 'Meta+k' : 'Ctrl+k', wrapLink, 'Inline link', { force: true }),
            // Apparently, ⌘; conflicts with a Safari's hotkey. Whatever.
            new Shortcut(isMac ? 'Meta+;' : 'Ctrl+;', insertDate, 'Insert date UTC', { force: true }),
        ]))
    }
}
