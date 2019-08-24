package input

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/chzyer/readline"
	"github.com/tanema/promptui/frmt"
	"github.com/tanema/promptui/screenbuf"
)

// ErrAbort is the error returned when confirm prompts are supplied "n"
var ErrAbort = errors.New("")

// ErrEOF is the error returned from prompts when EOF is encountered.
var ErrEOF = errors.New("^D")

// ErrInterrupt is the error returned from prompts when an interrupt (ctrl-c) is
// encountered.
var ErrInterrupt = errors.New("^C")

// Prompt represents a single line text field input with options for validation and input masks.
type Prompt struct {
	// Label is the value displayed on the command line prompt.
	//
	// The value for Label can be a simple string or a struct that will need to be accessed by dot notation
	// inside the templates. For example, `{{ .Name }}` will display the name property of a struct.
	Label interface{}

	// Default is the initial value for the prompt. This value will be displayed next to the prompt's label
	// and the user will be able to view or change it depending on the options.
	Default string

	// AllowEdit lets the user edit the default value. If false, any key press
	// other than <Enter> automatically clears the default value.
	AllowEdit bool

	// Validate is an optional function that fill be used against the entered value in the prompt to validate it.
	Validate func(string) error

	// Mask is an optional rune that sets which character to display instead of the entered characters. This
	// allows hiding private information like passwords.
	Mask rune

	// Templates can be used to customize the prompt output. If nil is passed, the
	// default templates are used. See the PromptTemplates docs for more info.
	Templates *PromptTemplates

	// IsConfirm makes the prompt ask for a yes or no ([Y/N]) question rather than request an input. When set,
	// most properties related to input will be ignored.
	IsConfirm bool

	// IsVimMode enables vi-like movements (hjkl) and editing.
	IsVimMode bool

	// the Pointer defines how to render the cursor.
	Pointer Pointer

	stdin  io.ReadCloser
	stdout io.WriteCloser
}

// PromptTemplates allow a prompt to be customized following stdlib
// text/template syntax. Custom state, colors and background color are available for use inside
// the templates and are documented inside the Variable section of the docs.
//
// Examples
//
// text/templates use a special notation to display programmable content. Using the double bracket notation,
// the value can be printed with specific helper functions. For example
//
// This displays the value given to the template as pure, unstylized text.
//	'{{ . }}'
//
// This displays the value colored in cyan
//	'{{ . | cyan }}'
//
// This displays the value colored in red with a cyan background-color
//	'{{ . | red | cyan }}'
//
// See the doc of text/template for more info: https://golang.org/pkg/text/template/
type PromptTemplates struct {
	// Prompt is a text/template for the prompt label displayed on the left side of the prompt.
	Prompt string

	// Prompt is a text/template for the prompt label when IsConfirm is set as true.
	Confirm string

	// Valid is a text/template for the prompt label when the value entered is valid.
	Valid string

	// Invalid is a text/template for the prompt label when the value entered is invalid.
	Invalid string

	// Success is a text/template for the prompt label when the user has pressed entered and the value has been
	// deemed valid by the validation function. The label will keep using this template even when the prompt ends
	// inside the console.
	Success string

	// Prompt is a text/template for the prompt label when the value is invalid due to an error triggered by
	// the prompt's validation function.
	ValidationError string
}

