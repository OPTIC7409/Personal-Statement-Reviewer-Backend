package main

import (
	"fmt"
	api "psr/cmd/api"
	secrets "psr/cmd/api/secrets"
	"psr/database"
	"psr/utils/helpful/discord"
)

func main() {
	secrets.InitializeEnvVariables()

	server := api.NewAPIServer(":3002", database.GetConnection())

	go func() {
		if err := server.Run(); err != nil {
			fmt.Println("Error starting server:", err)
		}
	}()

	discord.SendMessage(discord.StartLog, "[API] has started")

	select {}
}
