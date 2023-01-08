package internal

import (
	"fmt"
	"strings"

	"github.com/emersion/go-imap"
	"github.com/rivo/tview"
)

type State struct {
	app                *tview.Application
	selected_message   *MyMessage
	selected_folder    string
	selected_attribute string
	selected_pattern   string
	messages           []*MyMessage
	folders            []string
	Preselect          string
	Filter             string
}

// var attributes = []string{"from", "to", "replyto", "cc", "bcc", "subject"}

func EmailPatterns(addr *imap.Address) []string {
	res := []string{
		fmt.Sprintf("%s@%s", addr.MailboxName, addr.HostName),
		fmt.Sprintf("*@%s", addr.HostName),
	}
	if addr.PersonalName != "" {
		res = append(res,
			fmt.Sprintf("%s <%s@%s>", addr.PersonalName, addr.MailboxName, addr.HostName),
		)
	}
	return res
}

func PossibleFilters(message *MyMessage) []string {
	res := make([]string, 0)

	for _, addr := range message.Envelope.From {
		for _, pattern := range EmailPatterns(addr) {
			res = append(res, "from:"+pattern)
		}
	}
	for _, addr := range message.Envelope.To {
		for _, pattern := range EmailPatterns(addr) {
			res = append(res, "to:"+pattern)
		}
	}
	for _, addr := range message.Envelope.ReplyTo {
		for _, pattern := range EmailPatterns(addr) {
			res = append(res, "replyto:"+pattern)
		}
	}
	for _, addr := range message.Envelope.Cc {
		for _, pattern := range EmailPatterns(addr) {
			res = append(res, "cc:"+pattern)
		}
	}
	res = append(res, "subject:*"+message.Envelope.Subject+"*")
	return res
}

func (state *State) showAddRule(folders []string, messages []*MyMessage) tview.Primitive {
	SortMessages(messages)
	SortFolders(folders)

	form := tview.NewForm()

	form.AddTextView("", "Add new rule", 20, 1, false, false)

	messages_strs := MapArray(messages, func(idx int, msg **MyMessage) string { return MessageToString(idx, *msg) })

	form.AddDropDown("Select Message", messages_strs, GetIndex(messages, state.selected_message, 0), func(val string, idx int) {
		if state.selected_message != messages[idx] {
			state.selected_message = messages[idx]
			state.Render()
		}
	})

	form.AddDropDown("Move to folder", folders, GetIndex(folders, state.selected_folder, 0), func(val string, idx int) {
		if state.selected_folder != val {
			state.selected_folder = val
			state.Filter = state.Preselect + ":" + state.selected_folder
			state.Render()
		}
	})

	possible_filters := PossibleFilters(state.selected_message)
	form.AddDropDown("Predefineds", possible_filters, GetIndex(possible_filters, state.Preselect, 0), func(val string, idx int) {
		if state.Preselect != val {
			state.Preselect = val
			state.Filter = state.Preselect + ":" + state.selected_folder
			state.Render()
		}
	})

	form.AddTextView("Matching Messages", "", 80, 10, false, false)
	preview := form.GetFormItemByLabel("Matching Messages").(*tview.TextView)
	update_preview := func() {
		matching_messages := []string{}
		res := NewFilterRule(state.Filter)
		idx := 1
		if res == nil {
			matching_messages = append(matching_messages, "malformed rule")
		} else {
			for _, m := range state.messages {
				if res.Matches(m) {
					matching_messages = append(matching_messages, MessageToString(idx, m))
					idx += 1
				}
			}
			if len(matching_messages) == 0 {
				matching_messages = append(matching_messages, "- none -")
			}
		}
		preview.SetText(strings.Join(matching_messages, "\n"))
	}
	update_preview()

	form.AddTextArea("Filter", state.Filter, 80, 3, 160, func(val string) {
		if state.Filter != val {
			state.Filter = val
			update_preview()
			// state.Render()
		}
	})

	form.AddButton("Add Rule to Config", func() {
		AddNewRuleToConfig(state.Filter)
		state.app.Stop()
	})

	return form
}

func (s *State) Render() {
	cell := s.showAddRule(s.folders, s.messages)
	s.app.SetRoot(cell, true)
	return
}

func (r *MyApp) CmdAddRuleInteractive() error {
	conn := r.NewConnection()

	err := conn.connect()
	if err != nil {
		return err
	}
	defer conn.close()

	folders := conn.FoldersCached()

	err, messages := conn.yield_messages(r.Folder)
	if err != nil {
		return err
	}

	app := tview.NewApplication()

	state := State{messages: messages, folders: folders, app: app}
	state.Render()

	// app.EnableMouse(true)

	if err := app.Run(); err != nil {
		return err
	}

	return nil
}
