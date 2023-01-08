package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"imap_filter_go/internal"
)

func main() {
	type MyArgs struct {
		Account       string
		Verbose       bool
		LogFilename   string
		DebugImap     bool
		configActions []func(run *internal.MyApp)
		Action        func(run *internal.MyApp) error
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
				Name: "password",
				Action: func(ctx *cli.Context, v string) error {
					args.configActions = append(args.configActions, func(run *internal.MyApp) {
						run.ImapPassword = v
					})
					return nil
				},
			},
		},
		Action: func(ctx *cli.Context) error {
			if ctx.Args().Len() > 0 {
				return fmt.Errorf("invalid command: '%s'", ctx.Args().First())
			}
			return fmt.Errorf("no command")
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
				},
				Action: func(cCtx *cli.Context) error {
					if cCtx.Args().Len() > 0 {
						return fmt.Errorf("invalid argument: '%s'", cCtx.Args().First())
					}
					args.Action = func(run *internal.MyApp) error { return run.Execute(true, -1) }
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
				},
				Action: func(cCtx *cli.Context) error {
					if cCtx.Args().Len() > 0 {
						return fmt.Errorf("invalid argument: '%s'", cCtx.Args().First())
					}
					args.Action = func(run *internal.MyApp) error { return run.Execute(false, -1) }
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		internal.LError().Fatal(err)
	}
	cleanup := internal.SetupLogging(args.Verbose, args.LogFilename)
	defer cleanup()

	_, config, err := internal.ReadConfig(true)
	if err != nil {
		fmt.Println("error parsing config: ", err)
	}
	run := internal.NewMyApp(config, args.DebugImap)
	for _, action := range args.configActions {
		action(run)
	}

	err = args.Action(run)
	if err != nil {
		internal.LError().Fatal(err)
	}
	internal.LInfo().Println("done")
}
