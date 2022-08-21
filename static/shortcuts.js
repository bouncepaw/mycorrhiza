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

let notTextField = event => !(event.target instanceof Node && isTextField(event.target))

let allShortcuts = []

class ShortcutHandler {
    constructor(element, override, filter = () => true) {
        this.element = element

        // This shit is built like this: map is the whole shortcut graph. active
        // points to the node of the shortcut graph that's currently "selected".
        // When the user presses a key, the code finds the key node in the
        // active subgraph and sets active to the found node. On each key press
        // we get a narrower view until we reach a leaf and execute the action.
        this.map = {}
        this.active = this.map

        // This is the listing of all shortcuts for the help dialog.
        this.shortcuts = []

        this.override = override
        this.filter = filter
        this.timeout = null

        this.element.addEventListener('keydown', event => this.handleKeyDown(event))
    }

    register(shortcuts, action) {
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
            leaf.shortcut = chord
        }
    }

    add(shortcuts, action, description = null) {
        shortcuts = strToArr(shortcuts)
        this.register(shortcuts, action)
        this.shortcuts.push({ action, shortcuts, description })
    }

    bindElement(shortcut, element, description = null) {
        element = typeof element === 'string' ? $(element) : element
        if (!element) return
        this.add(shortcut, () => {
            if (isTextField(element)) {
                element.focus()
            } else {
                element.click()
            }
        }, description)
    }

    bindLink(shortcut, link, ...other) {
        this.add(shortcut, () => window.location.href = link, ...other)
    }

    bindCollection(trigger, elements, description) {
        this.shortcuts.push({
            shortcuts: [trigger + ' 1 — 9'],
            description,
        })

        if (typeof elements === 'string') {
            elements = $$(elements)
        } else if (Array.isArray(elements)) {
            elements = elements.map(el => typeof el === 'string' ? $(el) : el)
        }

        for (let i = 1; i <= elements.length && i < 10; i++) {
            let element = elements[i - 1]
            this.register(`${trigger} ${i}`, () => {
                if (isTextField(element)) {
                    element.focus()
                } else {
                    element.click()
                }
            })
        }
    }

    handleKeyDown(event) {
        if (event.defaultPrevented) return
        if (['Control', 'Alt', 'Shift', 'Meta'].includes(event.key)) return
        if (!this.filter(event)) return

        let shortcut = keyEventToShortcut(event)

        if (!this.active[shortcut]) {
            this.resetActive()
            return
        }

        this.active = this.active[shortcut]
        if (this.active.action) {
            event.stopPropagation()
            this.active.action(event)
            if (this.override) event.preventDefault()
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

rrh.shortcuts = {
    // Global shortcuts work everywhere.
    global: new ShortcutHandler(document, false),

    groupMap: {}, // this is used for name lookups
    groups: [], // this is used in the help dialog
    group(name, element = window, filter = () => true) {
        if (!this.groupMap[name]) {
            let group = new ShortcutHandler(element, false, filter)
            this.groupMap[name] = group
            this.groups.push({ name, group })
            return group
        }
        return this.groupMap[name]
    },
}

rrh.l10n('List of shortcuts', { ru: 'Горячие клавиши' })
rrh.l10n('Close this dialog', { ru: 'Закрыть диалог' })

class ShortcutsHelpDialog {
    constructor() {
        this.backdrop = rrh.html`<div class="dialog-backdrop" hidden></div>`
        this.backdrop.onclick = () => this.close()
        document.body.appendChild(this.backdrop)

        this.dialog = rrh.html`
            <div class="dialog shortcuts-help" tabindex="0" hidden>
                <div class="dialog__header">
                    <h1 class="dialog__title">${rrh.l10n('List of shortcuts')}</h1>
                    <button class="dialog__close-button" aria-label="${rrh.l10n('Close this dialog')}"></button>
                </div>
                <div class="dialog__content"></div>
            </div>
        `
        this.dialog.querySelector('.dialog__close-button').onclick = () => this.close()
        this.dialog.onkeydown = event => event.key === 'Escape' && this.close()
        document.body.appendChild(this.dialog)

        for (let { name, group } of rrh.shortcuts.groups) {
            if (group.shortcuts.length === 0) continue
            this.dialog.querySelector('.dialog__content').appendChild(rrh.html`
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

    open() {
        this.prevActiveElement = document.activeElement

        document.body.overflow = 'hidden'
        this.backdrop.hidden = false
        this.dialog.hidden = false
        this.dialog.focus()
    }

    close() {
        document.body.overflow = ''
        this.backdrop.hidden = true
        this.dialog.hidden = true

        if (this.prevActiveElement) {
            this.prevActiveElement.focus()
            this.prevActiveElement = null
        }
    }
}

rrh.l10n('Common', { ru: 'Общее' })
rrh.l10n('Hypha', { ru: 'Гифа' })
rrh.l10n('Editor', { ru: 'Редактор' })
rrh.l10n('Format', { ru: 'Форматирование' })

let helpDialog = null
let openHelp = () => {
    if (!helpDialog) helpDialog = new ShortcutsHelpDialog()
    helpDialog.open()
}

rrh.shortcuts.global.add(isMac ? 'Meta+/' : 'Ctrl+/', openHelp)

let common = rrh.l10nify(rrh.shortcuts.group('Common', window, notTextField))
common.bindCollection('g', '.top-bar__highlight-link', 'First 9 header links', 'Header link')
common.bindLink('g h', '/', 'Home')
common.bindLink('g l', '/list/', 'List of hyphae')
common.bindLink('g r', '/recent-changes/', 'Recent changes')
common.bindElement('g u', '.auth-links__user-link', 'Your profile′s hypha')
common.add('?', openHelp, rrh.l10n('Shortcut help'))

if (document.body.dataset.rrhAddr.startsWith('/hypha')) {
    let hypha = rrh.shortcuts.group('Hypha', window, notTextField)
    hypha.bindCollection('', 'article .wikilink', 'First 9 hypha′s links')
    hypha.bindElement(['p', 'Alt+ArrowLeft', 'Ctrl+Alt+ArrowLeft'], '.prevnext__prev', 'Next hypha')
    hypha.bindElement(['n', 'Alt+ArrowRight', 'Ctrl+Alt+ArrowRight'], '.prevnext__next', 'Previous hypha')
    hypha.bindElement(['s', 'Alt+ArrowUp', 'Ctrl+Alt+ArrowUp'], $$('.navi-title a').slice(1, -1).slice(-1)[0], 'Parent hypha')
    hypha.bindElement(['c', 'Alt+ArrowDown', 'Ctrl+Alt+ArrowDown'], '.subhyphae__link', 'First child hypha')
    hypha.bindElement(['e', isMac ? 'Meta+Enter' : 'Ctrl+Enter'], '.btn__link_navititle[href^="/edit/"]', 'Edit this hypha')
    hypha.bindElement('v', '.hypha-info__link[href^="/hypha/"]', 'Go to hypha′s page')
    hypha.bindElement('a', '.hypha-info__link[href^="/media/"]', 'Go to media management')
    hypha.bindElement('h', '.hypha-info__link[href^="/history/"]', 'Go to history')
    hypha.bindElement('r', '.hypha-info__link[href^="/rename/"]', 'Rename this hypha')
    hypha.bindElement('b', '.hypha-info__link[href^="/backlinks/"]', 'Backlinks')
}

if (document.body.dataset.rrhAddr.startsWith('/edit') && editTextarea) {
    let editor = rrh.shortcuts.group('Editor', window)
    editor.bindElement(isMac ? 'Meta+Enter' : 'Ctrl+Enter', $('.edit-form__save'), 'Save changes')
    editor.bindElement(isMac ? 'Meta+Shift+Enter' : 'Ctrl+Shift+Enter', $('.edit-form__preview'), 'Preview changes')

    let format = rrh.shortcuts.group('Format', editTextarea)
    format.override = true
    format.add(isMac ? 'Meta+b' : 'Ctrl+b', wrapBold, 'Bold')
    format.add(isMac ? 'Meta+i' : 'Ctrl+i', wrapItalic, 'Italic')
    format.add(isMac ? 'Meta+Shift+m' : 'Ctrl+M', wrapMonospace, 'Monospaced')
    format.add(isMac ? 'Meta+Shift+i' : 'Ctrl+I', wrapHighlighted, 'Highlight')
    format.add(isMac ? 'Meta+.' : 'Ctrl+.', wrapLifted, 'Superscript')
    format.add(isMac ? 'Meta+Comma' : 'Ctrl+Comma', wrapLowered, 'Subscript')
    format.add(isMac ? 'Meta+Shift+x' : 'Ctrl+X', wrapStrikethrough, 'Strikethrough')
    format.add(isMac ? 'Meta+k' : 'Ctrl+k', wrapLink, 'Inline link')
    // Apparently, ⌘; conflicts with a Safari's hotkey. Whatever.
    format.add(isMac ? 'Meta+;' : 'Ctrl+;', insertDate, 'Insert date UTC')
}
