package tracker

import "sync/atomic"

type entity struct {
	id string

	bn int64  // biggest nonce  atomic
	ln *queue // locked nonces
	rn *queue // released nonces
}

func newEntity(id string, bn int64) *entity {
	return &entity{
		id: id,
		bn: bn,
		ln: newQ(),
		rn: newQ(),
	}
}

func (e *entity) nonce() int {
	if e.rn.size() > 0 {
		n := e.rn.pop()
		e.ln.push(n)
		return n
	}

	n := atomic.AddInt64(&e.bn, 1)

	e.ln.push(int(n))
	return int(n)
}

func (e *entity) burn(n int) {
	e.ln.remove(n)
}

func (e *entity) release(n int) {
	e.ln.remove(n)
	e.rn.push(n)
}
