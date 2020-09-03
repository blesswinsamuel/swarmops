# Docker Swarm Operator

> GitOps operator for Docker Swarm.

## Getting started

```bash
swarm_operator start --git-repo git@github.com:blesswinsamuel/swarm-operator-example.git --private-key-file ~/.ssh/id_rsa --repo-dir /tmp/swarm-operator-repo
```

At startup, Swarm Operator generates a SSH key and logs the public key. Find the SSH public key by running:

```bash
swarm_operator getkey
```

## Sync swarm

```bash
swarm_operator sync
```
