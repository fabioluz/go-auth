package logs

import "context"

func (service *LogService) Get(userID string, pageSize int, after string) ([]Log, error) {
	return service.logRepository.Get(context.Background(), userID, pageSize, after)
}
