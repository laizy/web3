package web3

import (
	"errors"
	"sync"
)

type LogParser interface {
	ParseLog(log *Log) (*ParsedEvent, error)
}

type NilParser struct{}

func (self *NilParser) ParseLog(log *Log) (*ParsedEvent, error) {
	return nil, errors.New("can not parse log")
}

var parser LogParser = &NilParser{}
var mutex = &sync.RWMutex{}

func RegisterParser(p LogParser) {
	mutex.Lock()
	defer mutex.Unlock()
	parser = p
}

func GetParser() LogParser {
	mutex.RLock()
	defer mutex.RUnlock()
	return parser
}
