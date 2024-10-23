package secrets

import (
	"context"
	"fmt"
	"os"
	"psr/utils/helpful/discord"

	infisical "github.com/infisical/go-sdk"
)

var envVars map[string]string

func fetchEnvVariables() map[string]string {
	client := infisical.NewInfisicalClient(context.Background(), infisical.Config{
		SiteUrl: "https://app.infisical.com",
	})

	_, err := client.Auth().UniversalAuthLogin("79169051-4a9b-482f-9441-f8c61108637b", "1d2136979592c60e5753d34d40994230b39a6498b5670f4542521feaaad83bba")

	if err != nil {
		fmt.Printf("Database Authentication failed: %v", err)
		os.Exit(1)
	}

	secrets, err := client.Secrets().List(infisical.ListSecretsOptions{
		ProjectID:   "ed5d43cc-1aea-4f58-9e47-3c65d24026b0",
		Environment: "dev",
		SecretPath:  "/",
	})

	if err != nil {
		discord.SendMessage(discord.ErrorLog, fmt.Sprintf("Failed to fetch secrets: %v", err))
		fmt.Printf("Failed to fetch secrets: %v", err)
		os.Exit(1)
	}

	envVars := make(map[string]string)

	for _, secret := range secrets {
		envVars[secret.SecretKey] = secret.SecretValue
	}

	discord.SendMessage(discord.MessageLog, fmt.Sprintf("Loaded %d environment variables.", len(envVars)))

	return envVars
}

func InitializeEnvVariables() {
	envVars = fetchEnvVariables()
	fmt.Printf("Loaded %d environment variables.\n", len(envVars))
}

func GetEnvVariable(key string) string {
	value := envVars[key]
	if value == "" {
		fmt.Printf("Environment variable %s not set.\n", key)
	}
	return value
}
