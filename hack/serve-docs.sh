#!/usr/bin/env bash

# Copyright 2019 Pressinfra SRL
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -o errexit
set -o nounset
set -o pipefail

ROOT_DIR="$(dirname "${BASH_SOURCE[0]}")/.."

cd "${ROOT_DIR}"

PWD=$(pwd)

set -x

docker run --rm -p 1313:1313 \
    -v "$PWD/docs:/site/content" \
    -v "$PWD/public:/site/public" \
    -v "$PWD/.git:/site/.git" \
    gcr.io/pl-infra/hugo:latest \
        --config=/site/content/config.toml \
        --bind 0.0.0.0 -b http://localhost:1313/ \
        --enableGitInfo server
