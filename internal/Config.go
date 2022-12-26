package internal

import (
	"errors"
	"os"
	"os/exec"

	"github.com/pelletier/go-toml/v2"
)

func ConfigFilename() (string, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	config_filename := homedir + "/.config/imap_filter_go.toml"
	return config_filename, nil
}

type AccountConfig struct {
	Host         string
	Inbox        string
	Password     string
	PasswordEval string
	Port         int
	Ssl          bool
	User         string
	Rules        []FilterRule
}

func NewAccountConfigDefaults() *AccountConfig {
	return &AccountConfig{
		Host:         "www.example.com",
		Inbox:        "INBOX",
		Password:     "",
		PasswordEval: "",
		Port:         993,
		Ssl:          true,
		User:         "user",
		Rules:        make([]FilterRule, 0),
	}
}

func (res *AccountConfig) parse(data map[string]interface{}, normalize bool) {
	valb, ok := data["ssl"].(bool)
	if ok {
		res.Ssl = valb
	}
	vals, ok := data["host"].(string)
	if ok {
		res.Host = vals
	}
	vals, ok = data["inbox"].(string)
	if ok {
		res.Inbox = vals
	}
	vals, ok = data["password"].(string)
	if ok {
		res.Password = vals
	}
	vals, ok = data["password_eval"].(string)
	if ok {
		res.PasswordEval = vals
	}
	vali, ok := data["port"].(int)
	if ok {
		res.Port = vali
	}
	vals, ok = data["user"].(string)
	if ok {
		res.User = vals
	}
	// log.Println(data["rules"].([]string))
	valsa, ok := _parse_string_array(data["rules"])
	if ok {
		parse_rules := func(data []string) []FilterRule {
			res := make([]FilterRule, len(data))
			for idx, rule := range data {
				res[idx] = NewFilterRule(rule)
			}
			return res
		}
		encode_rules := func(data []FilterRule) []string {
			res := make([]string, len(data))
			for idx, rule := range data {
				res[idx] = rule.Encode()
			}
			return res
		}

		res.Rules = parse_rules(valsa)

		if normalize {
			data["rules"] = encode_rules(res.Rules)
		}
	}
}

func _parse_string_array(data interface{}) ([]string, bool) {
	valsa, ok := data.([]interface{})
	if !ok {
		return nil, false
	}

	res := make([]string, 0)
	for _, r := range valsa {
		r_, ok2 := r.(string)
		if !ok2 {
			return nil, false
		}
		res = append(res, r_)
	}
	return res, true
}

func readConfigRaw() (map[string]interface{}, error) {
	LVerbose().Println("read config")

	config_filename, err := ConfigFilename()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(config_filename); errors.Is(err, os.ErrNotExist) {
		LError().Println("could not find config file")
		return nil, nil
	}

	config_, err := os.Open(config_filename)
	if err != nil {
		return nil, err
	}

	d := toml.NewDecoder(config_)
	d.DisallowUnknownFields()
	s := map[string]interface{}{}
	err = d.Decode(&s)
	if err != nil {
		return nil, err
	}

	// log.Println(s)
	return s, nil
}

func WriteConfigRaw(data map[string]interface{}) error {
	LVerbose().Println("write config")

	config_filename, err := ConfigFilename()
	if err != nil {
		return err
	}

	f, err := os.Create(config_filename)
	if err != nil {
		return err
	}
	enc := toml.NewEncoder(f)
	if err != nil {
		return err
	}
	enc.SetArraysMultiline(true)
	enc.Encode(data)

	// def get_account_name() -> str:
	//     account_names = config.account_names()

	//     account = cast(str, args.account) if "account" in vars(args) else None
	//     if account is None and len(account_names) > 0:
	//         account = account_names[0]
	//         log.info('using default account "%s"' % account)

	//     if account is None:
	//         log.error('no accounts defined in "%s"' % CONFIG_PATH)
	//         exit(-1)
	//     else:
	//         return account

	// account_name = get_account_name()
	// account_cfg = config.account(account_name)

	return nil
}

func resolvePassword(password_eval string) (string, error) {
	LVerbose().Println("evaluate password: '" + password_eval + "'")

	cmd := exec.Command("sh", "-c", password_eval)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	res := string(out)
	// LInfo().Println("password: '" + res + "'")
	return res, err
}

func ReadConfig(normalize bool) (map[string]interface{}, *AccountConfig, error) {
	data, err := readConfigRaw()
	if err != nil {
		return data, nil, err
	}

	accounts_, valid := data["account"].(map[string]interface{})
	if !valid {
		return data, nil, errors.New("config format invalid")
	}

	accounts := make([]string, 0)
	for account := range accounts_ {
		accounts = append(accounts, account)
	}
	LVerbose().Println("accounts: ", accounts)

	if len(accounts) == 0 {
		return data, nil, errors.New("warning no account found.")
	}
	if len(accounts) > 1 {
		LWarning().Println("warning more than one account not supported. Picking first.")
	}
	account := accounts[0]

	settings, valid := accounts_[account].(map[string]interface{})
	if !valid {
		return data, nil, errors.New("config format invalid")
	}

	res := NewAccountConfigDefaults()
	res.parse(settings, normalize)

	// log.Println(res)

	if normalize {
		WriteConfigRaw(data)
	}

	if res.PasswordEval != "" {
		res.Password, err = resolvePassword(res.PasswordEval)
		if err != nil {
			return data, nil, err
		}
	}

	LVerbose().Printf("loaded %d rules\n", len(res.Rules))

	return data, res, nil
}
