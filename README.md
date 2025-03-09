# gocount
Go count


This is significantly faster than its equivalent in C lang thanks to goroutine
(benched on 2.402.919 files and 412.806 dirs)


async count in go:
real    0m0,854s
user    0m8,429s
sys     0m2,649s


sync count in go:
real    0m2,137s
user    0m1,296s
sys     0m0,905s

sync count in C:
real    0m1,108s
user    0m0,211s
sys     0m0,879s

