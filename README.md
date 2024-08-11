# isokinexp
[![Build + Test](https://github.com/behringer24/isokinexp/actions/workflows/go.yml/badge.svg)](https://github.com/behringer24/isokinexp/actions/workflows/go.yml)
[![Release build](https://github.com/behringer24/isokinexp/actions/workflows/release.yml/badge.svg)](https://github.com/behringer24/isokinexp/actions/workflows/release.yml)

isokinexp is a command line tool to import/export/sort files from isokinetic measuring devices from Ferstl like the Isomed 2000. Files are copied or moved into a file and filename structure /<name of person>/<name of person><Y-m-d>.txt

# Basic usage

```
Usage:
  isokinexp [flags]

Flags:
  -d, --delete          delete source files from input directory.
                          isokinexp will move files instead of copying
  -f, --filter string   filter names (use * as placeholder)
  -h, --help            help for isokinexp
  -i, --in string       path to import files from
  -o, --out string      path to export files to
  -v, --version         version for isokinexp
```

To copy files from an input path to an output path:

```
isokinexp -i import/path -o output/path
```

To also delete the source files you can use the -d or --delete option:

```
isokinexp -d -i import/path -o output/path
```

To filter the proband name you can use the --filter / -f Option. For an exact match of the name use:

```
isokinexp -i import/path -o output/path -f XYZ4711
```

Use the * / asterisk as a placeholder. Filter for all names starting with XYZ use:

```
isokinexp -i import/path -o output/path -f XYZ*
```

To filter for all names containing XYZ somewhere in the name use:

```
isokinexp -i import/path -o output/path -f *XYZ*
```

Use * wherever needed. Example:

```
isokinexp -i import/path -o output/path -f *XYZ*00*287
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
