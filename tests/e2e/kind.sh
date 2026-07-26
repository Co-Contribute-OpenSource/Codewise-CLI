#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
test_root="$(mktemp -d)"
trap 'rm -rf "$test_root"' EXIT

export CODEWISE_HOME="$test_root/home"
codewise="$test_root/codewise"
go build -o "$codewise" "$repo_root"

create_environment() {
  local name="$1"
  local strategy="$2"
  "$codewise" env create "$name"
  printf 'strategy: %s\n' "$strategy" >>"$CODEWISE_HOME/envs/$name/k8s.yaml"
}

kubectl create namespace codewise-e2e

# kubectl strategy: apply a real object and verify it through the API server.
kubectl_dir="$test_root/kubectl"
mkdir -p "$kubectl_dir/k8s"
printf '%s\n' \
  'apiVersion: v1' \
  'kind: ConfigMap' \
  'metadata:' \
  '  name: codewise-kubectl' \
  >"$kubectl_dir/k8s/configmap.yaml"
create_environment kubectl-e2e kubectl
sed -i 's/namespace: kubectl-e2e/namespace: codewise-e2e/' "$CODEWISE_HOME/envs/kubectl-e2e/k8s.yaml"
(cd "$kubectl_dir" && "$codewise" deploy run --env kubectl-e2e)
kubectl -n codewise-e2e get configmap codewise-kubectl
"$codewise" deploy status --env kubectl-e2e

# Helm strategy: install a minimal chart and validate Helm-owned state.
helm_dir="$test_root/helm"
mkdir -p "$helm_dir/helm/chart/templates"
printf '%s\n' \
  'apiVersion: v2' \
  'name: codewise-e2e' \
  'version: 0.1.0' \
  >"$helm_dir/helm/chart/Chart.yaml"
printf '%s\n' \
  'apiVersion: v1' \
  'kind: ConfigMap' \
  'metadata:' \
  '  name: codewise-helm' \
  >"$helm_dir/helm/chart/templates/configmap.yaml"
create_environment helm-e2e helm
sed -i 's/namespace: helm-e2e/namespace: codewise-e2e/' "$CODEWISE_HOME/envs/helm-e2e/k8s.yaml"
(cd "$helm_dir" && "$codewise" deploy run --env helm-e2e)
helm -n codewise-e2e status helm-e2e
"$codewise" deploy history --env helm-e2e
printf '%s\n' 'data:' '  revision: "2"' >>"$helm_dir/helm/chart/templates/configmap.yaml"
(cd "$helm_dir" && "$codewise" deploy run --env helm-e2e)
"$codewise" deploy rollback --env helm-e2e --revision 1

# GitOps strategy: install a minimal Application CRD, apply the generated
# Argo CD Application, and verify it is accepted by the API server.
kubectl apply -f - <<'EOF'
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: applications.argoproj.io
spec:
  group: argoproj.io
  names:
    kind: Application
    plural: applications
    singular: application
  scope: Namespaced
  versions:
  - name: v1alpha1
    served: true
    storage: true
    schema:
      openAPIV3Schema:
        type: object
        x-kubernetes-preserve-unknown-fields: true
EOF
kubectl create namespace argocd
create_environment gitops-e2e gitops
sed -i 's/namespace: gitops-e2e/namespace: codewise-e2e/' "$CODEWISE_HOME/envs/gitops-e2e/k8s.yaml"
sed -i 's#repo: ""#repo: https://github.com/example/repo.git#' "$CODEWISE_HOME/envs/gitops-e2e/gitops.yaml"
(cd "$test_root" && "$codewise" deploy run --env gitops-e2e)
kubectl -n argocd get application gitops-e2e
"$codewise" deploy status --env gitops-e2e

# Docker strategy: perform an actual local build without downloading a base.
docker_dir="$test_root/docker"
mkdir -p "$docker_dir"
printf '%s\n' 'FROM scratch' 'COPY marker /marker' >"$docker_dir/Dockerfile"
printf '%s\n' 'codewise-e2e' >"$docker_dir/marker"
(cd "$docker_dir" && "$codewise" docker build --tag codewise-e2e:test)
docker image inspect codewise-e2e:test >/dev/null
