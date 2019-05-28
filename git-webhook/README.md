git-webhook
===

It's a simple HTTP server that can be feed with github events and updates
WordPress sites references. It updates all WordPress resources with proper
annotations.

## Annotations

| Annotation                                    | Description                                                 |
| ---                                           | ---                                                         |
| `webhook.stack.presslabs.org/auto-update-ref` | Whenever to auto-update the git reference. Default `false`. |
| `webhook.stack.presslabs.org/follow-git-ref`  | The git reference to follow. Usually a branch name          |
