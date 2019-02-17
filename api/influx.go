package main

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"time"

	client "github.com/influxdata/influxdb1-client/v2"
)

const influxAddr = "http://influx:8086"
const influxDbName = "tcc_data"
const telemetryMeasurement = "telemetry_data"
const probeMeasurement = "probe_data"
const influxUsername = "tcc"
const influxPassword = "tcc_ufrj"
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

// InfluxInsert ...
func InfluxInsert(point Point) {
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

// CreateInfluxQuery ...
func CreateInfluxQuery(influxQuery InfluxQuery, measurement string) (res []client.Result, err error) {

	var query = "SELECT "

	var selectByField = influxQuery.SelectFields != nil && len(influxQuery.SelectFields) > 0
	var selectByTag = influxQuery.SelectTags != nil && len(influxQuery.SelectTags) > 0
	var selectMean = influxQuery.SelectMeanField != ""
	var hasSelect = selectByField || selectByTag || selectMean

	var whereLastMinutes = influxQuery.WhereLastXMinutes != ""
	var wherePeriod = influxQuery.WhereStartTime != "" && influxQuery.WhereEndTime != ""
	var whereField = influxQuery.Where != nil && len(influxQuery.Where) > 0
	var hasWhere = whereLastMinutes || wherePeriod || whereField

	var groupByMeanTime = influxQuery.SelectMeanInterval != ""
	var groupByTag = influxQuery.GroupByTag != ""
	var hasGroupBy = groupByMeanTime || groupByTag

	if hasSelect {
		if selectMean {
			query += "MEAN(" + influxQuery.SelectMeanField + ") "

		} else {
			if selectByField {
				for i := 0; i < len(influxQuery.SelectFields); i++ {
					query += influxQuery.SelectFields[i]
					query += "::field"
					if i < (len(influxQuery.SelectFields) - 1) {
						query += ", "
					}
				}
			}

			if selectByTag {
				if selectByField {
					query += ", "
				}

				for i := 0; i < len(influxQuery.SelectTags); i++ {
					query += influxQuery.SelectTags[i]
					query += "::tag"
					if i < (len(influxQuery.SelectTags) - 1) {
						query += ", "
					}
				}
			}
		}

	} else {
		query += "* "
	}

	query += " FROM \"" + measurement + "\" "

	if hasWhere {
		query += "WHERE "

		if whereField {
			var mapCount int
			for field, value := range influxQuery.Where {

				query += "\"" + field + "\""
				query += " = "

				if reflect.TypeOf(value).String() == "string" {
					query += "'" + value.(string) + "' "

				} else if reflect.TypeOf(value).String() == "float64" {
					query += strconv.Itoa(int((value.(float64)))) + " "
				}

				mapCount++

				if mapCount < (len(influxQuery.Where)) {
					query += "AND "
				}
			}

			if whereLastMinutes || wherePeriod {
				query += "AND "
			}
		}

		if wherePeriod {
			query += "time >= '" + influxQuery.WhereStartTime + "' AND "
			query += "time <= '" + influxQuery.WhereEndTime + "' "

		} else if whereLastMinutes {
			now := time.Now()

			minutesInt, error := strconv.Atoi(influxQuery.WhereLastXMinutes)
			if error != nil {
				fmt.Println(err)
			}

			nowSubtract := now.Add(time.Duration(-minutesInt) * time.Minute)
			query += "time >= '" + nowSubtract.UTC().Format("2006-01-02T15:04:05Z0700") + "' AND "
			query += "time <= '" + now.UTC().Format("2006-01-02T15:04:05Z0700") + "' "
		}
	}

	if hasGroupBy {
		query += "GROUP BY "

		if groupByMeanTime {
			query += "time(" + influxQuery.SelectMeanInterval + "m)"
			if groupByTag {
				query += ", "
			}
		}
		if groupByTag {
			query += influxQuery.GroupByTag
		}
	}

	res, err = queryDB(influxClient, query)

	fmt.Println(query)

	if err != nil {
		log.Fatal(err)
	}

	return res, nil
}
