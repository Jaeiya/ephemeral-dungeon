package utils

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

var style = lipgloss.NewStyle()

const _tagLen = 4

var colorStyleMap = map[string]lipgloss.Style{
	// ";gu;":  style.Underline(true).Foreground(ansi.BrightGreen),
	// ";dgu;": style.Underline(true).Foreground(ansi.Green),

	// Bright colors
	";rd;": style.Foreground(ansi.BrightRed),
	";gn;": style.Foreground(ansi.BrightGreen),
	";yw;": style.Foreground(ansi.BrightYellow),
	";be;": style.Foreground(ansi.BrightBlue),
	";ma;": style.Foreground(ansi.BrightMagenta),
	";cn;": style.Foreground(ansi.BrightCyan),
	";we;": style.Foreground(ansi.BrightWhite),

	// Dark colors
	";dr;": style.Foreground(ansi.Red),
	";dg;": style.Foreground(ansi.Green),
	";dy;": style.Foreground(ansi.Yellow),
	";db;": style.Foreground(ansi.Blue),
	";dm;": style.Foreground(ansi.Magenta),
	";dc;": style.Foreground(ansi.Cyan),
	";dw;": style.Foreground(ansi.White),
	";gy;": style.Foreground(ansi.BrightBlack),
}

var colorCodes = map[string]string{
	";rd;": "",
	";gn;": "",
	";yw;": "",
	";be;": "",
	";ma;": "",
	";cn;": "",
	";we;": "",
	"@rd;": "",
	"@gn;": "",
	"@yw;": "",
	"@be;": "",
	"@ma;": "",
	"@cn;": "",
	"@we;": "",
	";dr;": "",
	";dg;": "",
	";dy;": "",
	";db;": "",
	";dm;": "",
	";dc;": "",
	";dw;": "",
}

type StyledString struct {
	text      string
	formatted string // Cached rendered text
}

func NewStyledString(s string) StyledString {
	return StyledString{text: s}
}

func (ss *StyledString) Render() string {
	if len(ss.formatted) > 0 {
		fmt.Println("cached")
		return ss.formatted
	}

	var sb strings.Builder
	currentStyle := style.Foreground(ansi.White)
	start := 0

	for i := 0; i < len(ss.text); i++ {
		if ss.text[i] != ';' && ss.text[i] != '@' {
			continue
		}
		if i+_tagLen > len(ss.text) {
			break
		}
		if tagStyle, exists := colorStyleMap[ss.text[i:i+_tagLen]]; exists {
			sb.WriteString(currentStyle.Render(ss.text[start:i]))
			currentStyle = tagStyle
			i += _tagLen - 1 // -1 because loop does i++
			start = i + 1
		}
	}

	if start < len(ss.text) {
		sb.WriteString(currentStyle.Render(ss.text[start:]))
	}

	ss.formatted = sb.String()
	return ss.formatted
}

func (ss StyledString) Raw() string {
	sb := strings.Builder{}
	sb.Grow(len(ss.text))

	for i := 0; i < len(ss.text); i++ {
		if i+_tagLen <= len(ss.text) {
			if _, exists := colorCodes[ss.text[i:i+_tagLen]]; exists {
				i += _tagLen - 1 // -1 because loop does i++
				continue
			}
		}
		sb.WriteByte(ss.text[i])
	}

	return strings.ToLower(sb.String())
}
