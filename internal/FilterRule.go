package internal

import (
	"fmt"
	"strings"

	"github.com/emersion/go-imap"
	"github.com/gobwas/glob"
)

type FilterRuleRaw struct {
	raw string
}

func (r *FilterRuleRaw) Encode() string {
	return r.raw
}

func (r *FilterRuleRaw) ToString() string {
	return r.raw
}

func (r *FilterRuleRaw) Describe(m *MyMessage) string {
	return "invalid rule"
}

func (r *FilterRuleRaw) ShouldRun(m *MyMessage) bool {
	return false
}

func (r *FilterRuleRaw) Execute(mailbox *Mailbox, m *MyMessage) error {
	return nil
}

func (r *FilterRuleRaw) Matches(m *MyMessage) bool {
	return false
}

type FilterRuleMove struct {
	who      string
	pattern  string
	pattern_ glob.Glob
	folder   string
}

func (r *FilterRuleMove) Encode() string {
	return fmt.Sprintf("%s:%s:%s", r.who, r.pattern, r.folder)
}

func (r *FilterRuleMove) ToString() string {
	return fmt.Sprintf("%s:%s:%s", r.who, r.pattern, r.folder)
}

func (r *FilterRuleMove) Describe(m *MyMessage) string {
	from_ := m.Envelope.From
	var from imap.Address = imap.Address{PersonalName: "unknown", AtDomainList: "unknown", HostName: "unknown"}
	if len(from_) > 0 {
		from = *from_[0]
	}

	return fmt.Sprintf("moving to '%s': %s \"%s\" <%s>, %s", r.folder, m.Envelope.Date, from.PersonalName, from.AtDomainList+"@"+from.HostName, m.Envelope.Subject)
}

func (r *FilterRuleMove) ShouldRun(m *MyMessage) bool {
	return m.Folder != r.folder
}

func contains[K comparable](s []K, e K) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func (r *FilterRuleMove) Execute(mailbox *Mailbox, m *MyMessage) error {
	// mailbox._connection.Select(m.Folder, false)
	seqset := new(imap.SeqSet)
	seqset.AddNum(m.Uid)

	//         if not mailbox.has_folder(self.folder):
	//             log.info("warning: creating non-existing folder: %s" % self.folder)
	//             mailbox.create_folder(self.folder)
	if !contains(mailbox.FoldersCached(), r.folder) {
		mailbox._connection.Create(r.folder)
	}
	// LInfo().Println(1, seqset, m.Uid)
	err := mailbox._connection.UidMove(seqset, r.folder)
	if err != nil {
		LError().Println("error moving message: ", err)
		return err
	}
	return nil
}

func _matches_pattern1(addrs []*imap.Address, pattern glob.Glob) bool {
	for _, addr := range addrs {
		if _matches_pattern2(addr.PersonalName, pattern) {
			return true
		}
		email := fmt.Sprintf("%s@%s", addr.MailboxName, addr.HostName)
		if _matches_pattern2(email, pattern) {
			return true
		}
		format1 := fmt.Sprintf("%s <%s>", addr.PersonalName, email)
		if _matches_pattern2(format1, pattern) {
			return true
		}
		format2 := fmt.Sprintf("\"%s\" <%s>", addr.PersonalName, email)
		if _matches_pattern2(format2, pattern) {
			return true
		}
	}
	return false
}

func _matches_pattern2(data string, pattern glob.Glob) bool {
	if pattern.Match(data) {
		// log.Println(data, pattern)
		return true
	}
	return false
}

func _contains_lower(data []string, pattern string) bool {
	for _, d := range data {
		if strings.ToLower(d) == pattern {
			return true
		}
	}
	return false
}

func (r *FilterRuleMove) Matches(m *MyMessage) bool {
	if _contains_lower(m.Flags, "\\flagged") {
		// LVerbose().Println("flagged: ", m)
		// log.Println(r.who, m.Envelope, r.pattern_, m.Flags)
		return false
	}

	if r.who == "from" {
		return _matches_pattern1(m.Envelope.From, r.pattern_)
	}
	if r.who == "to" {
		return _matches_pattern1(m.Envelope.To, r.pattern_)
	}
	if r.who == "replyto" {
		return _matches_pattern1(m.Envelope.ReplyTo, r.pattern_)
	}
	if r.who == "cc" {
		return _matches_pattern1(m.Envelope.Cc, r.pattern_)
	}
	if r.who == "bcc" {
		return _matches_pattern1(m.Envelope.Bcc, r.pattern_)
	}
	if r.who == "subject" {
		return _matches_pattern2(m.Envelope.Subject, r.pattern_)
	}

	LError().Println("malformed rule: ", r.who, r.Encode())

	return false
}

type FilterRule interface {
	ToString() string
	Encode() string
	Describe(m *MyMessage) string
	ShouldRun(m *MyMessage) bool
	Execute(mailbox *Mailbox, m *MyMessage) error
	Matches(m *MyMessage) bool
}

func NewFilterRule(raw string) FilterRule {
	res := NewFilterRuleMove(raw)
	if res != nil {
		return res
	}

	LError().Println("malformed rule: ", raw)
	return &FilterRuleRaw{raw: raw}
}

func NewFilterRuleMove(raw string) *FilterRuleMove {
	parts := strings.Split(raw, ":")
	if len(parts) != 3 {
		return nil
	}

	who, pattern, folder := parts[0], parts[1], parts[2]
	who = strings.ToLower(who)

	// folder = "TEST/" + folder
	// pattern = strings.Replace(pattern, ".", "/", -1) // todo remove
	// folder = strings.Replace(folder, ".", "/", -1) // todo remove

	if !(who == "from" || who == "to" || who == "replyto" || who == "cc" || who == "bcc" || who == "subject") {
		return nil
	}

	pattern_, err := glob.Compile(pattern)
	if err != nil {
		return nil
	}

	//     if who not in ["from", "to", "replyto", "cc", "bcc", "subject"]:
	return &FilterRuleMove{who: who, pattern: pattern, folder: folder, pattern_: pattern_}
}
