# gocount
Go count


This is significantly faster than its equivalent in C lang thanks to goroutine

Notes 
> benched on 9.445.438 files and 1.400.174 directories
> it used max concurrent walkDir calls: **69403**





### Display the max number of concurrent calls to walkDir()

In order to get an idea of the resource usage

Tracks and displays the maximum number of concurrent calls to `walkDir`. 

A shared counter for the current active goroutines and one for the maximum observed value are updated in a thread-safe manner using a mutex. 

### Explanation

- **Concurrent Call Tracking:**  
  Two counters, `currentConcurrent` and `maxConcurrent`, are maintained. When `walkDir` starts, it increments `currentConcurrent` and, if needed, updates `maxConcurrent` if the current value is higher than any previously observed. When `walkDir` finishes (via a deferred function), it decrements `currentConcurrent`.

- **Synchronization:**  
  A mutex (`mu`) ensures that updates to the shared counters (for files, directories, and concurrency) are done safely across concurrent goroutines.

- **WaitGroup:**  
  A WaitGroup (`wg`) is used to wait for all goroutines to finish before printing the final results.

This displays the maximum number of concurrent calls to `walkDir` that occurred during the traversal.
