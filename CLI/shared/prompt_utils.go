package shared

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"golang.org/x/term"
)

type MenuOptions struct {
	Title    string
	Items    []string
	ExitText string
	SoftExit bool
}

var Borders = struct {
	// Horizontal bar
	H string
	// Vertical Bar
	V string
	// Top left corner
	Tlc string
	// Top right corner
	Trc string
	// bottom left corner
	Blc string
	// Bottom right corner
	Brc string
	// Top Junction
	Tj string
	// Left Junction
	Lj string
	// Right Junction
	Rj string
	// Bottom Junction
	Bj string
	// Cross Junction
	Cj string
}{
	H:   "─",
	V:   "│",
	Tlc: "╭",
	Trc: "╮",
	Blc: "╰",
	Brc: "╯",
	Tj:  "┬",
	Lj:  "├",
	Rj:  "┤",
	Bj:  "┴",
	Cj:  "┼",
}

var TerminalMap = map[string]string{
	";bk;": "\033[90m",
	";dg;": "\033[32m",
	";r;":  "\033[91m",
	";g;":  "\033[92m",
	";y;":  "\033[93m",
	";b;":  "\033[94m",
	";m;":  "\033[95m",
	";c;":  "\033[96m",
	";w;":  "\033[97m",
	";x;":  "\033[0m",
	";@;":  "\u25CF",
	";*;":  "\u2022",
	";u;":  "\033[4m",
	"\t":   "    ",
}

func PromptMenu(opt MenuOptions, r *bufio.Reader) (int, error) {
	title := opt.Title
	items := opt.Items
	exitText := "Exit"
	if opt.ExitText != "" {
		exitText = opt.ExitText
	}

	padding := 4
	spacing := 2
	borderLen := len(title) + padding + spacing

	ClearScreen()

	fmt.Printf(
		"\n\n\t\033[90m%s%s\033[90m%s\n\t%s  \033[97m%s\033[90m  %s\033[0m\n\t\033[90m%s%s\033[90m%s\n\n",
		Borders.Tlc,
		fmt.Sprintf("\033[90m%s\033[0m", strings.Repeat(Borders.H, borderLen-2)),
		Borders.Trc,
		Borders.V,
		title,
		Borders.V,
		Borders.Blc,
		fmt.Sprintf("\033[90m%s\033[0m", strings.Repeat(Borders.H, borderLen-2)),
		Borders.Brc,
	)

	items = append(items, exitText)
	itemCount := 0
	for i, item := range items {
		if item == "" {
			fmt.Println("")
			continue
		}
		itemCount += 1

		padding := " "
		if len(items) < 10 || itemCount > 9 {
			padding = ""
		}

		if i+1 < len(items) {
			fmt.Print(
				ColorString(fmt.Sprintf("\t\t  ;y;%d.%s ;x;%s\n", itemCount, padding, item)),
			)
		} else {
			fmt.Print(
				ColorString(fmt.Sprintf("\n\t\t  ;r;%d.%s %s;x;\n\n", itemCount, padding, item)),
			)
		}
	}

	input := PromptInput("Choose a number", r)
	if input == "q" || input == "quit" || input == "exit" {
		if opt.SoftExit {
			return itemCount, nil
		}
		os.Exit(0)
	}

	choice, err := ParseInt(input)
	if err != nil {
		return 0, fmt.Errorf("'%s' is not a valid choice; integers only", input)
	}

	if choice < 1 || choice > len(items) {
		return 0, fmt.Errorf("invalid choice; select a number between 1 and %d", len(items))
	}

	// Last item is always "Exit"
	if choice == itemCount {
		if opt.SoftExit {
			return itemCount, nil
		}
		fmt.Printf("\n\n  \033[95mBye Bye!!")
		os.Exit(0)
	}

	return choice, nil
}

func PromptInput(msg string, r *bufio.Reader) string {
	fmt.Printf(ColorString("\n;w;%s;x;\n> ;g;"), msg)
	s, err := r.ReadString('\n')
	if err != nil {
		panic(fmt.Errorf("failed to retrieve input from prompt::%w", err))
	}
	return strings.TrimSpace(s)
}

func PromptBackToMenu(r *bufio.Reader) {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	HideCursor()
	fmt.Print(ColorString("\n\n ;bk;Press any key to continue..."))
	byteInput := make([]byte, 1)
	os.Stdin.Read(byteInput)
	ShowCursor()
}

/*
ReplaceAll uses the token map to replace all keys with the
mapped value, within the given string.
*/
func ReplaceAll(s string, tokenMap map[string]string) string {
	sb := strings.Builder{}
	for i := 0; i < len(s); {
		var isSwapped bool
		for t, v := range tokenMap {
			if strings.HasPrefix(s[i:], t) {
				sb.WriteString(v)
				i += len(t)
				isSwapped = true
				break
			}
		}
		if !isSwapped {
			sb.WriteByte(s[i])
			i++
		}
	}

	return sb.String()
}

func ColorString(s string) string {
	return ReplaceAll(s, TerminalMap)
}

func ResetCursor() {
	fmt.Print("\033[1;1H")
}

func HideCursor() {
	fmt.Print("\033[?25l")
}

func ShowCursor() {
	fmt.Print("\033[?25h")
}

func ClearScreen() {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()

	default:
		fmt.Print("\033[3J\033[H\033[2J")
	}
}

func PrintError(err error) {
	fmt.Printf(ColorString("\n;w;[;r;ERROR;w;]: ;y;%s;x;\n"), err)
}
