package internal

import (
	"github.com/c-bata/go-prompt"
	"github.com/chzyer/readline"
)

func select_string_readline(values []string) (string, error) {
	values2 := make([]readline.PrefixCompleterInterface, len(values))
	for idx, message := range values {
		values2[idx] = readline.PcItem(message)
	}

	var completer = readline.NewPrefixCompleter(values2...)

	l, err := readline.NewEx(&readline.Config{
		AutoComplete: completer,
	})
	if err != nil {
		return "", err
	}
	defer l.Close()
	l.CaptureExitSignal()

	line, err := l.Readline()
	if err != nil {
		return "", err
	}
	return line, err
}

func select_string_prompt(values []string) (string, error) {
	s := []prompt.Suggest{}

	for _, message := range values {
		s = append(s,
			prompt.Suggest{Text: message})
	}

	completer := func(d prompt.Document) []prompt.Suggest {
		return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
	}

	t := prompt.Input("> ", completer)
	return t, nil
}

var select_string = select_string_prompt

// var select_string = select_string_readline

func select_item[K any](items []K, formatter func(idx int, item K) string) (*K, error) {
	strings := make([]string, len(items))
	for idx, item := range items {
		strings[idx] = formatter(idx, item)
	}

	res, err := select_string(strings)
	if err != nil {
		return nil, err
	}

	for idx, str := range strings {
		if str == res {
			return &items[idx], nil
		}
	}
	return nil, err
}
