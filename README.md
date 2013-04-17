log2csv
=======

*a simple tool convert Go gclog to csv*

### Usage

```
Usage1: log2csv -i gc.log -o gc.csv
Usage2: GCTRACE=1 your-go-program 2>&1 | log2csv -o gc.csv
  -h=false: Show Usage
  -i="": The input file (default: Stdin)
  -o="": The output file (default: Stdout)
  -t=false: Add timestamp at line head(Stdin input only)
```

Authors
-------

**Brandon Chen**

+ http://brandonc.me
+ http://github.com/parkghost


License
---------------------

Licensed under the Apache License, Version 2.0: http://www.apache.org/licenses/LICENSE-2.0