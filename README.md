# Docker Swarm Operator

> GitOps operator for Docker Swarm.

## Getting started

```bash
swarm_operator start --git-repo git@github.com:blesswinsamuel/swarm-operator-example.git --keys-dir /tmp/swarm-operator-keys --repo-dir /tmp/swarm-operator-repo
```

At startup, Swarm Operator generates a SSH key and logs the public key. Find the SSH public key by running:

```bash
swarm_operator getkey
```

## Sync swarm

```bash
swarm_operator sync
```
