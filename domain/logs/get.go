package logs

import "context"

func (service *LogService) GetLogs(userID string, pageSize int, after string) ([]Log, error) {
	return service.logRepository.GetLogs(context.Background(), userID, pageSize, after)
}
