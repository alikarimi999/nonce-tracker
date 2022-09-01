package tracker

import (
	"sync"
)

type Tracker struct {
	mux *sync.Mutex
	as  map[string]*entity
}

func NewTracker() *Tracker {
	return &Tracker{
		mux: &sync.Mutex{},
		as:  make(map[string]*entity),
	}
}

func (t *Tracker) AddEntity(id string, latest_nonce int64) {
	t.mux.Lock()
	defer t.mux.Unlock()

	for s := range t.as {
		if s == id {
			return
		}
	}
	t.as[id] = newEntity(id, latest_nonce)
}

func (t *Tracker) NewTicket(id string) *Ticket {
	t.mux.Lock()
	defer t.mux.Unlock()

	if a, ok := t.as[id]; ok {
		return &Ticket{Owner: id, Nonce: a.nonce()}
	}
	return &Ticket{Owner: "", Nonce: 0}
}

func (t *Tracker) BurnTicket(ti *Ticket) {
	t.mux.Lock()
	defer t.mux.Unlock()

	if a, ok := t.as[ti.Owner]; ok {
		a.burn(ti.Nonce)
	}
}

func (t *Tracker) ReleaseTicket(ti *Ticket) {
	t.mux.Lock()
	defer t.mux.Unlock()

	a, ok := t.as[ti.Owner]
	if ok {
		a.release(ti.Nonce)
	}

}
