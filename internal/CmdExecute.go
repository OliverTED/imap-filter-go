package internal

import (
	"os"
	"time"
)

func (r *MyApp) CmdExecute(dry_run bool, repeatAfterSeconds int) error {
	type TodoAction struct {
		rule    FilterRule
		message *MyMessage
	}

	if repeatAfterSeconds > 0 && repeatAfterSeconds < 60 {
		LWarning().Println("RepeatAfterSeconds to small; increasing to 60 seconds")
		repeatAfterSeconds = 60
	}

	for {
		conn := r.NewConnection()

		err := conn.connect()
		if err != nil {
			return err
		}
		defer conn.close()

		if r.DebugImap {
			conn._connection.SetDebug(os.Stdout)
		}

		err, messages := conn.yield_messages(r.Folder)
		if err != nil {
			return err
		}

		if len(messages) == 0 {
			LInfo().Println("no messages in the inbox, ", r.Folder)
			return nil
		}

		//         account = self.config.account(account_name)
		//         folders = folders or account.get_inbox_names()
		todo_actions := make([]TodoAction, 0)
		for _, message := range messages {
			possible_actions := make([]FilterRule, 0)
			for _, rule := range r.Rules {
				if !IsFlagged(message) && rule.Matches(message) && rule.ShouldRun(message) {
					possible_actions = append(possible_actions, rule)
				}
			}

			if len(possible_actions) == 0 {
				// log.Println("rule(s) matched no message")
				continue
			}

			if len(possible_actions) > 1 {
				LWarning().Println("warning: multiple rules for messages!")
				for _, rule := range possible_actions {
					LWarning().Println(rule.ToString())
				}
			}

			if len(possible_actions) >= 1 {
				action := possible_actions[0]
				todo_actions = append(todo_actions, TodoAction{rule: action, message: message})
			}
		}

		for _, action := range todo_actions {
			LInfo().Println(action.rule.Describe(action.message))
			if !dry_run {
				action.rule.Execute(&conn, action.message)
			}
		}

		if repeatAfterSeconds < 0 {
			break
		} else {
			time.Sleep(time.Duration(repeatAfterSeconds) * time.Second)
		}
	}

	return nil
}
