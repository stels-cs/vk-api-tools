package VkApi

import (
	"fmt"
	"strings"
)

const overhead = "return[%s];"
const overheadSize = len(overhead)
const maxPackSize = 15000

type ExecutePack struct {
	calls []string
	size  int
}

func (ep *ExecutePack) IsFull() bool {
	if len(ep.calls) >= 25 {
		return true
	}
	if ep.Size() >= maxPackSize {
		return true
	}
	return false
}

// Проверяет возможно ли добавить метод в пакет
func (ep *ExecutePack) CanAdd(method Method) bool {
	data, err := method.toExecute()
	if err != nil {
		return false
	}
	if len(data)+ep.Size() > maxPackSize {
		return false
	}
	return true
}

// Добавляет метод в пакет и возворящает индекс, если индекс = -1
// значи добавить метод не получилось, пакет уже полный
func (ep *ExecutePack) Add(method Method) (int, error) {
	if ep.IsFull() {
		return -1, nil
	}
	data, err := method.toExecute()
	if err != nil {
		return -1, err
	}
	if len(data)+ep.Size() > maxPackSize {
		return -1, nil
	}
	ep.calls = append(ep.calls, data)
	ep.size += len(data) + 1
	return len(ep.calls) - 1, nil
}

// Размер execute кода в символах
func (ep *ExecutePack) Size() int {
	return ep.size + overheadSize
}

// Код execute запроса
func (ep *ExecutePack) GetCode() string {
	str := strings.Join(ep.calls, ",")
	return fmt.Sprintf(overhead, str)
}

// Количество запросов в пакете
func (ep *ExecutePack) Count() int {
	return len(ep.calls)
}

// Удаляет все запросы из пакета
func (ep *ExecutePack) Clear() {
	ep.calls = make([]string, 0)
	ep.size = 0
}
