# SerMon
Service endpoints monitoring service. Providers an API for registering
and monitoring endpoints of HTTP services and analyse their network properties.

## Service Definition

A sample service definition is shown below, of which **name, endpoint, upCodes, metadata**
are the input fields while registering a service. Post registration, each service is monitored every 15 minutes
by default (can be changed using env variable) and metrics are collected. A maximum of 30 records of metrics
history is maintained.
```
{
  "id": 1,
  "name": "Google",
  "endpoint": "https://google.com",
  "upCodes": "200,302,301",
  "metadata": "google-web",
  "status": "UP",
  "metrics": {
    "timeStamp": "2020-07-23 12:22:11.3111039 +0530 IST m=+60.010645701",
    "dnsTime": 0.0489812,
    "connectTime": 0,
    "tlsTime": 0,
    "totalTime": 0.0489812
  },
  "timeSeriesMetrics": [
    {
      "timeStamp": "2020-07-23 10:04:45.9729621 +0530 IST m=+21.006139401",
      "dnsTime": 0.1149825,
      "connectTime": 0.0090198,
      "tlsTime": 0.0579789,
      "totalTime": 0.1149825
    },
    {
      "timeStamp": "2020-07-23 10:05:24.9844617 +0530 IST m=+60.017639001",
      "dnsTime": 0.0474716,
      "connectTime": 0,
      "tlsTime": 0,
      "totalTime": 0.0479937
    }
  ]
}
```

## Environment Variables

The following are the environment variables that can be used to tune the service

 * PORT - Service port (default=80)
 * POLL_INTERVAL - Polling interval (in minutes) for monitoring service (default=15)
 * DB_PATH - File path to SQLite DB file (default=sermon.db)
 