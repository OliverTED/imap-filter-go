package internal

import (
	"fmt"

	"github.com/emersion/go-imap"
	"golang.org/x/term"
)

func max_len(str string, max_len int) string {
	if len(str) > max_len {
		str = str[:max_len]
	}
	return str
}
func short_describe_emails(addrs []*imap.Address) string {
	res := ""
	for _, addr := range addrs {
		if res != "" {
			res = res + ","
		}
		res = res + short_describe_email(addr)
	}
	return res
}
func short_describe_email(addr *imap.Address) string {
	if len(addr.PersonalName) > 4 {
		return max_len(addr.PersonalName, 20)
	}
	return fmt.Sprintf("%s@%s", addr.MailboxName, addr.HostName)
}

func short_describe_message(idx int, m *MyMessage) string {
	return fmt.Sprintf(
		"|%03d/%s/%s/%s/%s|",
		idx,
		max_len(short_describe_emails(m.Envelope.From), 15),
		max_len(m.Envelope.Subject, 20),
		max_len(short_describe_emails(m.Envelope.Cc), 15),
		max_len(short_describe_emails(m.Envelope.ReplyTo), 15),
	)
}

func select_message(messages []*MyMessage) (*MyMessage, error) {
	fmt.Println("Select Message which should be filtered:")
	message, err := select_item(messages, short_describe_message)
	if message == nil || err != nil {
		return nil, err
	}
	fmt.Println("Building filter for message " + short_describe_message(0, *message))
	return *message, nil
}

func select_folder(folders []string) (string, error) {
	fmt.Println("Select Folder:")
	t, err := select_string(folders)
	if err != nil {
		return "", err
	}

	fmt.Println("Building filter for message " + t)
	return t, nil
}

func (r *MyApp) InteractiveAddRule() error {
	// oldState, err := term.MakeRaw(0)
	oldState, err := term.GetState(0)
	if err != nil {
		panic(err)
	}
	// restore terminal afterwards (needs 'reset' afterwards, as
	// characters are not echoed any more)
	defer term.Restore(0, oldState)

	conn := r.NewConnection()

	err = conn.connect()
	if err != nil {
		return err
	}
	defer conn.close()

	folders := conn.FoldersCached()

	folder := "INBOX"
	err, messages := conn.yield_messages(folder)
	if err != nil {
		return err
	}

	_, err = select_message(messages)
	if err != nil {
		return err
	}

	_, err = select_folder(folders)
	if err != nil {
		return err
	}

	return nil
}
