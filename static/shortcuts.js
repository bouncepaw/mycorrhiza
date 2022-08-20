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
let shortcutsGroup = null

class ShortcutHandler {
    constructor(element, override, filter = () => true) {
        this.element = element
        this.map = {}
        this.active = this.map
        this.override = override
        this.filter = filter
        this.timeout = null

        this.handleKeyDown = this.handleKeyDown.bind(this)
        this.resetActive = this.resetActive.bind(this)
        this.addEventListeners()
    }

    addEventListeners() {
        this.element.addEventListener('keydown', this.handleKeyDown)
    }

    add(text, action, description = null, shownInHelp = true) {
        let shortcuts = text.trim().split(',').map(shortcut => shortcut.trim().split(' '))

        if (shortcutsGroup && shownInHelp) {
            shortcutsGroup.push({
                action,
                shortcut: text,
                description,
            })
        }

        for (let shortcut of shortcuts) {
            let node = this.map
            for (let key of shortcut) {
                if (!node[key]) {
                    node[key] = {}
                }
                node = node[key]
                if (node.action) {
                    delete node.action
                    delete node.shortcut
                    delete node.description
                }
            }

            node.action = action
            node.shortcut = shortcut
            node.description = description
        }
    }

    group(...args) {
        if (typeof args[0] === 'string') this.fakeItem(args.shift())
        shortcutsGroup = []

        args[0].bind(this)()

        if (shortcutsGroup && shortcutsGroup.length) allShortcuts.push(shortcutsGroup)
        shortcutsGroup = null
    }

    bindElement(shortcut, element, ...other) {
        element = typeof element === 'string' ? $(element) : element
        if (!element) return
        this.add(shortcut, () => {
            if (isTextField(element)) {
                element.focus()
            } else {
                element.click()
            }
        }, ...other)
    }

    bindLink(shortcut, link, ...other) {
        this.add(shortcut, () => window.location.href = link, ...other)
    }

    bindCollection(prefix, elements, collectionDescription, itemDescription) {
        this.fakeItem(prefix + ' 1 – 9', collectionDescription)

        if (typeof elements === 'string') {
            elements = $$(elements)
        } else if (Array.isArray(elements)) {
            elements = elements.map(el => typeof el === 'string' ? $(el) : el)
        }

        for (let i = 1; i <= elements.length && i < 10; i++) {
            this.bindElement(`${prefix} ${i}`, elements[i - 1], `${itemDescription} #${i}`, false)
        }
    }

    fakeItem(shortcut, description = null) {
        let list = shortcutsGroup || allShortcuts
        list.push({
            shortcut: description ? shortcut : null,
            description: description || shortcut,
        })
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
        this.timeout = window.setTimeout(this.resetActive, 1500)
    }

    resetActive() {
        this.active = this.map
        if (this.timeout) {
            clearTimeout(this.timeout)
            this.timeout = null
        }
    }
}

const l10n = s => s

