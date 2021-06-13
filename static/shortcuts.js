class Shortcut {
    // turns the given event into a string representation of it.
    static fromEvent(event) {
        let elideShift = event.key.toUpperCase() === event.key && event.shiftKey;
        return (event.ctrlKey ? 'Ctrl+' : '') +
            (event.altKey ? 'Alt+' : '') +
            (event.metaKey ? 'Meta+' : '') +
            (!elideShift && event.shiftKey ? 'Shift+' : '') +
            (event.key === ',' ? 'Comma' : event.key === ' ' ? 'Space' : event.key);
    }

    // Some keys look better with cool symbols instead of their long and boring names.
    static prettify(shortcut, isMac) {
        let keys = shortcut.split('+');

        // Uh it places the cmd sign before the letter to follow the Mac conventions, I guess.
        if (isMac) {
            let cmdIdx = keys.indexOf('Meta');
            if (cmdIdx !== -1 && keys.length - cmdIdx > 2) {
                let tmp = keys[cmdIdx + 1];
                keys[cmdIdx + 1] = 'Meta';
                keys[cmdIdx] = tmp;
            }
        }

        let lastKey = keys[keys.length - 1];
        // Uhh add Shift if the letter is uppercase??
        if (!keys.includes('Shift') && lastKey.toUpperCase() === lastKey && lastKey.toLowerCase() !== lastKey) {
            keys.splice(keys.length - 1, 0, 'Shift');
        }

        return keys.map((key, i) => {
            // If last element and there is more than one element and it's a letter
            if (i === keys.length - 1 && i > 0 && key.length === 1) {
                // Show in upper case. ⌘K looks better ⌘k, no doubt.
                key = key.toUpperCase();
            }

            return `<kbd>${Shortcut.symbolifyKey(key, isMac)}</kbd>`;
        }).join(isMac ? '' : ' + ');
    }

    static symbolifyKey(key, isMac) {
        if (isMac) {
            switch (key) {
                case 'Ctrl': return '⌃';
                case 'Alt': return '⌥';
                case 'Shift': return '⇧';
                case 'Meta': return '⌘';
            }
        }

        switch (key) {
            case 'ArrowLeft': return '←';
            case 'ArrowRight': return '→';
            case 'ArrowTop': return '↑';
            case 'ArrowBottom': return '↓';
            case 'Comma': return ',';
            case 'Enter': return '↩';
            case ' ': return 'Space';
        }
        return key
    }
}

