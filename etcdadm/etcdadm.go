package etcdadm

import (
	"bufio"
	"encoding/base64"
	"io/ioutil"
	"log"
	"text/template"

	"github.com/tmc/scp"
	"golang.org/x/crypto/ssh"
)

func Init(session *ssh.Session, etcdConfig *EtcdConfig) {
	writeEtcdServiceEnvFile(etcdConfig)
	writeCertFiles(etcdConfig)
}

func writeCertFiles(etcdConfig *EtcdConfig) {
	writeRemoteFile(base64.StdEncoding.DecodeString(etcdConfig.ClientCert), etcdConfig.CertFile)
	writeRemoteFile(base64.StdEncoding.DecodeString(etcdConfig.ClientKey), etcdConfig.KeyFile)
}

func writeRemoteFile(data []byte, remoteFile string, session *ssh.Session) {
	tmpfile, err := ioutil.TempFile("", "pf9tmp")
	tmpfile.WriteString(string(data))
	defer tmpfile.Close()
	scp.CopyPath(tmpfile.Name, remoteFile, session)
}

func writeEtcdServiceEnvFile(etcdConfig *EtcdConfig) {
	envFileData := `
ETCD_NAME={{.Name}}
ETCD_STRICT_RECONFIG_CHECK={{.StrictReconfigCheck}}
ETCD_INITIAL_CLUSTER_TOKEN={{.InitialClusterToken}}
ETCD_INITIAL_CLUSTER_STATE={{.InitialClusterState}}
ETCD_INITIAL_CLUSTER={{.InitialCluster}}
ETCD_INITIAL_ADVERTISE_PEER_URLS={{.InitialAdvertisePeerUrls}}
ETCD_LISTEN_PEER_URLS={{.ListenPeerUrls}}
ETCD_ADVERTISE_CLIENT_URLS={{.AdvertiseClientUrls}}
ETCD_LISTEN_CLIENT_URLS={{.ListenClientUrls}}
ETCD_DATA_DIR={{.DataDir}}
ETCD_CERT_FILE={{.CertFile}}
ETCD_KEY_FILE={{.KeyFile}}
ETCD_TRUSTED_CA_FILE={{.TrustedCaFile}}
ETCD_PEER_KEY_FILE={{.PeerKeyFile}}
ETCD_PEER_CERT_FILE={{.PeerCertFile}}
ETCD_PEER_TRUSTED_CA_FILE={{.PeerTrustedCaFile}}
ETCD_CLIENT_CERT_AUTH={{.ClientCertAuth}}
ETCD_DEBUG={{.Debug}}`
	t := template.Must(template.New("envFileData").Parse(envFileData))
	tmpfile, err := ioutil.TempFile("", "etcdenv")
	if err != nil {
		log.Fatalf("Could not create tmp file for etcd env configuration")
	}
	//defer os.Remove(tmpfile.Name())
	defer tmpfile.Close()
	writer := bufio.NewWriter(tmpfile)
	err = t.Execute(writer, etcdConfig)
	writer.Flush()
	if err != nil {
		log.Fatalf("Could not write env file for etcd %v", err)
	}
}

func Join() {

}

func Upgrade() {

}

func Reset() {

}
