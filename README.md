# stdjson

stdjson converts an application output to **stdout** and **stderr** into json format with each line getting an index number.

## Install

stdjson is written in go and can therefore be installed via: ``go get github.com/mfmayer/mfgrep`` (go must be installed).

## Usage of stdjson

```
$ ./stdjson 
./stdjson: missing command
Usage: ./stdjson COMMAND ARGUMENTS...
 Try : ./stdjson ls -al
```

## Example

Reason for writing this tool was logging of service output to e.g. azure log analytics. Microsoft's log analytics agent is able to log any output (stdout & stderr) into the azure log analytics container logs and stores line by line with a timestamp in milliesecond resolution. Unfortunately, this can cause the order of output to be lost if the output is faster, e.g. in case of a panic when outputting a stack trace.

For this reason stdjson wraps every line of output into the following json format: ``{"i": <line index>,"d": <line data/content>}``

```
$ ./stdjson ls -al
{"d":"total 2688","i":"1"}
{"d":"drwxrwxr-x 3 mayema mayema    4096 Aug 18 13:07 .","i":"2"}
{"d":"drwxr-xr-x 8 mayema mayema    4096 Aug 18 12:57 ..","i":"3"}
{"d":"drwxrwxr-x 8 mayema mayema    4096 Aug 18 12:57 .git","i":"4"}
{"d":"-rw-rw-r-- 1 mayema mayema     269 Aug 18 12:57 .gitignore","i":"5"}
{"d":"-rw-rw-r-- 1 mayema mayema    1071 Aug 18 12:57 LICENSE","i":"6"}
{"d":"-rw-rw-r-- 1 mayema mayema    2322 Aug 18 13:07 main.go","i":"7"}
{"d":"-rw-rw-r-- 1 mayema mayema     135 Aug 18 12:57 README.md","i":"8"}
{"d":"-rwxrwxr-x 1 mayema mayema 2722213 Aug 18 13:07 stdjson","i":"9"}
```