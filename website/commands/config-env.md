# Config and environment commands

## `config init`

Create `~/.codewise/config.yaml` with defaults:

```bash
codewise config init
```

The command leaves an existing configuration unchanged.

## `config view`

```bash
codewise config view
```

## `env create <name>`

```bash
codewise env create dev
codewise env create staging --interactive
```

| Flag | Short | Description |
| --- | --- | --- |
| `--interactive` | `-i` | Prompt for namespace, context, repository, branch, image, and tag |

Environment names are Kubernetes-safe DNS labels: lowercase letters, numbers,
and hyphens, with a maximum length of 63 characters.

## `env list`

```bash
codewise env list
```

## `env show <name>`

Print the complete environment without opening several component files:

```bash
codewise env show staging
codewise env show staging --format json
```

## `env strategy <name> <strategy>`

Persist strategy selection instead of relying on files in the current working
directory:

```bash
codewise env strategy staging helm
codewise env strategy staging kubectl
codewise env strategy staging gitops
codewise env strategy staging auto
```

`auto` restores automatic selection. Explicit selection is recommended for CI
and important environments.

## `env delete <name>`

```bash
codewise env delete dev
codewise env delete dev --yes
```

| Flag | Description |
| --- | --- |
| `--yes` | Skip the confirmation prompt |
