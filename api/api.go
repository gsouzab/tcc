package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// GetSensors retorna todos os sensores
func GetSensors(w http.ResponseWriter, r *http.Request) {
	response := RedisGetAllSensors()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetSensor retorna um sensor especifico
func GetSensor(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	sensorID := params["id"]

	response := RedisGetSensor(sensorID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CreateSensor cria um sensor
func CreateSensor(w http.ResponseWriter, r *http.Request) {
	var sensor Sensor
	_ = json.NewDecoder(r.Body).Decode(&sensor)

	response := RedisSetSensor(sensor)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// DeleteSensor deleta um sensor
func DeleteSensor(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	sensorID := params["id"]
	response := RedisDeleteSensor(sensorID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

// InsertProbe ...
func InsertProbe(w http.ResponseWriter, r *http.Request) {
	var point Point
	var response Response

	_ = json.NewDecoder(r.Body).Decode(&point)
	point.Name = probeMeasurement

	InfluxInsert(point)

	w.Header().Set("Content-Type", "application/json")
	response.Success = true
	response.Msg = "Probe inserido com sucesso!"
	response.Data = point

	json.NewEncoder(w).Encode(response)
}

// InsertTelemetry ...
func InsertTelemetry(w http.ResponseWriter, r *http.Request) {
	var point Point
	var response Response

	_ = json.NewDecoder(r.Body).Decode(&point)
	point.Name = telemetryMeasurement

	InfluxInsert(point)

	w.Header().Set("Content-Type", "application/json")
	response.Success = true
	response.Msg = "Telemetria inserida com sucesso!"
	response.Data = point

	json.NewEncoder(w).Encode(response)
}

// QueryProbes ...
func QueryProbes(w http.ResponseWriter, r *http.Request) {

	var query InfluxQuery
	var response Response

	_ = json.NewDecoder(r.Body).Decode(&query)

	res, error := CreateInfluxQuery(query, probeMeasurement)

	if error != nil {
		response.Success = false
		response.Error = err.Error()
	} else {
		response.Success = true
		response.Data = res[0].Series
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// QueryTelemetry ...
func QueryTelemetry(w http.ResponseWriter, r *http.Request) {

	var query InfluxQuery
	var response Response

	_ = json.NewDecoder(r.Body).Decode(&query)

	res, error := CreateInfluxQuery(query, telemetryMeasurement)

	if error != nil {
		response.Success = false
		response.Error = err.Error()
	} else {
		response.Success = true
		response.Data = res[0].Series
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
