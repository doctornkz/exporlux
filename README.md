# Exporlux 
### Simple metrics pull-based Prometheus fetcher
#### Extremely usefull with Yandex-Tank InfluxDB plugin

### How to use:
#### Clone, change constants and compile:
```
git clone https://github.com/doctornkz/exporlux.git
cd exporlux
vim main.go
go build -o exporlux
```

#### Save file to reachable for YandeTank place:
```
cp exporlux /tmp
```

#### Add section to your yaml.load config:
```
shellexec:
  start: /tmp/exporlux &
  end: pkill exporlux
```

##### Note #1: Exporlux console logging works only with `short_only` Yandex-Tank console settings.
##### Note #2: If something goes wrong, exporlux will be closed with `os.Exit(0)`. 

#### Running:

After start exporlux will parse your load.yaml, sections `influx` , `phantom` and uses settings for workload.

```
$ yandex-tank -o "influx.tank_tag=SomeTag"
23:08:15 [INFO] Loading plugins...
23:08:15 [INFO] Grafana link: http://*****:3000/dashboard/db/sometag?var-uuid=e737eb93-112b-425b-af26-891217221367&from=-5m&to=now
23:08:15 [INFO] Performing test
23:08:15 [INFO] Configuring plugins...
23:08:15 [INFO] Testing connection to resolved address 10.*.*.44 and port 80
23:08:15 [INFO] Resolved ****.ru into 10.*.*.44:80
23:08:15 [INFO] Configuring StepperWrapper...
23:08:15 [INFO] Using cached stpd-file: ./logs/_9254ae20b87a52b4f267f84562a3f093.stpd
23:08:15 [INFO] rps_schedule is set. Overriding cached instances param from config: 1
23:08:15 [INFO] Preparing test...   0, speed:    13 Krps
23:08:16 [INFO] Checking tank resources...
23:08:16 [INFO] Starting test...
23:08:16 [INFO] using verbose histogram
23:08:16 [INFO] Plugin <yandextank.plugins.Phantom.plugin.Plugin object at 0x7f102cd3e550> required 0.007444 seconds to start
23:08:16 [INFO] Plugin <yandextank.plugins.Autostop.plugin.Plugin object at 0x7f102cd8de90> required 0.000016 seconds to start
23:08:16 [INFO] Executing: /tmp/exporlux &
23:08:16 [INFO] Plugin <yandextank.plugins.ShellExec.plugin.Plugin object at 0x7f102cd8d4d0> required 0.008112 seconds to start
23:08:16 [INFO] Plugin <yandextank.plugins.Influx.plugin.Plugin object at 0x7f102cd8d490> required 0.000043 seconds to start
23:08:16 [INFO] Plugin <yandextank.plugins.Console.plugin.Plugin object at 0x7f102cd8d750> required 0.000014 seconds to start
23:08:16 [INFO] Plugin <yandextank.plugins.RCAssert.plugin.Plugin object at 0x7f102cd54150> required 0.000003 seconds to start
23:08:16 [INFO] Plugin <yandextank.plugins.ResourceCheck.plugin.Plugin object at 0x7f102cd54f10> required 0.000003 seconds to start
23:08:16 [INFO] Plugin <yandextank.plugins.JsonReport.plugin.Plugin object at 0x7f102cd54cd0> required 0.000002 seconds to start
23:08:16 [INFO] Waiting for test to finish...
23:08:16 [INFO] Artifacts dir: /home/doctor/Work/go/src/github.com/doctornkz/exporlux/logs/2018-07-16_23-08-15.l8S9DC
23:08:16 [INFO] Config file load.yaml found
23:08:16 [INFO] Influx settings: {*** **** loaddb ****.ru 8086}
23:08:16 [INFO] Exporter settings: {****.ru:80 0}
23:08:16 [INFO] Exporter: http://*****.ru:1957
23:08:16 [INFO] Influx backend: http://*****:8086
23:08:16 [INFO] Metrics successfully sent from node ****.ru:80 to influx ****
....
23:08:19 [INFO] ts:1531771696   RPS:1   avg:13.64       min:13.64       max:13.64       q95:13.64
23:08:20 [INFO] ts:1531771697   RPS:1   avg:3.01        min:3.01        max:3.01        q95:3.01
23:08:21 [INFO] ts:1531771698   RPS:1   avg:4.08        min:4.08        max:4.08        q95:4.08
23:08:22 [INFO] Metrics successfully sent from node ****.ru:80 to influx ****
23:08:22 [INFO] ts:1531771699   RPS:1   avg:5.64        min:5.64        max:5.64        q95:5.64
23:08:23 [INFO] ts:1531771700   RPS:1   avg:3.53        min:3.53        max:3.53        q95:3.53
23:08:24 [INFO] ts:1531771701   RPS:1   avg:3.29        min:3.29        max:3.29        q95:3.29
23:08:25 [INFO] ts:1531771702   RPS:1   avg:2.67        min:2.67        max:2.67        q95:2.67
23:08:26 [INFO] ts:1531771703   RPS:1   avg:3.05        min:3.05        max:3.05        q95:3.05
23:08:27 [INFO] ts:1531771704   RPS:1   avg:2.96        min:2.96        max:2.96        q95:2.96
23:08:28 [INFO] Metrics successfully sent from node ****.ru:80 to influx ****
....<Ctrl+C Pressed>
^C23:12:09 [INFO] Do not press Ctrl+C again, the test will be broken otherwise
23:12:09 [INFO] Trying to shutdown gracefully...
23:12:09 [INFO] Finishing test...
23:12:09 [INFO] Stopping load generator and aggregator
23:12:09 [INFO] Terminating phantom process with PID 31191
23:12:11 [INFO] ts:1531771926   RPS:1   avg:3.58        min:3.58        max:3.58        q95:3.58
23:12:11 [INFO] ts:1531771927   RPS:1   avg:3.50        min:3.50        max:3.50        q95:3.50
23:12:11 [INFO] ts:1531771928   RPS:1   avg:7.48        min:7.48        max:7.48        q95:7.48
23:12:11 [INFO] Timestamps without stats:
23:12:11 [INFO] 1531771929
23:12:11 [INFO] ts:1531771929   RPS:1   avg:10.69       min:10.69       max:10.69       q95:10.69
23:12:11 [INFO] Stopping monitoring
23:12:11 [INFO] Executing: pkill exporlux
23:12:11 [INFO] Post-processing test...
23:12:11 [INFO] Artifacts dir: /home/doctor/Work/go/src/github.com/doctornkz/exporlux/logs/2018-07-16_23-08-15.l8S9DC
23:12:11 [WARNING] File not found to collect: validation_error.yaml
23:12:11 [INFO] Done graceful shutdown
23:12:11 [INFO] Close allocated resources...
23:12:11 [INFO] Done performing test with code 1
```
##### After test metrics appears in InfluxDB. Now you can improve your LoadTest dashboard.
##### That ugly hack works well with CI automation, especially inside Docker image e.g. Gitlab-CI pipelines.

### Enjoy!
