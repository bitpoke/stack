/*
Copyright 2018 Pressinfra SRL.

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
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	wordpressv1alpha1 "github.com/presslabs/wordpress-operator/pkg/apis/wordpress/v1alpha1"
)

var _ = Describe("Git repo Webhook", func() {
	var (
		// stop channel for controller manager
		stop chan struct{}
		// webhook server
		s *Server

		c client.Client
	)

	BeforeEach(func() {
		mgr, err := manager.New(cfg, manager.Options{})
		Expect(err).NotTo(HaveOccurred())

		Expect(err).To(Succeed())

		s, err = NewServer(mgr, ":0")
		Expect(err).To(BeNil())
		Expect(mgr.Add(s)).To(Succeed())

		// create new k8s client
		c, err = client.New(cfg, client.Options{})
		Expect(err).To(Succeed())

		stop = StartTestManager(mgr)
	})

	AfterEach(func() {
		close(stop)
	})

	When("receiving a GitHub push event", func() {
		var (
			wp                                           *wordpressv1alpha1.Wordpress
			wpName                                       string
			repo_name, repo_org                          string
			repo_org_name, repo_html_url, repo_clone_url string
			branch                                       string
			commitSha1                                   string
			webhookBody                                  []byte
			recorder                                     *httptest.ResponseRecorder
			req                                          *http.Request
		)

		BeforeEach(func() {
			wpName = fmt.Sprintf("wp-%d", rand.Int31())
			repo_name = fmt.Sprintf("repo-%d", rand.Int31())
			repo_org = "presslabs"
			repo_org_name = repo_org + "/" + repo_name
			repo_html_url = "https://github.com/" + repo_org_name
			repo_clone_url = repo_html_url + ".git"

			branch = fmt.Sprintf("branch-%d", rand.Int31())

			commitSha1 = "409e57f658ca22628468a509b86c5d41869509df"

			replicas := int32(1)
			wp = &wordpressv1alpha1.Wordpress{
				ObjectMeta: metav1.ObjectMeta{
					Name:      wpName,
					Namespace: "default",
					Annotations: map[string]string{
						"stack.presslabs.org/git-follow-name": branch,
					},
				},
				Spec: wordpressv1alpha1.WordpressSpec{
					Domains:  []wordpressv1alpha1.Domain{wordpressv1alpha1.Domain("test.com")},
					Image:    "wordpress",
					Replicas: &replicas,
					CodeVolumeSpec: &wordpressv1alpha1.CodeVolumeSpec{
						GitDir: &wordpressv1alpha1.GitVolumeSource{
							GitRef:     "7515154901d55a1",
							Repository: repo_clone_url,
						},
					},
				},
			}
			Expect(c.Create(context.TODO(), wp)).To(Succeed())
			wpKey := client.ObjectKey{Name: wp.Name, Namespace: wp.Namespace}
			// Wait for the resource to be cached
			Eventually(func() error {
				return s.Get(context.TODO(), wpKey, wp)
			}).Should(Succeed())

			var err error
			webhookBody, err = json.Marshal(map[string]interface{}{
				"ref":      "refs/heads/" + branch,
				"before":   "bd5bb8a087b0db65e18b9ec4c92d0df1fcb13345",
				"after":    commitSha1,
				"created":  false,
				"deleted":  false,
				"forced":   false,
				"base_ref": nil,
				"commits":  []map[string]interface{}{},
				"repository": map[string]interface{}{
					"id":        24887812,
					"name":      repo_name,
					"full_name": repo_org + "/" + repo_name,
					"private":   false,
					"owner": map[string]interface{}{
						"name":  repo_org,
						"email": "ping@presslabs.com",
						"login": repo_org,
						"id":    25032711,
					},
					"html_url":     repo_html_url,
					"fork":         false,
					"url":          repo_html_url,
					"created_at":   1412681978,
					"updated_at":   "2019-05-29T10:04:48Z",
					"pushed_at":    1559571474,
					"git_url":      "git://github.com/" + repo_org_name + ".git",
					"ssh_url":      "git@github.com:" + repo_org_name + ".git",
					"clone_url":    repo_clone_url,
					"homepage":     "https://www.presslabs.com/stack/",
					"organization": repo_org,
				},
			})
			Expect(err).To(BeNil())

			recorder = httptest.NewRecorder()

			req, err = http.NewRequest("POST", "/github", bytes.NewBuffer(webhookBody))
			Expect(err).To(BeNil())

			req.Header.Add("User-Agent", "GitHub-Hookshot/0a2cefb")
			req.Header.Add("X-GitHub-Delivery", "60f280a6-860a-11e9-91d1-e5b9ac6e1d3e")
			req.Header.Add("X-GitHub-Event", "push")
		})

		AfterEach(func() {
			Expect(c.Delete(context.TODO(), wp)).To(Succeed())
		})
		It("updates the ref of all the sites matching the repo/branch", func() {
			s.Handler.ServeHTTP(recorder, req)

			response := recorder.Result()

			responseBody, err := ioutil.ReadAll(response.Body)
			Expect(err).To(BeNil())

			Expect(response.StatusCode).To(Equal(http.StatusOK), string(responseBody))
			Expect(string(responseBody)).To(BeEmpty())

			Eventually(func() error {
				wpKey := client.ObjectKey{Name: wp.Name, Namespace: wp.Namespace}
				err := c.Get(context.TODO(), wpKey, wp)
				if err != nil {
					return err
				}
				if wp.Spec.CodeVolumeSpec.GitDir.GitRef != commitSha1 {
					return errors.New("commit sha not matching")
				}
				return nil
			}).Should(Succeed())
		})
		When("using signed payload", func() {
			BeforeEach(func() {
				// Add webhook signature

				secret := "asdasd123123"
				Expect(os.Setenv("WEBHOOK_SECRET", secret)).To(Succeed())

				sig := hmac.New(sha1.New, []byte(secret))
				_, err := sig.Write([]byte(webhookBody))
				Expect(err).To(BeNil())
				webhookBodySha1 := hex.EncodeToString(sig.Sum(nil))

				req.Header.Add("X-Hub-Signature", "sha1="+webhookBodySha1)
			})
			AfterEach(func() {
				Expect(os.Setenv("WEBHOOK_SECRET", ""))
			})
			It("updates the ref of all the sites matching the repo/branch", func() {
				s.Handler.ServeHTTP(recorder, req)

				response := recorder.Result()

				responseBody, err := ioutil.ReadAll(response.Body)
				Expect(err).To(BeNil())

				Expect(response.StatusCode).To(Equal(http.StatusOK), string(responseBody))
				Expect(string(responseBody)).To(BeEmpty())

				Eventually(func() error {
					wpKey := client.ObjectKey{Name: wp.Name, Namespace: wp.Namespace}
					err := c.Get(context.TODO(), wpKey, wp)
					if err != nil {
						return err
					}
					if wp.Spec.CodeVolumeSpec.GitDir.GitRef != commitSha1 {
						return errors.New("commit sha not matching")
					}
					return nil
				}).Should(Succeed())
			})
		})
	})
})
