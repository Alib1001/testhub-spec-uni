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
//
// Code generated from specification version 8.14.0: DO NOT EDIT

package esapi

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

func newIndicesGetDataLifecycleFunc(t Transport) IndicesGetDataLifecycle {
	return func(name []string, o ...func(*IndicesGetDataLifecycleRequest)) (*Response, error) {
		var r = IndicesGetDataLifecycleRequest{Name: name}
		for _, f := range o {
			f(&r)
		}

		if transport, ok := t.(Instrumented); ok {
			r.instrument = transport.InstrumentationEnabled()
		}

		return r.Do(r.ctx, t)
	}
}

// ----- API Definition -------------------------------------------------------

// IndicesGetDataLifecycle returns the data stream lifecycle of the selected data streams.
//
// This API is experimental.
//
// See full documentation at https://www.elastic.co/guide/en/elasticsearch/reference/master/data-streams-get-lifecycle.html.
type IndicesGetDataLifecycle func(name []string, o ...func(*IndicesGetDataLifecycleRequest)) (*Response, error)

// IndicesGetDataLifecycleRequest configures the Indices Get Data Lifecycle API request.
type IndicesGetDataLifecycleRequest struct {
	Name []string

	ExpandWildcards string
	IncludeDefaults *bool

	Pretty     bool
	Human      bool
	ErrorTrace bool
	FilterPath []string

	Header http.Header

	ctx context.Context

	instrument Instrumentation
}

// Do executes the request and returns response or error.
func (r IndicesGetDataLifecycleRequest) Do(providedCtx context.Context, transport Transport) (*Response, error) {
	var (
		method string
		path   strings.Builder
		params map[string]string
		ctx    context.Context
	)

	if instrument, ok := r.instrument.(Instrumentation); ok {
		ctx = instrument.Start(providedCtx, "indices.get_data_lifecycle")
		defer instrument.Close(ctx)
	}
	if ctx == nil {
		ctx = providedCtx
	}

	method = "GET"

	if len(r.Name) == 0 {
		return nil, errors.New("name is required and cannot be nil or empty")
	}

	path.Grow(7 + 1 + len("_data_stream") + 1 + len(strings.Join(r.Name, ",")) + 1 + len("_lifecycle"))
	path.WriteString("http://")
	path.WriteString("/")
	path.WriteString("_data_stream")
	path.WriteString("/")
	path.WriteString(strings.Join(r.Name, ","))
	if instrument, ok := r.instrument.(Instrumentation); ok {
		instrument.RecordPathPart(ctx, "name", strings.Join(r.Name, ","))
	}
	path.WriteString("/")
	path.WriteString("_lifecycle")

	params = make(map[string]string)

	if r.ExpandWildcards != "" {
		params["expand_wildcards"] = r.ExpandWildcards
	}

	if r.IncludeDefaults != nil {
		params["include_defaults"] = strconv.FormatBool(*r.IncludeDefaults)
	}

	if r.Pretty {
		params["pretty"] = "true"
	}

	if r.Human {
		params["human"] = "true"
	}

	if r.ErrorTrace {
		params["error_trace"] = "true"
	}

	if len(r.FilterPath) > 0 {
		params["filter_path"] = strings.Join(r.FilterPath, ",")
	}

	req, err := newRequest(method, path.String(), nil)
	if err != nil {
		if instrument, ok := r.instrument.(Instrumentation); ok {
			instrument.RecordError(ctx, err)
		}
		return nil, err
	}

	if len(params) > 0 {
		q := req.URL.Query()
		for k, v := range params {
			q.Set(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}

	if len(r.Header) > 0 {
		if len(req.Header) == 0 {
			req.Header = r.Header
		} else {
			for k, vv := range r.Header {
				for _, v := range vv {
					req.Header.Add(k, v)
				}
			}
		}
	}

	if ctx != nil {
		req = req.WithContext(ctx)
	}

	if instrument, ok := r.instrument.(Instrumentation); ok {
		instrument.BeforeRequest(req, "indices.get_data_lifecycle")
	}
	res, err := transport.Perform(req)
	if instrument, ok := r.instrument.(Instrumentation); ok {
		instrument.AfterRequest(req, "elasticsearch", "indices.get_data_lifecycle")
	}
	if err != nil {
		if instrument, ok := r.instrument.(Instrumentation); ok {
			instrument.RecordError(ctx, err)
		}
		return nil, err
	}

	response := Response{
		StatusCode: res.StatusCode,
		Body:       res.Body,
		Header:     res.Header,
	}

	return &response, nil
}

// WithContext sets the request context.
func (f IndicesGetDataLifecycle) WithContext(v context.Context) func(*IndicesGetDataLifecycleRequest) {
	return func(r *IndicesGetDataLifecycleRequest) {
		r.ctx = v
	}
}

// WithExpandWildcards - whether wildcard expressions should get expanded to open or closed indices (default: open).
func (f IndicesGetDataLifecycle) WithExpandWildcards(v string) func(*IndicesGetDataLifecycleRequest) {
	return func(r *IndicesGetDataLifecycleRequest) {
		r.ExpandWildcards = v
	}
}

// WithIncludeDefaults - return all relevant default configurations for the data stream (default: false).
func (f IndicesGetDataLifecycle) WithIncludeDefaults(v bool) func(*IndicesGetDataLifecycleRequest) {
	return func(r *IndicesGetDataLifecycleRequest) {
		r.IncludeDefaults = &v
	}
}

// WithPretty makes the response body pretty-printed.
func (f IndicesGetDataLifecycle) WithPretty() func(*IndicesGetDataLifecycleRequest) {
	return func(r *IndicesGetDataLifecycleRequest) {
		r.Pretty = true
	}
}

// WithHuman makes statistical values human-readable.
func (f IndicesGetDataLifecycle) WithHuman() func(*IndicesGetDataLifecycleRequest) {
	return func(r *IndicesGetDataLifecycleRequest) {
		r.Human = true
	}
}

// WithErrorTrace includes the stack trace for errors in the response body.
func (f IndicesGetDataLifecycle) WithErrorTrace() func(*IndicesGetDataLifecycleRequest) {
	return func(r *IndicesGetDataLifecycleRequest) {
		r.ErrorTrace = true
	}
}

// WithFilterPath filters the properties of the response body.
func (f IndicesGetDataLifecycle) WithFilterPath(v ...string) func(*IndicesGetDataLifecycleRequest) {
	return func(r *IndicesGetDataLifecycleRequest) {
		r.FilterPath = v
	}
}

// WithHeader adds the headers to the HTTP request.
func (f IndicesGetDataLifecycle) WithHeader(h map[string]string) func(*IndicesGetDataLifecycleRequest) {
	return func(r *IndicesGetDataLifecycleRequest) {
		if r.Header == nil {
			r.Header = make(http.Header)
		}
		for k, v := range h {
			r.Header.Add(k, v)
		}
	}
}

// WithOpaqueID adds the X-Opaque-Id header to the HTTP request.
func (f IndicesGetDataLifecycle) WithOpaqueID(s string) func(*IndicesGetDataLifecycleRequest) {
	return func(r *IndicesGetDataLifecycleRequest) {
		if r.Header == nil {
			r.Header = make(http.Header)
		}
		r.Header.Set("X-Opaque-Id", s)
	}
}
