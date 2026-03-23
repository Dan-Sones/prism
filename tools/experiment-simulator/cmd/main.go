package main

import (
	"encoding/json"
	"experiment-simulator/internal/parsers"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	fmt.Print(
		"  _____                      _                      _     ____  _                 _       _             \n" +
			" | ____|_  ___ __   ___ _ __(_)_ __ ___   ___ _ __ | |_  / ___|(_)_ __ ___  _   _| | __ _| |_ ___  _ __ \n" +
			" |  _| \\ \\/ / '_ \\ / _ \\ '__| | '_ " + "`" + " _ \\ / _ \\ '_ \\| __| \\___ \\| | '_ " + "`" + " _ \\| | | | |/ _" + "`" + " | __/ _ \\| '__|\n" +
			" | |___ >  <| |_) |  __/ |  | | | | | | |  __/ | | | |_   ___) | | | | | | | |_| | | (_| | || (_) | |   \n" +
			" |_____/_/\\_\\ .__/ \\___|_|  |_|_| |_| |_|\\___|_| |_|\\__| |____/|_|_| |_| |_|\\__,_|_|\\__,_|\\__\\___/|_|   \n" +
			"            |_|\n",
	)

	entries, err := os.ReadDir("resources")
	if err != nil {
		log.Fatal(err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		data, err := os.ReadFile(filepath.Join("resources", entry.Name()))
		if err != nil {
			log.Fatal(err)
		}

		parsed := parsers.ParseExperimentConfig(data)

		formatted, _ := json.MarshalIndent(parsed, "", "\t")
		fmt.Println(string(formatted))

	}

}
