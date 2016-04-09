package RemoteSyslog

import (
	"../..//op/go-logging"
	"time"
	"os"
	"errors"
	"net"
	"strconv"
)

type PapertrailBackend struct {
	ClientHostname         string
	Tag			string

	Hostname               string
	Network                string
	ConnectTimeout         int
	WriteTimeout           int
	Port                   int

	Logger                 *Logger

	connectTimeoutDuration time.Duration
	writeTimeoutDuration   time.Duration
	raddr                  string
}

func NewPapertrailBackend(config *PapertrailBackend) (*PapertrailBackend) {
	if config.ClientHostname == "" {
		hostname, err := os.Hostname()
		if err != nil {
			config.ClientHostname = "Unknown"
		} else {
			config.ClientHostname = hostname
		}
	}

	if config.Tag == "" {
		config.Tag = "logging"
	}

	if config.ConnectTimeout == 0 {
		config.connectTimeoutDuration = time.Duration(30) * time.Second
	} else {
		config.connectTimeoutDuration = time.Duration(config.ConnectTimeout) * time.Second
	}

	if config.WriteTimeout == 0 {
		config.writeTimeoutDuration = time.Duration(30) * time.Second
	} else {
		config.writeTimeoutDuration = time.Duration(config.WriteTimeout) * time.Second
	}

	switch config.Network {
	case "udp":
		config.Network = "udp"
	case "tcp":
		config.Network = "tcp"
	default:
		return errors.New("Invalid Network. Neither UDP nor TCP.")
	}

	config.raddr = net.JoinHostPort(config.Hostname, strconv.Itoa(config.Port))

	var err error
	config.Logger, err = Dial(config.ClientHostname, config.Network, config.raddr, nil, config.connectTimeoutDuration, config.writeTimeoutDuration, 99990)
	if err != nil {
		return err
	}

	return config
}

func (w *PapertrailBackend) Crit(line []byte) (error) {
	w.Logger.Packets <- Packet{
		SevCrit,
		LogLocal0,
		w.ClientHostname,
		w.Tag,
		time.Now(),
		line,
	}

	return nil
}

func (w *PapertrailBackend) Err(line []byte) (error) {
	w.Logger.Packets <- Packet{
		SevErr,
		LogLocal0,
		w.ClientHostname,
		w.Tag,
		time.Now(),
		line,
	}

	return nil
}

func (w *PapertrailBackend) Warning(line []byte) (error) {
	w.Logger.Packets <- Packet{
		SevWarning,
		LogLocal0,
		w.ClientHostname,
		w.Tag,
		time.Now(),
		line,
	}

	return nil
}

func (w *PapertrailBackend) Notice(line []byte) (error) {
	w.Logger.Packets <- Packet{
		SevNotice,
		LogLocal0,
		w.ClientHostname,
		w.Tag,
		time.Now(),
		line,
	}

	return nil
}

func (w *PapertrailBackend) Info(line []byte) (error) {
	w.Logger.Packets <- Packet{
		SevInfo,
		LogLocal0,
		w.ClientHostname,
		w.Tag,
		time.Now(),
		line,
	}

	return nil
}

func (w *PapertrailBackend) Debug(line []byte) (error) {
	w.Logger.Packets <- Packet{
		SevDebug,
		LogLocal0,
		w.ClientHostname,
		w.Tag,
		time.Now(),
		line,
	}

	return nil
}

func (b *PapertrailBackend) Log(level logging.Level, calldepth int, rec *logging.Record) error {
	line := rec.Formatted(calldepth + 1)
	switch level {
	case logging.CRITICAL:
		return b.Crit(line)
	case logging.ERROR:
		return b.Err(line)
	case logging.WARNING:
		return b.Warning(line)
	case logging.NOTICE:
		return b.Notice(line)
	case logging.INFO:
		return b.Info(line)
	case logging.DEBUG:
		return b.Debug(line)
	default:
	}

	panic("Unhandled log level")
}
