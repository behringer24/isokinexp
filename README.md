# isokinexp
[![Build + Test](https://github.com/behringer24/isokinexp/actions/workflows/go.yml/badge.svg)](https://github.com/behringer24/isokinexp/actions/workflows/go.yml)
[![Release build](https://github.com/behringer24/isokinexp/actions/workflows/release.yml/badge.svg)](https://github.com/behringer24/isokinexp/actions/workflows/release.yml)

isokinexp is a command to import export files from isokinetic measuring devices.

# Basic usage

```
Usage:
  isokinexp [flags]

Flags:
  -d, --delete       delete source files from input directory.
                     isokinexp will move files instead of copying
  -h, --help         help for isokinexp
  -i, --in string    path to import files from
  -o, --out string   path to export files to
  -v, --version      version for isokinexp
```

To import files from an input path to an output path:

```
isokinexp -i import/path -o output/path
```

to also delete the source files you can use the -d or --delete option

```
isokinexp -d -i import/path -o output/path
```

# Installation
## from pre build release

Go to https://github.com/behringer24/isokinexp/releases to download the latest release file for your architecture. Extract the file to your harddrive.

## Execute on Windows

Execute the isokinexp.exe directly from the folder you extracted the .exe file to.Alternatively you can copy the .exe file to a folder that is available via your %PATH% variable.

# Build and install from source

Checkout the main branch from git@github.com:behringer24/isokinexp.git

```
git clone git@github.com:behringer24/isokinexp.git
cd ./isokinexp
go build -v ./...
go install
```