// Run executes the prompt. Its displays the label and default value if any, asking the user to enter a value.
// Run will keep the prompt alive until it has been canceled from the command prompt or it has received a valid
// value. It will return the value and an error if any occurred during the prompt's execution.
func (p *Prompt) Run() (string, error) {
	p.prepareTemplates()

	c := &readline.Config{
		Stdin:          p.stdin,
		Stdout:         p.stdout,
		EnableMask:     p.Mask != 0,
		MaskRune:       p.Mask,
		HistoryLimit:   -1,
		VimMode:        p.IsVimMode,
		UniqueEditLine: true,
	}

	err := c.Init()
	if err != nil {
		return "", err
	}

	rl, err := readline.NewEx(c)
	if err != nil {
		return "", err
	}
	// we're taking over the cursor,  so stop showing it.
	rl.Write([]byte(hideCursor))
	sb := screenbuf.New(rl)

	validFn := func(x string) error {
		return nil
	}
	if p.Validate != nil {
		validFn = p.Validate
	}

	var inputErr error
	input := p.Default
	if p.IsConfirm {
		input = ""
	}
	eraseDefault := input != "" && !p.AllowEdit
	cur := NewCursor(input, p.Pointer, eraseDefault)

	listen := func(input []rune, pos int, key rune) ([]rune, int, bool) {
		_, _, keepOn := cur.Listen(input, pos, key)
		err := validFn(cur.Get())
		var prompt []byte

		if err != nil {
			prompt = frmt.Render(p.Templates.Invalid, p.Label)
		} else {
			prompt = frmt.Render(p.Templates.Valid, p.Label)
			if p.IsConfirm {
				prompt = frmt.Render(p.Templates.Prompt, p.Label)
			}
		}

		echo := cur.Format()
		if p.Mask != 0 {
			echo = cur.FormatMask(p.Mask)
		}

		prompt = append(prompt, []byte(echo)...)
		sb.Reset()
		sb.Write(prompt)
		if inputErr != nil {
			validation := frmt.Render(p.Templates.ValidationError, inputErr)
			sb.Write(validation)
			inputErr = nil
		}
		sb.Flush()
		return nil, 0, keepOn
	}

	c.SetListener(listen)

	for {
		_, err = rl.Readline()
		inputErr = validFn(cur.Get())
		if inputErr == nil {
			break
		}

		if err != nil {
			break
		}
	}

	if err != nil {
		switch err {
		case readline.ErrInterrupt:
			err = ErrInterrupt
		case io.EOF:
			err = ErrEOF
		}
		if err.Error() == "Interrupt" {
			err = ErrInterrupt
		}
		sb.Reset()
		sb.WriteString("")
		sb.Flush()
		rl.Write([]byte(showCursor))
		rl.Close()
		return "", err
	}

	echo := cur.Format()
	if p.Mask != 0 {
		echo = cur.FormatMask(p.Mask)
	}

	prompt := frmt.Render(p.Templates.Success, p.Label)
	prompt = append(prompt, []byte(echo)...)

	if p.IsConfirm {
		lowerDefault := strings.ToLower(p.Default)
		if strings.ToLower(cur.Get()) != "y" && (lowerDefault != "y" || (lowerDefault == "y" && cur.Get() != "")) {
			prompt = frmt.Render(p.Templates.Invalid, p.Label)
			err = ErrAbort
		}
	}

	sb.Reset()
	sb.Write(prompt)
	sb.Flush()
	rl.Write([]byte(showCursor))
	rl.Close()

	return cur.Get(), err
}

func (p *Prompt) prepareTemplates() {
	tpls := p.Templates
	if tpls == nil {
		tpls = &PromptTemplates{}
	}

	if p.IsConfirm {
		if tpls.Confirm == "" {
			confirm := "y/N"
			if strings.ToLower(p.Default) == "y" {
				confirm = "Y/n"
			}
			tpls.Confirm = fmt.Sprintf(`{{ iconQ | bold }} {{ . | bold }}? {{ "[%s]" | faint }} `, confirm)
		}
	} else {
		if tpls.Prompt == "" {
			tpls.Prompt = fmt.Sprintf(`{{iconQ | bold }} {{ . | bold }}{{":" | bold}} `)
		}
	}

	if tpls.Valid == "" {
		tpls.Valid = fmt.Sprintf(`{{iconGood | bold}} {{ . | bold }}{{":" | bold}} `)
	}

	if tpls.Invalid == "" {
		tpls.Invalid = fmt.Sprintf(`{{ iconBad | bold }} {{ . | bold }}{{":" | bold}} `)
	}

	if tpls.ValidationError == "" {
		tpls.ValidationError = `{{ ">>" | red }} {{ . | red }}`
	}

	if tpls.Success == "" {
		tpls.Success = fmt.Sprintf(`{{ . | faint }}{{":" | faint }}`)
	}

	p.Templates = tpls
}
