// Copyright (c) 2018 Nikita Chisnikov
// Distributed under the MIT/X11 software license

package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func getMyPath() string {
	ex, err := os.Executable()
	if err != nil {
		log.Fatal("getMyPath():", err.Error())
	}
	exPath := filepath.Dir(ex)
	return exPath
}

func handlePath(path, cfname, dfname string) *BlockInfo {
	if path == "" {
		return nil
	}
	path, err := filepath.Abs(path)
	if err != nil {
		log.Fatal("filepath.Abs():", err.Error())
	}
	path += "/"
	pathD, err := exec.LookPath(path + dfname)
	if err != nil {
		log.Fatal("LookPath("+dfname+"):", err.Error())
	}
	pathC, err := exec.LookPath(path + cfname)
	if err != nil {
		log.Fatal("LookPath("+cfname+"):", err.Error())
	}
	chsumD, err := checksum(pathD)
	if err != nil {
		log.Fatal("checksum("+dfname+"):", err.Error())
	}
	chsumC, err := checksum(pathC)
	if err != nil {
		log.Fatal("checksum("+cfname+"):", err.Error())
	}
	return &BlockInfo{Path: path, Cli: &CoreFile{Name: cfname, Chsum: chsumC},
		Daemon: &CoreFile{Name: dfname, Chsum: chsumD}}
}

func checkDaemonFile(bi *BlockInfo) error {
	var chsum string
	var err error
	if bi == nil {
		log.Fatal("checkFiles(): BlockInfo nil:")
	}
	chsum, err = checksum(bi.Path + bi.Daemon.Name)
	if err != nil {
		log.Fatal("checksum(daemon):", err.Error())
	}
	if chsum != bi.Daemon.Chsum {
		log.Fatal("daemon checksum changed!")
	}
	log.Println(bi.Daemon.Name, "checksum match")
	return err
}

func checkCliFile(bi *BlockInfo) error {
	var chsum string
	var err error
	if bi == nil {
		log.Fatal("checkFiles(): BlockInfo nil:")
	}
	chsum, err = checksum(bi.Path + bi.Cli.Name)
	if err != nil {
		log.Fatal("checksum(cli):", err.Error())
	}
	if chsum != bi.Cli.Chsum {
		log.Fatal("cli checksum changed!")
	}
	log.Println(bi.Cli.Name, "checksum match")
	return err
}
