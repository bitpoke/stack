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
	"context"
	"net/http"
	"os"
	"strings"

	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/driver/github"
	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"

	"github.com/presslabs/stack/git-webhook/pkg/webhook/git"

	wordpressv1alpha1 "github.com/presslabs/wordpress-operator/pkg/apis/wordpress/v1alpha1"
)

type Server struct {
	client.Client
	client.FieldIndexer
	http.Server
	Mux *http.ServeMux
	Log logr.Logger
}

// Add creates a new Git Webhook Server and adds it to the Manager with default RBAC.
// The Manager will set fields on the Server and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	s, err := NewServer(mgr, ":8080")
	if err != nil {
		return err
	}
	return mgr.Add(s)
}

func NewServer(mgr manager.Manager, addr string) (*Server, error) {
	s := &Server{
		Client:       mgr.GetClient(),
		FieldIndexer: mgr.GetFieldIndexer(),
		Server: http.Server{
			Addr: addr,
		},
		Log: logf.Log.WithName("webhook-server"),
	}
	s.Mux = http.NewServeMux()
	s.Server.Handler = s.Mux

	s.Mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	s.Mux.HandleFunc("/github", s.githubWebhook)

	err := s.IndexField(
		&wordpressv1alpha1.Wordpress{}, "git.followed_ref", func(in runtime.Object) []string {
			wp := in.(*wordpressv1alpha1.Wordpress)

			if wp.Spec.CodeVolumeSpec == nil || wp.Spec.CodeVolumeSpec.GitDir == nil {
				return []string{}
			}

			if len(wp.Annotations["stack.presslabs.org/git-follow-name"]) == 0 {
				return []string{}
			}

			index, err := git.GitRepoFollowedRef(
				wp.Spec.CodeVolumeSpec.GitDir.Repository,
				wp.Annotations["stack.presslabs.org/git-follow-name"],
			)

			if err != nil {
				return []string{}
			}

			return []string{index}
		},
	)
	if err != nil {
		return s, err
	}

	return s, nil
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
	ghClient := github.NewDefault()

	secret := func(webhook scm.Webhook) (secret string, err error) {
		return os.Getenv("WEBHOOK_SECRET"), nil
	}

	webhook, err := ghClient.Webhooks.Parse(r, secret)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	switch event := webhook.(type) {
	case *scm.PushHook:
		if strings.HasPrefix(event.Ref, "refs/heads/") {
			repo, branch, ref := event.Repo.Clone, strings.TrimPrefix(event.Ref, "refs/heads/"), event.After
			s.Log.Info("received push event",
				"repo", repo,
				"branch", branch,
				"ref", ref,
			)
			// nolint errcheck
			go s.updateRef(repo, branch, ref)
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

func (s *Server) updateRef(repo, branch, ref string) error {
	wpList := wordpressv1alpha1.WordpressList{}

	cachedRef, err := git.GitRepoFollowedRef(repo, branch)
	if err != nil {
		s.Log.Error(err, "couldn't obtain git followed ref cache index")
		return err
	}

	listsOptions := &client.ListOptions{FieldSelector: fields.OneTermEqualSelector("git.followed_ref", cachedRef)}
	err = s.Client.List(context.TODO(), listsOptions, &wpList)
	if err != nil {
		s.Log.Error(err, "couldn't fetch wordpress resources")
		return err
	}

	if len(wpList.Items) == 0 {
		s.Log.Info("no wordpresses matched the given ref", "ref", cachedRef)
	}
	for _, wp := range wpList.Items {
		if wp.Spec.CodeVolumeSpec == nil || wp.Spec.CodeVolumeSpec.GitDir == nil {
			s.Log.Info("wp.Spec.CodeVolumeSpec[.GitDir] is nil")
			continue
		}

		followedRef, err := git.GitRepoFollowedRef(
			wp.Spec.CodeVolumeSpec.GitDir.Repository,
			wp.Annotations["stack.presslabs.org/git-follow-name"],
		)
		if err != nil {
			s.Log.Error(err, "couldn't obtain git followed ref index")
			continue
		}

		if followedRef != cachedRef {
			s.Log.Info("cachedRef not matching followedRef")
			continue
		}

		wp.Spec.CodeVolumeSpec.GitDir.GitRef = ref
		err = s.Client.Update(context.TODO(), &wp)
		if err != nil {
			s.Log.Error(err, "couldn't update wordpress git ref")
			continue
		}
	}
	return nil
}
