import subprocess
import argparse
import json
import shutil


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

    return p.parse_args()


def dir_to_target_path(v: str) -> str:
    v = v.replace("\\", "/").replace(":", "")
    if v.startswith('/'):
        return v[1:]
    return v


def sync(target: str, dir: str, excludes: list[str], dry: bool, limit: str):
    cmd = ["rclone", "sync", "-v", dir,
           f"{target}:{dir_to_target_path(dir)}"]
    cmd.extend([f"--exclude={p}" for p in excludes])
    if dry:
        cmd.append('--dry-run')
    if limit:
        cmd.extend(['--bwlimit', limit])
    subprocess.call(cmd)


def get_configs() -> map:
    cmd = ["rclone", "config", "dump"]
    output = subprocess.check_output(cmd)
    return json.loads(output)


def main():
    args = parse_args()

    if not shutil.which("rclone"):
        print("Error: rclone is not installed")
        exit(1)

    if not get_configs().get(args.target):
        print("Error: could not find any rclone "
              f"config with the name {args.target}")
        exit(1)

    with open(args.list, 'r') as f:
        for line in f.readlines():
            line = line.replace("\n", "").replace("\r", "")
            if len(line) == 0 or line[0] == "#":
                continue
            split = line.replace("\n", "").replace("\r", "").split(",")
            dir = split[0]
            excludes = split[1:]
            print(f"----- BACKING UP {dir}")
            sync(args.target, dir, excludes, args.dry, args.limit)


if __name__ == "__main__":
    main()
