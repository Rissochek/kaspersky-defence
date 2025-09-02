package utils

import (
	"log"
	"os"
	"strconv"
)

func GetKeyFromEnv(key string) int {
	value, exists := os.LookupEnv(key)
	defaultValues := map[string]int{
		"WORKERS":    2,
		"QUEUE_SIZE": 32,
		"TASKS": 10,
	}

	if !exists {
		log.Printf("%v value is not set in environment. Using default value: %v", key, defaultValues[key])
		return defaultValues[key]
	}

	valueInt, err := strconv.Atoi(value)
	if err != nil {
		log.Printf("failed to convert %v to int: %v. Using default value: %v", key, err, defaultValues[key])
		return defaultValues[key]
	}

	log.Printf("%v using env value = %v", key, valueInt)
	return valueInt
}