
export default class Editor {

    #textarea: HTMLTextAreaElement;

    constructor(textarea: HTMLTextAreaElement) {
        this.#textarea = textarea;
    }

    getContent(): string {
        return this.#textarea.value;
    }

    bold() {
        this.#coverBySymbols('**');
    }

    code() {
        this.#coverBySymbols('`');
    }

    italic() {
        this.#coverBySymbols('_');
    }

    #coverBySymbols(symbol: string)
    {
        const textarea = this.#textarea;
        const start = textarea.selectionStart;
        const end = textarea.selectionEnd;
        const text = textarea.value;

        // Получаем текст до и после выделения
        const beforeSelection = text.substring(0, start);
        const afterSelection = text.substring(end);

        // Проверяем, обернут ли текст в **
        const isBold = beforeSelection.endsWith(symbol) && afterSelection.startsWith(symbol);

        const symbolLen = symbol.length;

        if (isBold) {
            // Убираем ** и корректируем позицию курсора
            textarea.value = beforeSelection.slice(0, -symbolLen) + text.substring(start, end) + afterSelection.slice(symbolLen);
            textarea.setSelectionRange(start - symbolLen, end - symbolLen);
        } else {
            // Добавляем ** и устанавливаем курсор между ними
            textarea.value = beforeSelection + `${symbol}${text.substring(start, end)}${symbol}` + afterSelection;
            if (start === end) {
                // Если текст не выделен, ставим курсор между **
                textarea.setSelectionRange(start + symbolLen, start + symbolLen);
            } else {
                // Если текст выделен, выделение остается на тексте
                textarea.setSelectionRange(start + symbolLen, end + symbolLen);
            }
        }

        textarea.focus();
    }

    heading() {
        const textarea = this.#textarea;
        const start = textarea.selectionStart;
        const end = textarea.selectionEnd;
        const text = textarea.value;

        // Получаем текст до выделения
        const beforeSelection = text.substring(0, start);

        // Определяем начало строки, чтобы проверить наличие '### '
        const lineStart = beforeSelection.lastIndexOf('\n') + 1;
        const lineText = text.substring(lineStart, start);

        const hasHeading = lineText.startsWith('### ');

        if (hasHeading) {
            // Убираем '### ' и корректируем позицию курсора
            textarea.value = beforeSelection.slice(0, lineStart) + lineText.slice(4) + text.substring(start);
            textarea.setSelectionRange(start - 4, end - 4);
        } else {
            // Добавляем '### ' перед началом строки
            textarea.value = beforeSelection.slice(0, lineStart) + '### ' + lineText + text.substring(start);
            textarea.setSelectionRange(start + 4, end + 4);
        }

        textarea.focus();
    }

    quotes() {
        const textarea = this.#textarea;
        const start = textarea.selectionStart;
        const end = textarea.selectionEnd;
        const text = textarea.value;

        if (start === end) {
            // Если текст не выделен
            if (text.trim() === '') {
                // Если textarea пустая, добавляем '> ' без переноса строк
                textarea.value = '> ';
                textarea.setSelectionRange(2, 2); // Устанавливаем курсор после '> '
            } else {
                // Если в textarea есть текст, добавляем две пустые строки и '> '
                const beforeCursor = text.substring(0, start);
                const afterCursor = text.substring(start);
                const newText = beforeCursor + '\n\n> ' + afterCursor;
                textarea.value = newText;
                textarea.setSelectionRange(start + 4, start + 4); // Устанавливаем курсор после '> '
            }
        } else {
            // Если текст выделен, обрабатываем каждую строку
            const selectedText = text.substring(start, end);
            const lines = selectedText.split('\n');

            const isQuoted = lines.every(line => line.startsWith('> '));

            const processedText = lines.map(line => {
                if (line.startsWith('> ')) {
                    // Убираем '> ' если уже есть
                    return line.slice(2);
                } else {
                    // Добавляем '> ' если его нет
                    return '> ' + line;
                }
            }).join('\n');

            textarea.value = text.substring(0, start) + processedText + text.substring(end);

            // Корректируем выделение, чтобы '> ' не выделялось
            const adjustment = isQuoted ? -2 : 2;
            const newStart = start + (isQuoted ? 0 : 2);
            const newEnd = end + adjustment * lines.length;
            textarea.setSelectionRange(newStart, newEnd);
        }

        textarea.focus();
    }

    link() {
        const textarea = this.#textarea;
        const start = textarea.selectionStart;
        const end = textarea.selectionEnd;
        const text = textarea.value;

        if (start === end) {
            // Если текст не выделен
            const beforeCursor = text.substring(0, start);
            const afterCursor = text.substring(start);

            if (beforeCursor.endsWith('[') && afterCursor.startsWith('](url)')) {
                // Если курсор находится между '[' и '](url)', удаляем конструкцию
                textarea.value = beforeCursor.slice(0, -1) + afterCursor.slice(6);
                textarea.setSelectionRange(start - 1, start - 1);
            } else {
                // Вставляем '[](url)' и перемещаем курсор между [ и ]
                const newText = beforeCursor + '[](url)' + afterCursor;
                textarea.value = newText;
                textarea.setSelectionRange(start + 1, start + 1);
            }
        } else {
            // Если текст выделен
            const selectedText = text.substring(start, end);
            const beforeSelection = text.substring(0, start);
            const afterSelection = text.substring(end);

            // Вставляем '[](url)' и перемещаем выделенный текст между [ и ], а слово url становится выделенным
            const newText = beforeSelection + '[' + selectedText + '](url)' + afterSelection;
            textarea.value = newText;
            textarea.setSelectionRange(start + selectedText.length + 3, start + selectedText.length + 6);
        }

        textarea.focus();
    }
}