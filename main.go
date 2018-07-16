package main

import (
	"bufio"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	// YAML Formatter
	"gopkg.in/yaml.v2"

	// Cool logger && nice textFormat wrapper
	"github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"

	// Influx client
	client "github.com/influxdata/influxdb/client/v2"
)

var (
	err error
	u   *url.URL
	log *logrus.Logger
)

type metrics struct {
	timestamp int64
	metric    map[string]interface{}
}

// LoadYAML - new type config yandex-tank
type LoadYAML struct {
	Influx struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Database string `yaml:"database"`
		Address  string `yaml:"address"`
		Port     int    `yaml:"port"`
	}
	Phantom struct {
		Address string `yaml:"address"`
		Port    int    `yaml:"port"`
	}
	Console struct {
		Enabled   bool `yaml:"enabled"`
		ShortOnly bool `yaml:"short_only"`
	}
}

var config = struct {
	conf        LoadYAML
	timeKey     []string
	picks       []string
	urlExporter string
	urlInflux   string
}{
	timeKey: []string{"time"},
	picks: []string{
		"cp_user",
		"cp_sys",
		"cp_iowait",
		"memory_total",
		"memory_free",
		"memory_buffers",
		"memory_cached",
		"interface_bytes_in",
		"interface_bytes_out",
		"hdd_busy_time",
		"hdd_operations_read",
		"hdd_operations_written",
	},
}

const (
	exporterPort       = "1957"
	loadYAMLfilename   = "load.yaml"
	monitoringInterval = 5 * time.Second
)

func init() {
	log = logrus.New()

	// 22:06:19 [INFO] Plugin <yandextank.plugins.Autostop.plugin.Plugin object at 0x7f0533efee90> required 0.000025 seconds to start
	log.Formatter = &easy.Formatter{
		TimestampFormat: "15:04:05",
		LogFormat:       "%time% [%lvl%] %msg%\n",
	}

	log.SetLevel(logrus.DebugLevel)

	// Reading load config
	yamlFile, err := ioutil.ReadFile(loadYAMLfilename)
	check(err)
	err = yaml.Unmarshal(yamlFile, &config.conf)
	check(err)
	log.Printf("Config file %s found", loadYAMLfilename)
	log.Printf("Influx settings: %v ", config.conf.Influx)
	log.Printf("Exporter settings: %v ", config.conf.Phantom)
	// Exporter configuration
	// Parse ip:port, [ip]:port, ip:*** port:***, (IPv4 Only)
	config.urlExporter = "http://" + phantomToExporter(config.conf.Phantom.Address) + ":" + exporterPort
	log.Printf("Exporter: %v", config.urlExporter)

	// Influx configuration
	//http://localhost:8086
	config.urlInflux = "http://" +
		config.conf.Influx.Address + ":" +
		strconv.Itoa(config.conf.Influx.Port)
	log.Printf("Influx backend: %v", config.urlInflux)

}

func phantomToExporter(ph string) string {
	// IPv4 only
	socket := strings.Split(ph, ":")
	if len(socket) < 2 {
		socket[1] = "80"
	}
	socket[0] = strings.Replace(socket[0], "[", "", 1)
	socket[0] = strings.Replace(socket[0], "]", "", 1)
	return socket[0]
}

func main() {
	for {
		batch, err := metricReader()
		check(err)
		// log.Printf("Metrics received... ")
		err = influxUploader(batch)
		check(err)
		if config.conf.Console.Enabled && config.conf.Console.ShortOnly {
			log.Printf("Metrics successfully sent from node %s to influx %s", config.conf.Phantom.Address, config.conf.Influx.Address)
		}
		time.Sleep(monitoringInterval)
	}
}

func influxUploader(b metrics) error {
	//startTime := time.Now()

	// Create a new HTTPClient
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     config.urlInflux,
		Username: config.conf.Influx.Username,
		Password: config.conf.Influx.Password,
	})
	check(err)
	defer c.Close()

	// Create a new point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  config.conf.Influx.Database,
		Precision: "s",
	})
	check(err)

	// TODO: Make configurable
	tags := map[string]string{"system": "common"}

	pt, err := client.NewPoint("monitoring", tags, b.metric, time.Unix(b.timestamp, 0))
	check(err)
	bp.AddPoint(pt)

	// Write the batch
	err = c.Write(bp)
	check(err)

	// Close client resources
	err = c.Close()
	check(err)

	//log.Printf("metricWriter, %f ms", timeSpent(startTime))
	return nil
}

func check(err error) {
	startTime := time.Now()
	if err != nil {
		log.Printf("check, %f ms", timeSpent(startTime))
		log.Printf("Error: %v, gracefull exit", err)
		log.Printf("Gracefull exit")
		os.Exit(0)
	}
}

func metricReader() (metrics, error) {
	//startTime := time.Now()
	var b metrics
	mtrmap := make(map[string]interface{})
	// Metrics to regexp:
	rgxPicks, err := regexpComposer(config.picks)
	check(err)
	rgxTime, err := regexpComposer(config.timeKey)
	check(err)
	u, err := url.Parse(config.urlExporter + "/metrics")
	check(err)
	parameters := url.Values{}
	u.RawQuery = parameters.Encode()

	client := http.Client{Timeout: time.Duration(600) * time.Second}
	req, _ := http.NewRequest("GET", u.String(), nil)
	resp, err := client.Do(req)
	check(err)

	defer resp.Body.Close()
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		str := scanner.Text()
		if checkKey(str, rgxTime) {
			_, v, err := stringToList(str)
			check(err)
			b.timestamp, err = strconv.ParseInt(v, 10, 64)
			check(err)
		} else {
			if checkKey(str, rgxPicks) {
				k, v, err := stringToList(str)
				check(err)
				mtrmap[k] = stringToDigital(v)
			}
		}

	}
	b.metric = mtrmap

	err = scanner.Err()
	check(err)
	// log.Printf("metricReader, %f ms", timeSpent(startTime))
	if b.timestamp == 0 {
		return b, errors.New("Zero timestamp. Check exporter version")
	}

	return b, nil
}

func stringToDigital(s string) interface{} {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return s
	}
	return f
}

func stringToList(s string) (string, string, error) {
	list := strings.Split(s, " ")
	if len(list) < 2 {
		return "", "", errors.New("Incorrect key-value string: " + s)
	}
	var key string
	var delimiter string
	for k, v := range list {
		if k > 0 {
			delimiter = " "
		}
		if (len(list) - k) == 1 {
			break
		}
		key = key + delimiter + v
	}
	return key, list[len(list)-1], nil
}

func regexpComposer(list []string) (string, error) {
	if len(list) == 0 {
		return "", errors.New("Empty metrics")
	}
	var out string
	for _, v := range list {
		out = out + "(^" + v + ".*)"
	}
	out = strings.Replace(out, ")(", ")|(", -1)
	return out, nil
}

func checkKey(s string, rgx string) bool {
	m, err := regexp.MatchString(rgx, s)
	check(err)
	return m
}

func timeSpent(t time.Time) float64 {
	return float64((time.Now().UnixNano() - t.UnixNano()) / 1000000)
}
