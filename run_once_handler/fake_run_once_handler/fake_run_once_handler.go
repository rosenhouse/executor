package fake_run_once_handler

import (
	"github.com/cloudfoundry-incubator/runtime-schema/models"
	"sync"
)

type FakeTaskHandler struct {
	numberOfCalls   int
	handledTasks map[string]string
	mutex           *sync.RWMutex
	cancel          <-chan struct{}
}

func New() *FakeTaskHandler {
	return &FakeTaskHandler{
		handledTasks: make(map[string]string),
		mutex:           &sync.RWMutex{},
	}
}

func (handler *FakeTaskHandler) Task(runOnce *models.Task, executorId string, cancel <-chan struct{}) {
	handler.mutex.Lock()
	defer handler.mutex.Unlock()

	handler.cancel = cancel

	_, present := handler.handledTasks[runOnce.Guid]
	if !present {
		handler.numberOfCalls++
		handler.handledTasks[runOnce.Guid] = executorId
	}
}

func (handler *FakeTaskHandler) NumberOfCalls() int {
	handler.mutex.Lock()
	defer handler.mutex.Unlock()

	return handler.numberOfCalls
}

func (handler *FakeTaskHandler) HandledTasks() map[string]string {
	handler.mutex.RLock()
	defer handler.mutex.RUnlock()

	handled := map[string]string{}

	for k, v := range handler.handledTasks {
		handled[k] = v
	}

	return handled
}

func (handler *FakeTaskHandler) GetCancel() <-chan struct{} {
	handler.mutex.Lock()
	defer handler.mutex.Unlock()

	return handler.cancel
}
