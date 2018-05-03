package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
)

var addr = flag.String("Server", "localhost:6379", "Redis server address")
var client *redis.Client

//Sensor eh a estrutura que armazena um sensor.
type Sensor struct {
	Name        string `json:"name,omitempty"`
	Latitude    string `json:"latitude,omitempty"`
	Longitude   string `json:"longitude,omitempty"`
	Description string `json:"description,omitempty"`
	Mac         string `json:"mac,omitempty"`
}

//Response eh a estrutura que armazena a resposta de uma requisicao.
type Response struct {
	Sucess bool        `json:"sucess,omitempty"`
	Msg    string      `json:"msg,omitempty"`
	Error  string      `json:"error,omitempty"`
	Data   interface{} `json:"data,omitempty"`
}

//RedisInit eh o metodo que inicializa a conexao com o cliente Redis
func RedisInit(indexDatabase int) *redis.Client {
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",            // no password set
		DB:       indexDatabase, // use default DB
	})

	return client
}

//RedisGetAllSensors eh o metodo que retorna todos os sensores armazenados no Redis
func RedisGetAllSensors() Response {
	var response Response
	response.Sucess = true
	var sensors []Sensor

	keys, err := client.Keys("*").Result()
	if err != nil {
		response.Sucess = false
		response.Error = err.Error()
		return response
	}

	for _, key := range keys {
		m, err := client.HGetAll(key).Result()
		if err != nil {
			response.Sucess = false
			response.Error = err.Error()
			return response
		}

		var sensor Sensor
		for k, v := range m {
			switch k {
			case "Name":
				sensor.Name = v
			case "Latitude":
				sensor.Latitude = v
			case "Longitude":
				sensor.Longitude = v
			case "Description":
				sensor.Description = v
			case "Mac":
				sensor.Mac = v
			default:
			}
		}

		sensors = append(sensors, sensor)
	}

	response.Msg = "Consulta realizada com sucesso!"
	response.Data = sensors

	return response
}

//RedisSetSensor eh o metodo que insere um sensor no Redis
func RedisSetSensor(sensor Sensor) Response {
	var response Response
	response.Sucess = true

	ok, err := client.HMSet(sensor.Mac, map[string]interface{}{
		"Name":        sensor.Name,
		"Latitude":    sensor.Latitude,
		"Longitude":   sensor.Longitude,
		"Description": sensor.Description,
		"Mac":         sensor.Mac,
	}).Result()

	if err != nil {
		response.Sucess = false
		response.Error = err.Error()
		return response
	}

	response.Msg = ok
	response.Data = sensor

	return response
}

//RedisDeleteSensor eh o metodo que deleta um sensor do Redis
func RedisDeleteSensor(sensorID string) Response {
	var response Response
	response.Sucess = true

	status, err := client.Del(sensorID).Result()

	if err != nil {
		response.Sucess = false
		response.Error = err.Error()
		return response
	}

	if status == 1 {
		response.Msg = "Sensor deletado com sucesso!"
	} else {
		response.Msg = "O sensor não existe!"
	}

	return response
}

//GetSensors retorna todos os sensores
func GetSensors(w http.ResponseWriter, r *http.Request) {
	response := RedisGetAllSensors()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

//CreateSensor cria um sensor
func CreateSensor(w http.ResponseWriter, r *http.Request) {
	var sensor Sensor
	_ = json.NewDecoder(r.Body).Decode(&sensor)

	response := RedisSetSensor(sensor)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

//DeleteSensor deleta um sensor
func DeleteSensor(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	sensorID := params["id"]
	response := RedisDeleteSensor(sensorID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	flag.Parse()

	router := mux.NewRouter()
	RedisInit(0)

	router.HandleFunc("/sensors", GetSensors).Methods("GET")
	router.HandleFunc("/sensors", CreateSensor).Methods("POST")
	router.HandleFunc("/sensor/{id}", DeleteSensor).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8000", router))
}
