package store

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/smahjoub/events-api/objects"
)

// IEventStore is the database interface for storing Events
type IEventStore interface {
	Get(ctx context.Context, in *objects.GetRequest) (*objects.Event, error)
	List(ctx context.Context, in *objects.ListRequest) ([]*objects.Event, error)
	Create(ctx context.Context, in *objects.CreateRequest) error
	UpdateDetails(ctx context.Context, in *objects.UpdateDetailsRequest) error
	Cancel(ctx context.Context, in *objects.CancelRequest) error
	Reschedule(ctx context.Context, in *objects.RescheduleRequest) error
	Delete(ctx context.Context, in *objects.DeleteRequest) error
}

func init() {
	rand.Seed(time.Now().UTC().Unix())
}

// GenerateUniqueID will returns a time based sortable unique id
func GenerateUniqueID() string {
	word := []byte("0987654321")
	rand.Shuffle(len(word), func(i, j int) {
		word[i], word[j] = word[j], word[i]
	})
	now := time.Now().UTC()
	return fmt.Sprintf("%010v-%010v-%s", now.Unix(), now.Nanosecond(), string(word))
}
