package ttl

import (
	"fmt"
	"log/slog"
	"time"
)

type Checker interface {
	CheckForTTL() (int, error)
}

func TTLChecker(logger *slog.Logger, checker Checker) error {

	ticker := time.NewTicker(time.Minute * 1)
	for {
		select {
		case <-ticker.C:
			id, err := checker.CheckForTTL()
			if err != nil {
				logger.Error("no expired segments")
				continue
			}
			logger.Info(fmt.Sprintf("the segment of user %d has been expired", id))
		}
	}

}
