# pathfinder

Pathfinder is a command-line tool that helps you navigate between directories without needing to specify the entire path. Simply specify the name of the directory you want to switch to, and Pathfinder will handle the rest. It also supports specifying entire paths if needed.

## Features

- Quickly switch from one directory to another without entering the full path
- Pathfinder also supports full paths
- Use the `-b` flag to get back to the previous directory you were in
- Uses a cache to store frequently visited directories, so you can move between them faster

## Installation

> [!WARNING]
> Pathfinder is only tested in Linux right now.

Pathfinder has a big limitation which can be overcomed only through a small step of methods that must be done to get pathfinder working.

> [!IMPORTANT]
> I will make a install script to make help you install pathfinder easily. Till then please follow the below steps

- Clone this repository
```
git clone https://github.com/vyshnav-vinod/pathfinder.git
```

- Make sure you have go installed and run (you should be inside the cloned folder)
```
go build -o pathfinder
```

> [!NOTE]
> You can also use `go install`, I mentioned this method to make it easier for you to find the path of the executable 

- Open `pathfinder.sh` in any text editor and replace `pfexecpath` in Line 4 with the full path of `pathfinder` (The executable that was produced using the `go build`) command

- Open `~/.bashrc` in any text editor and add a alias to run `pathfinder.sh` at the end of the file
```
alias pf='. /path/to/pathfinder.sh'
```
or
```
alias pf='source /path/to/pathfinder.sh'
```

- Restart your terminal or just reload bashrc
```
source ~/.bashrc
```

You are good to go. Check below usage and examples to start using pathfinder. If you were stuck or unable to install, please raise an [issue](https://github.com/vyshnav-vinod/pathfinder/issues) and i will help you.

## Usage

```bash
pf [directory name/path] (flags)
```

## Flags

- `-b, --back` : Move back to the previous directory from where `pathfinder` was called.
- `-i, --ignore` : Ignore searching for the folder in the current directory
- `-h, --help` : Display the help message

## Examples

- Go to the folder named `dirname`
```bash
pf dirname
```

- Go to the folder named `dirname` using full path
```bash
pf ~/Desktop/dirname
```

- Go back to previous directory
```bash
pf -b
```

## Issues

Pathfinder is still in early stages. You can raise an [issue](https://github.com/vyshnav-vinod/pathfinder/issues) and i will be glad to help. You can also contribute to pathfinder by raising a PR. But first, it would be better if you raised an issue regarding what you will be adding.
