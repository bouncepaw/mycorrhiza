const editTextarea = document.getElementsByClassName('edit-form__textarea')[0]

function placeCursor(position, el = editTextarea) {
    el.selectionEnd = position
    el.selectionStart = el.selectionEnd
}

function getSelectedText(el = editTextarea) {
    const [start, end] = [el.selectionStart, el.selectionEnd]
    const text = el.value
    return text.substring(start, end)
}

function textInserter(text, cursorPosition = null, el = editTextarea) {
    return function() {
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
        const [start, end] = [el.selectionStart, el.selectionEnd]
        if (postfix == null) {
            postfix = prefix
        }
        text = getSelectedText(el)
        result = prefix + text + postfix
        el.setRangeText(result, start, end, 'select')
        el.focus()
        placeCursor(end + cursorPosition)
    }
}

const wrapBold = selectionWrapper(2, '**'),
    wrapItalic = selectionWrapper(2, '//'),
    wrapMonospace = selectionWrapper(1, '`'),
    wrapHighlighted = selectionWrapper(2, '!!'),
    wrapLifted = selectionWrapper(1, '^'),
    wrapLowered = selectionWrapper(2, ',,'),
    wrapStrikethrough = selectionWrapper(2, '~~'),
    wrapLink = selectionWrapper(2, '[[', ']]')

const insertHorizontalBar = textInserter('----\n'),
    insertImgBlock = textInserter('img {\n\t\n}\n', 7),
    insertTableBlock = textInserter('table {\n\t\n}\n', 9),
    insertRocket = textInserter('=> '),
    insertXcl = textInserter('<= '),
    insertHeading2 = textInserter('## '),
    insertHeading3 = textInserter('### '),
    insertCodeblock = textInserter('```\n\n```\n', 4)

function insertDate() {
    let date = new Date().toISOString().split('T')[0]
    textInserter(date)()
}

function insertUserlink() {
    const userlink = document.querySelector('.header-links__entry_user a')
    const userHypha = userlink.getAttribute('href').substring(7) // no /hypha/
    textInserter('[[' + userHypha + ']]')()
}
