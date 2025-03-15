package message

import (
	"context"
	"fmt"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	mdparser "github.com/gomarkdown/markdown/parser"
	"github.com/microcosm-cc/bluemonday"
	"github.com/nessai1/aiinterview/internal/utils"
	"regexp"
	"strings"
)

type Parser struct {
	sanitizer *bluemonday.Policy

	highlighter *Highlighter
}

func (p *Parser) Parse(content []byte) ([]byte, error) {
	extensions := mdparser.CommonExtensions | mdparser.AutoHeadingIDs | mdparser.NoEmptyLineBeforeBlock
	mdp := mdparser.NewWithExtensions(extensions)

	// Регулярное выражение для поиска блоков кода
	codeBlockRegex := regexp.MustCompile("(?s)```(\\w+)(.*?)\\n?```")

	// Карта для хранения временных замен
	replacements := make(map[string]string)

	// Заменяем блоки кода на временные UUID
	processedText := codeBlockRegex.ReplaceAllStringFunc(string(content), func(match string) string {

		// TODO: хуй знает пока как правильно обработать. Мне лень
		id, _ := utils.GenerateUUIDv7()
		// Извлекаем язык и код
		matches := codeBlockRegex.FindStringSubmatch(match)
		language := matches[1]

		//innerCode := strings.Replace(matches[2], "\\n", "", 1)
		//innerCode = strings.ReplaceAll(innerCode, "\\n", "\n")
		resultCode, _ := p.highlighter.Highlight([]byte(matches[2]), context.TODO(), language)

		// Сохраняем оригинальный блок кода
		replacements[id] = strings.ReplaceAll(string(resultCode), "\\n", "\n")
		// Возвращаем временную замену
		return fmt.Sprintf("{#{%s}#}", id)
	})

	processedText = strings.ReplaceAll(processedText, "\\n", "<br>")

	doc := mdp.Parse([]byte(processedText))

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	processedText = string(p.sanitizer.SanitizeBytes(markdown.Render(doc, renderer)))

	// Вставляем блоки кода обратно
	for id, originalCode := range replacements {
		processedText = strings.ReplaceAll(processedText, fmt.Sprintf("{#{%s}#}", id), originalCode)
	}

	return []byte(processedText), nil
}

func NewParser(highlighter *Highlighter) *Parser {
	p := Parser{
		sanitizer:   createDefaultSanitizer(),
		highlighter: highlighter,
	}

	return &p
}

func createDefaultSanitizer() *bluemonday.Policy {
	return bluemonday.UGCPolicy().AllowElementsMatching(regexp.MustCompile(`> `))
}
