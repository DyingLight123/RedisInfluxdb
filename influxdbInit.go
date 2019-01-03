package RedisInfluxdb

import (
	"errors"
	"github.com/influxdata/influxdb/client/v2"
	"time"
)

func ConnInfluxdb(addr string, user string, password string) (client.Client, error) {
	if addr == "" {
		return nil, errors.New("influxdbaddr can not be empty")
	}
	if user == "" {
		return nil, errors.New("influxdbusername can not be empty")
	}
	cli, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     addr,
		Username: user,
		Password: password,
	})
	if err != nil {
		return nil, err
	}
	return cli, nil
}

func QueryDB(cli client.Client, cmd string, data string) (res []client.Result, err error) {
	q := client.Query{
		Command:  cmd,
		Database: data,
	}
	if response, err := cli.Query(q); err == nil {
		if response.Error() != nil {
			return nil, response.Error()
		}
		res = response.Results
	} else {
		return nil, err
	}
	return res, nil
}

func WritesPoints(cli client.Client, field map[string]string, data string, measurement string) error {
	//t := time.Now()
	if data == "" {
		return errors.New("influxdbdata can not be empty")
	}
	if measurement == "" {
		return errors.New("influxdbmeasurement can not be empty")
	}
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  data,
		Precision: "s",
	})
	if err != nil {
		return err
	}
	for key, value := range field {
		tags := map[string]string{"key": key}
		fields := map[string]interface{}{
			"value": value,
		}
		pt, err := client.NewPoint(
			measurement,
			tags,
			fields,
			time.Now(),
		)
		if err != nil {
			return err
		}
		bp.AddPoint(pt)
	}
	if err := cli.Write(bp); err != nil {
		return err
	}
	//elapsed := time.Since(t)
	//fmt.Println(elapsed)
	return nil
}
