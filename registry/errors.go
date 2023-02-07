package registry

import (
	"fmt"
	"sync"

	"github.com/laizy/web3"
	"github.com/laizy/web3/abi"
	"github.com/laizy/web3/utils"
)

type ErrorRegistry struct {
	errors        map[[4]byte]*abi.Error
	contractNames map[web3.Address]string
	lock          sync.RWMutex
}

func NewErrorRegistry() *ErrorRegistry {
	errorRegister := &ErrorRegistry{}
	for _, e := range abi.DefaultError() {
		errorRegister.Register(e)
	}
	return errorRegister
}

func (self *ErrorRegistry) Register(e *abi.Error) {
	self.lock.Lock()
	defer self.lock.Unlock()
	if len(self.errors) == 0 {
		self.errors = make(map[[4]byte]*abi.Error)
	}
	var id [4]byte
	copy(id[:], e.ID())
	if errors := self.errors[id]; errors != nil {
		utils.EnsureTrue(errors.Name == e.Name)
		return
	}
	self.errors[id] = e
}

func (self *ErrorRegistry) RegisterFromAbi(abi *abi.ABI) {
	for _, e := range abi.Errors {
		self.Register(e)
	}
}

func (self *ErrorRegistry) ParseError(info []byte) (string, error) {
	if len(info) < 4 {
		return "", fmt.Errorf("short info")
	}

	var id [4]byte
	copy(id[:], info[:4])

	e := self.errors[id]
	if e == nil {
		return "", fmt.Errorf("can not parse error with sig: %x", id)
	}

	errInterface, err := abi.Decode(e.Inputs, info[:])
	if err != nil {
		return "", fmt.Errorf("can not parse err when decode: %s", err)
	}
	Innertemplate := "%v"
	inner := ""
	for i, v := range e.Inputs.TupleElems() {
		inner += fmt.Sprintf(Innertemplate, errInterface.(map[string]interface{})[abi.NameToKey(v.Name, i)])
		if i+1 != len(e.Inputs.TupleElems()) {
			inner += ","
		}
	}
	return e.Name + "(" + inner + ")", nil
}
