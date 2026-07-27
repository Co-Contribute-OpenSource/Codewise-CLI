# Deploy commands

Deployment commands operate on profiles created with `codewise env create`.

## `deploy plan`

```bash
codewise deploy plan --env staging
```

Displays the selected strategy and command. GitOps plans include the generated Argo CD Application manifest.

## `deploy explain`

```bash
codewise deploy explain --env staging
```

Explains why the strategy was selected, required tools and access, the exact
external command, each execution phase, and the strategy-specific recovery
procedure. It does not contact the cluster.

## `deploy run`

```bash
codewise deploy run --env staging --dry-run
codewise deploy run --env staging
```

| Flag | Description |
| --- | --- |
| `--env` | Environment profile to deploy; required |
| `--dry-run` | Preview without cluster access or resource mutation |

## Observe a deployment

```bash
codewise deploy status --env staging
codewise deploy logs --env staging
codewise deploy logs --env staging --follow
codewise deploy history --env staging
```

## Roll back

```bash
codewise deploy rollback --env staging
codewise deploy rollback --env staging --revision 2
```

| Flag | Description |
| --- | --- |
| `--env` | Environment profile; required |
| `--revision` | Positive Helm revision to restore; required |

History and automatic rollback are supported only for Helm environments.
kubectl environments must reapply known-good manifests; GitOps environments
must revert desired state in Git.

## Preview environments

```bash
codewise deploy preview --pr 42
codewise deploy preview --image ghcr.io/example/myapp:pr-42
codewise deploy preview --pr 42 --keep
```

By default, Codewise removes the temporary environment profile after the deployment finishes. `--keep` preserves it.
