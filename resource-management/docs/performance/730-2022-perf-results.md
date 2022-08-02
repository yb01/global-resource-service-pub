
### Performance test results for the 0.1.0 release:

|   Test   | Total Nodes | Regions|Nodes per Region| RPs per Region | Nodes per RP | Schedulers| Nodes per scheduler list | Notes | Register<br>Latency<br>(ms) | List<br>Latency<br>(ms) | Watch<br>P50(ms) | P90(ms) | P99(ms) |
|:--------:| :---: | :---:|----:|----:|----:|----:|----:| ----:|----:|----:|----:|----:| ----:|
|  test-1  | 1m  | 5 | 200K | 10 | 20K | 20 | 25K| disable metric, daily data change pattern|301|871|108|175|211|
| test-1.1 | 1m  | 5 | 200K | 10 | 20K | 20 | 25K| disable metric, RP down data change pattern|374|1012|1021|1137|1156|
|  test-2  | 1m  | 5 | 200K | 10 | 20K | 20 | 25K| enable metric, daily data change pattern|298|1097|116|181|201|
| test-2.1 | 1m  | 5 | 200K | 10 | 20K | 20 | 25K| enable metric,RP down data change pattern|359|1012|1002|1074|1093|
|  test-3  | 1m  | 5 | 200K | 10 | 20K | 20 | 50K| disable metric, daily data change pattern|369|1766|109|173|217|
| test-3.1 | 1m  | 5 | 200K | 10 | 20K | 20 | 50K| disable metric, RP down data change pattern|337|1679|877|1174|1200|
|  test-4  | 1m  | 5 | 200K | 10 | 20K | 40 | 25K| disable metric, daily data change pattern|135|811|92|161|195|
|  test-5  | 1m  | 5 | 200K | 10 | 20K | 20 | 25K| disable metric, rp down, all in one region|209|641|451|513|529|

*Service is gce n1-standard-32 VM with 32 core and 120GB Ram, 500GB SSD, premium network. 
*Scheduler and the resource simulators are can with n1-standard-8 VMs*


#### Regions each component deployed


#### Service:

|        Region |             Location |
|--------------:|---------------------:|
| us Central-1a | Council bluffs, IOWA |



#### Simulators:

|        Region |             Location |
|--------------:|---------------------:|
| us Central-1a | Council bluffs, IOWA |
|    us east1-b |    Moncks COrner, SC |
|    us west2-a |               LA, CA |
|    us west3-c | Salt Lake city, Utah |
|    us west4-a |    las Vegas, Nevada |


#### Schedulers:

|     Region |             Location |
|-----------:|---------------------:|
| us west3-b | Salt Lake city, Utah |
| us east4-b |    Ashburn, Virginia |
