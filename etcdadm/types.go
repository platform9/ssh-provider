package etcdadm

type EtcdConfig struct {
	Name                     string `json:"name"`
	StrictReconfigCheck      bool   `json:"strictReconfigCheck"`
	InitialClusterState      string `json:"initialClusterState"`
	InitialCluster           string `json:"initialCluster"`
	InitialClusterToken      string `json:"initialClusterToken"`
	InitialAdvertisePeerUrls string `json:"initialAdvertisePeerUrls"`
	ListenPeerUrls           string `json:"listenPeerUrls"`
	AdvertiseClientUrls      string `json:"advertiseClientUrls"`
	ListenClientUrls         string `json:"listenClientUrls"`
	DataDir                  string `json:"dataDir"`
	CertFile                 string `json:"certFile"`
	KeyFile                  string `json:"keyFile"`
	TrustedCaFile            string `json:"trustedCaFile"`
	ClientCert               string `json:"clientCertFile"`
	ClientKey                string `json:"clientKeyFile"`
	ClientTrustedCa          string `json:"clientTrustedCaFile"`
	PeerKeyFile              string `json:"peerKeyFile"`
	PeerCertFile             string `json:"peerCertFile"`
	PeerTrustedCaFile        string `json:"peerTrustedCaFile"`
	PeerKey                  string `json:"peerKey"`
	PeerCert                 string `json:"peerCert"`
	PeerTrustedCa            string `json:"peerTrustedCa"`
	ClientCertAuth           bool   `json:"clientCertAuth"`
	Debug                    bool   `json:"debug"`
}
