---
title: Contributing
linktitle: Contributing
description: "How to contribute to Stack"
categories: []
keywords: ['stack', 'docs', 'wordpress', 'kubernetes', 'contributing', 'development', 'issues']
draft: false
aliases: []
slug: contributing
toc: true
related: true
---

## Issues

Issues are being tracked at https://github.com/presslabs/stack/issues.  
They can range from bug reports to questions about development, installation process and other stuff related to the project.

## Development

> ###### NOTE
>
> Before making a change to this repository, please open an issue to discuss why the change is needed, and how it should be implemented if necessary.

1. Clone the repository:  
Optionally you could [fork the repo](https://github.com/presslabs/stack/fork) first.  
`git clone https://github.com/presslabs/stack.git && cd stack`

2. Install dependencies:  
`make dependencies`

3. Implement your changes. Remember to lint your code with `make lint`.

4. Deploy to your cluster with [skaffold](https://skaffold.dev/docs/getting-started/#installing-skaffold):  

    If you want to develop locally, use [minikube](https://github.com/kubernetes/minikube#installation) to get a locally running cluster. You can follow the steps described [here](install-stack-on-minikube.md), except for the ones explaining how to install the stack.

    If you want to have your changes constantly deployed run:  
`skaffold dev`

    If you want to manually deploy your changes run:  
`skaffold run`

5. Test your changes.  
    There's not much written code in stack, except for a few components (default-backend, git-webhook...) which might have their own Makefile and tests that can be run with `make test`.  
    So most of the times testing will mean just running the stack with your changes and poking around.

    If you need to run a site for your use-case, read the [related documentation](docs/running-wordpress-on-kubernetes.md).

6. Open a pull request at https://github.com/presslabs/stack/compare.
    We'll review the pull request and assist you on getting it merged.
