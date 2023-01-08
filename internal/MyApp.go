package internal

type MyApp struct {
	ImapUser     string
	ImapHost     string
	ImapSsl      bool
	ImapPort     int
	ImapPassword func() (string, error)
	Rules        []FilterRule
	Folder       string
	DebugImap    bool
	RawConfig    *Config
	AccountName  string
}

func (r *MyApp) NewConnection() (*Mailbox, error) {
	password, err := r.ImapPassword()
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

func NewMyApp(rawconfig *Config, config *AccountConfig, DebugImap bool) *MyApp {
	r := MyApp{}

	r.RawConfig = rawconfig
	r.AccountName = config.AccountName

	r.DebugImap = DebugImap
	r.ImapUser = config.User
	r.ImapHost = config.Host
	r.ImapPort = config.Port
	r.ImapSsl = config.Ssl
	r.ImapPassword = config.GetPassword
	r.Folder = config.Inbox
	r.Rules = config.Rules

	return &r
}
