// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

//Config ...
type Config struct {
	Path       string
	ListenAddr string
}

//DefaultConfig ...
var DefaultConfig = Config{
	Path:       "/prom",
	ListenAddr: "8000",
}
