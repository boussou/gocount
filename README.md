# gocount
Go count


This is significantly faster than its equivalent in C lang thanks to goroutine.   
(benched on 2.402.919 files and 412.806 dirs)

```


async count in go no goroutine quota:
time gocount-nolimit ~/www
real    0m3,079s
user    0m27,006s
sys     0m9,052s


async count in go with a limit of 100 goroutines:
time gocount ~/www
real    0m3,310s
user    0m33,378s
sys     0m10,612s

sync count in C:
time countfiles ~/www
/home/nadir/www contains 9.445.438 files and 1.400.174 directories

```

## Implementation a quota of goroutines 

I wanted to be able to set a **maximum concurrent calls** to walkDir().  
A way to limit the number of concurrent calls to walkDir by using a buffered channel as a semaphore. 

This pattern allows to limit concurrent goroutine usage while still traversing the directory tree efficiently.

### Explanation
- **Semaphore (`sem` channel):**  
  A buffered channel of capacity `maxConcurrent` is used to control the number of concurrent `walkDir` calls. 
  
  Before processing a directory, a goroutine writes an empty struct into `sem` (acquiring a token). Once the work in that invocation is finished, the token is released by reading from the channel (in the deferred function). 
  
  This ensures that no more than `maxConcurrent` goroutines are actively traversing directories at any given time.

- **WaitGroup:**  
  The WaitGroup (`wg`) ensures that the main function waits until all directory traversals are complete.

- **Mutex:**  
  A Mutex (`mu`) protects the shared counters (`totalFiles` and `totalDirs`) during concurrent updates.
