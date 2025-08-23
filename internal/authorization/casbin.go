package authorization

import (
	casbin "github.com/casbin/casbin/v2"
	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
)

// NewEnforcer membuat model dan policy dari file.
// Parameter modelPath dan policyPath adalah lokasi file .conf dan .csv.
func NewEnforcer(modelPath, policyPath string) (*casbin.Enforcer, error) {
	adapter := fileadapter.NewAdapter(policyPath)
	enforcer, err := casbin.NewEnforcer(modelPath, adapter) 
	if err != nil {
		return nil, err
	}
	if err := enforcer.LoadPolicy(); err != nil {
		return nil, err
	}
	return enforcer, nil
}

// - Agar Casbin dapat berjalan, Anda membutuhkan dua file konfigurasi di proyek Anda:
// - rbac_model.conf – mendefinisikan struktur model otorisasi (bagaimana permintaan, policy, dan role dicocokkan).
// - rbac_policy.csv – berisi daftar policy (aturan) yang menetapkan role mana yang boleh melakukan aksi apa terhadap resource tertentu.