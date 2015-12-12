// Copyright 2015 Felipe A. Cavani. All rights reserved.
// Use of this source code is governed by Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"net"
	"net/smtp"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/fcavani/e"
	"github.com/fcavani/log"
	"github.com/fcavani/net/dns"
	mysmtp "github.com/fcavani/net/smtp"
	flags "github.com/jessevdk/go-flags"
	"gopkg.in/ini.v1"

	"monlite"
)

const appname = "monlite"

type options struct {
	Conf  string `short:"c" long:"configuration" description:"Configuration file." required:"true" default:"/etc/monlite.ini"`
	Log   string `short:"l" long:"log" description:"File to log to."`
	Level string `short:"v" long:"level" description:"Log level for log going to stderr. The options are: nolog, panic, fatal, error, warning, info, debug and protocol." default:"debug"`
}

func main() {
	println("Starting monlite...")
	// Configuration
	var opts options
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal("can't parse the command line options:", err)
	}
	cfg, err := ini.Load(opts.Conf)
	if err != nil {
		log.Fatal("Error reading configuratio file:", opts.Conf)
	}

	// Log stuff
	println("Log...")
	name := appname
	pid := os.Getpid()
	pidstr := strconv.FormatInt(int64(pid), 10)
	name = name + " (" + pidstr + ")"

	fnull, err := os.OpenFile("/dev/null", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		log.Fatal("open log file failed:", err)
	}
	defer fnull.Close()
	stderrBack := log.NewWriter(fnull).F(log.DefFormatter)
	if opts.Level != "" && opts.Level != "nolog" {
		level, err := log.ParseLevel(opts.Level)
		if err != nil {
			log.Fatal("Invalid log level.")
		}
		stderrBack = log.Filter(
			log.NewWriter(os.Stderr).F(log.DefFormatter),
			log.Op(log.Ge, "level", level),
		)
	}

	logfile := cfg.Section("log").Key("file").MustString("/var/log/monlite.log")
	if opts.Log != "" {
		logfile = opts.Log
	}
	loglevel := cfg.Section("log").Key("level").MustString("debug")
	level, err := log.ParseLevel(loglevel)
	if err != nil {
		log.Fatal("Invalid log level.")
	}
	var fileBack log.LogBackend
	if logfile != "" {
		f, err := os.OpenFile(logfile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
		if err != nil {
			log.Fatalf("open log file %v failed: %v", logfile, err)
		}
		defer f.Close()
		fileBack = log.Filter(
			log.NewWriter(f).F(log.DefFormatter),
			log.Op(log.Ge, "level", level),
		)
	}

	log.Log = log.New(stderrBack, true).Domain(name).InfoLevel()
	if fileBack != nil {
		log.Log = log.New(
			log.NewMulti(
				stderrBack,
				log.DefFormatter,
				fileBack,
				log.DefFormatter,
			),
			true,
		).Domain(name).InfoLevel()
	}

	log.Println("Log Ok!")

	dns.SetLookupHostFunction(net.LookupHost)

	log.Println("Configuration...")

	cfgMail := cfg.Section("mail")

	smtpTimeout, err := cfgMail.Key("timeout").Int()
	if err != nil {
		log.Fatal("invalid smtp timeout")
	}

	server := cfgMail.Key("smtp").String()
	s := strings.Split(server, ":")
	if len(s) > 0 {
		server = s[0]
	}

	auth := smtp.PlainAuth(
		"",
		cfgMail.Key("account").String(),
		cfgMail.Key("password").String(),
		server,
	)

	hostname := "your system"
	if hn, err := os.Hostname(); err == nil {
		hostname = hn
	}

	mons := make([]*monlite.Monitor, 0)
	for _, sec := range cfg.Sections() {
		if !strings.HasPrefix(sec.Name(), "service.") {
			continue
		}
		name := strings.TrimPrefix(sec.Name(), "service.")
		log.DebugLevel().Printf("Adding monitor %v for url: %v", name, sec.Key("url").String())
		to, err := sec.Key("timeout").Int()
		if err != nil {
			log.Fatalf("invalid value in timeout for %v", name)
		}
		p, err := sec.Key("periode").Int()
		if err != nil {
			log.Fatalf("invalid value in timeout for %v", name)
		}
		sleep, err := sec.Key("sleep").Int()
		if err != nil {
			log.Fatalf("invalid value in sleep for %v", name)
		}
		fails, err := sec.Key("fails").Int()
		if err != nil {
			log.Fatalf("invalid value in fails for %v", name)
		}
		mons = append(mons, &monlite.Monitor{
			Name:    name,
			Url:     sec.Key("url").String(),
			Timeout: time.Duration(to) * time.Second,
			Periode: time.Duration(p) * time.Second,
			Sleep:   time.Duration(sleep) * time.Second,
			Fails:   fails,
			OnFail: func(m *monlite.Monitor) error {

				body := "Mime-Version: 1.0\n"
				body += "Content-Type: text/plain; charset=utf-8\n"
				body += "From:" + cfgMail.Key("from").String() + "\n"
				body += "To:" + cfgMail.Key("to").String() + "\n"
				body += "Subject: [" + hostname + "] " + "Monitor fail for " + m.Name + "\n"
				body += "Hi! This is " + hostname + ".\n\n"
				body += "Monitor fail for " + m.Name + " " + m.Url + "\n\n"
				body += "Tank you, our lazy boy.\n"
				body += time.Now().Format(time.RFC1123Z)

				err := mysmtp.SendMail(
					cfgMail.Key("smtp").String(),
					auth,
					cfgMail.Key("from").String(),
					[]string{cfgMail.Key("to").String()},
					cfgMail.Key("helo").String(),
					[]byte(body),
					time.Duration(smtpTimeout)*time.Second,
					false,
				)
				if err != nil {
					return e.Forward(err)
				}
				return nil
			},
			OnUnFail: func(m *monlite.Monitor) error {

				body := "Mime-Version: 1.0\n"
				body += "Content-Type: text/plain; charset=utf-8\n"
				body += "From:" + cfgMail.Key("from").String() + "\n"
				body += "To:" + cfgMail.Key("to").String() + "\n"
				body += "Subject: [" + hostname + "] " + "Monitor ok for " + m.Name + "\n"
				body += "Hi! This is " + hostname + ".\n\n"
				body += "Monitor ok for " + m.Name + " " + m.Url + "\n\n"
				body += "Tank you, our lazy boy.\n"
				body += time.Now().Format(time.RFC1123Z)

				err := mysmtp.SendMail(
					cfgMail.Key("smtp").String(),
					auth,
					cfgMail.Key("from").String(),
					[]string{cfgMail.Key("to").String()},
					cfgMail.Key("helo").String(),
					[]byte(body),
					time.Duration(smtpTimeout)*time.Second,
					false,
				)
				if err != nil {
					return e.Forward(err)
				}
				return nil
			},
		})
	}

	log.Println("Starting monitors...")

	for _, m := range mons {
		err := m.Start()
		if err != nil {
			log.Fatalf("Failed to start monitor for %v. Error: %v", m.Name, err)
		}
	}

	log.Println("Monitors ok!")

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	<-sig

	log.Println("Stop monitors...")

	for _, m := range mons {
		log.DebugLevel().Printf("Stop monitor %v", m.Name)
		err := m.Stop()
		if err != nil {
			log.Errorf("Failed to start monitor for %v. Error: %v", m.Name, err)
		}
	}
	log.Println("End.")
}
