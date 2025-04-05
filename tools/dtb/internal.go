package dtb

import "time"

func DoWithAttempts(callback func() error, maxAttempts int, delay time.Duration) error {
	// function executes callback several times
	var err error
	for maxAttempts > 0 {
		if err = callback(); err != nil {
			time.Sleep(delay)
			maxAttempts--
			continue
		}
		return nil
	}
	return err
}
