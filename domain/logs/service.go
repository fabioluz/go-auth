package logs

type LogService struct {
	logRepository LogRepository
}

func NewLogService(logRepository LogRepository) *LogService {
	return &LogService{
		logRepository,
	}
}