(() => {
    const $ = document.querySelector.bind(document);
    const $$ = (...args) => Array.prototype.slice.call(document.querySelectorAll(...args));

    // Some things look different on Mac.
    // Note that the ⌘ command key is called Meta in JS for some reason.
    const isMac = /Macintosh/.test(window.navigator.userAgent);

    function isTextField(element) {
        let name = element.nodeName.toLowerCase();
        return name === 'textarea' ||
            name === 'select' ||
            (name === 'input' && !['submit', 'reset', 'checkbox', 'radio'].includes(element.type)) ||
            element.isContentEditable;
    }

    let notTextField = event => !(event.target instanceof Node && isTextField(event.target));

    // The whole shortcut table for current page. It is used for generating the dialog.
    let allShortcuts = [];
    // Temporary variable for building a shortcut group.
    let shortcutsGroup = null;

    // Advanced stuff.
    class ShortcutHandler {
        constructor(element, filter = () => true) {
            this.element = element;
            this.map = {};
            this.active = this.map;
            this.filter = filter;
            this.timeout = null;

            this.handleKeyDown = this.handleKeyDown.bind(this);
            this.resetActive = this.resetActive.bind(this);
            this.element.addEventListener('keydown', this.handleKeyDown);
        }

        add(text, action, description = null) {
            let shortcuts = text.split(',').map(shortcut => shortcut.trim().split(' '));

            if (shortcutsGroup) {
                shortcutsGroup.push({
                    action,
                    shortcut: text,
                    description,
                })
            }

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

        groupStart() {
            shortcutsGroup = [];
        }

        groupEnd() {
            if (shortcutsGroup && shortcutsGroup.length) allShortcuts.push(shortcutsGroup);
            shortcutsGroup = null;
        }

        // A dirty and shameful hack for inserting non-generated entries into the table.
        fakeItem(shortcut, description = null) {
            // So it's a boolean, right?
            let list = shortcutsGroup || allShortcuts;
            // And we push something into a boolean. I give up.
            list.push({
                shortcut: description ? shortcut : null,
                description: description || shortcut,
            });
        }

        handleKeyDown(event) {
            if (event.defaultPrevented) return;
            if (['Control', 'Alt', 'Shift', 'Meta'].includes(event.key)) return;
            if (!this.filter(event)) return;

            let shortcut = Shortcut.fromEvent(event);

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

    let prevActiveElement = null;
    let shortcutsListDialog = null;

    function openShortcutsReference() {
        if (!shortcutsListDialog) { // I guess the dialog is reused for second and subsequent invocations.
            let wrap = document.createElement('div');
            wrap.className = 'dialog-wrap';
            shortcutsListDialog = wrap;

            let dialog = document.createElement('div');
            dialog.className = 'dialog shortcuts-modal';
            dialog.tabIndex = 0;
            wrap.appendChild(dialog);

            let dialogHeader = document.createElement('div');
            dialogHeader.className = 'dialog__header';
            dialog.appendChild(dialogHeader);

            let title = document.createElement('h1');
            title.className = 'dialog__title';
            title.textContent = 'List of shortcuts';
            dialogHeader.appendChild(title);

            let closeButton = document.createElement('button');
            closeButton.className = 'dialog__close-button';
            closeButton.setAttribute('aria-label', 'Close this dialog'); // a11y gang
            dialogHeader.appendChild(closeButton);

            for (let item of allShortcuts) {
                if (item.description && !item.shortcut) {
                    let heading = document.createElement('h2');
                    heading.className = 'shortcuts-group-heading';
                    heading.textContent = item.description;
                    dialog.appendChild(heading);

                } else {
                    let list = document.createElement('ul');
                    list.className = 'shortcuts-group';

                    for (let shortcut of item) {
                        let listItem = document.createElement('li');
                        listItem.className = 'shortcut-row';
                        list.appendChild(listItem);

                        let descriptionColumn = document.createElement('div')
                        descriptionColumn.className = 'shortcut-row__description';
                        descriptionColumn.textContent = shortcut.description;
                        listItem.appendChild(descriptionColumn);

                        let shortcutColumn = document.createElement('div');
                        shortcutColumn.className = 'shortcut-row__keys';
                        shortcutColumn.innerHTML = shortcut.shortcut.split(',')
                            .map(shortcuts => shortcuts.trim().split(' ').map((sc) => Shortcut.prettify(sc, isMac)).join(' '))
                            .join(' or ');
                        listItem.appendChild(shortcutColumn);
                    }

                    dialog.appendChild(list);
                }
            }

            let handleClose = (event) => {
                event.preventDefault();
                event.stopPropagation();
                closeShortcutsReference();
            };

            let dialogShortcuts = new ShortcutHandler(dialog, notTextField);

            dialogShortcuts.add('Escape', handleClose);
            closeButton.addEventListener('click', handleClose);
            wrap.addEventListener('click', handleClose);

            dialog.addEventListener('click', event => event.stopPropagation());

            document.body.appendChild(wrap);
        }

        document.body.overflow = 'hidden';
        shortcutsListDialog.hidden = false;
        prevActiveElement = document.activeElement;
        shortcutsListDialog.children[0].focus();
    }

    function closeShortcutsReference() {
        if (shortcutsListDialog) {
            document.body.overflow = '';
            shortcutsListDialog.hidden = true;

            if (prevActiveElement) {
                prevActiveElement.focus();
                prevActiveElement = null;
            }
        }
    }

    window.addEventListener('load', () => {
        let globalShortcuts = new ShortcutHandler(document, notTextField);

        // Global shortcuts

        let bindElement = bindElementFactory(globalShortcuts);
        let bindLink = bindLinkFactory(globalShortcuts);

        // * Common shortcuts
        globalShortcuts.fakeItem('Common');

        // Nice indentation here
        globalShortcuts.groupStart();
            globalShortcuts.fakeItem('g 1 – 9', 'First 9 header links');
            bindLink('g h', '/', 'Home');
            bindLink('g l', '/list/', 'List of hyphae');
            bindLink('g r', '/recent-changes/', 'Recent changes');
            bindElement('g u', '.header-links__entry_user .header-links__link', 'Your profile′s hypha');
        globalShortcuts.groupEnd();

        let headerLinks = $$('.header-links__link');
        for (let i = 1; i <= headerLinks.length && i < 10; i++) {
            bindElement(`g ${i}`, headerLinks[i-1], `Header link #${i}`);
        }

        // * Hypha shortcuts
        if (typeof editTextarea === 'undefined') {
            globalShortcuts.fakeItem('Hypha');

            globalShortcuts.groupStart();
                globalShortcuts.fakeItem('1 – 9', 'First 9 hypha′s links');
                bindElement('p, Alt+ArrowLeft, Ctrl+Alt+ArrowLeft', '.prevnext__prev', 'Next hypha');
                bindElement('n, Alt+ArrowRight, Ctrl+Alt+ArrowRight', '.prevnext__next', 'Previous hypha');
                bindElement('s, Alt+ArrowTop, Ctrl+Alt+ArrowTop', $$('.navi-title a').slice(1, -1).slice(-1)[0], 'Parent hypha');
                bindElement('e, Ctrl+Enter', '.hypha-tabs__link[href^="/edit/"]', 'Edit this hypha');
            globalShortcuts.groupEnd();

            let hyphaLinks = $$('article .wikilink');
            for (let i = 1; i <= hyphaLinks.length && i < 10; i++) {
                bindElement(i.toString(), hyphaLinks[i-1], `Hypha link #${i}`);
            }
        }

        // * Editor shortcuts
        if (typeof editTextarea !== 'undefined') {
            let editorShortcuts = new ShortcutHandler(editTextarea);
            let bindElement = bindElementFactory(editorShortcuts);

            let shortcuts = [
                // Inspired by MS Word, Pages, Google Docs and Telegram desktop clients.
                // And by myself, too.

                // Win+Linux    Mac              Action              Description
                ['Ctrl+b',      'Meta+b',        wrapBold,           'Format: Bold'],
                ['Ctrl+i',      'Meta+i',        wrapItalic,         'Format: Italic'],
                ['Ctrl+M',      'Meta+Shift+m',  wrapMonospace,      'Format: Monospaced'],
                ['Ctrl+I',      'Meta+Shift+i',  wrapHighlighted,    'Format: Highlight'],
                ['Ctrl+.',      'Meta+.',        wrapLifted,         'Format: Superscript'],
                ['Ctrl+Comma',  'Meta+Comma',    wrapLowered,        'Format: Subscript'],
                // Strikethrough conflicts with 1Password on my machine but
                // I'm probably the only Mycorrhiza user who uses 1Password. -handlerug
                ['Ctrl+X',      'Meta+Shift+x',  wrapStrikethrough,  'Format: Strikethrough'],
                ['Ctrl+k',      'Meta+k',        wrapLink,           'Format: Link'],
            ];

            editorShortcuts.fakeItem('Editor');

            editorShortcuts.groupStart();
            for (let shortcut of shortcuts) {
                if (isMac) {
                    editorShortcuts.add(shortcut[1], ...shortcut.slice(2))
                } else {
                    editorShortcuts.add(shortcut[0], ...shortcut.slice(2))
                }
            }
            editorShortcuts.groupEnd();

            editorShortcuts.groupStart();
            bindElement(isMac ? 'Meta+Enter' : 'Ctrl+Enter', $('.edit-form__save'), 'Save changes');
            editorShortcuts.groupEnd();

            // Help shortcut
            editorShortcuts.add(isMac ? 'Meta+/' : 'Ctrl+/', openShortcutsReference);
        }

        // * Meta shortcuts
        globalShortcuts.add('?', openShortcutsReference);
    });
})();
