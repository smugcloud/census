# Census CLI Tool

Tool which interacts with the Census API to retrieve two possible sets of data.

## Usage

```
Usage: census [params] [comma separated list of states]

e.g. census --averages oregon,washington,california

  -averages
    	Return average income below poverty across
	the states specified.
  -csv
    	Print CSV output of all state information.
```

## Building

The included Makefile will produce builds for both macOS and Linux (64bit).

Pre-built binaries are included in the `builds` directory.