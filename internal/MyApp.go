package internal

type MyApp struct {
	ImapUser     string
	ImapHost     string
	ImapSsl      bool
	ImapPort     int
	ImapPassword string
	Rules        []FilterRule
	Folder       string
	DebugImap    bool
}

func (r *MyApp) NewConnection() Mailbox {
	mb := Mailbox{
		user:     r.ImapUser,
		password: r.ImapPassword,
		ssl:      r.ImapSsl,
		host:     r.ImapHost,
		port:     r.ImapPort,
	}

	return mb
}

func NewMyApp(config *AccountConfig, DebugImap bool) *MyApp {
	r := MyApp{}

	r.DebugImap = DebugImap
	r.ImapUser = config.User
	r.ImapHost = config.Host
	r.ImapPort = config.Port
	r.ImapSsl = config.Ssl
	r.ImapPassword = config.Password
	r.Folder = config.Inbox
	r.Rules = config.Rules

	return &r
}
