# gocount

Go count - simplistic synchronous walking equivalent to code written in C lang.

This branch to keep track of it for benchmark comparisons   
(and have a quick access / without a tag)


This is **2x slower** than its equivalent in C lang.  
(roughly 193%)

```
sync count in go:
real    0m2,137s
user    0m1,296s
sys     0m0,905s

sync count in C:
real    0m1,108s
user    0m0,211s
sys     0m0,879s
```