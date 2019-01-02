package RedisInfluxdb

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

/*func (this *RedisInfluxdb) GetInfluxdb() (*map[string]interface{}, error) {
	conn, err := ConnInfluxdb(this.InfluxdbAddr, this.InfluxdbUser, this.InfluxdbPassword)
	if err != nil {
		return nil, err
	}
	cmd := fmt.Sprintf("select * from %s where time >= %s and time < %s tz('Asia/Shanghai')",
		"map", "'" + t1+ "'", "'" + t2 + "'")
	res, err := QueryDB(conn, cmd)
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
			log.Fatal(err)
		}
		m := new(MapResults)
		m.MapTime = row[0].(string)
		m.Key = row[1].(string)

		x := row[2].(string)
		n := make(map[string]interface{})
		json.Unmarshal([]byte(x), &n)
		//fmt.Println(n)
		m.Value = n

		switch row[2].(type) {
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

		}
		fmt.Println(m)
		results = append(results, m)
	}
	return results, nil
}*/
