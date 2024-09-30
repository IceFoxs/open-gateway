package markdown

import (
	"strconv"
	"strings"
)

const (
	WARNING_COLOR     string = "#EA9F00"
	ERROR_COLOR       string = "#FF0000"
	NOTICE_COLOR      string = "#5AB030"
	CONTENT_COLOR     string = "#664B4B"
	CONTENT_RED_COLOR string = "#FF0000"
)

type Markdown struct {
	stringBuilder *strings.Builder
}

func WrapColor(color string, text string) string {
	return "<font color=" + color + ">" + text + "</font>  "
}

func WrapFont(text string) string {
	return "<font  size=2 >" + text + "</font>  "
}

func WrapFontColor(color, text string) string {
	return "<font color=" + color + " size=2>" + text + "</font>  "
}
func NewMarkdown(opts ...MarkdownOption) *Markdown {
	markdown := &Markdown{
		stringBuilder: new(strings.Builder),
	}

	for _, opt := range opts {
		opt(markdown)
	}

	return markdown
}

func (m *Markdown) Builder() string {
	return m.stringBuilder.String()
}
func WithTitle1(text string) MarkdownOption {
	return func(config *Markdown) {
		config.stringBuilder.WriteString("# ")
		config.stringBuilder.WriteString(text)
		config.stringBuilder.WriteString("\n")
	}
}

func WithTitle2(text string) MarkdownOption {
	return func(config *Markdown) {
		config.stringBuilder.WriteString("## ")
		config.stringBuilder.WriteString(text)
		config.stringBuilder.WriteString("\n")
	}
}
func WithTitle3(text string) MarkdownOption {
	return func(config *Markdown) {
		config.stringBuilder.WriteString("#### ")
		config.stringBuilder.WriteString(text)
		config.stringBuilder.WriteString("\n")
	}
}

func WithTitle4(text string) MarkdownOption {
	return func(config *Markdown) {
		config.stringBuilder.WriteString("#### ")
		config.stringBuilder.WriteString(text)
		config.stringBuilder.WriteString("\n")
	}
}

func WithTitle5(text string) MarkdownOption {
	return func(config *Markdown) {
		config.stringBuilder.WriteString("##### ")
		config.stringBuilder.WriteString(text)
		config.stringBuilder.WriteString("\n")
	}
}

func WithTitle6(text string) MarkdownOption {
	return func(config *Markdown) {
		config.stringBuilder.WriteString("###### ")
		config.stringBuilder.WriteString(text)
		config.stringBuilder.WriteString("\n")
	}
}
func WithUl(text []string) MarkdownOption {
	return func(config *Markdown) {
		if len(text) == 0 {
			return
		}
		for _, t := range text {
			config.stringBuilder.WriteString("- ")
			config.stringBuilder.WriteString(t)
			config.stringBuilder.WriteString("\n")
		}
	}
}

func WithText(text string) MarkdownOption {
	return func(config *Markdown) {
		config.stringBuilder.WriteString("\n")
		config.stringBuilder.WriteString(text)
		config.stringBuilder.WriteString("\n")
	}
}

func WithLi(text []string) MarkdownOption {
	return func(config *Markdown) {
		if len(text) == 0 {
			return
		}
		for i, t := range text {
			config.stringBuilder.WriteString(strconv.Itoa(i) + ". ")
			config.stringBuilder.WriteString(t)
			config.stringBuilder.WriteString("\n")
		}
	}
}

// MarkdownOption ...
type MarkdownOption func(*Markdown)
