package Vk

import (
	"strings"
	"fmt"
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

func (ep *ExecutePack) Add(method ApiMethod) (int, error) {
	if ep.IsFull() {
		return -1, nil
	}
	data, err:= method.toExecute()
	if err != nil {
		return -1, err
	}
	if len(data) + ep.Size() > maxPackSize {
		return -1, nil
	}
	ep.calls = append(ep.calls, data)
	ep.size += len(data) + 1
	return len(ep.calls) - 1, nil
}

func (ep *ExecutePack) Size() int {
	return ep.size + overheadSize
}

func (ep *ExecutePack) GetCode() string {
	str := strings.Join(ep.calls, ",")
	return fmt.Sprintf(overhead, str)
}

func (ep *ExecutePack) Count() int {
	return len(ep.calls)
}

func (ep *ExecutePack) Clear() {
	ep.calls = make([]string, 0)
	ep.size = 0
}