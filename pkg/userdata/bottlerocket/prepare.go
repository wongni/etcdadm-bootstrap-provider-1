package bottlerocket

import (
	"path/filepath"

	"github.com/go-logr/logr"
	"github.com/mrajashree/etcdadm-bootstrap-provider/pkg/userdata"
	"sigs.k8s.io/cluster-api/bootstrap/kubeadm/api/v1alpha3"
)

const (
	orgCertsPath = "/etc/etcd/pki"
	newCertsPath = "/var/lib/etcd/pki"
)

func prepare(input *userdata.BaseUserData) {
	input.Header = cloudConfigHeader
	input.WriteFiles = append(input.WriteFiles, input.AdditionalFiles...)
	input.SentinelFileCommand = sentinelFileCommand
	patchCertPaths(input)
}

func patchCertPaths(input *userdata.BaseUserData) {
	// hacky. new array bc I need to update the original object and these are not pointers
	files := make([]v1alpha3.File, 0, len(input.WriteFiles))
	for _, file := range input.WriteFiles {
		if filepath.Dir(file.Path) == orgCertsPath {
			file.Path = filepath.Join(newCertsPath, filepath.Base(file.Path))
		}
		files = append(files, file)
	}

	input.WriteFiles = files
}

func logIgnoredFields(input *userdata.BaseUserData, log logr.Logger) {
	if len(input.PreEtcdadmCommands) > 0 {
		log.Info("Ignoring PreEtcdadmCommands. Not supported with bottlerocket")
	}
	if len(input.PostEtcdadmCommands) > 0 {
		log.Info("Ignoring PostEtcdadmCommands. Not supported with bottlerocket")
	}
	if input.NTP != nil {
		log.Info("Ignoring NTP. Not supported with bottlerocket")
	}
	if input.DiskSetup != nil {
		log.Info("Ignoring DiskSetup. Not supported with bottlerocket")
	}
	if len(input.Mounts) > 0 {
		log.Info("Ignoring Mounts. Not supported with bottlerocket")
	}
}
