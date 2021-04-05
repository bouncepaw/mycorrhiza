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
    wrapStrikethrough = selectionWrapper(2, '~~')

function insertHorizontalBar() {
    insertTextAtCursor('----\n')
}

function insertImgBlock() {
    insertTextAtCursor('img {\n\t\n}\n', 7)
}
