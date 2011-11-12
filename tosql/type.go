// Copyright 2010  The "GotoSQL" Authors
//
// Use of this source code is governed by the BSD 2-Clause License
// that can be found in the LICENSE file.
//
// This software is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES
// OR CONDITIONS OF ANY KIND, either express or implied. See the License
// for more details.

package tosql

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
