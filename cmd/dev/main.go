package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"

	"go.uber.org/zap"
	"google.golang.org/api/option"

	authsession "github.com/vinylhousegarage/idpproxy/internal/auth/session"
	firesessionstore "github.com/vinylhousegarage/idpproxy/internal/auth/sessionstore/firestore"
	authuser "github.com/vinylhousegarage/idpproxy/internal/auth/user"
	authcodeservice "github.com/vinylhousegarage/idpproxy/internal/authcode/service"
	authcodestore "github.com/vinylhousegarage/idpproxy/internal/authcode/store"
	"github.com/vinylhousegarage/idpproxy/internal/config"
	"github.com/vinylhousegarage/idpproxy/internal/deps"
	firebaseutil "github.com/vinylhousegarage/idpproxy/internal/firebase"
	"github.com/vinylhousegarage/idpproxy/internal/kms"
	"github.com/vinylhousegarage/idpproxy/internal/oauth/github/callback"
	githubme "github.com/vinylhousegarage/idpproxy/internal/oauth/github/me"
	githubstore "github.com/vinylhousegarage/idpproxy/internal/oauth/github/store"
	"github.com/vinylhousegarage/idpproxy/internal/router"
	"github.com/vinylhousegarage/idpproxy/internal/server"
	"github.com/vinylhousegarage/idpproxy/public"
)

const (
	sessionCollectionName = "sessions"
	sessionTTL            = 24 * time.Hour
	sessionIDBytes        = 32
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic("failed to initialize logger: " + err.Error())
	}
	defer func() {
		_ = logger.Sync()
	}()

	ctx := context.Background()

	firebaseCfg, err := config.LoadFirebaseConfig()
	if err != nil {
		logger.Fatal(
			"failed to load Firebase config",
			zap.Error(err),
		)
	}

	app, err := firebaseutil.NewFirebaseApp(
		ctx,
		option.WithCredentialsJSON(firebaseCfg.CredentialsJSON),
	)
	if err != nil {
		logger.Fatal(
			"failed to initialize Firebase App",
			zap.Error(err),
		)
	}

	authClient, err := firebaseutil.NewAuthClient(ctx, app, logger)
	if err != nil {
		logger.Fatal(
			"failed to initialize Firebase Auth client",
			zap.Error(err),
		)
	}

	firestoreClient, err := firebaseutil.NewFirestoreClient(
		ctx,
		app,
		logger,
	)
	if err != nil {
		logger.Fatal(
			"failed to initialize Firestore client",
			zap.Error(err),
		)
	}
	defer func() {
		if err := firestoreClient.Close(); err != nil {
			logger.Error(
				"failed to close Firestore client",
				zap.Error(err),
			)
		}
	}()

	serviceAccountCfg := config.LoadServiceAccountConfig()

	kmsClient, err := kms.NewClient(
		ctx,
		serviceAccountCfg.ImpersonateSA,
	)
	if err != nil {
		logger.Fatal(
			"failed to initialize KMS client",
			zap.Error(err),
		)
	}
	defer func() {
		if err := kmsClient.Close(); err != nil {
			logger.Error(
				"failed to close KMS client",
				zap.Error(err),
			)
		}
	}()

	githubTokenKMSCfg, err := config.LoadGitHubTokenKMSConfig()
	if err != nil {
		logger.Fatal(
			"failed to load GitHub token KMS config",
			zap.Error(err),
		)
	}

	kmsAdapter, err := kms.NewAdapter(
		kmsClient,
		githubTokenKMSCfg.KeyName,
		[]byte(config.GitHubTokenKMSAAD),
	)
	if err != nil {
		logger.Fatal(
			"failed to initialize KMS adapter",
			zap.Error(err),
		)
	}

	tokenEncryptor, err := githubstore.NewKMSTokenEncryptor(
		githubTokenKMSCfg.KeyName,
		kmsAdapter,
	)
	if err != nil {
		logger.Fatal(
			"failed to initialize GitHub token encryptor",
			zap.Error(err),
		)
	}

	githubTokenRepo := githubstore.NewFirestoreGitHubTokenRepo(
		firestoreClient,
		tokenEncryptor,
	)

	sessionRepo := firesessionstore.NewRepository(
		firestoreClient,
		sessionCollectionName,
	)

	sessionUsecase := &authsession.Usecase{
		Repo:        sessionRepo,
		Now:         time.Now,
		TTL:         sessionTTL,
		IDGenerator: generateSessionID,
	}

	googleDeps := deps.NewGoogleDeps(authClient, logger)

	githubCfg, err := config.LoadGitHubDevOAuthConfig()
	if err != nil {
		logger.Fatal(
			"failed to load GitHub config",
			zap.Error(err),
		)
	}

	githubOAuthDeps := deps.NewGitHubOAuthDeps(githubCfg, logger)

	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	githubAPICfg := config.LoadGitHubAPIConfig()
	githubAPIDeps := deps.NewGitHubAPIDeps(
		githubAPICfg,
		httpClient,
		logger,
	)

	githubMeService := githubme.NewService(
		githubTokenRepo,
		httpClient,
	)

	githubMeHandler := githubme.NewGitHubMeHandler(
		githubMeService,
	)

	githubUserService := authuser.NewService(authClient)

	proxyCodeService := authcodeservice.NewService(
		authcodestore.NewMemoryStore(),
	)

	githubCallbackHandler := callback.NewGitHubCallbackHandler(
		githubOAuthDeps,
		githubAPIDeps,
		githubUserService,
		proxyCodeService,
		githubTokenRepo,
		sessionUsecase,
		githubOAuthDeps.Config.ClientID,
	)

	systemDeps := deps.NewSystemDeps(
		config.GoogleOIDCMetadataURL,
		httpClient,
		logger,
	)

	d := router.NewRouterDeps(
		public.PublicFS,
		githubAPIDeps,
		githubOAuthDeps,
		githubCallbackHandler,
		githubMeHandler,
		googleDeps,
		logger,
		systemDeps,
	)

	r := router.NewRouter(d)

	logger.Info(
		"starting idpproxy (dev)",
		zap.String("addr", ":"+config.GetPort()),
	)

	server.StartServer(r, logger)
}

func generateSessionID() (string, error) {
	b := make([]byte, sessionIDBytes)

	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}
