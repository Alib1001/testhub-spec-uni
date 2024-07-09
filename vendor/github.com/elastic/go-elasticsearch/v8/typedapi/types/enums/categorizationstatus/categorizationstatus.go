// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

// Code generated from the elasticsearch-specification DO NOT EDIT.
// https://github.com/elastic/elasticsearch-specification/tree/07bf82537a186562d8699685e3704ea338b268ef

// Package categorizationstatus
package categorizationstatus

import "strings"

// https://github.com/elastic/elasticsearch-specification/blob/07bf82537a186562d8699685e3704ea338b268ef/specification/ml/_types/Model.ts#L83-L86
type CategorizationStatus struct {
	Name string
}

var (
	Ok = CategorizationStatus{"ok"}

	Warn = CategorizationStatus{"warn"}
)

func (c CategorizationStatus) MarshalText() (text []byte, err error) {
	return []byte(c.String()), nil
}

func (c *CategorizationStatus) UnmarshalText(text []byte) error {
	switch strings.ReplaceAll(strings.ToLower(string(text)), "\"", "") {

	case "ok":
		*c = Ok
	case "warn":
		*c = Warn
	default:
		*c = CategorizationStatus{string(text)}
	}

	return nil
}

func (c CategorizationStatus) String() string {
	return c.Name
}
