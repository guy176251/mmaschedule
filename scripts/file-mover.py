#!/usr/bin/env python3

import argparse
import logging
import os
import shutil
import sys
from pathlib import Path
from subprocess import run
from time import sleep

logging.basicConfig(level=logging.INFO)


class Args:
    file: Path
    dest: Path
    postcmd: str | None
    precmd: str | None


parser = argparse.ArgumentParser()

parser.add_argument(
    "--file",
    help="File to watch",
    type=Path,
    required=True,
)

parser.add_argument(
    "--dest",
    help="Directory to copy file to when it changes.",
    type=Path,
    required=True,
)

parser.add_argument(
    "--precmd",
    help="Shell command to run before moving file.",
)

parser.add_argument(
    "--postcmd",
    help="Shell command to run after moving file.",
)

args = parser.parse_args(namespace=Args())

if not args.dest.is_dir():
    logging.error(f"Directory {args.dest} does not exist or is invalid.")
    sys.exit(1)

logging.info(f"Watching file {args.file}, will move to {args.dest} when changed.")

mtime = 0

try:
    mtime = int(os.stat(args.dest / args.file.name).st_mtime)
except FileNotFoundError:
    pass


def run_script(script: str):
    logging.info(f"Running script {script!r}")
    result = run(script, shell=True, capture_output=True)
    logging.info(f"Command ended with code {result.returncode}")


while True:
    try:
        sleep(60)
        new_mtime = int(os.stat(args.file).st_mtime)
        if mtime != new_mtime:
            logging.info("File changed, moving to destination")
            if isinstance(args.precmd, str):
                run_script(args.precmd)
            shutil.copy(args.file, args.dest)
            mtime = new_mtime
            if isinstance(args.postcmd, str):
                run_script(args.postcmd)
    except KeyboardInterrupt:
        print("\nExiting")
        break
    except FileNotFoundError:
        pass
    except Exception as err:
        logging.error(f"Operation failed: {err}", exc_info=True)
