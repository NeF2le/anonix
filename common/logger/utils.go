package logger

import (
	"encoding/base64"
	"log/slog"
)

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}

func Base64(key string, value []byte) slog.Attr {
	v64 := base64.StdEncoding.EncodeToString(value)
	return slog.Attr{
		Key:   key,
		Value: slog.StringValue(v64),
	}
}
