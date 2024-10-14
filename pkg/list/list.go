package list

import (
	"bufio"
	"io"
	"strings"
)

type Entry struct {
	Path string
	Args []string
}

func Unmarshal(r io.Reader) (res []Entry, err error) {
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, err
		}

		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			continue
		}

		if line[0] == '#' {
			continue
		}

		split := strings.Split(line, ",")
		path := split[0]
		args := split[1:]

		res = append(res, Entry{
			Path: path,
			Args: buildExcludeFlags(args),
		})
	}

	return res, nil
}

func buildExcludeFlags(args []string) []string {
	res := make([]string, 0, len(args))
	for _, arg := range args {
		arg = strings.TrimSpace(arg)
		if strings.HasPrefix(arg, "max-size=") {
			res = append(res, "--"+arg)
		} else {
			res = append(res, "--exclude="+arg)
		}
	}
	return res
}
