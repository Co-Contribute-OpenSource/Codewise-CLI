# Production readiness

Codewise is a client-side CLI. Install it on a workstation or CI runner; do not
deploy the CLI itself as a Kubernetes workload.

## Support status

Codewise is currently beta software. Use a disposable cluster first, review
every `deploy plan`, and keep the underlying Docker, Kubernetes, Helm, and Git
commands available for recovery.

The CI compatibility baseline is:

| Component | Validated version or range |
| --- | --- |
| Go | 1.25.x |
| Linux | GitHub-hosted Ubuntu runner |
| kind | 0.32.0 |
| Helm | 3.19.0 |
| Kubernetes | The default Kubernetes version bundled with kind 0.32.0 |
| Cosign | 3.0.x for release verification |
| Node.js | 22.x for building the documentation |

Docker must support BuildKit-compatible `docker build`. kubectl should remain
within one minor version of the target Kubernetes API server. macOS and Windows
binaries are built on release tags, while cluster end-to-end tests currently
run on Linux.

## Required tools by command

| Workflow | Required tools and access |
| --- | --- |
| Encoding and generation | Codewise only |
| `docker build` / `docker push` | Docker CLI and a reachable daemon; registry credentials for push |
| Raw Kubernetes deployment | kubectl, kubeconfig, reachable cluster, namespace permissions |
| Helm deployment | Helm, kubectl, kubeconfig, reachable cluster, release permissions |
| GitOps deployment | kubectl, Argo CD Application CRD, `argocd` namespace access, repository credentials configured in Argo CD |
| Release production | GoReleaser, Syft, Cosign, GitHub OIDC permissions |

## Select the strategy explicitly

Set `strategy` in the environment's `k8s.yaml` so operational commands do not
depend on which files happen to exist in the current directory:

```yaml
namespace: staging
context: my-cluster
strategy: helm # helm, kubectl, or gitops
```

Existing environments without this property retain automatic detection:
GitOps when a repository is configured, Helm when `helm/chart` exists, and
kubectl otherwise.

## Safe deployment sequence

```bash
kubectl config current-context
codewise deploy plan --env staging
codewise deploy run --env staging --dry-run
codewise deploy run --env staging
codewise deploy status --env staging
```

The dry run does not contact the cluster or create a namespace. It validates
and prints the command or GitOps manifest Codewise intends to apply.

## Rollback and recovery

Rollback is deliberately strategy-specific.

### Helm

```bash
codewise deploy history --env staging
codewise deploy rollback --env staging --revision 2
codewise deploy status --env staging
```

If Codewise is unavailable, use the equivalent Helm recovery path:

```bash
helm history RELEASE -n NAMESPACE
helm rollback RELEASE REVISION -n NAMESPACE
kubectl rollout status deployment/NAME -n NAMESPACE
```

### kubectl

Raw kubectl deployments do not have a Codewise revision history. Reapply a
known-good manifest or use Kubernetes rollout history for an individual
Deployment:

```bash
git checkout GOOD_COMMIT -- k8s/
kubectl apply -f k8s/ -n NAMESPACE
kubectl rollout undo deployment/NAME -n NAMESPACE
```

### GitOps

Revert the desired-state Git commit and let Argo CD reconcile it:

```bash
git revert BAD_COMMIT
git push
kubectl -n argocd get application ENVIRONMENT
```

Do not use `codewise deploy rollback` for GitOps. Codewise rejects that request
instead of mutating live state outside Git.

## Failure recovery

When a deployment fails:

1. Preserve the complete Codewise output. External tool errors are streamed.
2. Verify the current context and namespace.
3. Inspect workloads and events.
4. Recover with the strategy-specific procedure above.

```bash
kubectl config current-context
kubectl get pods -n NAMESPACE -o wide
kubectl get events -n NAMESPACE --sort-by=.lastTimestamp
kubectl describe pod POD -n NAMESPACE
codewise deploy logs --env staging
```

Codewise will not create a namespace when the namespace lookup fails because of
authentication, connectivity, or authorization. Only a confirmed NotFound
response triggers namespace creation.

## Upgrade Codewise

1. Read the release notes.
2. Download the archive for your operating system.
3. Verify its SHA-256 checksum.
4. Verify the Sigstore bundle.
5. Replace the binary and run smoke checks.

```bash
sha256sum --check checksums.txt
cosign verify-blob \
  --bundle checksums.txt.sigstore.json \
  --certificate-identity-regexp 'https://github.com/aryansharma9917/Codewise-CLI/' \
  --certificate-oidc-issuer https://token.actions.githubusercontent.com \
  checksums.txt

codewise version
codewise deploy plan --env staging
```

Environment YAML remains human-readable. Back it up before changing fields
across a major Codewise release.
