package callers

import (
	"time"
)

func Retry(operation func() error, maxRetries uint, baseDelay time.Duration) error {
	var err error

	for n := uint(0); n < maxRetries; n++ {
		if n > 0 {
			time.Sleep(baseDelay * time.Duration(1<<n))
		}

		err = operation()
		if err == nil {
			return nil
		}
	}
	return err
}
