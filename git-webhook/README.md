git-webhook
===

It's a simple HTTP server that can be feed with github events and updates
WordPress sites references. It updates all WordPress resources with proper
annotations.

## Annotations

| Annotation                             | Description                                                 |
| ---                                    | ---                                                         |
| `stack.presslabs.org/git-follow-name`  | The git reference to follow. Usually a branch name          |

## Quirks
* When setting up a GitHub webhook you must choose to receive JSON data payloads, otherwise the git-webhook server won't know how to decode the message.
