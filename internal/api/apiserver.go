package api

import (
	"database/sql"
	"os"

	sqlstore "github.com/artemiyKew/todo-list-rest-api/internal/store/SQLStore"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Start(config *Config) error {
	db, err := newDB(config.DataBaseURL)
	if err != nil {
		return err
	}
	defer db.Close()
	store := sqlstore.New(db)
	srv := newServer(store)
	logger, err := configureLogger()
	if err != nil {
		return err
	}
	srv.logger = logger

	return fasthttp.ListenAndServe(config.BindAddr, srv.router.Handler())
}

func newDB(dbURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func configureLogger() (*zap.Logger, error) {
	filename := "logs.log"
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder

	// Create file and console encoders
	fileEncoder := zapcore.NewJSONEncoder(config)
	consoleEncoder := zapcore.NewConsoleEncoder(config)

	logFile, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	// Create writers for file and console
	fileWriter := zapcore.AddSync(logFile)
	consoleWriter := zapcore.AddSync(os.Stdout)

	// Set the log level
	defaultLogLevel := zapcore.InfoLevel

	// Create cores for writing to the file and console
	fileCore := zapcore.NewCore(fileEncoder, fileWriter, defaultLogLevel)
	consoleCore := zapcore.NewCore(consoleEncoder, consoleWriter, defaultLogLevel)

	// Combine cores
	core := zapcore.NewTee(fileCore, consoleCore)
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return logger, nil
}
