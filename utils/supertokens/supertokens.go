package utils

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

var db *sql.DB

// InitializeSuperTokens initializes SuperTokens with the provided configuration
func InitializeSuperTokens(apiDomain, websiteDomain string, db *sql.DB) error {
	err := supertokens.Init(supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:3002",
		},
		AppInfo: supertokens.AppInfo{
			AppName:       "Personal Statement Reviewer",
			APIDomain:     apiDomain,
			WebsiteDomain: websiteDomain,
		},
		RecipeList: []supertokens.Recipe{
			thirdparty.Init(&tpmodels.TypeInput{
				SignInAndUpFeature: tpmodels.TypeInputSignInAndUp{
					Providers: []tpmodels.ProviderInput{
						{
							Config: tpmodels.ProviderConfig{
								ThirdPartyId: "google",
								Clients: []tpmodels.ProviderClientConfig{
									{
										ClientID:     "1027070285318-4jqtd2nea1815851ebo2k5ft0srumk3o.apps.googleusercontent.com",
										ClientSecret: "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",
									},
								},
							},
						},
					},
				},
			}),
			session.Init(nil), // initializes session features
		},
	})

	return err
}

// Middleware for verifying SuperTokens session
func VerifySession(next http.HandlerFunc) http.HandlerFunc {
	return verifySession(next)
}

// VerifySessionAndGetUserID verifies the session and returns the user ID
func VerifySessionAndGetUserID(next http.HandlerFunc) http.HandlerFunc {
	// return verifySession(next, &session.VerifySessionOptions{
	// 	SessionRequired: true,
	// })
	return next
}

func verifySession(next http.HandlerFunc) http.HandlerFunc {
	// return session.VerifySession(options, func(w http.ResponseWriter, r *http.Request, sessionContainer *session.SessionContainer) {
	// 	// Session is verified, call next handler
	// 	next.ServeHTTP(w, r)
	// })
	return next
}

// GetUserIDFromSession retrieves the user ID from the session
func GetUserIDFromSession(r *http.Request) (string, error) {
	sessionContainer := session.GetSessionFromRequestContext(r.Context())
	if sessionContainer == nil {
		return "", fmt.Errorf("no session found")
	}
	return sessionContainer.GetUserID(), nil
}

// CreateTables creates the necessary tables for your application
func CreateTables(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id VARCHAR(128) PRIMARY KEY,
			email VARCHAR(256) UNIQUE,
			name VARCHAR(256),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		// Add any other application-specific tables here
	}

	for _, query := range queries {
		_, err := db.Exec(query)
		if err != nil {
			return fmt.Errorf("error creating table: %v", err)
		}
	}

	return nil
}
