package syslog

import (
	"flag"
	"os"

	"gopkg.in/yaml.v3"
)

type SyslogConfig struct {
	QueueSize int          `yaml:"queue_size,omit_empty"`
	Host      string       `yaml:"host"`
	Port      int          `yaml:"port,omit_empty"`
	Protocol  string       `yaml:"protocol,omit_empty"`
	Facility  int          `yaml:"facility,omit_empty"`
	Formatter RfcFormatter `yaml:"formatter,omit_empty"`
	Rfc       int          `yaml:"rfc,omit_empty"`
	TLS       TLS          `yaml:"tls,omit_empty"`
	BasicAuth BasicAuth    `yaml:"basic_auth,omit_empty"`
}

type RfcFormatter struct {
	RfcNum int `yaml:"rfcNum,omit_empty"`
}

type TLS struct {
	CaCrt              interface{} `yaml:"ca.crt,omit_empty"`
	CertCrt            interface{} `yaml:"cert.crt,omit_empty"`
	CertKey            interface{} `yaml:"cert.key,omit_empty"`
	ServerName         interface{} `yaml:"server_name,omit_empty"`
	InsecureSkipVerify bool        `yaml:"insecure_skip_verify,omit_empty"`
}
type BasicAuth struct {
	Username interface{} `yaml:"username"`
	Password interface{} `yaml:"password"`
}

const (
	DEF_QUEUE_SIZE_VALUE           = 1000
	DEF_INSECURE_SKIP_VERIFY_VALUE = false
	DEF_RFC_NUM_VALUE              = 3164
	DEF_SYSLOG_SERVER_PORT         = 514
	DEF_SYSLOG_FACILITY            = 16
	DEF_SYSLOG_PROTOCOL            = "udp"
)

var (
	sysCfg       SyslogConfig
	syslogConfig = flag.String("syslogConfig", "", "Configuration file for syslog")
)

func Init() {
	if *syslogConfig != "" {
		cfgData, err := os.ReadFile(*syslogConfig)
		if err != nil {
			panic(err)
		}

		err = yaml.Unmarshal(cfgData, &sysCfg)
		if err != nil {
			panic(err)
		}
		if sysCfg.QueueSize == 0 {
			sysCfg.QueueSize = DEF_QUEUE_SIZE_VALUE
		}

		if sysCfg.Port == 0 {
			sysCfg.Port = DEF_SYSLOG_SERVER_PORT
		}

		if sysCfg.Protocol == "" {
			sysCfg.Protocol = DEF_SYSLOG_PROTOCOL
		}

	}

}

type SyslogWriter struct {
	In *chan string
}

func GetSyslogWriter() *SyslogWriter {
	inChan := make(chan string, sysCfg.QueueSize)

	return &SyslogWriter{
		In: &inChan,
	}
}

func (w *SyslogWriter) Write(b []byte) (int, error) {
	lMsg := string(b)

	select {
	case *w.In <- lMsg:
		//message sent
	default:
		<-*w.In
		*w.In <- lMsg
	}

	return 0, nil
}
