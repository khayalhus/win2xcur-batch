# win2xcur-batch
Converts and renames cursors from Windows format (*.cur, *.ani) to Xcursor format.

## Prerequisites

The program has only been tested on Linux Mint.

Install [Go](https://go.dev/doc/install)

This project is a helper program for [win2xcur](https://github.com/quantum5/win2xcur). Original project needs to be installed.

## Usage

Place your `.ani` or `.cur` files in a directory inside directory `Unzipped`. Cursor directory might be renamed if it contains whitespace as `win2xcur` converter does not support them in folder path. Your hierarchy should look something like this.

```sh
❯ tree
.
├── converter.go
├── go.mod
├── LICENSE
├── map.json
├── README.md
└── Unzipped
    └── CursorFolder
        ├── 01-Normal.ani
        ├── 02-Link.ani
        ├── 03-Loading.ani
        ├── 04-Help.ani
        ├── 05-Text Select Alt.ani
        ├── 05-Text Select.ani
        ├── 06-Handwriting.ani
        ├── 07-Precision.ani
        ├── 08-Unavailable.ani
        ├── 09-Location Select.ani
        ├── 10-Person Select.ani
        ├── 11-Vertical Resize.ani
        ├── 12-Horizontal Resize.ani
        ├── 13-Diagonal Resize 1.ani
        ├── 14-Diagonal Resize 2.ani
        ├── 15-Move.ani
        ├── 16-Alternate Select.ani
        └── Installer.inf
```

Execute the program.

```sh
go build
./win2xcur-batch
```

Look into `Sorted` directory for final output. Place the cursor directory in `/usr/share/icons` or `/home/$USER/.icons`