# gbenchcmp

Utility for comparison of [googlebenchmark](https://github.com/google/benchmark) results.

Usage:
```
$ go get github.com/dvyukov/gbenchcmp
$ gbenchcmp old.res new.res

Benchmark                  Time(ns): old        new     diff  CPU(ns): old         new     diff
===============================================================================================
contraction_NxNxN_48T/10             404        433   +7.18%           404         433   +7.18%
contraction_NxNxN_48T/64            8843       8834   -0.10%          8848        8840   -0.09%
contraction_NxNxN_48T/512        1220311     455426  -62.68%      10044706    15969236  +58.98%
contraction_NxNxN_48T/4k        96134407   84270601  -12.34%    3402981956  3901162020  +14.64%
contraction_NxNxN_48T/5k       181780696  152923999  -15.87%    6337035902  7136203158  +12.61%
```

The tool can also show CPU load and choose last/best results (if several results are present in the input files).
