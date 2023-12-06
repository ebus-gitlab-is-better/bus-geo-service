package data

import (
	"bus-geo-service/internal/biz"
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type Bus struct {
	ID       uint    `redis:"id"`
	DriverID string  `redis:"driver_id"`
	Battery  uint    `redis:"battery"`
	Lat      float32 `redis:"lat"`
	Lon      float32 `redis:"lon"`
}

type busRepo struct {
	data *Data
}

func NewBusRepo(data *Data) biz.BusRepo {
	return &busRepo{data: data}
}

// GetById implements biz.BusRepo.
func (*busRepo) GetById(context.Context) (*biz.Bus, error) {
	panic("unimplemented")
}

// List implements biz.BusRepo.
func (*busRepo) List(context.Context) ([]*biz.Bus, error) {
	panic("unimplemented")
}

// Update implements biz.BusRepo.
func (r *busRepo) Update(ctx context.Context, bus *biz.Bus) error {
	busBytes, err := json.Marshal(bus)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("bus:%d", bus.ID)

	err = r.data.rdb.Set(ctx, key, busBytes, 30*time.Minute).Err()
	if err != nil {
		return err
	}

	r.data.cent.Publish(context.TODO(), "bus:public", busBytes)

	return nil
}
