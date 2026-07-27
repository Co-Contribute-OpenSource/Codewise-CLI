package env

type K8sConfig struct {
	Namespace string `json:"namespace" yaml:"namespace"`
	Context   string `json:"context" yaml:"context"`
	Strategy  string `json:"strategy,omitempty" yaml:"strategy,omitempty"`
}

type HelmConfig struct {
	Release string `json:"release" yaml:"release"`
	Chart   string `json:"chart" yaml:"chart"`
	Values  string `json:"values" yaml:"values"`
}

type GitOpsConfig struct {
	Repo   string `json:"repo" yaml:"repo"`
	Path   string `json:"path" yaml:"path"`
	Branch string `json:"branch" yaml:"branch"`
}

type ValuesConfig struct {
	Image struct {
		Repository string `json:"repository" yaml:"repository"`
		Tag        string `json:"tag" yaml:"tag"`
	} `json:"image" yaml:"image"`
}

type Env struct {
	Name   string
	K8s    K8sConfig
	Helm   HelmConfig
	GitOps GitOpsConfig
	Values ValuesConfig
}
