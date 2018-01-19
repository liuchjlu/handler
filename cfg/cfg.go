package cfg

type HttpConfig struct {
	Enabled bool   `json:"enabled"`
	Listen  string `json:"listen"`
}

type RpcConfig struct {
	Enabled bool   `json:"enabled"`
	Listen  string `json:"listen"`
}

// type K8sAddr struct {
// 	Address	string 	`json:"address"`
// }

// type EtcdAddr struct {
// 	Address	string 	`json:"address"`
// }
var (
	K8sAddr		string
	EtcdAddr	string
)