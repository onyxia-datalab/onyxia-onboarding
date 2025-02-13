package domain

type Quota struct {
	MemoryRequest           string
	CPURequest              string
	MemoryLimit             string
	CPULimit                string
	StorageRequest          string
	MaxPods                 string
	EphemeralStorageRequest string
	EphemeralStorageLimit   string
	GPURequest              string
	GPULimit                string
}

type Quotas struct {
	Enabled      bool
	Default      Quota
	UserEnabled  bool
	User         Quota
	GroupEnabled bool
	Group        Quota
}
