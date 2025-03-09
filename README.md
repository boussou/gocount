# gocount

Go count all files and directories.

This is significantly faster than its equivalent in C lang thanks to goroutine.  

Implementation a quota of goroutines 
(applying a quota of 100 highly limits resource usage without impacting perf)

Notes 
> benched on 9.445.438 files and 1.400.174 directories)  
> quota of 100 goroutines used  
> without quota used max concurrent walkDir calls: **69403**

> 


```
async count in go no goroutine quota:
time gocount-nolimit ~/www
real    0m3,079s
user    0m27,006s
sys     0m9,052s
Note: it used Max Concurrent walkDir calls: 69403


async count in go with a limit of 100 goroutines:
time gocount ~/www
real    0m3,310s
user    0m33,378s
sys     0m10,612s

sync count in C:
time countfiles ~/www
~/www contains 9.445.438 files and 1.400.174 directories

Using pure linux commandline 
time find ~/www |wc -l     # => 10.845.613 files
real    0m18,903s
user    0m2,989s
sys     0m8,056s
```


### Implementation a quota of goroutines : Explanation


This code is intented to set a **maximum concurrent calls** to walkDir().  
A way to limit the number of concurrent calls to walkDir is by using a buffered channel as a semaphore. 

This pattern allows to limit concurrent goroutine usage while still traversing the directory tree efficiently.


- **Semaphore (`sem` channel):**  
  A buffered channel of capacity `maxConcurrent` is used to control the number of concurrent `walkDir` calls. 
  
  Before processing a directory, a goroutine writes an empty struct into `sem` (acquiring a token). Once the work in that invocation is finished, the token is released by reading from the channel (in the deferred function). 
  
  This ensures that no more than `maxConcurrent` goroutines are actively traversing directories at any given time.

- **WaitGroup:**  
  The WaitGroup (`wg`) ensures that the main function waits until all directory traversals are complete.

- **Mutex:**  
  A Mutex (`mu`) protects the shared counters (`totalFiles` and `totalDirs`) during concurrent updates.
