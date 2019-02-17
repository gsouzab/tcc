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
	Success bool        `json:"success,omitempty"`
	Msg     string      `json:"msg,omitempty"`
	Error   string      `json:"error,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// Point eh a estrutura que armazena a resposta de uma requisicao.
type Point struct {
	Name   string                 `json:"name,omitempty"`
	Tags   map[string]string      `json:"tags,omitempty"`
	Fields map[string]interface{} `json:"fields,omitempty"`
	Time   time.Time              `json:"time,omitempty"`
}

// ProbeMessage object
type ProbeMessage struct {
	SSID      string `json:"ssid,omitempty"`
	DstMac    string `json:"dstMac,omitempty"`
	SrcMac    string `json:"srcMac,omitempty"`
	Sensor    string `json:"sensor,omitempty"`
	Channel   string `json:"channel,omitempty"`
	RSSI      string `json:"rssi,omitempty"`
	CreatedAt int64  `json:"createdAt,omitempty"`
}

// TelemetryMessage object
type TelemetryMessage struct {
	Sensor    string `json:"sensor,omitempty"`
	CO2       string `json:"co2,omitempty"`
	Temp      string `json:"temp,omitempty"`
	Hum       string `json:"hum,omitempty"`
	Latitude  string `json:"latitude,omitempty"`
	Longitude string `json:"longitude,omitempty"`
	CreatedAt int64  `json:"createdAt,omitempty"`
}

// InfluxQuery object
type InfluxQuery struct {
	SelectFields       []string               `json:"selectFields,omitempty"`
	SelectTags         []string               `json:"selectTags,omitempty"`
	SelectMeanField    string                 `json:"selectMeanField,omitempty"`
	SelectMeanInterval string                 `json:"selectMeanInterval,omitempty"`
	Where              map[string]interface{} `json:"where,omitempty"`
	WhereStartTime     string                 `json:"whereStartTime,omitempty"`
	WhereEndTime       string                 `json:"whereEndTime,omitempty"`
	WhereLastXMinutes  string                 `json:"whereLastXMinutes,omitempty"`
	GroupByTag         string                 `json:"groupByTag,omitempty"`
}
