package graceful

import (
	"os"
	"sync"
)

type shutdownHandler struct {
	C         chan os.Signal
	mutex     sync.Mutex
	callbacks []ShutdownFunc
}

func newHandler(notify chan os.Signal) *shutdownHandler {
	return &shutdownHandler{
		C: notify,
	}
}

func (h *shutdownHandler) add(fn ShutdownFunc) {
	h.mutex.Lock()
	h.callbacks = append(h.callbacks, fn)
	h.mutex.Unlock()
}

func (h *shutdownHandler) get() []ShutdownFunc {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	cur := make([]ShutdownFunc, len(h.callbacks))
	copy(cur, h.callbacks)

	return cur
}

func (h *shutdownHandler) clear() {
	h.mutex.Lock()
	h.callbacks = make([]ShutdownFunc, 0)
	h.mutex.Unlock()
}
