package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli/v2"

	"imap_filter_go/internal"
)

func main() {
	type MyArgs struct {
		Account       string
		Verbose       bool
		LogFilename   string
		DebugImap     bool
		Interactive   bool
		Rule          string
		configActions []func(run *internal.MyApp)
		Action        func(run *internal.MyApp) error
		ExtraRules    []internal.FilterRule
	}

	args := MyArgs{
		configActions: make([]func(run *internal.MyApp), 0),
	}

	app := &cli.App{
		Name:  "imap-filter-go",
		Usage: "Filter your IMAP mailbox and sort your emails.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name: "account",
				Action: func(ctx *cli.Context, v string) error {
					args.Account = v
					return nil
				},
			},
			&cli.StringFlag{
				Name: "host",
				Action: func(ctx *cli.Context, v string) error {
					args.configActions = append(args.configActions, func(run *internal.MyApp) {
						run.ImapHost = v
					})
					return nil
				},
			},
			&cli.BoolFlag{
				Name:  "ssl",
				Value: false,
				Action: func(ctx *cli.Context, v bool) error {
					args.configActions = append(args.configActions, func(run *internal.MyApp) {
						run.ImapSsl = v
					})
					return nil
				},
			},
			&cli.BoolFlag{
				Name:  "verbose",
				Value: false,
				Action: func(ctx *cli.Context, v bool) error {
					args.Verbose = v
					return nil
				},
			},
			&cli.StringFlag{
				Name:  "logfilename",
				Value: "",
				Action: func(ctx *cli.Context, v string) error {
					args.LogFilename = v
					return nil
				},
			},
			&cli.BoolFlag{
				Name:  "debug-imap",
				Value: false,
				Action: func(ctx *cli.Context, v bool) error {
					args.DebugImap = v
					return nil
				},
			},
			&cli.IntFlag{
				Name: "port",
				Action: func(ctx *cli.Context, v int) error {
					args.configActions = append(args.configActions, func(run *internal.MyApp) {
						run.ImapPort = v
					})
					return nil
				},
			},
			&cli.StringFlag{
				Name: "user",
				Action: func(ctx *cli.Context, v string) error {
					args.configActions = append(args.configActions, func(run *internal.MyApp) {
						run.ImapUser = v
					})
					return nil
				},
			},
			&cli.StringFlag{
				Name: "password-eval",
				Action: func(ctx *cli.Context, v string) error {
					args.configActions = append(args.configActions, func(run *internal.MyApp) {
						run.ImapPasswordFunc = func() (string, error) {
							return internal.ResolvePassword(v)
						}
					})
					return nil
				},
			},
		},
		Action: func(ctx *cli.Context) error {
			if ctx.Args().Len() > 0 {
				return fmt.Errorf("invalid command: '%s'", strings.Join(ctx.Args().Slice(), " "))
			}
			cli.ShowAppHelp(ctx)
			return fmt.Errorf("missing command")
		},
		Commands: []*cli.Command{
			{
				Name:  "dry-run",
				Usage: "run, but only print actions",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name: "folder",
						Action: func(ctx *cli.Context, v string) error {
							args.configActions = append(args.configActions, func(run *internal.MyApp) {
								run.Folder = v
							})
							return nil
						},
					},
					&cli.StringFlag{
						Name: "rule",
						Action: func(ctx *cli.Context, v string) error {
							rule_, err := internal.NewFilterRule(v)
							if err != nil {
								return err
							}
							args.ExtraRules = append(args.ExtraRules, rule_)
							return nil
						},
					},
				},
				Action: func(ctx *cli.Context) error {
					if ctx.Args().Len() > 0 {
						return fmt.Errorf("invalid command: '%s'", strings.Join(ctx.Args().Slice(), " "))
					}
					args.Action = func(run *internal.MyApp) error {
						return run.CmdExecute(run.Connect, run.Folder, true, -1, append(run.Rules, args.ExtraRules...))
					}
					return nil
				},
			},
			{
				Name:  "execute",
				Usage: "execute imap filter",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name: "folder",
						Action: func(ctx *cli.Context, v string) error {
							args.configActions = append(args.configActions, func(run *internal.MyApp) {
								run.Folder = v
							})
							return nil
						},
					},
					&cli.StringFlag{
						Name: "rule",
						Action: func(ctx *cli.Context, v string) error {
							rule_, err := internal.NewFilterRule(v)
							if err != nil {
								return err
							}
							args.ExtraRules = append(args.ExtraRules, rule_)
							return nil
						},
					},
				},
				Action: func(ctx *cli.Context) error {
					if ctx.Args().Len() > 0 {
						return fmt.Errorf("invalid command: '%s'", strings.Join(ctx.Args().Slice(), " "))
					}
					args.Action = func(run *internal.MyApp) error {
						return run.CmdExecute(run.Connect, run.Folder, false, -1, append(run.Rules, args.ExtraRules...))
					}
					return nil
				},
			},
			{
				Name: "add-rule",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name: "folder",
						Action: func(ctx *cli.Context, v string) error {
							args.configActions = append(args.configActions, func(run *internal.MyApp) {
								run.Folder = v
							})
							return nil
						},
					},
					&cli.BoolFlag{
						Name:  "interactive",
						Value: false,
						Action: func(ctx *cli.Context, v bool) error {
							args.Interactive = v
							return nil
						},
					},
					&cli.StringFlag{
						Name: "rule",
						Action: func(ctx *cli.Context, v string) error {
							args.Rule = v
							return nil
						},
					},
				},
				Action: func(ctx *cli.Context) error {
					if ctx.Args().Len() > 0 {
						return fmt.Errorf("invalid command: '%s'", strings.Join(ctx.Args().Slice(), " "))
					}
					if args.Interactive && args.Rule != "" {
						return fmt.Errorf("--interactive not valid with --rule")
					}
					if args.Interactive {
						args.Action = func(run *internal.MyApp) error {
							return run.CmdAddRuleInteractive(run.Connect)
						}
					} else if args.Rule != "" {
						args.Action = func(run *internal.MyApp) error {
							return run.CmdAddRule(args.Rule)
						}
					} else {
						return fmt.Errorf("need either --interactive or --rule")
					}

					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		internal.LError().Fatal("ERROR: ", err)
		os.Exit(1)
	}
	cleanup := internal.SetupLogging(args.Verbose, args.LogFilename)
	defer cleanup()

	if args.Action == nil {
		args.Action = func(run *internal.MyApp) error {
			fmt.Println(app.Usage)
			return nil
		}
	}

	rawconfig, accountconfig, err := internal.ReadConfig(false)
	if err != nil {
		fmt.Println("error parsing config: ", err)
		os.Exit(1)
	}

	run := internal.NewMyApp(rawconfig, accountconfig, args.DebugImap)
	for _, action := range args.configActions {
		action(run)
	}

	err = args.Action(run)
	if err != nil {
		internal.LError().Fatal("ERROR: ", err)
	}
	internal.LInfo().Println("done")
}
