/*
Copyright 2022 Daisuke Taniwaki.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHPA(t *testing.T) {
	hour, minute, err := parseHours("07:20")
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	if !assert.Equal(t, hour, int32(7)) {
		t.FailNow()
	}
	if !assert.Equal(t, minute, int32(20)) {
		t.FailNow()
	}

	hour, minute, err = parseHours("07:20am")
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	if !assert.Equal(t, hour, int32(7)) {
		t.FailNow()
	}
	if !assert.Equal(t, minute, int32(20)) {
		t.FailNow()
	}

	hour, minute, err = parseHours("07:20pm")
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	if !assert.Equal(t, hour, int32(19)) {
		t.FailNow()
	}
	if !assert.Equal(t, minute, int32(20)) {
		t.FailNow()
	}
}
