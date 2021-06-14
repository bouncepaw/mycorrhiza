(() => {
    const $ = document.querySelector.bind(document);
    const $$ = (...args) => Array.prototype.slice.call(document.querySelectorAll(...args));

    const isMac = /Macintosh/.test(window.navigator.userAgent);

    function keyEventToShortcut(event) {
        let elideShift = event.key.toUpperCase() === event.key && event.shiftKey;
        return (event.ctrlKey ? 'Ctrl+' : '') +
            (event.altKey ? 'Alt+' : '') +
            (event.metaKey ? 'Meta+' : '') +
            (!elideShift && event.shiftKey ? 'Shift+' : '') +
            (event.key === ',' ? 'Comma' : event.key === ' ' ? 'Space' : event.key);
    }

    function prettifyShortcut(shortcut) {
        let keys = shortcut.split('+');

        if (isMac) {
            let cmdIdx = keys.indexOf('Meta');
            if (cmdIdx !== -1 && keys.length - cmdIdx > 2) {
                let tmp = keys[cmdIdx + 1];
                keys[cmdIdx + 1] = 'Meta';
                keys[cmdIdx] = tmp;
            }
        }

        let lastKey = keys[keys.length - 1];
        if (!keys.includes('Shift') && lastKey.toUpperCase() === lastKey && lastKey.toLowerCase() !== lastKey) {
            keys.splice(keys.length - 1, 0, 'Shift');
        }

        for (let i = 0; i < keys.length; i++) {
            if (isMac) {
                switch (keys[i]) {
                    case 'Ctrl': keys[i] = '⌃'; break;
                    case 'Alt': keys[i] = '⌥'; break;
                    case 'Shift': keys[i] = '⇧'; break;
                    case 'Meta': keys[i] = '⌘'; break;
                }
            }

            if (i === keys.length - 1 && i > 0 && keys[i].length === 1) {
                keys[i] = keys[i].toUpperCase();
            }

            switch (keys[i]) {
                case 'ArrowLeft': keys[i] = '←'; break;
                case 'ArrowRight': keys[i] = '→'; break;
                case 'ArrowUp': keys[i] = '↑'; break;
                case 'ArrowDown': keys[i] = '↓'; break;
                case 'Comma': keys[i] = ','; break;
                case 'Enter': keys[i] = '↩'; break;
                case ' ': keys[i] = 'Space'; break;
            }

            keys[i] = `<kbd>${keys[i]}</kbd>`;
        }

        return keys.join(isMac ? '' : ' + ');
    }

    function isTextField(element) {
        let name = element.nodeName.toLowerCase();
        return name === 'textarea' ||
            name === 'select' ||
            (name === 'input' && !['submit', 'reset', 'checkbox', 'radio'].includes(element.type)) ||
            element.isContentEditable;
    }

    let notTextField = event => !(event.target instanceof Node && isTextField(event.target));

    let allShortcuts = [];
    let shortcutsGroup = null;

    class ShortcutHandler {
        constructor(element, override, filter = () => true) {
            this.element = element;
            this.map = {};
            this.active = this.map;
            this.override = override;
            this.filter = filter;
            this.timeout = null;

            this.handleKeyDown = this.handleKeyDown.bind(this);
            this.resetActive = this.resetActive.bind(this);
            this.addEventListeners();
        }

        addEventListeners() {
            this.element.addEventListener('keydown', this.handleKeyDown);
        }

        add(text, action, description = null, shownInHelp = true) {
            let shortcuts = text.trim().split(',').map(shortcut => shortcut.trim().split(' '));

            if (shortcutsGroup && shownInHelp) {
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

        groupStart(title = null) {
            if (title) this.fakeItem(title);
            shortcutsGroup = [];
        }

        groupEnd() {
            if (shortcutsGroup && shortcutsGroup.length) allShortcuts.push(shortcutsGroup);
            shortcutsGroup = null;
        }

        fakeItem(shortcut, description = null) {
            let list = shortcutsGroup || allShortcuts;
            list.push({
                shortcut: description ? shortcut : null,
                description: description || shortcut,
            });
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
                if (this.override) event.preventDefault();
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

    function bindCollectionFactory(handler) {
        return (prefix, elements, collectionDescription, itemDescription) => {
            handler.fakeItem(prefix + ' 1 – 9', collectionDescription);

            if (typeof elements === 'string') {
                elements = $$(elements);
            } else if (Array.isArray(elements)) {
                elements = elements.map(el => typeof el === 'string' ? $(el) : el);
            }

            for (let i = 1; i <= elements.length && i < 10; i++) {
                bindElementFactory(handler)(`${prefix} ${i}`, elements[i-1], `${itemDescription} #${i}`, false);
            }
        }
    }

    class ShortcutsHelpDialog {
        constructor() {
            let template = $('#dialog-template');
            this.wrap = template.content.firstElementChild.cloneNode(true);
            this.wrap.classList.add('shortcuts-help');
            this.wrap.hidden = true;

            let handleClose = (event) => {
                event.preventDefault();
                event.stopPropagation();
                this.close();
            };

            this.shortcuts = new ShortcutHandler(this.wrap, false);
            this.shortcuts.add('Escape', handleClose, null, false);

            this.wrap.querySelector('.dialog__title').textContent = 'List of shortcuts';
            this.wrap.querySelector('.dialog__close-button').addEventListener('click', handleClose);

            this.wrap.addEventListener('click', handleClose);
            this.wrap.querySelector('.dialog').addEventListener('click', event => event.stopPropagation());

            let shortcutsGroup;
            let shortcutsGroupTemplate = document.createElement('div');
            shortcutsGroupTemplate.className = 'shortcuts-group';

            for (let item of allShortcuts) {
                if (item.description && !item.shortcut) {
                    shortcutsGroup = shortcutsGroupTemplate.cloneNode();
                    this.wrap.querySelector('.dialog__content').appendChild(shortcutsGroup);

                    let heading = document.createElement('h2');
                    heading.className = 'shortcuts-group-heading';
                    heading.textContent = item.description;
                    shortcutsGroup.appendChild(heading);

                } else {
                    let list = document.createElement('ul');
                    list.className = 'shortcuts-list';

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
                            .map(shortcuts => shortcuts.trim().split(' ').map(prettifyShortcut).join(' '))
                            .join(' <span class="kbd-or">or</span> ');
                        listItem.appendChild(shortcutColumn);
                    }

                    shortcutsGroup.appendChild(list);
                }
            }

            document.body.appendChild(this.wrap);
        }

        open() {
            this.prevActiveElement = document.activeElement;

            document.body.overflow = 'hidden';
            this.wrap.hidden = false;
            this.wrap.children[0].focus();
        }

        close() {
            document.body.overflow = '';
            this.wrap.hidden = true;

            if (this.prevActiveElement) {
                this.prevActiveElement.focus();
                this.prevActiveElement = null;
            }
        }
    }

    window.addEventListener('load', () => {
        let helpDialog = null;
        let openHelp = () => {
            if (!helpDialog) helpDialog = new ShortcutsHelpDialog();
            helpDialog.open();
        };

        // Global shortcuts work everywhere.
        let globalShortcuts = new ShortcutHandler(document, false);
        globalShortcuts.add('?, ' + (isMac ? 'Meta+/' : 'Ctrl+/'), openShortcutsReference);

        // Common shortcuts work everywhere except on text fields.
        let commonShortcuts = new ShortcutHandler(document, false, notTextField);

        let bindElement = bindElementFactory(commonShortcuts);
        let bindLink = bindLinkFactory(commonShortcuts);
        let bindCollection = bindCollectionFactory(commonShortcuts);

        // Common shortcuts
        commonShortcuts.groupStart('Common');
            bindCollection('g', '.header-links__link', 'First 9 header links', 'Header link');
            bindLink('g h', '/', 'Home');
            bindLink('g l', '/list/', 'List of hyphae');
            bindLink('g r', '/recent-changes/', 'Recent changes');
            bindElement('g u', '.header-links__entry_user .header-links__link', 'Your profile′s hypha');
        commonShortcuts.groupEnd();

        if (typeof editTextarea === 'undefined') {
            // Hypha shortcuts
            commonShortcuts.groupStart('Hypha');
                bindCollection('', 'article .wikilink', 'First 9 hypha′s links', 'Hypha link');
                bindElement('p, Alt+ArrowLeft, Ctrl+Alt+ArrowLeft', '.prevnext__prev', 'Next hypha');
                bindElement('n, Alt+ArrowRight, Ctrl+Alt+ArrowRight', '.prevnext__next', 'Previous hypha');
                bindElement('s, Alt+ArrowUp, Ctrl+Alt+ArrowUp', $$('.navi-title a').slice(1, -1).slice(-1)[0], 'Parent hypha');
                bindElement('c, Alt+ArrowDown, Ctrl+Alt+ArrowDown', '.subhyphae__link', 'First child hypha');
                bindElement('e, Ctrl+Enter', '.hypha-tabs__link[href^="/edit/"]', 'Edit this hypha');
            commonShortcuts.groupEnd();

        } else {
            // Hypha editor shortcuts. These work only on editor's text area.
            let editorShortcuts = new ShortcutHandler(editTextarea, true);
            let bindElement = bindElementFactory(editorShortcuts);

            let shortcuts = [
                // Win+Linux    Mac                  Action              Description
                ['Ctrl+b',      'Meta+b',            wrapBold,           'Format: Bold'],
                ['Ctrl+i',      'Meta+i',            wrapItalic,         'Format: Italic'],
                ['Ctrl+M',      'Meta+Shift+m',      wrapMonospace,      'Format: Monospaced'],
                ['Ctrl+I',      'Meta+Shift+i',      wrapHighlighted,    'Format: Highlight'],
                ['Ctrl+.',      'Meta+Shift+.',      wrapLifted,         'Format: Superscript'],
                ['Ctrl+Comma',  'Meta+Shift+Comma',  wrapLowered,        'Format: Subscript'],
                ['Ctrl+X',      'Meta+Shift+x',      wrapStrikethrough,  'Format: Strikethrough'],
                ['Ctrl+k',      'Meta+k',            wrapLink,           'Format: Link'],
            ];

            editorShortcuts.groupStart('Editor');
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
        }
    });
})();
