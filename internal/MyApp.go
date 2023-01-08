package internal

import "os"

type MyApp struct {
	ImapUser         string
	ImapHost         string
	ImapSsl          bool
	ImapPort         int
	ImapPasswordFunc func() (string, error)
	Rules            []FilterRule
	Folder           string
	DebugImap        bool
	RawConfig        *Config
	AccountName      string
}

func (r *MyApp) newConnection() (*Mailbox, error) {
	password, err := r.ImapPasswordFunc()
	if err != nil {
		return nil, err
	}
	mb := Mailbox{
		user:     r.ImapUser,
		password: password,
		ssl:      r.ImapSsl,
		host:     r.ImapHost,
		port:     r.ImapPort,
	}

	return &mb, nil
}
func (r *MyApp) Connect() (*Mailbox, error) {
	conn, err := r.newConnection()
	if err != nil {
		return nil, err
	}

	err = conn.connect()
	if err != nil {
		return nil, err
	}

	if r.DebugImap {
		conn._connection.SetDebug(os.Stdout)
	}
	return conn, nil
}

func NewMyApp(rawconfig *Config, config *AccountConfig, DebugImap bool) *MyApp {
	r := MyApp{}

	r.RawConfig = rawconfig
	r.AccountName = config.AccountName

	r.DebugImap = DebugImap
	r.ImapUser = config.User
	r.ImapHost = config.Host
	r.ImapPort = config.Port
	r.ImapSsl = config.Ssl
	r.ImapPasswordFunc = config.GetPassword
	r.Folder = config.Inbox
	r.Rules = config.Rules

	return &r
}
