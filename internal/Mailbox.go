package internal

import (
	"errors"

	"fmt"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

type Mailbox struct {
	host        string
	port        int
	ssl         bool
	user        string
	password    string
	sep         string
	_connection *client.Client

	_folders *[]string
}

func (m *Mailbox) connect() error {
	url := fmt.Sprintf("%s:%d", m.host, m.port)
	LVerbose().Printf("Connecting to %s\n", url)

	conn, err := client.DialTLS(url, nil)
	m._connection = conn
	if err != nil {
		return err
		// log.Fatal(err)
	}
	LVerbose().Println("Connected")

	if err := m._connection.Login(m.user, m.password); err != nil {
		return err
		// log.Fatal(err)
	}
	LVerbose().Println("Logged in")

	return nil
}

func (m *Mailbox) close() error {
	err := m._connection.Logout()

	err = m._connection.Close()
	if err != nil {
		return err
	}

	m._connection = nil
	return nil
}

type MyMessage struct {
	Uid      uint32
	Envelope *imap.Envelope
	Flags    []string
	Folder   string
}

func (m MyMessage) String() string {
	// mid := m.envelope.MessageId
	from := m.Envelope.From[0]
	subject := m.Envelope.Subject

	return fmt.Sprintf("%s/%s", from, subject)
}

func (m *Mailbox) yield_messages(folder string) (error, []*MyMessage) {
	if m._connection == nil {
		return errors.New("error: connection closed"), nil
	}
	LVerbose().Println("IMAP Change to folder:", folder)
	mbox, err := m._connection.Select(folder, false)
	if err != nil {
		return fmt.Errorf("IMAP: error changing folder '%s' %w", folder, err), nil
	}
	// log.Println("Flags for INBOX:", mbox.Flags)

	from := uint32(1)
	to := mbox.Messages

	seqset := new(imap.SeqSet)
	seqset.AddRange(from, to)

	messages := make(chan *imap.Message)
	done := make(chan error, 1)
	go func() {
		done <- m._connection.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope, imap.FetchFlags, imap.FetchUid}, messages)
	}()

	res := make([]*MyMessage, 0)

	for val := range messages {
		res = append(res,
			&MyMessage{Envelope: val.Envelope, Folder: folder, Flags: val.Flags, Uid: val.Uid},
		)
	}

	if err := <-done; err != nil {
		return err, nil
	}

	LVerbose().Println("IMAP: query messages: ", len(res))

	return nil, res
}

func (m *Mailbox) list_folders() []string {
	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)

	go func() {
		done <- m._connection.List("", "*", mailboxes)
	}()

	res := make([]string, 0)

	for val := range mailboxes {
		res = append(res, val.Name)
	}

	if err := <-done; err != nil {
		LError().Fatal(err)
	}

	LVerbose().Println("IMAP: list_folders: ", len(res))

	return res
}

func (m *Mailbox) FoldersCached() []string {
	if m._folders == nil {
		res := m.list_folders()
		m._folders = &res
	}
	return *m._folders
}
