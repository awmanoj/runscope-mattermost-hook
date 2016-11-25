# runscope\_mattermost\_hook

This is an example of how to integrate runscope with Mattermost. 

Example is a simple web service which you need to run to receive the callback on runscope events and do HTTP POST to Mattermost incoming webhook URL. 

## Quick Start

* Add your `mattermost` incoming webhook URL in app.development.ini 
* Build the binary: 
```
$ go build
```
* Run the service: 
```
$ nohup ./runscope-mattermost-hook 2>log.txt 1>&2&
```

## More details: 

```
/**************************************************************************** 
 * 
 * following is the callback data that runscope sends to custom webhook.
 *
 *****************************************************************************
{
   "test_id": "76598752-cbda-4e1d-820f-6274a62f74ff",
   "test_name": "Buckets Test",
   "test_run_id": "9c15aa62-21f0-48f2-a819-c99bdf8e4543",
   "team_id": "6b9c7f65-9e11-4f77-85ad-e6ee7a28232d",
   "team_name": "Acme Inc.",
   "environment_uuid": "98290cfc-a008-4ab7-9ea4-8906f12b228f",
   "environment_name": "Staging Settings",
   "bucket_name": "Rocket Sled",
   "bucket_key": "7xzcnsgbwox2",
   "test_url": "https://www.runscope.com/radar/7xzcnsgbwox2/76598752-cbda-4e1d-820f-6274a62f74ff",
   "test_run_url": "https://www.runscope.com/radar/7xzcnsgbwox2/76598752-cbda-4e1d-820f-6274a62f74ff/results/9c15aa62-21f0-48f2-a819-c99bdf8e4543",
   "trigger_url": "https://api.runscope.com/radar/09039249-fdfd-4e1d-820f-6274a62f74ff/trigger",
   "result": "fail",
   "started_at": 1384281308.548077,
   "finished_at": 1384281310.680218,
   "agent": null,
   "region": "us1",
   "region_name": "US East - Northern Virginia",
   "initial_variables": {}, 
   "requests": [{
     "url": "https://api.runscope.com/",
     "variables": {
        "fail": 0,
        "total": 1,
        "pass": 1
     },
     "assertions": {
        "fail": 0,
        "total": 2,
        "pass": 2
     },         
     "scripts": {
        "fail": 0,
        "total": 1,
        "pass": 1
     },
     "result": "pass",
     "method": "GET",
     "response_time_ms": 123,
     "response_size_bytes": 2048,
     "response_status_code": 200,
     "note": "Root URL"
  }]
}
```
