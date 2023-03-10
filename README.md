# rclone Backup Script

A very simple python script to backup important local directories using [rclone](https://rclone.org/).

## Usage

The script works with simple CSV lists to define what to backup and what not.

> example.csv
```
/home/user/importantdata,**/exclude1/**

# Upload photos and exclude raw files and
# exclude files that are larger than 1 GiB
/mnt/data/photos,**.ARW,**.CRW,max-size=1G
```

In the first column, you specify the full path to the location to backup. The following columns define exclude filter patterns. [Here](https://rclone.org/filtering/) you can find detailed documentation how to specify filters. There are also some special flags like `max-size` which you can use to specify additional filters passed to the rclone command.

You also need to configure an rclone target for your backups. Simply use `rclone config` to use the setup assistent or edit the rcon config file. This can be found via the command `rclone config file` (with this handy command you can directly edit it: `vim $(rclone config file | tail -1)`). The following example shows a target configuration for Backblaze B2.

```conf
[backblaze-b2]
type = b2
account = *****
key = *****

[backups]
type = alias
remote = backblaze-b2:my-backups
```

Then, simply execute the script passing the list file as well as the rcon target.

```
./backups/main.py -l example.csv -t backups
```

There are also other useful command line flags. Simply use the `--help` flag to get more info.

```
❯ ./backup/main.py --help
usage: main.py [-h] -l LIST -t TARGET [--dry] [--limit LIMIT] [--mode {sync,copy}] [--no-colors]

options:
  -h, --help            show this help message and exit
  -l LIST, --list LIST  The list of backup resources.
  -t TARGET, --target TARGET
                        The rclone target config.
  --dry                 Only print files to be backed up without transferring them.
  --limit LIMIT         Bandwidth limit.
  --mode {sync,copy}    The rclone replication mode to use.
  --no-colors           Supress colorful output
```

## Docker Image

You can also use the provided Docker image to do automated backups of your servers, for example.


Therefore, you need to bind the Rclone config as well as the backup list and the corresponding directories which should be backed up into the container.

```
docker run \
  --rm \
  -v $HOME/rclone.conf:/root/.config/rclone/rclone.conf:ro \
  -v $HOME/list.csv:/app/list.csv:ro \
  -v $HOME/myfiles:/backup/myfiles:ro \
  ghcr.io/zekrotja/rclone-backup -l list.csv -t backups
```