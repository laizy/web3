/*
 * Copyright (C) 2018 The ontology Authors
 * This file is part of The ontology library.
 *
 * The ontology is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The ontology is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with The ontology.  If not, see <http://www.gnu.org/licenses/>.
 */
package schema

import (
	"errors"
	"github.com/ontio/ontology/core/store/common"
)

var ErrNotFound = errors.New("not found")

// Store iterator for iterate store
type StoreIterator interface {
	Next() bool //Next item. If item available return true, otherwise return false
	//Prev() bool           //previous item. If item available return true, otherwise return false
	First() bool //First item. If item available return true, otherwise return false
	//Last() bool           //Last item. If item available return true, otherwise return false
	//Seek(key []byte) bool //Seek key. If item available return true, otherwise return false
	Key() []byte   //Return the current item key
	Value() []byte //Return the current item value
	Release()      //Close iterator
	Error() error  // Error returns any accumulated error.
}

// PersistStore of ledger
type PersistStore interface {
	Get(key []byte) ([]byte, error)
	BatchPut(key []byte, value []byte)              //Put a key-value pair to batch
	BatchDelete(key []byte)                         //Delete the key in batch
	NewIterator(prefix []byte) common.StoreIterator //Return the iterator of store
}
