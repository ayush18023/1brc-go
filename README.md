# 1brc-go 

The One Billion Row Challenge (1BRC) is a fun exploration of how far modern Java can be pushed for aggregating one billion rows from a text file. 

The Challenge was originally coined for Java based compilers but here i have implemented in golang .

The original repo can be found at https://github.com/gunnarmorling/1brc

The Implementation is in different methods running for machine with specs below
### Specifications:
**Processor:** Intel(R) Core(TM) i5-8265U CPU @ 1.60GHz, 1800 Mhz, 4 Core(s), 8 Logical Processor(s)<br />
**Physical Memory(RAM):** 8GB

### Method 1: Hashmap but no Parallization
In this method hashmap is used where key is hashed to index in a array. No go routine is used  the program executes serially. The program takes **7m 35s** to execute approximately
### Method 2: Go's Maps with pointers but no Parallization
This method is similar to first but instead of using hashmap golang's inbuilt map is used which maps key to pointer of values. The execution time reduced to around **5m** 
### Method 3: Parallization(Batch Processing) with Hashmaps
Here the execution is parellized by batch processing along with custom hashmaps. Lines are chunked together in different sizes of multiple 1024 bytes. The execution time is reduced to **4m 45s**
### Method 4: Parallization(Batch Processing) with Go Maps
This method is same as above only maps are used. Execution time is tested for different chunk lengths along with the cores in use.<br/>
chunk len: 64x1024, 128x1024, 256x1024 <br/>
| 8 Cores |  3m54s  |  3m22s   |   3m23s  |<br/>
| 6 Cores |  3m54s  |  3m21s   |   3m49s  |<br/>
### Further Optimizations
1. Using Buffers instead of scanner.Scan()
2. Different methods to convert to tempratures to Int
3. Profiling the program
4. Improving Hashing algorithm to use instead of go maps
5. Execution time can be further reduced by increasing CPU Cores to increase parallization