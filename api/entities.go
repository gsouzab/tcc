package main

import "time"

// Sensor eh a estrutura que armazena um sensor.
type Sensor struct {
	Name        string `json:"name,omitempty"`
	Latitude    string `json:"latitude,omitempty"`
	Longitude   string `json:"longitude,omitempty"`
	Description string `json:"description,omitempty"`
	Mac         string `json:"mac,omitempty"`
}

// Response eh a estrutura que armazena a resposta de uma requisicao.
type Response struct {
	Success bool        `json:"success"`
	Msg     string      `json:"msg"`
	Error   string      `json:"error"`
	Data    interface{} `json:"data"`
}

// Point eh a estrutura que armazena a resposta de uma requisicao.
type Point struct {
	Name   string                 `json:"name"`
	Tags   map[string]string      `json:"tags"`
	Fields map[string]interface{} `json:"fields"`
	Time   time.Time              `json:"time"`
}

// TelemetryMessage object
type TelemetryMessage struct {
	Sensor    string `json:"sensor"`
	CO2       int    `json:"co2"`
	Temp      int    `json:"temp"`
	Hum       int    `json:"hum"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
	CreatedAt int64  `json:"createdAt"`
}