class ShortcutsHelpDialog {
    constructor() {
        this.backdrop = rrh.html`<div class="dialog-backdrop" hidden></div>`
        this.backdrop.onclick = () => this.close()
        document.body.appendChild(this.backdrop)

        this.dialog = rrh.html`
            <div class="dialog shortcuts-help" tabindex="0" hidden>
                <div class="dialog__header">
                    <h1 class="dialog__title">${l10n('List of shortcuts')}</h1>
                    <button class="dialog__close-button" aria-label="${l10n('Close this dialog')}"></button>
                </div>
                <div class="dialog__content"></div>
            </div>
        `
        this.dialog.querySelector('.dialog__close-button').onclick = () => this.close()
        document.body.appendChild(this.dialog)

        this.shortcuts = new ShortcutHandler(this.dialog, false)
        this.shortcuts.add('Escape', () => this.close(), null, false)

        let shortcutsGroup
        for (let item of allShortcuts) {
            if (item.description && !item.shortcut) {
                shortcutsGroup = rrh.html`
                    <div class="shortcuts-group">
                        <h2 class="shortcuts-group-heading">${item.description}</h2>
                    </div>
                `
                this.dialog.querySelector('.dialog__content').appendChild(shortcutsGroup)
            } else if (shortcutsGroup) {
                shortcutsGroup.appendChild(rrh.html`
                    <ul class="shortcuts-list">
                        ${item.map(({ description, shortcut }) => `
                            <li class="shortcut-row">
                                <div class="shortcut-row__description">${description}</div>
                                <div class="shortcut-row__keys">
                                    ${shortcut.split(',')
                                        .map(shortcuts => shortcuts
                                            .trim()
                                            .split(' ')
                                            .map(prettifyShortcut)
                                            .join(' '))
                                        .join(' <span class="kbd-or">or</span> ')}
                                </div>
                            </li>
                        `)}
                    </ul>
                `)
            }
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

window.addEventListener('load', () => {
    let helpDialog = null
    let openHelp = () => {
        if (!helpDialog) helpDialog = new ShortcutsHelpDialog()
        helpDialog.open()
    }

    let onEditPage = typeof editTextarea !== 'undefined'

    // Global shortcuts work everywhere.
    let globalShortcuts = new ShortcutHandler(document, false)
    globalShortcuts.add(isMac ? 'Meta+/' : 'Ctrl+/', openHelp)

    // Page shortcuts work everywhere except on text fields.
    let pageShortcuts = new ShortcutHandler(document, false, notTextField)
    pageShortcuts.add('?', openHelp, null, false)

    // Common shortcuts
    pageShortcuts.group('Common', function () {
        this.bindCollection('g', '.top-bar__highlight-link', 'First 9 header links', 'Header link')
        this.bindLink('g h', '/', 'Home')
        this.bindLink('g l', '/list/', 'List of hyphae')
        this.bindLink('g r', '/recent-changes/', 'Recent changes')
        this.bindElement('g u', '.auth-links__user-link', 'Your profile′s hypha')
    })

    if (!onEditPage) {
        // Hypha shortcuts
        pageShortcuts.group('Hypha', function () {
            this.bindCollection('', 'article .wikilink', 'First 9 hypha′s links', 'Hypha link')
            this.bindElement('p, Alt+ArrowLeft, Ctrl+Alt+ArrowLeft', '.prevnext__prev', 'Next hypha')
            this.bindElement('n, Alt+ArrowRight, Ctrl+Alt+ArrowRight', '.prevnext__next', 'Previous hypha')
            this.bindElement('s, Alt+ArrowUp, Ctrl+Alt+ArrowUp', $$('.navi-title a').slice(1, -1).slice(-1)[0], 'Parent hypha')
            this.bindElement('c, Alt+ArrowDown, Ctrl+Alt+ArrowDown', '.subhyphae__link', 'First child hypha')

            this.bindElement('e, ' + (isMac ? 'Meta+Enter' : 'Ctrl+Enter'), '.btn__link_navititle[href^="/edit/"]', 'Edit this hypha')
            this.bindElement('v', '.hypha-info__link[href^="/hypha/"]', 'Go to hypha′s page')
            this.bindElement('a', '.hypha-info__link[href^="/media/"]', 'Go to media management')
            this.bindElement('h', '.hypha-info__link[href^="/history/"]', 'Go to history')
            this.bindElement('r', '.hypha-info__link[href^="/rename/"]', 'Rename this hypha')
            this.bindElement('b', '.hypha-info__link[href^="/backlinks/"]', 'Backlinks')
        })

    } else {
        // Hypha editor shortcuts. These work only on editor's text area.
        let editorShortcuts = new ShortcutHandler(editTextarea, true)

        let shortcuts = [
            // Win+Linux    Mac                  Action              Description
            ['Ctrl+b', 'Meta+b', wrapBold, 'Format: Bold'],
            ['Ctrl+i', 'Meta+i', wrapItalic, 'Format: Italic'],
            ['Ctrl+M', 'Meta+Shift+m', wrapMonospace, 'Format: Monospaced'],
            ['Ctrl+I', 'Meta+Shift+i', wrapHighlighted, 'Format: Highlight'],
            ['Ctrl+.', 'Meta+.', wrapLifted, 'Format: Superscript'],
            ['Ctrl+Comma', 'Meta+Comma', wrapLowered, 'Format: Subscript'],
            ['Ctrl+X', 'Meta+Shift+x', wrapStrikethrough, 'Format: Strikethrough'],
            ['Ctrl+k', 'Meta+k', wrapLink, 'Format: Inline link'],
            // Apparently, ⌘; conflicts with a Safari's hotkey. Whatever.
            ['Ctrl+;', 'Meta+;', insertDate, 'Insert date UTC'],
        ]

        editorShortcuts.group('Editor', function () {
            for (let shortcut of shortcuts) {
                if (isMac) {
                    this.add(shortcut[1], ...shortcut.slice(2))
                } else {
                    this.add(shortcut[0], ...shortcut.slice(2))
                }
            }
        })

        editorShortcuts.group(function () {
            this.bindElement(isMac ? 'Meta+Enter' : 'Ctrl+Enter', $('.edit-form__save'), 'Save changes')
            this.bindElement(isMac ? 'Meta+Shift+Enter' : 'Ctrl+Shift+Enter', $('.edit-form__preview'), 'Preview changes')
        })
    }
})
