#!/usr/bin/env python

import subprocess
import argparse
import json
import shutil
from log import disable_colors, error, highlight, success, fail


def parse_args():
    p = argparse.ArgumentParser()

    p.add_argument('-l', '--list', type=str, required=True,
                   help="The list of backup resources.")
    p.add_argument('-t', '--target', type=str, required=True,
                   help="The rclone target config.")
    p.add_argument('--dry', action='store_true',
                   help="Only print files to be backed up without "
                        "transferring them.")
    p.add_argument('--limit', type=str,
                   help="Bandwidth limit.")
    p.add_argument('--mode', type=str, default="sync",
                   choices=["sync", "copy"],
                   help="The rclone replication mode to use.")
    p.add_argument("--no-colors", action='store_true',
                   help="Supress colorful output")

    return p.parse_args()


def dir_to_target_path(v: str) -> str:
    v = v.replace("\\", "/").replace(":", "")
    if v.startswith('/'):
        return v[1:]
    return v


def join_target(target, dir):
    if ':' in target:
        return f"{target}/{dir_to_target_path(dir)}"
    return f"{target}:{dir_to_target_path(dir)}"


def build_exclude_flags(excludes: list[str]) -> list[str]:
    flags = []
    for exclude in excludes:
        if exclude.startswith('max-size='):
            flags.append(f"--max-size={exclude[len('max-size='):]}")
        else:
            flags.append(f"--exclude={exclude}")
    return flags


def sync(
    mode: str, target: str, dir: str,
    excludes: list[str], dry: bool, limit: str,
) -> int:
    cmd = ["rclone", mode, "-v", dir, join_target(target, dir)]
    cmd.extend(build_exclude_flags(excludes))
    if dry:
        cmd.append('--dry-run')
    if limit:
        cmd.extend(['--bwlimit', limit])
    return subprocess.call(cmd)


def get_configs() -> map:
    cmd = ["rclone", "config", "dump"]
    output = subprocess.check_output(cmd)
    return json.loads(output)


def main() -> int:
    args = parse_args()

    if args.no_colors:
        disable_colors()

    if not shutil.which("rclone"):
        error("rclone is not installed")
        return 1

    cfg = args.target.split(':')[0]
    if not get_configs().get(cfg):
        error("could not find any rclone "
              f"config with the name '{cfg}'")
        return 2

    results = []
    with open(args.list, 'r') as f:
        for line in f.readlines():
            line = line.replace("\n", "").replace("\r", "")
            if len(line) == 0 or line[0] == "#":
                continue
            split = line.replace("\n", "").replace("\r", "").split(",")
            dir = split[0]
            excludes = split[1:]
            highlight(f"=== BACKING UP {dir}")
            exit_code = sync(args.mode, args.target, dir,
                             excludes, args.dry, args.limit)
            results.append((exit_code, dir))

    results_successful = [p for (c, p) in results if c == 0]
    results_failed = [p for (c, p) in results if c != 0]

    print()

    if len(results_failed) == 0:
        success(
            f"All {len(results_successful)} backup targets "
            "finished successfully.")
        return 0

    fail(f"{len(results_failed)} backup tasks had "
         f"errors and {len(results_successful)} ended successfully.")

    print("\nSuccessful tasks:")
    for dir in results_successful:
        print(f"  {dir}")

    fail("\nFailed tasks:")
    for dir in results_failed:
        fail(f"  {dir}")
    return 3


if __name__ == "__main__":
    exit(main())
