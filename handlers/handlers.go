package handlers

import (
	"io/ioutil"
	"net/http"

	"github.com/smahjoub/events-api/errors"
	"github.com/smahjoub/events-api/objects"
	"github.com/smahjoub/events-api/store"
)

// IEventHandler is implement all the handlers
type IEventHandler interface {
	Get(w http.ResponseWriter, r *http.Request)
	List(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	UpdateDetails(w http.ResponseWriter, r *http.Request)
	Cancel(w http.ResponseWriter, r *http.Request)
	Reschedule(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	store store.IEventStore
}

// NewEventHandler return current IEventHandler implementation
func NewEventHandler(store store.IEventStore) IEventHandler {
	return &handler{store: store}
}

func (h *handler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		WriteError(w, errors.ErrValidEventIDIsRequired)
		return
	}
	evt, err := h.store.Get(r.Context(), &objects.GetRequest{ID: id})
	if err != nil {
		WriteError(w, err)
		return
	}
	WriteResponse(w, &objects.EventResponseWrapper{Event: evt})
}

func (h *handler) List(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	// after
	after := values.Get("after")
	// name
	name := values.Get("name")
	// limit
	limit, err := IntFromString(w, values.Get("limit"))
	if err != nil {
		return
	}
	// list events
	list, err := h.store.List(r.Context(), &objects.ListRequest{
		Limit: limit,
		After: after,
		Name:  name,
	})
	if err != nil {
		WriteError(w, err)
		return
	}
	WriteResponse(w, &objects.EventResponseWrapper{Events: list})
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		WriteError(w, errors.ErrUnprocessableEntity)
		return
	}
	evt := &objects.Event{}
	if Unmarshal(w, data, evt) != nil {
		return
	}
	if err := checkSlot(evt.Slot); err != nil {
		WriteError(w, err)
		return
	}
	if err = h.store.Create(r.Context(), &objects.CreateRequest{Event: evt}); err != nil {
		WriteError(w, err)
		return
	}
	WriteResponse(w, &objects.EventResponseWrapper{Event: evt})
}

func (h *handler) UpdateDetails(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		WriteError(w, errors.ErrUnprocessableEntity)
		return
	}
	req := &objects.UpdateDetailsRequest{}
	if Unmarshal(w, data, req) != nil {
		return
	}

	// check if event exist
	if _, err := h.store.Get(r.Context(), &objects.GetRequest{ID: req.ID}); err != nil {
		WriteError(w, err)
		return
	}

	if err = h.store.UpdateDetails(r.Context(), req); err != nil {
		WriteError(w, err)
		return
	}
	WriteResponse(w, &objects.EventResponseWrapper{})
}

func (h *handler) Cancel(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		WriteError(w, errors.ErrValidEventIDIsRequired)
		return
	}

	// check if event exist
	if _, err := h.store.Get(r.Context(), &objects.GetRequest{ID: id}); err != nil {
		WriteError(w, err)
		return
	}

	if err := h.store.Cancel(r.Context(), &objects.CancelRequest{ID: id}); err != nil {
		WriteError(w, err)
		return
	}
	WriteResponse(w, &objects.EventResponseWrapper{})
}

func (h *handler) Reschedule(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		WriteError(w, errors.ErrUnprocessableEntity)
		return
	}
	req := &objects.RescheduleRequest{}
	if Unmarshal(w, data, req) != nil {
		return
	}
	if err := checkSlot(req.NewSlot); err != nil {
		WriteError(w, err)
		return
	}

	// check if event exist
	if _, err := h.store.Get(r.Context(), &objects.GetRequest{ID: req.ID}); err != nil {
		WriteError(w, err)
		return
	}

	if err = h.store.Reschedule(r.Context(), req); err != nil {
		WriteError(w, err)
		return
	}
	WriteResponse(w, &objects.EventResponseWrapper{})
}

func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		WriteError(w, errors.ErrValidEventIDIsRequired)
		return
	}

	// check if event exist
	if _, err := h.store.Get(r.Context(), &objects.GetRequest{ID: id}); err != nil {
		WriteError(w, err)
		return
	}

	if err := h.store.Delete(r.Context(), &objects.DeleteRequest{ID: id}); err != nil {
		WriteError(w, err)
		return
	}
	WriteResponse(w, &objects.EventResponseWrapper{})
}
