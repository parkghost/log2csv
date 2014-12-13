log2csv
=======

[![Build Status](https://travis-ci.org/parkghost/log2csv.png)](https://travis-ci.org/parkghost/log2csv) 

*a simple tool convert Go gclog to csv*

### Usage

```
Usage1: log2csv -i gc.log -o gc.csv
Usage2: GODEBUG=gctrace=1 your-go-program 2>&1 | log2csv -o gc.csv
       (GO version below 1.2) GOGCTRACE=1 your-go-program 2>&1 | log2csv -o gc.csv
  -i="": The input file (Default: Stdin)
  -o="": The output file
  -t=false: Insert timestamp at the beginning of each line (Required Stdin input)
  -h   : show help usage
```

#### Output

```csv
numgc,nproc,pause,sweep,mark,wait,heap0,heap1,obj,nmalloc,nfree,nspan,nbgsweep,npausesweep,nhandoff,nhandoffcnt,nsteal,nstealcnt,nprocyield,nosyield,nsleep
1,1,5,0,186,0,0,0,18,19,1,0,0,0,0,0,0,0,0,0,0
2,1,4,1,316,1,0,0,208,209,1,12,0,0,0,0,0,0,0,0,0
3,1,1,0,360,1,0,0,565,616,51,33,0,0,0,0,0,0,0,0,0
4,2,755,4,10745,2,0,0,2361,2608,247,52,8,0,1,2,6,114,16,10,14
5,2,2463,2,3395,13,0,1,6298,6676,378,90,0,0,0,0,7,163,16,1,0
6,1,21,14,2858,1,1,3,15684,18565,2881,188,1,0,0,0,0,0,0,0,0
7,2,6,6,1621,2,1,2,14259,28924,14665,376,1,0,0,0,1,17,16,10,1
8,2,5,9,1668,1,1,2,14536,39712,25176,376,1,0,1,2,2,20,19,4,0
9,2,1,7,3114,1,1,2,13972,49530,35558,376,0,0,0,0,7,220,16,1,0
10,2,1,3185,12,17,1,3,17867,62532,44665,376,0,0,0,0,7,220,16,1,0
11,2,24,5,1418,1,1,3,19292,77210,57918,484,6,0,0,0,1,3,13,0,0
12,2,2,2844,1,154,1,3,18071,90905,72834,492,7,0,0,0,7,226,16,10,1
13,2,1,6,1614,197,2,4,20139,106514,86375,492,324,0,0,0,1,13,16,10,1
14,2,1,11,1424,1,1,3,20220,122275,102055,510,5,0,0,0,1,9,16,1,0
15,2,2,728,2689,14,1,3,13151,131182,118031,510,19,0,0,0,6,110,16,1,0
16,2,1,6,3139,0,1,3,19185,146019,126834,510,1,0,0,0,7,200,16,1,0
17,2,1,3002,11,146,2,4,20707,162326,141619,510,116,0,1,2,7,226,16,10,1
```

Authors
-------

**Brandon Chen**

+ http://brandonc.me
+ http://github.com/parkghost


License
---------------------

This project is licensed under the MIT license