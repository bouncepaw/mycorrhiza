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

        group(...args) {
            if (typeof args[0] === 'string') this.fakeItem(args.shift());
            shortcutsGroup = [];

            args[0].bind(this)();

            if (shortcutsGroup && shortcutsGroup.length) allShortcuts.push(shortcutsGroup);
            shortcutsGroup = null;
        }

        bindElement(shortcut, element, ...other) {
            element = typeof element === 'string' ? $(element) : element;
            if (!element) return;
            this.add(shortcut, () => {
                if (isTextField(element)) {
                    element.focus();
                } else {
                    element.click();
                }
            }, ...other);
        }

        bindLink(shortcut, link, ...other) {
            this.add(shortcut, () => window.location.href = link, ...other);
        }

        bindCollection(prefix, elements, collectionDescription, itemDescription) {
            this.fakeItem(prefix + ' 1 – 9', collectionDescription);

            if (typeof elements === 'string') {
                elements = $$(elements);
            } else if (Array.isArray(elements)) {
                elements = elements.map(el => typeof el === 'string' ? $(el) : el);
            }

            for (let i = 1; i <= elements.length && i < 10; i++) {
                this.bindElement(`${prefix} ${i}`, elements[i-1], `${itemDescription} #${i}`, false);
            }
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
                event.stopPropagation();
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

    class ShortcutsHelpDialog {
        constructor() {
            let template = $('#dialog-template');
            let clonedTemplate = template.content.cloneNode(true);
            this.backdrop = clonedTemplate.children[0];
            this.dialog = clonedTemplate.children[1];

            this.dialog.classList.add('shortcuts-help');
            this.dialog.hidden = true;
            this.backdrop.hidden = true;

            document.body.appendChild(this.backdrop);
            document.body.appendChild(this.dialog);

            this.close = this.close.bind(this);

            this.dialog.querySelector('.dialog__title').textContent = 'List of shortcuts';
            this.dialog.querySelector('.dialog__close-button').addEventListener('click', this.close);
            this.backdrop.addEventListener('click', this.close);

            this.shortcuts = new ShortcutHandler(this.dialog, false);
            this.shortcuts.add('Escape', this.close, null, false);

            let shortcutsGroup;
            let shortcutsGroupTemplate = document.createElement('div');
            shortcutsGroupTemplate.className = 'shortcuts-group';

            for (let item of allShortcuts) {
                if (item.description && !item.shortcut) {
                    shortcutsGroup = shortcutsGroupTemplate.cloneNode();
                    this.dialog.querySelector('.dialog__content').appendChild(shortcutsGroup);

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
        }

        open() {
            this.prevActiveElement = document.activeElement;

            document.body.overflow = 'hidden';
            this.backdrop.hidden = false;
            this.dialog.hidden = false;
            this.dialog.focus();
        }

        close() {
            document.body.overflow = '';
            this.backdrop.hidden = true;
            this.dialog.hidden = true;

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
        globalShortcuts.add('?, ' + (isMac ? 'Meta+/' : 'Ctrl+/'), openHelp);

        // Page shortcuts work everywhere except on text fields.
        let pageShortcuts = new ShortcutHandler(document, false, notTextField);
        pageShortcuts.add('?', openHelp, null, false);

        // Common shortcuts
        pageShortcuts.group('Common', function () {
            this.bindCollection('g', '.header-links__link', 'First 9 header links', 'Header link');
            this.bindLink('g h', '/', 'Home');
            this.bindLink('g l', '/list/', 'List of hyphae');
            this.bindLink('g r', '/recent-changes/', 'Recent changes');
            this.bindElement('g u', '.header-links__entry_user .header-links__link', 'Your profile′s hypha');
        });

        if (typeof editTextarea === 'undefined') {
            // Hypha shortcuts
            pageShortcuts.group('Hypha', function () {
                this.bindCollection('', 'article .wikilink', 'First 9 hypha′s links', 'Hypha link');
                this.bindElement('p, Alt+ArrowLeft, Ctrl+Alt+ArrowLeft', '.prevnext__prev', 'Next hypha');
                this.bindElement('n, Alt+ArrowRight, Ctrl+Alt+ArrowRight', '.prevnext__next', 'Previous hypha');
                this.bindElement('s, Alt+ArrowUp, Ctrl+Alt+ArrowUp', $$('.navi-title a').slice(1, -1).slice(-1)[0], 'Parent hypha');
                this.bindElement('c, Alt+ArrowDown, Ctrl+Alt+ArrowDown', '.subhyphae__link', 'First child hypha');
                this.bindElement('e, Ctrl+Enter', '.hypha-tabs__link[href^="/edit/"]', 'Edit this hypha');
            });

        } else {
            // Hypha editor shortcuts. These work only on editor's text area.
            let editorShortcuts = new ShortcutHandler(editTextarea, true);

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

            editorShortcuts.group('Editor', function () {
                for (let shortcut of shortcuts) {
                    if (isMac) {
                        this.add(shortcut[1], ...shortcut.slice(2))
                    } else {
                        this.add(shortcut[0], ...shortcut.slice(2))
                    }
                }
            });

            editorShortcuts.group(function () {
                this.bindElement(isMac ? 'Meta+Enter' : 'Ctrl+Enter', $('.edit-form__save'), 'Save changes');
            });
        }
    });
})();
