package registry

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/laizy/web3"
	"github.com/laizy/web3/abi"
	"github.com/laizy/web3/utils"
)

type EventRegistry struct {
	events        map[web3.Hash]*abi.Event
	contractNames map[web3.Address]string
	lock          sync.RWMutex
}

func NewEventRegistry() *EventRegistry {
	return &EventRegistry{}
}

func (self *EventRegistry) RegisterContractAlias(c web3.Address, name string) {
	self.lock.Lock()
	defer self.lock.Unlock()
	if len(self.contractNames) == 0 {
		self.contractNames = map[web3.Address]string{}
	}
	self.contractNames[c] = name
}

func (self *EventRegistry) Register(e *abi.Event) {
	self.lock.Lock()
	defer self.lock.Unlock()
	if len(self.events) == 0 {
		self.events = map[web3.Hash]*abi.Event{}
	}
	if event := self.events[e.ID()]; event != nil {
		utils.EnsureTrue(event.Name == e.Name)
		return
	}
	self.events[e.ID()] = e
}

func (self *EventRegistry) RegisterFromAbi(abi *abi.ABI) {
	for _, e := range abi.Events {
		self.Register(e)
	}
}

func (self *EventRegistry) RegisterFromHumanString(eventStr string) {
	e := abi.MustNewEvent(eventStr)
	self.Register(e)
}

func (self *EventRegistry) ParseLog(log *web3.Log) (*web3.ParsedEvent, error) {
	if len(log.Topics) == 0 {
		return nil, errors.New("no topic found")
	}
	e := self.GetEvent(log.Topics[0])
	if e == nil {
		return nil, fmt.Errorf("can not parse log with sig: %s", log.Topics[0].String())
	}
	val, err := abi.ParseLog(e.Inputs, log)
	if err != nil {
		return nil, err
	}
	sig := e.DetailedSig()
	addr := log.Address.String()
	if name := self.contractNames[log.Address]; name != "" {
		addr = name
	}
	return &web3.ParsedEvent{
		Contract: addr,
		Sig:      sig,
		Values:   val,
	}, nil
}

func (self *EventRegistry) GetEvent(id web3.Hash) *abi.Event {
	self.lock.RLock()
	defer self.lock.RUnlock()

	return self.events[id]
}

func (self *EventRegistry) DumpLog(log *web3.Log) string {
	decoded, err := self.ParseLog(log)
	if err != nil {
		return err.Error()
	}

	buf, err := json.MarshalIndent(decoded, "", "  ")
	utils.Ensure(err)

	return string(buf)
}
