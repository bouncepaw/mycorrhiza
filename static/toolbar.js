const editTextarea = document.getElementsByClassName('edit-form__textarea')[0]

function placeCursor(position, el = editTextarea) {
    el.selectionEnd = position
    el.selectionStart = el.selectionEnd
}

function selectRange(left, right, el = editTextarea) {
    el.selectionEnd = right
    el.selectionStart = left
}

function getSelectedText(el = editTextarea) {
    const [start, end] = [el.selectionStart, el.selectionEnd]
    const text = el.value
    return text.substring(start, end)
}

function textInserter(text, cursorPosition = null, el = editTextarea) {
    return function() {
        window.hyphaChanged = true
        const [start, end] = [el.selectionStart, el.selectionEnd]
        el.setRangeText(text, start, end, 'select')
        el.focus()
        if (cursorPosition == null) {
            placeCursor(end + text.length)
        } else {
            placeCursor(end + cursorPosition)
        }
    }
}

function selectionWrapper(cursorPosition, prefix, postfix = null, el = editTextarea) {
    return function() {
        window.hyphaChanged = true
        let [start, end] = [el.selectionStart, el.selectionEnd]
        if (postfix == null) {
            postfix = prefix
        }
        let text = getSelectedText(el)
        let removing = false
        let result
        if (text.startsWith(prefix) && text.endsWith(postfix)) {
            // selection is decorated, so we just cut it
            removing = true
            result = text.substring(cursorPosition, text.length - cursorPosition)
        } else if ( (prefix == el.value.slice(start-cursorPosition, start)) &&
                    (postfix == el.value.slice(end, end+cursorPosition)) ) {
            // selection is surrounded by decorations
            removing = true
            result = text
            start -= cursorPosition
            end += cursorPosition
        } else {
            // no decorations, so we add them
            result = prefix + text + postfix
        }
        el.setRangeText(result, start, end, 'select')
        el.focus()
        if (removing) {
            selectRange(start, end-cursorPosition*2)
        } else {
            selectRange(start+cursorPosition, end+cursorPosition)
        }
    }
}

const wrapBold = selectionWrapper(2, '**'),
    wrapItalic = selectionWrapper(2, '//'),
    wrapMonospace = selectionWrapper(1, '`'),
    wrapHighlighted = selectionWrapper(2, '++'),
    wrapLifted = selectionWrapper(2, '^^'),
    wrapLowered = selectionWrapper(2, ',,'),
    wrapStrikethrough = selectionWrapper(2, '~~'),
    wrapUnderline = selectionWrapper(2, '__'),
    wrapLink = selectionWrapper(2, '[[', ']]')

const insertHorizontalBar = textInserter('\n----\n'),
    insertImgBlock = textInserter('\nimg {\n   \n}\n', 10),
    insertTableBlock = textInserter('\ntable {\n   \n}\n', 12),
    insertRocket = textInserter('\n=> '),
    insertXcl = textInserter('\n<= '),
    insertHeading2 = textInserter('\n## '),
    insertHeading3 = textInserter('\n### '),
    insertCodeblock = textInserter('\n```\n\n```\n', 5),
    insertBulletedList = textInserter('\n* '),
    insertNumberedList = textInserter('\n*. ')

function insertDate() {
    let date = new Date().toISOString().split('T')[0]
    textInserter(date)()
}

function insertTimeUTC() {
	let time = new Date().toISOString().substring(11, 19) + " UTC"
	textInserter(time)()
}

function insertUserlink() {
    const userlink = document.querySelector('.auth-links__user-link')
    const userHypha = userlink.getAttribute('href').substring(7) // no /hypha/
    textInserter('[[' + userHypha + ']]')()
}
