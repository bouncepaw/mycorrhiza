const $ = document.querySelector.bind(document);
const $$ = document.querySelectorAll.bind(document);

function keyEventToShortcut(event) {
    let elideShift = event.key.toUpperCase() === event.key && event.shiftKey;
    return (event.ctrlKey ? 'Ctrl+' : '') +
        (event.altKey ? 'Alt+' : '') +
        (event.metaKey ? 'Meta+' : '') +
        (!elideShift && event.shiftKey ? 'Shift+' : '') +
        event.key;
}

function isTextField(element) {
    let name = element.nodeName.toLowerCase();
    return name === 'textarea' ||
        name === 'select' ||
        (name === 'input' && !['submit', 'reset', 'checkbox', 'radio'].includes(element.type)) ||
        element.isContentEditable;
}

class ShortcutHandler {
    constructor(element, filter = () => {}) {
        this.element = element;
        this.map = {};
        this.active = this.map;
        this.filter = filter;
        this.timeout = null;

        this.handleKeyDown = this.handleKeyDown.bind(this);
        this.resetActive = this.resetActive.bind(this);
        this.addEventListeners();
    }

    addEventListeners() {
        this.element.addEventListener('keydown', this.handleKeyDown);
    }

    add(text, action) {
        let shortcuts = text.split(',').map(shortcut => shortcut.trim().split(' '));

        for (let shortcut of shortcuts) {
            let node = this.map;
            for (let key of shortcut) {
                if (!node[key]) {
                    node[key] = {};
                }
                node = node[key];
                if (node.action) {
                    delete node.action;
                }
            }

            node.action = action;
        }
    }

    handleKeyDown(event) {
        if (event.defaultPrevented) return;
        if (['Control', 'Alt', 'Shift', 'Meta'].includes(event.key)) return;
        if (!this.filter(event)) return;

        let shortcut = keyEventToShortcut(event);

        if (!this.active[shortcut]) {
            this.resetActive();
            return;
        }

        this.active = this.active[shortcut];
        if (this.active.action) {
            this.active.action(event);
            this.resetActive();
            return;
        }

        if (this.timeout) clearTimeout(this.timeout);
        this.timeout = window.setTimeout(this.resetActive, 1500);
    }

    resetActive() {
        this.active = this.map;
        if (this.timeout) {
            clearTimeout(this.timeout)
            this.timeout = null;
        }
    }
}

let notFormField = event => !(event.target instanceof Node && isTextField(event.target));
let globalShortcuts = new ShortcutHandler(document, notFormField);

globalShortcuts.add('p', () => alert('hello p'));
globalShortcuts.add('h', () => alert('hi h!'));
globalShortcuts.add('g h', () => alert('hi g h!!!'));