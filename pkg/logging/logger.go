package logging

import (
	"log"
	"log/slog"
	"os"
)



func  SetupLogger(modeLog string) *slog.Logger{
	var logger *slog.Logger

	switch modeLog {
	case "debug": 
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "jsonDebug":	
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "jsonInfo":	
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		log.Fatal("not init modelog: ", modeLog)
		
	}

	return logger
}

func Err(err error) slog.Attr {
	return slog.Attr{
		Key: "error",
		Value: slog.StringValue(err.Error()),
	}
}