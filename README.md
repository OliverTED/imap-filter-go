## Imap-Filter: Filter an IMAP Mailbox from cli/mutt/etc

'imap-filter-go' is a simple cli tool to log into an imap server and rule based filter the email messages (such as moving all messages from a certain sender to a specific folder)

### Install

Install via golang's install directive:

```
go install github.com/OliverTED/imap-filter-go@latest
```

The executable will be in '~/go/bin'.

### Quickstart


To quickly apply a email filter:

```
imap_filter_go --verbose \
    --host imap.myserver.com --ssl --port 993 \
    --user testuser@myserver.com \
    --password-eval 'echo mypassword' \
    execute --dry-run \
    --rule 'from:*@dailycodingproblem.com:INBOX/learn-coding' \
    --rule 'from:*@indiehackers.com:INBOX/indiehackers'
```


This will move all emails where from matches
'*@dailycodingproblem.com' to 'INBOX/learn-coding' and anything where
from matches '*@indiehackers.com' to 'INBOX/indiehackers'.  Actually
it will only print what it would do if you did not specify 'dry-run'.

All those parameters can be stored inside the configuration
'~/.config/imap_filter_go.toml':

```
[account.default]
host = 'imap.myserver.com
inbox = 'INBOX'
password_eval = 'echo mypassword'
ssl = true
user = 'testuser@myserver.com'
port = 993
rules = [
  'from:*@dailycodingproblem.com:INBOX/learn-coding',
  'from:*@indiehackers.com:INBOX/indiehackers',
]
```



### Mutt

Put these lines into your 'muttrc' to have 'Mod+f' start the imap-filter and 'Mod-Shift-f' add a new rule:

```
macro index \ef "<shell-escape>term -e bash -c \"imap_filter_go execute; read\" &<enter>"
macro index \eF "<shell-escape>term -e bash -c \"imap_filter_go add-rule --interactive\" &<enter>"
```

### Development


Sources can be downloaded from github.  Then ```go install ...``` installs the current sources onto your system.  Push requests are welcome.
