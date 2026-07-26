package generator

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func ScaffoldProject(projectName string, withDocker, withDeployment bool) error {
	if projectName == "" {
		return fmt.Errorf("please provide a project name using --project")
	}

	basePath := filepath.Join(".", projectName)

	// Folder structure
	dirs := []string{
		"cmd", "pkg", "internal", "configs", "scripts", "templates",
	}

	for _, dir := range dirs {
		fullPath := filepath.Join(basePath, dir)
		if err := os.MkdirAll(fullPath, 0755); err != nil {
			return fmt.Errorf("create %s: %w", fullPath, err)
		}
	}

	// Dockerfile
	if withDocker {
		dockerfile := `FROM golang:1.25-alpine
WORKDIR /app
COPY . .
RUN go build -o main .
CMD ["./main"]`

		err := os.WriteFile(filepath.Join(basePath, "Dockerfile"), []byte(dockerfile), 0644)
		if err != nil {
			return fmt.Errorf("write Dockerfile: %w", err)
		}
		fmt.Println("📦 Dockerfile created.")
	}

	// Kubernetes deployment.yaml
	if withDeployment {
		k8sPath := filepath.Join(basePath, "k8s")
		if err := os.MkdirAll(k8sPath, 0755); err != nil {
			return fmt.Errorf("create Kubernetes directory: %w", err)
		}

		deployment := `apiVersion: apps/v1
kind: Deployment
metadata:
  name: ` + projectName + `
spec:
  replicas: 1
  selector:
    matchLabels:
      app: ` + projectName + `
  template:
    metadata:
      labels:
        app: ` + projectName + `
    spec:
      containers:
      - name: ` + projectName + `
        image: ` + projectName + `:latest
        ports:
        - containerPort: 8080`

		err := os.WriteFile(filepath.Join(k8sPath, "deployment.yaml"), []byte(deployment), 0644)
		if err != nil {
			return fmt.Errorf("write deployment.yaml: %w", err)
		}
		fmt.Println("📄 k8s/deployment.yaml created.")
	}

	if err := setupGitRepo(basePath); err != nil {
		return err
	}

	fmt.Println("✅ Project scaffolded successfully.")
	return nil
}

func setupGitRepo(basePath string) error {
	cmd := exec.Command("git", "init")
	cmd.Dir = basePath
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("initialize Git repository: %w", err)
	}

	gitignore := `# Binaries
*.exe
*.dll
*.so
*.dylib
*.test
*.out

# Vendor
/vendor/

# Logs
*.log

# IDEs and editors
.vscode/
.idea/
*.swp

# Build
/build/
bin/
`

	err := os.WriteFile(filepath.Join(basePath, ".gitignore"), []byte(gitignore), 0644)
	if err != nil {
		return fmt.Errorf("write .gitignore: %w", err)
	}

	fmt.Println("🔧 Git repo initialized with .gitignore")
	return nil
}
