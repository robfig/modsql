// Copyright 2010  The "SQLModel" Authors
//
// Use of this source code is governed by the Simplified BSD License
// that can be found in the LICENSE file.
//
// This software is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES
// OR CONDITIONS OF ANY KIND, either express or implied. See the License
// for more details.

package sqlmodel


type output int

// Output where to write
const (
	FILEOUT output = iota // To file
	STDOUT                // To standard output

	_SQL_OUTPUT   = "model.sql"
	_MODEL_OUTPUT = "model.go" // Go definitions related to each SQL table
)

var errors bool

