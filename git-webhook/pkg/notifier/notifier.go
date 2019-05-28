/*
Copyright 2019 Pressinfra SRL.

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

package notifier

import (
	"sync"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

var defaultNotifier *Notifier
var once sync.Once

func DefaultNotifier() *Notifier {
	once.Do(func() {
		defaultNotifier = NewNotifier()
	})

	return defaultNotifier
}

type repoBranch struct {
	repo   string
	branch string
}

type Notifier struct {
	mu           sync.Mutex
	repoRegistry map[repoBranch]string
	registry     map[repoBranch]map[client.ObjectKey]struct{}
	revRegistry  map[client.ObjectKey]repoBranch
}

func NewNotifier() *Notifier {
	return &Notifier{
		repoRegistry: make(map[repoBranch]string),
		registry:     make(map[repoBranch]map[client.ObjectKey]struct{}),
		revRegistry:  make(map[client.ObjectKey]repoBranch),
	}
}

func (n *Notifier) RegisterSite(key client.ObjectKey, repo, branch string) {
	n.mu.Lock()
	defer n.mu.Unlock()

	oldRepoBranch, found := n.revRegistry[key]
	rb := repoBranch{repo: repo, branch: branch}

	if n.registry[rb] == nil {
		n.registry[rb] = make(map[client.ObjectKey]struct{})
	}
	if !found {
		n.registry[rb][key] = struct{}{}
		n.revRegistry[key] = rb
	} else if oldRepoBranch != rb {
		n.revRegistry[key] = rb
	}
}

func (n *Notifier) UnregisterSite(key client.ObjectKey) {
	n.mu.Lock()
	defer n.mu.Unlock()

	rb, found := n.revRegistry[key]
	if !found {
		return
	}

	delete(n.revRegistry, key)
	delete(n.registry[rb], key)
}

func (n *Notifier) TriggerUpdate(repo, branch, ref string) {
	n.mu.Lock()
	defer n.mu.Unlock()

	rb := repoBranch{repo: repo, branch: branch}
	n.repoRegistry[rb] = ref
}

func (n *Notifier) GetSiteRef(key client.ObjectKey) (string, bool) {
	n.mu.Lock()
	defer n.mu.Unlock()

	rb, found := n.revRegistry[key]
	if !found {
		return "", false
	}

	ref, found := n.repoRegistry[rb]
	return ref, found
}
