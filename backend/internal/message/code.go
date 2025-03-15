package message

import (
	"context"
	"fmt"
	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"go.uber.org/zap/buffer"
	"io"
)

type Highlighter struct {
}

var ErrNoLangFound = fmt.Errorf("no language found by code")

func (h *Highlighter) Highlight(code []byte, ctx context.Context, lang string) ([]byte, error) {
	b := buffer.Buffer{}
	err := h.HighlightToWriter(&b, code, ctx, lang)
	if err != nil {
		return nil, fmt.Errorf("cannot highlight to writer: %w", err)
	}

	return b.Bytes(), nil
}

func (h *Highlighter) HighlightToWriter(w io.Writer, code []byte, _ context.Context, lang string) error {
	lexer := lexers.Get(lang)
	if lexer == nil {
		return ErrNoLangFound
	}

	style := styles.Get("github-dark")
	formatter := html.New(html.WithClasses(true))

	iterator, err := lexer.Tokenise(nil, string(code))
	if err != nil {
		return fmt.Errorf("cannot tokenise code for lang '%s': %w", lang, err)
	}

	err = formatter.Format(w, style, iterator)
	if err != nil {
		return fmt.Errorf("cannot format code to writer: %w", err)
	}

	return nil
}

func NewHighlighter() *Highlighter {
	return &Highlighter{}
}
