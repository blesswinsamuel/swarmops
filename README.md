# Docker Swarm Operator

> GitOps operator for Docker Swarm.

## Getting started

### Start server

```bash
swarmops serve \
    --git-repo git@github.com:blesswinsamuel/swarmops-stack-example.git \
    --git-branch master \
    --private-key-file ~/.ssh/id_rsa \
    --repo-dir /tmp/swarm-operator-repo \
    --stack-file stack.dev.yaml \
    --sync-interval 1m \
    --port 8080
```

### Deploy (one-off command)

```bash
swarmops deploy --stack-file stack.dev.yaml
```

### Sync swarm

```bash
curl localhost:8080/api/sync?force=true
```
