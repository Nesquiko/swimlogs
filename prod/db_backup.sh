#!/usr/bin/env bash
#
# took inspiration from https://github.com/shakibamoshiri/bash-CLI-template/blob/master/subcommand.sh

backups_dir="/opt/swimlogs/backups"
scriptname=$(basename "$0")

usage() {
	echo "Usage: $scriptname {help|backup|restore}"
	echo -e "\t help: show this help"
	echo -e "\t backup: creates a backup of database in $backups_dir"
	echo -e "\t restore: restores the latest backup in $backups_dir"
	exit "${1:-0}"
}

backup() {
	if [ ! -d "$backups_dir" ]; then
		mkdir -p "$backups_dir"
	fi

	backups_count=$(find "$backups_dir" -type f | wc -l)

	if [ "$backups_count" -gt 5 ]; then
		oldest_backup=$(find "$backups_dir" -type f -name "*.sql" | sort -t'_' -k1.29,1.33 -k1.26,1.27 -k1.23,1.24 | head -1)
		rm "$oldest_backup"
	fi

	datetime=$(date +%d-%m-%Y"_"%H_%M_%S)
	/usr/bin/docker exec -t db pg_dump --clean --exclude-table=schema_migrations \
		-U swimlogs >"$backups_dir"/"$datetime".sql
}

restore() {
	if [ ! -d "$backups_dir" ]; then
		echo "Backups directory does not exist"
		exit 1
	fi

	latest_backup=$(find "$backups_dir" -type f -name "*.sql" | sort -t'_' -k1.29,1.33 -k1.26,1.27 -k1.23,1.24 | tail -1)

	/usr/bin/docker <"$latest_backup" exec -i db psql -U swimlogs
}

function main() {
	if [ $# -lt 1 ]; then
		usage 1
	fi

	case ${1} in
	help)
		usage 0
		;;
	backup | restore)
		"$1"
		;;
	*)
		echo "unknown command: $1"
		usage 1
		;;
	esac
}

main "$@"
