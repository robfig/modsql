// Copyright 2010  The "GoSQL" Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gosql

type sqlType uint8

const (
	_ sqlType = iota
	Integer
	Float
	Text
	Blob
	Boolean
)

var (
	sqlType_str = map[sqlType]string{
		Integer: "Integer",
		Float:   "Float",
		Text:    "Text",
		Blob:    "Blob",
		Boolean: "Boolean",
	}

	sqlType_goType = map[sqlType]string{
		Integer: "int",
		Float:   "float32",
		Text:    "string",
		Blob:    "[]byte",
		Boolean: "bool",
	}
)

func (t sqlType) String() string {
	return sqlType_str[t]
}

func (t sqlType) Go() string {
	return sqlType_goType[t]
}
