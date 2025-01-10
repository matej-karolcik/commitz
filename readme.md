# commitz
================

A Git commit tool that assists in generating commit messages and README files.

## Installation
-------------

To use commitz, you need to have Go installed on your system. You can download the latest version of Go from the official website: <https://golang.org/dl/>

Once you have Go installed, run the following command:
```bash
go install github.com/matej-karolcik/commitz@latest
```
This will install commitz and its dependencies.

## Usage
-----

To use commitz, simply run the following command in your terminal:
```bash
commitz
```
This will start the commitz tool and prompt you to enter a message for your commit.

Alternatively, you can also use commitz with a custom prefix like a ticket number by running the following command:
```bash
commitz TICKET-123
```
