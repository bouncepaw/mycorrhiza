(() => {
    const $ = document.querySelector.bind(document);
    const $$ = (...args) => Array.prototype.slice.call(document.querySelectorAll(...args));

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

        add(text, action, description = null) {
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
                        delete node.shortcut;
                        delete node.description;
                    }
                }

                node.action = action;
                node.shortcut = shortcut;
                node.description = description;
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

    function bindElementFactory(handler) {
        return (shortcut, element, ...other) => {
            element = typeof element === 'string' ? $(element) : element;
            if (!element) return;
            handler.add(shortcut, () => {
                if (isTextField(element)) {
                    element.focus();
                } else {
                    element.click();
                }
            }, ...other);
        };
    }

    function bindLinkFactory(handler) {
        return (shortcut, link, ...other) => handler.add(shortcut, () => window.location.href = link, ...other);
    }

    window.addEventListener('load', () => {
        let notFormField = event => !(event.target instanceof Node && isTextField(event.target));
        let globalShortcuts = new ShortcutHandler(document, notFormField);

        let bindElement = bindElementFactory(globalShortcuts);
        let bindLink = bindLinkFactory(globalShortcuts);

        bindElement('p, Alt+ArrowLeft', '.prevnext__prev', 'Next hypha');
        bindElement('n, Alt+ArrowRight', '.prevnext__next', 'Previous hypha');
        bindElement('s, Alt+ArrowTop', $$('.navi-title a').slice(1, -1).slice(-1)[0], 'Parent hypha');

        bindLink('g h', '/', 'Home');
        bindLink('g l', '/list/', 'List of hyphae');
        bindLink('g r', '/recent-changes/', 'Recent changes');

        bindElement('g u', '.header-links__entry_user .header-links__link', 'Your profile′s hypha')

        let headerLinks = $$('.header-links__link');
        for (let i = 1; i <= headerLinks.length && i < 10; i++) {
            bindElement(`g ${i}`, headerLinks[i-1], `Header link #${i}`);
        }

        let hyphaLinks = $$('article .wikilink');
        for (let i = 1; i <= hyphaLinks.length && i < 10; i++) {
            bindElement(i.toString(), hyphaLinks[i-1], `Hypha link #${i}`);
        }
    });
})();