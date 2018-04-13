package jsongen

import (
	"encoding/json"
	"log"
	"math/rand"
	"time"

	"github.com/Shopify/sarama"
)

//WrapJSON wraps and returns generated at generateJSON() random JSON data
func WrapJSON() sarama.StringEncoder {
	start := time.Now()

	data := make(map[string]interface{})
	for i := 0; i < 200; i++ {
		num := rand.Int63()
		data["field"+string(i)] = num
	}

	exportData, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	t := time.Now()
	elapsed := t.Sub(start)
	log.Println("GenerateJSON time elapsed: ", elapsed)
	return sarama.StringEncoder(exportData)
}
