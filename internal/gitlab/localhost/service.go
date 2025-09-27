package localhost

type Service struct {
	events EventsRepo
}

type EventsRepo interface {
}

func NewService(events EventsRepo) *Service {
	return &Service{
		events: events,
	}
}
