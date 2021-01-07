package service

import (
	"context"
	"github.com/NYTimes/gizmo/server/kit"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/nokrPOC/internal/config"
	"github.com/nokrPOC/internal/service/modules"
)

type HermesService struct {
	//helixClient *buisness.BusinessApiClient
	pubSub *modules.PubSubHandler
	config *config.HermesServiceConfig
	Logger
}

func NewService(ctx context.Context, config *config.HermesServiceConfig) (*HermesService, error) {
	//helixClient, err := emailClient(config)
	//if err != nil {
	//	_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
	//	os.Exit(2)
	//}

	lg, lgClose, err := kit.NewLogger(context.Background(), config.ProjectID)
	// if logger fails to initialize use default logger and return new Service
	if err != nil {
		lg = log.NewLogfmtLogger(log.StdlibWriter{})
		lgClose = nil
		lg.Log("level", level.WarnValue(), "message", "Unable to start up logger defaulting to standard logger")
	}
	// logger is attached to the service but request scoped logger should be obtained from
	// the context.
	logger := serviceLog{lg, lgClose}

	pubSub := modules.NewPubSubService(config)
	return &HermesService{pubSub, config, &logger}, nil
}

// set up a new pub sub service that listens for messages

// emailClient returns a client for helix email sending. The client is
// authenticated if credentials are given.
//func emailClient(config *config.HermesServiceConfig) (*business.BusinessApiClient, error) {
//	ctx := context.Background()
//	// IMPORTANT! Need these two specific scopes for auth
//	scopes := []string{
//		"https://www.googleapis.com/auth/userinfo.email",
//		"https://www.googleapis.com/auth/userinfo.profile",
//	}
//	var client *http.Client
//	var err error
//	if config.Helix.HelixCredentials == "" {
//		client, err = google.DefaultClient(ctx, scopes...)
//		if err != nil {
//			return nil, fmt.Errorf("unable to get default client: %v", err)
//		}
//	} else {
//		conf, err := google.JWTConfigFromJSON([]byte(config.Helix.HelixCredentials), scopes...)
//		if err != nil {
//			return nil, fmt.Errorf("unable to parse credentials: %v, %v", config.Helix.HelixCredentials, err)
//		}
//		client = oauth2.NewClient(ctx, conf.TokenSource(ctx))
//	}
//	return business.NewBusinessServiceHTTPClient(config.Helix.HelixHost, log.NewJSONLogger(os.Stdout), kithttp.SetClient(client)), nil
//}
