package main

import (
	"fmt"
	"net/http"
	"path/filepath"
)

var dbUrlRemote = ""
var dbAuthToken = ""

func main() {

	// Doesn't work because not an absolute path (should work)
	config := NewClientConfig("file://./data/sqlite.db", nil)
	_, err := NewClient(config)
	if err != nil {
		fmt.Printf("Error Local DB Relative File Path: %s\n", err)
	}

	// Works
	absPath, _ := filepath.Abs("./data/sqlite.db")
	config = NewClientConfig("file://"+absPath, nil)
	dbLocal, err := NewClient(config)
	if err != nil {
		fmt.Printf("Error Local Abs Path: %s\n", err)
	}

	// Works at the beginning but times out after ~10 seconds
	config = NewClientConfig(dbUrlRemote, &dbAuthToken)
	dbRemoteNoEmbed, err := NewClient(config)
	if err != nil {
		fmt.Printf("Error Remote Not Embedded: %s\n", err)
	}

	// Works - sometimes randomly (happend 3-4 Times with 100x restarts) leaves a sqlite.db which is empty after restarting. This breaks the whole app and there's no way of recovering except for deleting the file.
	config = NewClientConfig(dbUrlRemote, &dbAuthToken)
	dbRemoteEmbed, err := NewEmbeddedClient(config)
	if err != nil {
		fmt.Printf("Error Remote Embedded: %s", err)
	}

	allClientsByName := map[string]*Client{
		"local":  dbLocal,
		"remote": dbRemoteNoEmbed,
		"embed":  dbRemoteEmbed,
	}

	// start simple http server with one get route
	server := &http.Server{
		Addr: ":8080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("Recovered from panic: %s\n", r)
				}
			}()
			for name, client := range allClientsByName {
				rows, err := client.Query("SELECT name FROM sqlite_master WHERE type='table' ORDER BY name;")
				if err != nil {
					fmt.Printf("Client: %s, Error: %s\n", name, err)
				}
				defer rows.Close()

				rowsForClient := "Rows for client" + name + ":\n"
				for rows.Next() {
					var val string
					err = rows.Scan(&val)
					if err != nil {
						fmt.Printf("Client: %s, Error: %s\n", val, err)
					}
					rowsForClient += val + "\n"
				}
				w.Write([]byte(rowsForClient + "\n\n"))
			}
		}),
	}

	fmt.Println("Server started at :8080")
	err = server.ListenAndServe()
	if err != nil {
		fmt.Printf("Error ListenAndServe: %s", err)
		panic(err)
	}
}
