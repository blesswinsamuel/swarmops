# Docker Swarm Operator

> GitOps operator for Docker Swarm.

## Getting started

```bash
docker_swarm_gitops \
    --git-repo git@github.com:blesswinsamuel/docker-swarm-gitops-stack-example.git \
    --git-branch master \
    --private-key-file ~/.ssh/id_rsa \
    --repo-dir /tmp/swarm-operator-repo \
    --stack-file stack.dev.yaml \
    --sync-interval 1m \
    --port 8080
```

## Sync swarm

(not implemented)

```bash
swarm_operator sync
```
