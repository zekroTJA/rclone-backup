package main

import (
	"fmt"
	"os"

	"github.com/alexflint/go-arg"
	"github.com/logrusorgru/aurora/v4"
	"github.com/zekrotja/rclone-backup/pkg/list"
	"github.com/zekrotja/rclone-backup/pkg/rclone"
)

type Args struct {
	List    string          `arg:"-l,--list,required" help:"The rule list of resources to backup"`
	Target  string          `arg:"-t,--target,required" help:"The rclone target config to backup to"`
	Dry     bool            `arg:"-d" help:"Only print files to be backed up without transferring them"`
	Limit   *string         `arg:"-l,--limit" help:"Bandwidth limit for file upload"`
	Mode    rclone.SyncMode `arg:"-m,--mode" default:"sync" help:"The rclone replication mode to be used"`
	NoColor bool            `arg:"--no-colors" help:"Suppress color output"`

	RclonePath *string `arg:"--rclone-path,env:RB_RCLONE_PATH" help:"The path to the rclone binary"`
}

func main() {
	var args Args
	arg.MustParse(&args)

	au := aurora.New(aurora.WithColors(!args.NoColor))

	printErr := func(format string, args ...any) {
		_, _ = fmt.Fprintf(os.Stderr, "%s %s\n", au.Bold(au.Red("error:")),
			fmt.Sprintf(format, args...))
	}

	listFile, err := os.Open(args.List)
	if err != nil {
		printErr("failed to open list file: %s", err)
		os.Exit(1)
	}

	lst, err := list.Unmarshal(listFile)
	if err != nil {
		printErr("failed to parse list: %s", err)
		os.Exit(1)
	}

	rc, err := rclone.New(args.RclonePath)
	if err != nil {
		printErr("failed to initialize rclone: %s", err)
		os.Exit(1)
	}

	for _, entry := range lst {
		fmt.Println(au.Cyan(fmt.Sprintf("Backing up %s ...", entry.Path)))
		err = rc.Sync(args.Mode, entry.Path, args.Target, entry.Args, args.Dry, args.Limit)
		if err != nil {
			printErr("failed to sync: %s", err)
		}
	}
}
