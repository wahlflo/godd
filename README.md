# godd
``godd`` is a simple cli tool for copying binary data similar to the ``dd`` command.

## Install godd

```
git clone https://github.com/wahlflo/godd.git
cd godd
sudo GOBIN=/usr/local/bin/ go install cmd/godd.go
```

## Usage

```
[bob@xps godd]$ godd -h
Usage: godd -source SOURCE -destination DESTINATION [OPTIONS]
Options:
  -buffer-size int
    	number of chunks which are buffered (default 500)
  -chunk-size int
    	size of chunks the tool reads / writes in KB (default 12)
  -destination string
    	path to the destination file / disk
  -h	show help
  -if string
    	path to the source file / disk (short form)
  -of string
    	path to the destination file / disk (short form)
  -source string
    	path to the source file / disk

```

If no input file is given the standard input is used as input stream and if no output file is given the standard output is used as destination.

The advantage over the traditional ``dd`` command is that it is a little bit faster and outputs the bottleneck (reading or writting) and an ETA out of the box.

### Example 1
```
> godd -if /dev/sda1 -of ./test.raw
[+] running since: 0d 00h 00m 01s     transferred data: 1GB     speed: 1604.08MB/s     buffer: 0% (reading is the bottleneck)     ETA: 0d 00h 04m 52s
[+] running since: 0d 00h 00m 11s     transferred data: 11GB     speed: 1031.73MB/s     buffer: 0% (reading is the bottleneck)     ETA: 0d 00h 07m 24s
[+] running since: 0d 00h 00m 21s     transferred data: 20GB     speed: 1006.19MB/s     buffer: 4% (reading is the bottleneck)     ETA: 0d 00h 07m 26s
[+] running since: 0d 00h 00m 31s     transferred data: 30GB     speed: 1012.11MB/s     buffer: 0% (reading is the bottleneck)     ETA: 0d 00h 07m 13s
...
```

The ``buffer`` statistic indicates how much the buffer is filled. A low filled buffer indicates that the reading the bottleneck and a filled buffer indicates that writing is the bottleneck.
