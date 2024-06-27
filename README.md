# ip-addr-counter

## Commands

1. Execute the application with 10000 concurrency to get a count of unique IP address:

    ```bash
    $ make run f=ip_addresses c=10000
    ```
   
    Output:
    ```bash
    2024/06/27 19:58:00 count of unique IP addresses: 1000000000
    2024/06/27 19:58:00 execution time: 3m31.888071083s
    ```

2. Run the unit tests:

    ```bash
    $ make test
    ```

3. Run the benchmarks:

    ```bash
    $ make benchmark
    ```

## Implementation

The input that contains the IP addresses (one per line) is divided into chunks according to `concurrency` parameter.
The application makes sure that the chunks are handled properly, i.e. each chunk starts on the new line with an IP address, not in random file parts.
Each file chunk is processed in it's own goroutine, meaning that the entire file is processed in `concurrency` number of goroutines.

The unique IP addresses are marked as bits in the Bitmap with the size of 2 ^ 32 (total number of all possible IPv4 addresses).
Bitmap is fully preallocated (512 MB) and therefore has a fixed size of used RAM (not counting the Go runtime memory usage).

Note: Go runtime memory usage might be optimised by adjusting the Garbage Collector's environment variable `GOGC` that defines how often the GC is executed.

All Bitmap operations are implemented in concurrent safe manner using atomics to achieve the highest execution speed.
