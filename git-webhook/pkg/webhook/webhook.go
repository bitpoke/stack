/*
Copyright 2019 Pressinfa SRL

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

package webhook

import (
	"net/http"
	"strings"

	"sigs.k8s.io/controller-runtime/pkg/manager"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"

	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/driver/github"
	"github.com/go-logr/logr"

	"github.com/presslabs/stack/git-webhook/pkg/notifier"
)

type Server struct {
	http.Server
	Mux *http.ServeMux
	Log logr.Logger
}

// Add creates a new Git Webhook Controller and adds it to the Manager with default RBAC.
// The Manager will set fields on the Controller and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	s := NewServer(":8080")
	return mgr.Add(s)
}

func NewServer(addr string) *Server {
	s := &Server{
		Server: http.Server{
			Addr: addr,
		},
		Log: logf.Log.WithName("webhook-server"),
	}
	s.Mux = http.NewServeMux()
	s.Server.Handler = s.Mux

	s.Mux.HandleFunc("/github", s.githubWebhook)

	return s
}

func (s *Server) Start(stop <-chan struct{}) error {
	s.Log.Info("webhook server is listening", "address", s.Addr)

	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}

	<-stop
	return nil
}

func (s *Server) githubWebhook(w http.ResponseWriter, r *http.Request) {
	client := github.NewDefault()
	secret := func(webhook scm.Webhook) (secret string, err error) {
		return "", nil
	}

	webhook, err := client.Webhooks.Parse(r, secret)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	switch event := webhook.(type) {
	case *scm.PushHook:
		if strings.HasPrefix(event.Ref, "refs/heads/") {
			repo, branch, ref := event.Repo.Clone, strings.TrimPrefix(event.Ref, "refs/heads/"), event.Commit.Sha
			notifier.DefaultNotifier().TriggerUpdate(repo, branch, ref)
			s.Log.Info("received push event",
				"repo", repo,
				"branch", branch,
				"ref", ref,
			)
		}
	case *scm.TagHook:
	case *scm.BranchHook:
	case *scm.IssueHook:
	case *scm.IssueCommentHook:
	case *scm.PullRequestHook:
	case *scm.PullRequestCommentHook:
	case *scm.ReviewCommentHook:
	}
}
