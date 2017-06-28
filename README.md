# mann

mann is a simple utility to store and index helpful bash/zsh snippets.

## Installation

Compile the go binary:

    user# go build main.go

In your `.bashrc` or `.zshrc` file, save the last command on the terminal
into a file in the mann directory.

Example for bash:

    PROMPT_COMMAND="history | tail -1 | cut -c 8- > ~/.mann/last_command"

## Usage

    user# sed -i s/hello/world/gi myfile.txt
    user# mann -a Simple string replace on file
    user# rpm -q --queryformat '%{VERSION}' centos-release
    7
    user# mann -a Retrieve CentOS release number
    user# mann -l
    rpm
    sed
    user# mann -l rpm
    # Retrieve CentOS release number
    rpm -q --queryformat '%{VERSION}' centos-release