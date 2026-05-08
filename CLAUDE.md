# Cert Manager Operator for OpenShift

OpenShift operator that manages the lifecycle of cert-manager, deploying and configuring it via a cluster-scoped `CertManager` CR named `cluster`. Built with Go 1.23, operator-sdk, controller-runtime, Jsonnet, and Makefile.

## Structure

- `api/` - CertManager API types
- `pkg/controller/` - Operator controllers (cainjector, controller, webhook)
- `bindata/` - Upstream cert-manager CRDs and deployment manifests (generated)
- `bundle/` - OLM bundle manifests
- `hack/` - Dev and CI scripts
- `vendor/` - Go dependencies

## Key Files

- Operator entrypoint: `main.go`
- Core controller logic: `pkg/controller/deployment/cert_manager_controller_deployment.go`
- All dev commands: `Makefile`
- Linter config: `.golangci.yaml`

## Commands

```bash
make build              # build operator binary
make test               # unit tests
make test-e2e           # e2e tests (requires stable cluster — see below)
make verify             # lint + vet
make deploy             # install CRDs and manifests into connected cluster
make local-run          # run operator locally against cluster
make update             # regenerate bindata/ after cert-manager version bump
```

## Critical Context

**Namespaces are hardcoded**: operator runs in `cert-manager-operator`, operand in `cert-manager`. These cannot be changed.

**Local dev requires a live OpenShift cluster** with `oc` connected. Local run order:
1. `make deploy` — installs CRDs/manifests
2. `oc scale --replicas=0 deploy --all -n cert-manager-operator` — stop in-cluster operator
3. `make local-run`

**E2E tests**: operator and operand must be fully stable first. Run `make test-e2e-wait-for-stable-state` before `make test-e2e`.

**Upgrading cert-manager**: bump `CERT_MANAGER_VERSION` in `Makefile`, run `make update`, then review `bindata/` for inconsistencies.

**Non-x86 builds**: add `--platform linux/amd64` to container build commands.

**Module path**: `github.com/openshift/cert-manager-operator` (not the fork repo name).

## More Info

See [BOOKMARKS.md](BOOKMARKS.md) for architecture docs, metrics guide, OLM bundle process, and upstream references.
