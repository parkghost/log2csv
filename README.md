log2csv
=======

[![Build Status](https://travis-ci.org/parkghost/log2csv.png)](https://travis-ci.org/parkghost/log2csv) 
[![Coverage Status](https://coveralls.io/repos/parkghost/log2csv/badge.svg)](https://coveralls.io/r/parkghost/log2csv)

*a simple tool convert Go gctrace to csv*

### Usage

```
Usage1: log2csv -i gc.log -o gc.csv
Usage2: GODEBUG=gctrace=1 your-go-program 2>&1 | log2csv -o gc.csv
  -i="stdin": The input file
  -o="stdout": The output file
  -t=true: Add timestamp at line head (the input file must be `stdin`)
```

### Installation
```
go get github.com/parkghost/log2csv/cmd/log2csv
```

### Example
```
$ GODEBUG=gctrace=1 godoc -http=:6060  2>&1 | log2csv -o gc.csv
```

###### Output
```
$ head gc.csv
unixtime,numgc,nproc,seq,sweep,mark,wait,heap0,heap1,obj,nmalloc,nfree,goroutines,nspan,nbgsweep,npausesweep,nhandoff,nhandoffcnt,nsteal,nstealcnt,nprocyield,nosyield,nsleep
1421931830.043593,1,1,4,0,123,0,0,0,21,21,0,2,16,0,0,0,0,0,0,0,0,0
1421931830.043797,2,1,0,0,101,0,0,0,40,41,1,3,20,0,0,0,0,0,0,0,0,0
1421931830.043930,3,1,1,0,121,0,0,0,128,144,16,4,30,0,0,0,0,0,0,0,0,0
1421931830.048720,4,1,1,1,129,0,0,0,139,190,51,4,31,0,0,0,0,0,0,0,0,0
1421931830.048840,5,1,0,0,222,0,0,0,179,259,80,4,31,0,0,0,0,0,0,0,0,0
1421931830.049519,6,1,2,0,595,2,0,0,956,1093,137,4,51,0,0,0,0,0,0,0,0,0
1421931830.055207,7,1,2,0,517,0,0,0,2600,3340,740,33,75,1,0,0,0,0,0,0,0,0
1421931830.057955,8,1,1,0,692,0,0,0,4248,6200,1952,116,105,0,0,0,0,0,0,0,0,0
1421931830.059537,9,1,1,0,670,0,0,1,4086,7423,3337,119,113,0,0,0,0,0,0,0,0,0
```

Author
-------

**Brandon Chen**

+ http://brandonc.me
+ http://github.com/parkghost
 
License
---------------------

This project is licensed under the MIT license
