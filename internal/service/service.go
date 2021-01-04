package service

import (
	"context"
	"fmt"
	"github.com/nokrPOC/internal/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"os"
	"github.com/nytm/messaging-helix-business-api/business"
)

type HermesService struct {
	helixClient *buisness.BusinessApiClient
	config *config.HermesServiceConfig
}

func NewService(ctx context.Context, config *config.HermesServiceConfig) (*HermesService, error) {
	helixClient := emailClient(config)
	return &HermesService{helixClient, config}, nil
}

// emailClient returns a client for helix email sending. The client is
// authenticated if credentials are given.
func emailClient(config *config.HermesServiceConfig) (*business.BusinessApiClient, error) {
	ctx := context.Background()
	// IMPORTANT! Need these two specific scopes for auth
	scopes := []string{
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/userinfo.profile",
	}
	var client *http.Client
	var err error
	if config.HelixCredentials == "" {
		client, err = google.DefaultClient(ctx, scopes...)
		if err != nil {
			return nil, fmt.Errorf("unable to get default client: %v", err)
		}
	} else {
		conf, err := google.JWTConfigFromJSON([]byte(config.HelixCredentials), scopes...)
		if err != nil {
			return nil, fmt.Errorf("unable to parse credentials: %v, %v", config.HelixCredentials, err)
		}
		client = oauth2.NewClient(ctx, conf.TokenSource(ctx))
	}
	return business.NewBusinessServiceHTTPClient(config.HelixHost, log.NewJSONLogger(os.Stdout), kithttp.SetClient(client)), nil
}
