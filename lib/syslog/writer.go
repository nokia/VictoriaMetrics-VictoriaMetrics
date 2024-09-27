// Current File defines the syslogWriter and the applicable methods
package syslog

import (
	"crypto/tls"
	"fmt"
	"net"
	"os"
)

type syslogWriter struct {
	conn      serverConn
	formatter formatter
	framer    framer
	sysCfg    *config
	tlsConfig *tls.Config
}

type serverConn interface {
	writeString(framer framer, formatter formatter, priority int64, hostname string, s string) error
	close() error
}

// basicDialer connects to the syslog server
func (w *syslogWriter) basicDialer() (serverConn, error) {
	c, err := net.Dial(w.sysCfg.SyslogConfig.Protocol, fmt.Sprintf("%s:%d", w.sysCfg.SyslogConfig.RemoteHost, w.sysCfg.SyslogConfig.Port))
	var sc serverConn
	if err == nil {
		sc = &netConn{conn: c}
	}
	return sc, err
}

func (w *syslogWriter) tlsDialer() (serverConn, error) {
	c, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", w.sysCfg.SyslogConfig.RemoteHost, w.sysCfg.SyslogConfig.Port), w.tlsConfig)
	var sc serverConn
	if err == nil {
		sc = &netConn{conn: c}
	} else {
		fmt.Printf("Error establishing the tls connection %v\n", err)
	}
	return sc, err
}

// connect updates the syslogWriter with a new serverConn
func (w *syslogWriter) connect() (serverConn, error) {
	var (
		conn serverConn
		err  error
	)
	if w.sysCfg.SyslogConfig.Protocol == "tcp+tls" {
		conn, err = w.tlsDialer()
	} else {
		conn, err = w.basicDialer()
	}
	if err == nil {
		w.conn = conn
		return conn, nil
	} else {
		return nil, err
	}
}

// send forwards the log message to the syslog server
func (w *syslogWriter) send(logLevel, msg string) (int, error) {
	priority := (w.sysCfg.SyslogConfig.Facility << 3) | logLevelMap[logLevel]

	var err error
	if w.conn != nil {
		err = w.conn.writeString(w.framer, w.formatter, priority, w.getHostname(), msg)
		if err == nil {
			return len(msg), nil
		}
	}
	//Establishes a new connection with the syslog server
	_, err = w.connect()

	if err != nil {
		return 0, err
	}

	err = w.conn.writeString(w.framer, w.formatter, priority, w.getHostname(), msg)
	if err != nil {
		return 0, err
	}
	return len(msg), nil
}

// getHostname returns the hostname field to be used in the log message
func (w *syslogWriter) getHostname() string {
	hostname := w.sysCfg.SyslogConfig.Hostname
	if hostname == "" {
		hostname, _ = os.Hostname()
	}
	return hostname
}

// logSender observes the buffered channel for log data to be written to the syslog server
func (w *syslogWriter) logSender() {
	for logEntry := range logChan {
		_, err := w.send(logEntry.LogLevel, logEntry.Msg)
		for err != nil {
			_, err = w.send(logEntry.LogLevel, logEntry.Msg)
		}
	}
}
