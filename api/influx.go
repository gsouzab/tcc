package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/influxdata/influxdb/client/v2"
)

const influxAddr = "http://influxdb:8086"
const influxDbName = "probes"
const probeMeasurement = "probe"
const telemetryMeasurement = "telemetry_data"
const influxUsername = "influx"
const influxPassword = "123456"
const precision = "s"

var influxClient client.Client
var batchPoints client.BatchPoints
var err error

// InfluxInit ...
func InfluxInit() error {
	influxClient, err = client.NewHTTPClient(client.HTTPConfig{
		Addr:     influxAddr,
		Username: influxUsername,
		Password: influxPassword,
	})

	if err != nil {
		log.Fatal(err)
	}

	return err
}

func createNewPointBatch() error {
	batchPoints, err = client.NewBatchPoints(client.BatchPointsConfig{
		Database:  influxDbName,
		Precision: precision,
	})

	if err != nil {
		log.Fatal(err)
	}

	return err
}

func createAndInsertNewPoint(point Point) error {
	pt, err := client.NewPoint(point.Name, point.Tags, point.Fields, point.Time)
	if err != nil {
		log.Fatal(err)
	}

	batchPoints.AddPoint(pt)
	return err
}

func writeBatchPointToDB() error {
	if err := influxClient.Write(batchPoints); err != nil {
		log.Fatal(err)
	}

	return err
}

// InfluxInsertProbe ...
func InfluxInsertProbe(point Point) {
	if err := createNewPointBatch(); err != nil {
		log.Fatal(err)
	}

	if err := createAndInsertNewPoint(point); err != nil {
		log.Fatal(err)
	}

	if err := writeBatchPointToDB(); err != nil {
		log.Fatal(err)
	}
}

func queryDB(clnt client.Client, cmd string) (res []client.Result, err error) {
	q := client.Query{
		Command:  cmd,
		Database: influxDbName,
	}
	if response, err := clnt.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	} else {
		return res, err
	}
	return res, nil
}

// InfluxQueryLastTelemetryLimit ...
func InfluxQueryLastTelemetryLimit(limitString string) (res []client.Result, err error) {
	limitInt, err := strconv.Atoi(limitString)
	if err != nil {
		limitInt = 10
	}

	q := fmt.Sprintf("SELECT * FROM %s LIMIT %d", telemetryMeasurement, limitInt)
	res, err = queryDB(influxClient, q)

	if err != nil {
		log.Fatal(err)
	}

	/*
		for i, row := range res[0].Series[0].Values {
			t, err := time.Parse(time.RFC3339, row[0].(string))
			if err != nil {
				log.Fatal(err)
			}
			val := row[1].(string)
			log.Printf("[%2d] %s: %s\n", i, t.Format(time.Stamp), val)
		}
	*/

	return res, nil
}
