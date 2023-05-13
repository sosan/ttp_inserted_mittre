package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

const (
	PONG = "PONG"
)

func main() {
	godotenv.Load(".env")
	ctx := context.Background()

	uri := os.Getenv("DB_URI")
	opt, _ := redis.ParseURL(uri)
	redisClient := redis.NewClient(opt)
	defer redisClient.Close()

	status, err := redisClient.Ping(ctx).Result()
	if (err != nil) && (status != PONG) {
		panic(err)
	}
	// con round trips a saco ! :-(
	content := ReadFromFile("ttp.json")
	log.Printf("TAMAÑO DEL ARRAY DEL CONTENIDO %d", len(content))
	countInserted := 0
	lengthArrary := len(content)
	for i := 0; i < lengthArrary; i++ {
		result, err := redisClient.HSet(ctx, "scanning-malware:ttp", *content[i].ID , *content[i].Name ).Result()
		if err != nil {
			log.Printf("error=%v en id %s name %s", err.Error(), *content[i].ID , *content[i].Name)
			continue
		}
		if result != 1 {
			log.Printf("No se ha insertado correctamente para id=%s name=%s", *content[i].ID , *content[i].Name )
			continue
		}
		countInserted++
	}

	log.Printf("FIN Insertados %d Tamaño array %d", countInserted, lengthArrary)

}

func ReadFromFile(name string) ContainerProps {
	fileName := filepath.Join("./", name)
	fileContentByte, err := os.ReadFile(fileName)

	if err != nil {
		return ContainerProps{}
	}

	var dataParsed ContainerProps
	err = json.Unmarshal(fileContentByte, &dataParsed)
	if err != nil {
		return ContainerProps{}
	}
	return dataParsed
}

type ContainerProps []ContainerProp

type ContainerProp struct {
	ID          *string `json:"id,omitempty"`
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}
