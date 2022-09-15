# Setting up global resource service test/demo environment in Google Cloud
This document gives brief introduction on how to set up global resource service test/demo environment in Google Cloud.

### **Pre-requirements**
1.  Request a VM in Google Cloud Platform, recommended configuration: ubuntu 20.04, 8 cpu or above, disk size 100GB or up.

2. Run "gcloud version" to ensure the Google Gloud SDK is updated (recommended Google Cloud SDK version is 396.0.0 and up). Please refer to https://cloud.google.com/sdk/docs/downloads-apt-get or https://cloud.google.com/sdk/docs/downloads-versioned-archives to upgrade your google cloud SDK. 

Step 3. Follow [dev environment setup instruction](dev-node-setup.md) to set up developer environment on this newly created GCE instance.
### **Automatically start testing**
1. Set correct env
```
export GRS_INSTANCE_PREFIX=[yourPreferInstanceName] SIM_NUM=5 CLIENT_NUM=2 SERVER_NUM=1
export SERVER_ZONE=us-central1-a SIM_ZONE=us-central1-a,us-east1-b,us-west2-a,us-west4-a,us-west3-c CLIENT_ZONE=us-west3-b,us-east4-b
export SIM_REGIONS="Beijing,Shanghai,Wulan,Guizhou,Reserved1" SIM_RP_NUM=10 NODES_PER_RP=20000 SCHEDULER_REQUEST_MACHINE=25000 SCHEDULER_REQUEST_LIMIT=26000 SCHEDULER_NUM=20
```
2. Simulator data pattern can be set seperately, total number should be same as SIM_NUM" 
```
export SIM_DATA_PATTERN=Daily,Outage,Daily,Daily,Daily SIM_WAIT_DOWN_TIME=0,6,0,0,0
```
1. If there is simulator data pattern is RP outage pattern, then set  SIM_DOWN_RP_NUM=0,10,0,0,0". 
```
export SIM_DOWN_RP_NUM=0,10,0,0,0
```
4. You can disable metrics to avoid possible performance lost by set "SERVICE_EXTRA_ARGS=--enable_metrics=false"
```
export SERVICE_EXTRA_ARGS="--enable_metrics=false"
```
5. You can only setup test environment and no testing started automaticly by set "AUTORUN_E2E=false"
```
export AUTORUN_E2E=false 
```
6. Once all env varible set correctly, run command below to start test env and testing
```
./hack/test-setup.sh
```


### **Collect logs**
1. Run on dev machines
2. Set any log env if different with default, and run test-logcollect.sh
```
./hack/test-logcollect.sh
```
3. By default, all logs copy to machine: sonyadev4:~/grs/logs/${SERVER_NUM}se${SIM_NUM}si${CLIENT_NUM}cl
```
DIR_ROOT=${DIR_ROOT:-"~"}
SIM_LOG_DIR=${SIM_LOG_DIR:-"${DIR_ROOT}/logs"}
SERVER_LOG_DIR=${SERVER_LOG_DIR:-"${DIR_ROOT}/logs"}
CLIENT_LOG_DIR=${CLIENT_LOG_DIR:-"${DIR_ROOT}/logs"}
DES_LOG_DIR=${DES_LOG_DIR:-"${DIR_ROOT}/grs/logs/${SERVER_NUM}se${SIM_NUM}si${CLIENT_NUM}cl"}
DES_LOG_INSTANCE=${DES_LOG_INSTANCE:-"sonyadev4"}
DES_LOG_INSTANCE_ZONE=${DES_LOG_INSTANCE_ZONE:-"us-central1-a"}
```

### **Manually start testing**
You can manually start testing on test environment if you use "AUTORUN_E2E=false" to start test env.
1. Start service

> you can get information as below after run "./hack/test-setup.sh" with "AUTORUN_E2E=false" 
```
You can start service using args: --master_ip=sonya-grs-server-us-central1-a-mig-nhmt --resource_urls=34.67.91.151:9119,35.199.189.187:9119,34.94.213.21:9119,34.125.157.142:9119,34.106.241.209:9119
```

> run on server machine: ${GRS_INSTANCE_PREFIX}-server
```
cd /home/sonyali/go/src/global-resource-service 
mkdir -p ~/logs
/usr/local/go/bin/go run resource-management/cmds/service-api/service-api.go  --master_ip=sonya-grs-server-us-central1-a-mig-nhmt --resource_urls=34.67.91.151:9119,35.199.189.187:9119,34.94.213.21:9119,34.125.157.142:9119,34.106.241.209:9119 -v=3 > ~/logs/sonya-grs-server-us-central1-a-mig-nhmt.log 2>&1
```

2. Start simulator
> run on simulator machine: ${GRS_INSTANCE_PREFIX}-sim
```
cd /home/sonyali/go/src/global-resource-service
mkdir -p ~/logs
/usr/local/go/bin/go run resource-management/test/resourceRegionMgrSimulator/main.go  --region_name=Wulan --rp_num=10 --nodes_per_rp=20000 --master_port=9119 --data_pattern=Outage --wait_time_for_make_rp_down=5 -v=3  > ~/logs/sonya-grs-sim-us-west2-a-2.log 2>&1
```

> each simulator should have different "region_name". available "region_name": 
```
    	Beijing   
	Shanghai  
	Wulan     
	Guizhou   
	Reserved1 
	Reserved2 
	Reserved3 
	Reserved4 
	Reserved5 
```

3. Run e2e testing

> you can get information as below after run "./hack/test-setup.sh" with "AUTORUN_E2E=false"
```
You can start scheduler using args: --service_url=34.172.122.124:8080
```

> run on client machines: ${GRS_INSTANCE_PREFIX}-client, $num is the number of scheduler that you want to start on this machine
```
cd /home/sonyali/go/src/global-resource-service
mkdir -p ~/logs 
$ for i in {1..$num}; do sleep 1; /usr/local/go/bin/go run resource-management/test/e2e/singleClientTest.go --service_url=34.172.122.124:8080 --request_machines=25000 --action=watch --repeats=1 --limit=26000 -v=6 > ~/logs/sonya-grs-client-us-east4-b-1.log.$i 2>&1 & done
```


### **Tear down test env**
```
./hack/test-teardown.sh
```
Note: Connection to GCP can be terminated unexpectedly, before tearing down the test environment, make sure your environment variables are set to be the same when you setup environment.