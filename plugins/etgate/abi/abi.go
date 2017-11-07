// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package abi

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/core/types"
)

// The ABI holds information about a contract's context and available
// invokable methods. It will allow you to type check function calls and
// packs data accordingly.
type ABI struct {
	Constructor Method
	Methods     map[string]Method
	Events      map[string]Event
}

// JSON returns a parsed ABI interface and error if it failed.
func JSON(reader io.Reader) (ABI, error) {
	dec := json.NewDecoder(reader)

	var abi ABI
	if err := dec.Decode(&abi); err != nil {
		return ABI{}, err
	}

	return abi, nil
}

// Pack the given method name to conform the ABI. Method call's data
// will consist of method_id, args0, arg1, ... argN. Method id consists
// of 4 bytes and arguments are all 32 bytes.
// Method ids are created from the first 4 bytes of the hash of the
// methods string signature. (signature = baz(uint32,string32))
func (abi ABI) Pack(name string, args ...interface{}) ([]byte, error) {
	// Fetch the ABI of the requested method
	var method Method

	if name == "" {
		method = abi.Constructor
	} else {
		m, exist := abi.Methods[name]
		if !exist {
			return nil, fmt.Errorf("method '%s' not found", name)
		}
		method = m
	}
	arguments, err := method.pack(args...)
	if err != nil {
		return nil, err
	}
	// Pack up the method ID too if not a constructor and return
	if name == "" {
		return arguments, nil
	}
	return append(method.Id(), arguments...), nil
}

// these variable are used to determine certain types during type assertion for
// assignment.
var (
	r_interSlice = reflect.TypeOf([]interface{}{})
	r_hash       = reflect.TypeOf(common.Hash{})
	r_bytes      = reflect.TypeOf([]byte{})
	r_byte       = reflect.TypeOf(byte(0))
)

func (abi ABI) Unpack(v interface{}, name string, log types.Log) error {
    var event = abi.Events[name]

    if len(log.Data) == 0 {
        return fmt.Errorf("abi: unmarshalling empty input")
    }

    valueOf := reflect.ValueOf(v)
    if reflect.Ptr != valueOf.Kind() {
        return fmt.Errorf("abi: Unpack(non-pointer %T)", v)
    }

    var (
        value = valueOf.Elem()
        typ = value.Type()
    )

    if reflect.Struct != value.Kind() {
        return fmt.Errorf("abi: Unpack(non-struct %T)", value)
    }

    var inter interface{}
    var iIndexed int
    if event.Anonymous {
        iIndexed = 0
    } else {
        iIndexed = 1
    }
    iNormal := 0

    for _, input := range event.Inputs {
        var reflectValue reflect.Value

        if input.Indexed {
            switch input.Type.Elem.T {
            case IntTy, UintTy:
                inter = readInteger(input.Type.Kind, log.Topics[iIndexed][:])
            case BoolTy:
                var err error
                inter, err = readBool(log.Topics[iIndexed][:])
                if err != nil {
                    return err
                }
            default:
                return fmt.Errorf("abi: unsupported indexed argument type")
            }
            reflectValue = reflect.ValueOf(inter)
            iIndexed++
        } else {
            marshalledValue, err := toGoType(iNormal, input, log.Data)
            if err != nil {
                return err
            }
            reflectValue = reflect.ValueOf(marshalledValue)
            iNormal++
        }

        var j int
        for j = 0; j < typ.NumField(); j++ {
            field := typ.Field(j)
            if field.Name == strings.ToUpper(input.Name[:1])+input.Name[1:] {
                if err := set(value.Field(j), reflectValue, input); err != nil {
                    return err
                }
                break
            }
        }
        if j == typ.NumField() {
            return fmt.Errorf("abi: Unpack(field not found %s)", input.Name)
        }
    }
    return nil
}

func (abi *ABI) UnmarshalJSON(data []byte) error {
	var fields []struct {
		Type      string
		Name      string
		Constant  bool
		Indexed   bool
		Anonymous bool
		Inputs    []Argument
		Outputs   []Argument
	}

	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}

	abi.Methods = make(map[string]Method)
	abi.Events = make(map[string]Event)
	for _, field := range fields {
		switch field.Type {
		case "constructor":
			abi.Constructor = Method{
				Inputs: field.Inputs,
			}
		// empty defaults to function according to the abi spec
		case "function", "":
			abi.Methods[field.Name] = Method{
				Name:    field.Name,
				Const:   field.Constant,
				Inputs:  field.Inputs,
				Outputs: field.Outputs,
			}
		case "event":
			abi.Events[field.Name] = Event{
				Name:      field.Name,
				Anonymous: field.Anonymous,
				Inputs:    field.Inputs,
			}
		}
	}

	return nil
}
