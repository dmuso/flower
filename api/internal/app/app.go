package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"flower/api/internal/domain/organisation"
	"flower/api/internal/domain/planning"
	"flower/api/internal/domain/project"
	"flower/api/internal/domain/story"
	"flower/api/internal/domain/tenancy"
	"flower/api/internal/domain/user"
	"flower/api/internal/handlers"
	"flower/api/internal/platform/auth"
	clk "flower/api/internal/platform/clock"
	"flower/api/internal/platform/config"
	"flower/api/internal/platform/db"
	"flower/api/internal/platform/email"
	"flower/api/internal/platform/middleware"
	"flower/api/internal/ports"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type App struct {
	Config *config.Config
	DB     *sql.DB
	Router *gin.Engine
	Logger *zap.Logger
	server *http.Server

	rootCtx    context.Context
	rootCancel context.CancelFunc
}

type Option func(*options)

type options struct {
	Clock ports.Clock
}

func WithClock(clock ports.Clock) Option {
	return func(o *options) {
		o.Clock = clock
	}
}

func New(cfg *config.Config, opts ...Option) (*App, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config is required")
	}

	o := options{Clock: clk.System{}}
	for _, opt := range opts {
		opt(&o)
	}
	if o.Clock == nil {
		return nil, fmt.Errorf("clock is required")
	}

	logger, err := newLogger(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	dbConn, err := db.Connect(cfg.Database)
	if err != nil {
		logger.Error("failed to connect to database",
			zap.String("component", "app"),
			zap.String("operation", "new"),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Migrate(cfg.Database); err != nil {
		if closeErr := dbConn.Close(); closeErr != nil {
			logger.Error("failed to close database after migration error",
				zap.String("component", "app"),
				zap.String("operation", "new"),
				zap.Error(closeErr),
			)
		}
		logger.Error("failed to run migrations",
			zap.String("component", "app"),
			zap.String("operation", "new"),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	if cfg.Environment == "production" || cfg.Environment == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.Logger("/health"))
	router.Use(middleware.CORS(cfg.FrontendOrigin))

	userRepo := user.NewRepository(dbConn)
	orgRepo := organisation.NewRepository(dbConn)
	projectRepo := project.NewRepository(dbConn)
	access := tenancy.NewService(dbConn)
	dir := &directory{orgs: orgRepo, projects: projectRepo}
	mailer := email.NewOutbox()
	userSvc := user.NewService(userRepo, o.Clock, mailer, cfg.FrontendOrigin, dir)
	orgSvc := organisation.NewService(orgRepo, userRepo)
	projectSvc := project.NewService(projectRepo, access)
	planner := planning.NewPlanner(o.Clock, projectRepo.WindowSettings)
	storyRepo := story.NewRepository(dbConn)
	storySvc := story.NewService(storyRepo, access, planner)

	remember := func(c *gin.Context, projectID string) error {
		u := middleware.CurrentUser(c)
		if u == nil {
			return fmt.Errorf("remember project: no session user")
		}
		return userSvc.RememberProject(c.Request.Context(), u.UserID, u.SessionID, projectID)
	}

	router.Use(middleware.Session(userSvc.LookupSession, auth.HashToken))

	if err := handlers.SetupRoutes(router, &handlers.Dependencies{
		Version:       cfg.Version,
		Pinger:        dbConn,
		Users:         user.NewHandler(userSvc, cfg.CookieSecure),
		Organisations: organisation.NewHandler(orgSvc),
		Projects:      project.NewHandler(projectSvc, remember),
		Stories:       story.NewHandler(storySvc, remember),
	}); err != nil {
		if closeErr := dbConn.Close(); closeErr != nil {
			logger.Error("failed to close database after route setup error",
				zap.String("component", "app"),
				zap.String("operation", "new"),
				zap.Error(closeErr),
			)
		}
		return nil, fmt.Errorf("failed to setup routes: %w", err)
	}

	rootCtx, rootCancel := context.WithCancel(context.Background())

	return &App{
		Config:     cfg,
		DB:         dbConn,
		Router:     router,
		Logger:     logger,
		rootCtx:    rootCtx,
		rootCancel: rootCancel,
	}, nil
}

func (a *App) Start() error {
	if a == nil {
		return fmt.Errorf("app is nil")
	}
	if a.Config == nil {
		return fmt.Errorf("app config is nil")
	}

	addr := fmt.Sprintf(":%s", a.Config.APIPort)
	a.server = &http.Server{
		Addr:              addr,
		Handler:           a.Router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	a.Logger.Info("starting Flower API",
		zap.String("component", "app"),
		zap.String("operation", "start"),
		zap.String("addr", addr),
		zap.String("version", a.Config.Version),
	)

	go func() {
		<-a.rootCtx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := a.server.Shutdown(shutdownCtx); err != nil {
			a.Logger.Error("failed to shut down HTTP server",
				zap.String("component", "app"),
				zap.String("operation", "shutdown"),
				zap.Error(err),
			)
		}
	}()

	if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("failed to listen on %s: %w", addr, err)
	}
	return nil
}

func (a *App) Cancel() {
	if a == nil || a.rootCancel == nil {
		return
	}
	a.rootCancel()
}

func (a *App) Close() error {
	if a == nil {
		return fmt.Errorf("app is nil")
	}
	a.Cancel()
	if a.DB != nil {
		if err := a.DB.Close(); err != nil {
			return fmt.Errorf("failed to close database: %w", err)
		}
	}
	if a.Logger != nil {
		if err := a.Logger.Sync(); err != nil {
			return fmt.Errorf("failed to sync logger: %w", err)
		}
	}
	return nil
}

func newLogger(cfg *config.Config) (*zap.Logger, error) {
	level := zapcore.InfoLevel
	if err := level.UnmarshalText([]byte(cfg.LogLevel)); err != nil {
		return nil, fmt.Errorf("invalid LOG_LEVEL %q: %w", cfg.LogLevel, err)
	}

	zapCfg := zap.NewProductionConfig()
	if cfg.Environment == "development" || cfg.Environment == "test" {
		zapCfg = zap.NewDevelopmentConfig()
	}
	zapCfg.Level = zap.NewAtomicLevelAt(level)

	logger, err := zapCfg.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build logger: %w", err)
	}
	zap.ReplaceGlobals(logger)
	return logger, nil
}
