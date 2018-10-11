package main

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"time"

	"github.com/influxdata/influxdb/client/v2"
)

const influxAddr = "http://influx:8086"
const influxDbName = "tcc_data"
const measurement = "probe"
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

// InfluxQueryTelemetry ...
func InfluxQueryTelemetry(queryTelemetry TelemetryQuery) (res []client.Result, err error) {

	var query = "SELECT "

	var selectByField = queryTelemetry.SelectFields != nil && len(queryTelemetry.SelectFields) > 0
	var selectByTag = queryTelemetry.SelectTags != nil && len(queryTelemetry.SelectTags) > 0
	var selectMean = queryTelemetry.SelectMeanField != ""
	var hasSelect = selectByField || selectByTag || selectMean

	var whereLastMinutes = queryTelemetry.WhereLastXMinutes != ""
	var wherePeriod = queryTelemetry.WhereStartTime != "" && queryTelemetry.WhereEndTime != ""
	var whereField = queryTelemetry.Where != nil && len(queryTelemetry.Where) > 0
	var hasWhere = whereLastMinutes || wherePeriod || whereField

	var groupByMeanTime = queryTelemetry.SelectMeanInterval != ""
	var groupByTag = queryTelemetry.GroupByTag != ""
	var hasGroupBy = groupByMeanTime || groupByTag

	if hasSelect {
		if selectMean {
			query += "MEAN(" + queryTelemetry.SelectMeanField + ") "

		} else {
			if selectByField {
				for i := 0; i < len(queryTelemetry.SelectFields); i++ {
					query += queryTelemetry.SelectFields[i]
					query += "::field"
					if i < (len(queryTelemetry.SelectFields) - 1) {
						query += ", "
					}
				}
			}

			if selectByTag {
				if selectByField {
					query += ", "
				}

				for i := 0; i < len(queryTelemetry.SelectTags); i++ {
					query += queryTelemetry.SelectTags[i]
					query += "::tag"
					if i < (len(queryTelemetry.SelectTags) - 1) {
						query += ", "
					}
				}
			}
		}

	} else {
		query += "* "
	}

	query += " FROM \"" + telemetryMeasurement + "\" "

	if hasWhere {
		query += "WHERE "

		if whereField {
			var mapCount int
			for field, value := range queryTelemetry.Where {

				query += "\"" + field + "\""
				query += " = "

				if reflect.TypeOf(value).String() == "string" {
					query += "'" + value.(string) + "' "

				} else if reflect.TypeOf(value).String() == "float64" {
					query += strconv.Itoa(int((value.(float64)))) + " "
				}

				mapCount++

				if mapCount < (len(queryTelemetry.Where)) {
					query += "AND "
				}
			}

			if whereLastMinutes || wherePeriod {
				query += "AND "
			}
		}

		if wherePeriod {
			query += "time >= '" + queryTelemetry.WhereStartTime + "' AND "
			query += "time <= '" + queryTelemetry.WhereEndTime + "' "

		} else if whereLastMinutes {
			now := time.Now()

			minutesInt, error := strconv.Atoi(queryTelemetry.WhereLastXMinutes)
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
			query += "time(" + queryTelemetry.SelectMeanInterval + "m)"
			if groupByTag {
				query += ", "
			}
		}
		if groupByTag {
			query += queryTelemetry.GroupByTag
		}
	}

	res, err = queryDB(influxClient, query)

	fmt.Println(query)

	if err != nil {
		log.Fatal(err)
	}

	return res, nil
}
