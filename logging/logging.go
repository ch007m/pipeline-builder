package logging

import (
	"log/slog"
	"os"
)

var Logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
