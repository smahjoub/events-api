package store

import (
	"context"
	"log"
	"os"

	"github.com/smahjoub/events-api/errors"
	"github.com/smahjoub/events-api/objects"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type pg struct {
	db *gorm.DB
}

// NewPostgresEventStore returns a postgres implementation of Event store
func NewPostgresEventStore(conn string) IEventStore {
	// create database connection
	db, err := gorm.Open(postgres.Open(conn),
		&gorm.Config{
			Logger: logger.New(
				log.New(os.Stdout, "", log.LstdFlags),
				logger.Config{
					LogLevel: logger.Info,
					Colorful: true,
				},
			),
		},
	)
	if err != nil {
		panic("Enable to connect to database: " + err.Error())
	}
	if err := db.AutoMigrate(&objects.Event{}); err != nil {
		panic("Enable to migrate database: " + err.Error())
	}
	// return store implementation
	return &pg{db: db}
}

func (p *pg) Get(ctx context.Context, in *objects.GetRequest) (*objects.Event, error) {
	evt := &objects.Event{}
	// take event where id == uid from database
	err := p.db.WithContext(ctx).Take(evt, "id = ?", in.ID).Error
	if err == gorm.ErrRecordNotFound {
		// not found
		return nil, errors.ErrEventNotFound
	}
	return evt, err
}

func (p *pg) List(ctx context.Context, in *objects.ListRequest) ([]*objects.Event, error) {
	if in.Limit == 0 || in.Limit > objects.MaxListLimit {
		in.Limit = objects.MaxListLimit
	}
	query := p.db.WithContext(ctx).Limit(in.Limit)
	if in.After != "" {
		query = query.Where("id > ?", in.After)
	}
	if in.Name != "" {
		query = query.Where("name ilike ?", "%"+in.Name+"%")
	}
	list := make([]*objects.Event, 0, in.Limit)
	err := query.Order("id").Find(&list).Error
	return list, err
}

func (p *pg) Create(ctx context.Context, in *objects.CreateRequest) error {
	if in.Event == nil {
		return errors.ErrObjectIsRequired
	}
	in.Event.ID = GenerateUniqueID()
	in.Event.Status = objects.Original
	in.Event.CreatedOn = p.db.NowFunc()
	return p.db.WithContext(ctx).
		Create(in.Event).
		Error
}

func (p *pg) UpdateDetails(ctx context.Context, in *objects.UpdateDetailsRequest) error {
	evt := &objects.Event{
		ID:          in.ID,
		Name:        in.Name,
		Description: in.Description,
		Website:     in.Website,
		Address:     in.Address,
		PhoneNumber: in.PhoneNumber,
		UpdatedOn:   p.db.NowFunc(),
	}
	return p.db.WithContext(ctx).Model(evt).
		Select("name", "description", "website", "address", "phone_number", "updated_on").
		Updates(evt).
		Error
}

func (p *pg) Cancel(ctx context.Context, in *objects.CancelRequest) error {
	evt := &objects.Event{
		ID:          in.ID,
		Status:      objects.Cancelled,
		CancelledOn: p.db.NowFunc(),
	}
	return p.db.WithContext(ctx).Model(evt).
		Select("status", "cancelled_on").
		Updates(evt).
		Error
}

func (p *pg) Reschedule(ctx context.Context, in *objects.RescheduleRequest) error {
	evt := &objects.Event{
		ID:            in.ID,
		Slot:          in.NewSlot,
		Status:        objects.Rescheduled,
		RescheduledOn: p.db.NowFunc(),
	}
	return p.db.WithContext(ctx).Model(evt).
		Select("status", "start_time", "end_time", "rescheduled_on").
		Updates(evt).
		Error
}

func (p *pg) Delete(ctx context.Context, in *objects.DeleteRequest) error {
	evt := &objects.Event{ID: in.ID}
	return p.db.WithContext(ctx).Model(evt).
		Delete(evt).
		Error
}
