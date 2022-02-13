/*
Copyright 2022 Daisuke Taniwaki..

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
	"fmt"
	"sync"

	"github.com/robfig/cron/v3"
	"k8s.io/apimachinery/pkg/types"
)

type LendingEvent string

const (
	LendingStart LendingEvent = "LendingStart"
	LendingEnd   LendingEvent = "LendingEnd"
)

type CronItem struct {
	Cron string
	Job  cron.Job
}

type Cron struct {
	cron    *cron.Cron
	entries map[string][]cron.EntryID
	lock    sync.RWMutex
}

func NewCron() *Cron {
	return &Cron{
		cron: cron.New(),
		lock: sync.RWMutex{},
	}
}

func (c *Cron) Add(namespacedName types.NamespacedName, cronItems []CronItem) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	resourceName := getResourceEntryName(namespacedName)
	entries, ok := c.entries[resourceName]
	if !ok {
		entries = []cron.EntryID{}
	}
	for _, item := range cronItems {
		entryId, err := c.cron.AddJob(item.Cron, item.Job)
		if err != nil {
			return err
		}
		entries = append(entries, entryId)
	}
	c.entries[resourceName] = entries
	return nil
}

func (c *Cron) Clear(namespacedName types.NamespacedName) {
	c.lock.Lock()
	defer c.lock.Unlock()

	resourceName := getResourceEntryName(namespacedName)
	entries, ok := c.entries[resourceName]
	if ok && entries != nil {
		for entryId := range entries {
			c.cron.Remove(cron.EntryID(entryId))
		}
	}
	delete(c.entries, resourceName)
}

func (c *Cron) Start() {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.cron.Start()
}

func (c *Cron) Stop() {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.cron.Stop()
}

func getResourceEntryName(namespacedName types.NamespacedName) string {
	return fmt.Sprintf("%s/%s", namespacedName.Namespace, namespacedName.Name)
}
