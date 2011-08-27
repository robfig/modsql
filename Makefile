# Copyright 2010  The "SQLModel" Authors
#
# Use of this source code is governed by the BSD-2 Clause license
# that can be found in the LICENSE file.
#
# This software is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES
# OR CONDITIONS OF ANY KIND, either express or implied. See the License
# for more details.

include $(GOROOT)/src/Make.inc

TARG=github.com/kless/sqlmodel
GOFILES=\
	column.go\
	main.go\
	metadata.go\
	table.go\
	type.go\

include $(GOROOT)/src/Make.pkg

