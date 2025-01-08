# magma-agent
![Build Status](https://github.com/andrewjouffray/magma-agent/actions/workflows/go.yml/badge.svg)

magma agent to track system directories, compute cryptographic hashes of files within those tracked directories and save them as json 'snapshots'. This is meant as a way to help developers track changes that they may have comited to on-premise systems that are not entirely tracked by a git repo.

## Requirements

- Go v 1.23

## Build

```
# after cloning
cd magma-agent
mkdir build
go build -o ./build/magma
sudo chmod +x ./install.sh
sudo ./install.sh
```

## Tests

```
# 'magma init' must be ran before testing
cd magma-agent
go test -cover ./...
```

## Usage
There a currently only 4 commands
- "init": Initializes the magma directory.
- "track [path]": Adds a new path to the track file.
- "untrack [path]": Removes a path from the track file.
- "snap [tag1] [tag2] ...": Creates a new cryptographic snapshot for all tracked files and directories.
