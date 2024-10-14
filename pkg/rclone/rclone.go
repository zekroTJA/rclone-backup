package rclone

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type SyncMode string

const (
	ModeSync SyncMode = "sync"
	ModeCopy SyncMode = "copy"
)

func (t *SyncMode) UnmarshalText(b []byte) error {
	switch string(b) {
	case "sync":
		*t = ModeSync
	case "copy":
		*t = ModeCopy
	default:
		return fmt.Errorf("invalid sync mode: %s; options are: sync, copy", string(b))
	}
	return nil
}

type Rclone struct {
	rclonePath string
}

func New(path *string) (t *Rclone, err error) {
	t = new(Rclone)

	if path != nil {
		t.rclonePath = *path
	} else {
		t.rclonePath, err = exec.LookPath("rclone")
		if err != nil {
			return nil, err
		}
	}

	stat, err := os.Stat(t.rclonePath)
	if err != nil {
		return nil, err
	}
	if stat.IsDir() {
		return nil, errors.New("path to rclone is a directory")
	}

	return t, nil
}

func (t *Rclone) Sync(
	mode SyncMode,
	source string,
	target string,
	excludeArgs []string,
	dry bool,
	limit *string,
) (err error) {
	cmd := exec.Command(
		t.rclonePath, string(mode),
		"--verbose",
		source,
		joinTarget(target, source))

	cmd.Args = append(cmd.Args, excludeArgs...)

	if dry {
		cmd.Args = append(cmd.Args, "--dry-run")
	}

	if limit != nil {
		cmd.Args = append(cmd.Args, "--bwlimit", *limit)
	}

	cmd.Stdout = os.Stdout

	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func dirToTargetPath(v string) string {
	v = strings.NewReplacer("\\", "/", ":", "").Replace(v)
	if strings.HasPrefix(v, "/") {
		return v[1:]
	}
	return v
}

func joinTarget(target, dir string) string {
	if strings.ContainsRune(target, ':') {
		return fmt.Sprintf("%s/%s", target, dirToTargetPath(dir))
	}
	return fmt.Sprintf("%s:%s", target, dirToTargetPath(dir))
}
