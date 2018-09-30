package main

import (
	"github.com/go-redis/redis"
)

const redisAddr = "redis:6379"
const redisPassword = ""

var redisClient *redis.Client

// RedisInit eh o metodo que inicializa a conexao com o cliente Redis
func RedisInit(indexDatabase int) {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       indexDatabase,
	})
}

// RedisGetAllSensors eh o metodo que retorna todos os sensores armazenados no Redis
func RedisGetAllSensors() Response {
	var response Response
	response.Success = true
	var sensors = make([]Sensor, 0)

	keys, err := redisClient.Keys("*").Result()
	if err != nil {
		response.Success = false
		response.Error = err.Error()
		return response
	}

	for _, key := range keys {
		persistedSensor, err := redisClient.HGetAll(key).Result()
		if err != nil {
			response.Success = false
			response.Error = err.Error()
			return response
		}

		var sensor = Sensor{
			Name:        persistedSensor["Name"],
			Latitude:    persistedSensor["Latitude"],
			Longitude:   persistedSensor["Longitude"],
			Description: persistedSensor["Description"],
			Mac:         persistedSensor["Mac"],
		}

		sensors = append(sensors, sensor)
	}

	response.Msg = "Consulta realizada com sucesso!"
	response.Data = sensors

	return response
}

// RedisGetSensor eh o metodo que retorna um sensor especifico armazenado no Redis
func RedisGetSensor(key string) Response {
	var response Response
	response.Success = true

	persistedSensor, err := redisClient.HGetAll(key).Result()
	if err != nil {
		response.Success = false
		response.Error = err.Error()
		return response
	}

	var sensor = Sensor{
		Name:        persistedSensor["Name"],
		Latitude:    persistedSensor["Latitude"],
		Longitude:   persistedSensor["Longitude"],
		Description: persistedSensor["Description"],
		Mac:         persistedSensor["Mac"],
	}

	response.Msg = "Consulta realizada com sucesso!"
	response.Data = sensor

	return response
}

// RedisSetSensor eh o metodo que insere um sensor no Redis
func RedisSetSensor(sensor Sensor) Response {
	var response Response
	response.Success = true

	ok, err := redisClient.HMSet(sensor.Mac, map[string]interface{}{
		"Name":        sensor.Name,
		"Latitude":    sensor.Latitude,
		"Longitude":   sensor.Longitude,
		"Description": sensor.Description,
		"Mac":         sensor.Mac,
	}).Result()

	if err != nil {
		response.Success = false
		response.Error = err.Error()
		return response
	}

	response.Msg = ok
	response.Data = sensor

	return response
}

// RedisDeleteSensor eh o metodo que deleta um sensor do Redis
func RedisDeleteSensor(sensorID string) Response {
	var response Response
	response.Success = true

	status, err := redisClient.Del(sensorID).Result()

	if err != nil {
		response.Success = false
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
