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

function insertTextAtCursor(text, cursorPosition = null, el = editTextarea) {
    const [start, end] = [el.selectionStart, el.selectionEnd]
    el.setRangeText(text, start, end, 'select')
    el.focus()
    if (cursorPosition == null) {
        placeCursor(end + text.length)
    } else {
        placeCursor(end + cursorPosition)
    }
}

function wrapSelection(prefix, postfix = null, el = editTextarea) {
    const [start, end] = [el.selectionStart, el.selectionEnd]
    if (postfix == null) {
        postfix = prefix
    }
    text = getSelectedText(el)
    result = prefix + text + postfix
    el.setRangeText(result, start, end, 'select')
    el.focus()
    placeCursor(end + (prefix + postfix).length)
}

function insertDate() {
    let date = new Date().toISOString().split('T')[0]
    insertTextAtCursor(date)
}

function wrapBold() {
    wrapSelection('**')
}

function wrapItalic() {
    wrapSelection('//')
}

function wrapMonospace() {
    wrapSelection('`')
}

function wrapHighlighted() {
    wrapSelection('!!')
}

function wrapLifted() {
    wrapSelection('^')
}

function wrapLowered() {
    wrapSelection(',,')
}

function wrapStroked() {
    wrapSelection('~~')
}

function insertHorizontalBar() {
    insertTextAtCursor('----\n')
}

function insertImgBlock() {
    insertTextAtCursor('img {\n\t\n}\n', 7)
}
