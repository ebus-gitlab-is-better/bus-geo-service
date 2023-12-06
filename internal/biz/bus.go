package biz

import "context"

type Bus struct {
	ID       uint    `redis:"id"`
	DriverID string  `redis:"driver_id"`
	Battery  uint    `redis:"battery"`
	Lat      float32 `redis:"lat"`
	Lon      float32 `redis:"lon"`
	// NextRoute *string `redis:"next_route"`
}

type BusRepo interface {
	Update(context.Context, *Bus) error
	List(context.Context) ([]*Bus, error)
	GetById(context.Context) (*Bus, error)
}

type BusUseCase struct {
	repo BusRepo
}

func NewBusUseCase(repo BusRepo) *BusUseCase {
	return &BusUseCase{repo: repo}
}

func (uc *BusUseCase) Update(ctx context.Context, bus *Bus) error {
	return uc.repo.Update(ctx, bus)
}

func (uc *BusUseCase) List(ctx context.Context) ([]*Bus, error) {
	return uc.repo.List(ctx)
}

func (uc *BusUseCase) GetById(ctx context.Context) (*Bus, error) {
	return uc.repo.GetById(ctx)
}
