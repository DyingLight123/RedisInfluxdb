package RedisInfluxdb

import (
	"encoding/json"
	"fmt"
	"github.com/influxdata/platform/kit/errors"
	"github.com/robfig/cron"
	g "github.com/soniah/gosnmp"
	"log"
	"strconv"
	"time"
)

type RedisInfluxdb struct {
	RedisAddr           string
	RedisPassword       string
	RedisKey            string
	InfluxdbAddr        string
	InfluxdbUser        string
	InfluxdbPassword    string
	InfluxdbData        string
	InfluxdbMeasurement string
}

type MapResults struct {
	MapTime string
	Key     string
	Value   map[string]interface{}
}

var pause = make(chan string, 1)

func (this *RedisInfluxdb) RedisWriteToInfluxdb() error {
	cli, err := ConnRedis(this.RedisAddr, this.RedisPassword)
	defer cli.Close()
	if err != nil {
		return err
	}
	field, err := cli.HGetAll(this.RedisKey).Result()
	if err != nil {
		return err
	}
	conn, err := ConnInfluxdb(this.InfluxdbAddr, this.InfluxdbUser, this.InfluxdbPassword)
	if err != nil {
		return err
	}
	err = WritesPoints(conn, field, this.InfluxdbData, this.InfluxdbMeasurement)
	if err != nil {
		return err
	}
	return nil
}

func (this *RedisInfluxdb) GetRedis() (*map[string]string, error) {
	cli, err := ConnRedis(this.RedisAddr, this.RedisPassword)
	defer cli.Close()
	if err != nil {
		return nil, err
	}
	field, err := cli.HGetAll(this.RedisKey).Result()
	if err != nil {
		return nil, err
	}
	return &field, nil
}

func (this *RedisInfluxdb) RefreshRedis(number int) {
	c := cron.New()
	c.AddFunc("@every "+"10s", func() {
		err := AddRedisData(this.RedisAddr, this.RedisPassword, this.RedisKey, number)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("redis已更新！")
	})
	c.Start()
	<-pause
	fmt.Println("continue")
	c.Stop()
	//time.AfterFunc(30 * time.Second, c.Stop)
}

func (this *RedisInfluxdb) PauseRedis() {
	fmt.Println("pause")
	pause <- "continue"
}

/*func (this *RedisInfluxdb) RefreshRedis1(number int, pause chan string) {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		if len(pause) == 1 {
			wg.Done()
			return
		}
		time.Sleep(10 * time.Second)
		fmt.Println(number)
		err := AddRedisData(this.RedisAddr, this.RedisPassword, this.RedisKey, number)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("redis已更新！")
	}()

	wg.Wait()
	fmt.Println(<-pause)
	//time.AfterFunc(30 * time.Second, c.Stop)
}*/

func (this *RedisInfluxdb) GetInfluxdb() ([]*MapResults, error) {
	conn, err := ConnInfluxdb(this.InfluxdbAddr, this.InfluxdbUser, this.InfluxdbPassword)
	if err != nil {
		return nil, err
	}
	if this.InfluxdbData == "" {
		return nil, errors.New("data can not be empty")
	}
	if this.InfluxdbMeasurement == "" {
		return nil, errors.New("measurement can not be empty")
	}
	//cmd := fmt.Sprintf("select * from %s where time >= %s and time < %s tz('Asia/Shanghai')",
	//	this.InfluxdbMeasurement, "'" + t1+ "'", "'" + t2 + "'")
	cmd := fmt.Sprintf("select * from %s tz('Asia/Shanghai')",
		this.InfluxdbMeasurement)
	res, err := QueryDB(conn, cmd, this.InfluxdbData)
	if err != nil {
		return nil, err
	}
	if len(res[0].Series) == 0 {
		return nil, nil
	}
	results := make([]*MapResults, 0)
	for _, row := range res[0].Series[0].Values {
		_, err := time.Parse(time.RFC3339, row[0].(string))
		if err != nil {
			log.Println(err)
		}
		m := new(MapResults)
		m.MapTime = row[0].(string)
		m.Key = row[1].(string)

		x := row[2].(string)
		n := make(map[string]interface{})
		json.Unmarshal([]byte(x), &n)
		//fmt.Println(n)
		m.Value = n

		/*switch row[2].(type) {
		case json.Number:
			str := string(row[2].(json.Number))
			m.Value, err = strconv.Atoi(str)
			if err != nil {
				log.Fatal(err)
			}
			//fmt.Println("执行了我")
		case string:
			m.Value, err = strconv.Atoi(row[2].(string))
			if err != nil {
				log.Fatal(err)
			}
		case :

		}*/
		//fmt.Println(m)
		results = append(results, m)
	}
	return results, nil
}

func (this *RedisInfluxdb) GetSnmpIfo(snmpaddr string, port string, oids string) (*map[string]interface{}, error) {
	port1, err := strconv.Atoi(port)
	if err != nil {
		return nil, err
	}
	params := &g.GoSNMP{
		Target:    snmpaddr,
		Port:      uint16(port1),
		Community: "public",
		Version:   g.Version2c,
		Timeout:   time.Duration(2) * time.Second,
	}
	err = params.Connect()
	if err != nil {
		//log.Fatalf("Connect() err: %v", err)
		return nil, err
	}
	defer params.Conn.Close()

	//oids := ".1.3.6.1.2.1.2.2.1.10"
	result, err := params.WalkAll(oids) // Get() accepts up to g.MAX_OIDS
	if err != nil {
		//log.Fatalf("Get() err: %v", err)
		return nil, err
	}
	m := make(map[string]interface{}, 0)
	for _, v := range result {
		//a := strings.Split(v.Name, ".")
		//m[a[10]] = v.Value
		m[v.Name] = v.Value
		/*switch v.Value {
		case g.Integer:
			m[v.Name] = v.Value.(int)
		default:
			m[v.Name] = v.Value
		}*/
	}
	return &m, nil
}
