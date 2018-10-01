// Copyright (c) 2018 Nikita Chisnikov
// Distributed under the MIT/X11 software license

package main

import (
	"io"
	"log"
	"flag"
	"time"
	"bytes"
	"io/ioutil"
	"os/exec"
	"encoding/json"
	"runtime/debug"
)

type BlockInfo struct {
	Path	string		`json:"path" binding:"required"`
	Cli	*CoreFile	`json:"cli" binding:"required"`
	Daemon	*CoreFile	`json:"daemon" binding:"required"`
	Wpass	string		`json:"wpass" binding:"required"`
}

type CoreFile struct {
	Name	string		`json:"name" binding:"required"`
	Chsum	string		`json:"chsum" binding:"required"`
}

const C_CLI_FNAME =	"blocknetdx-cli"
const C_DAEMON_FNAME =	"blocknetdxd"

func saveInfo(path, wpwd, spwd string) {
	var bi, bi_old *BlockInfo
	var err error
	var js []byte

	bi_old = loadInfo(spwd)
	if path != "" {
		log.Println("checking blocknet core files..")
		if bi = handlePath(path, C_CLI_FNAME, C_DAEMON_FNAME); bi == nil { return }
		log.Println("blocknet path:", bi.Path)
		log.Println(bi.Cli.Name, "checksum done!")
		log.Println(bi.Daemon.Name, "checksum done!")
	} else { bi = bi_old }
	if wpwd != "" {
		if bi == nil { bi = &BlockInfo{} }
		bi.Wpass = wpwd
	} else {
		if bi_old != nil { bi.Wpass = bi_old.Wpass }
	}
	js, err = json.Marshal(bi)
	bi = nil; bi_old = nil
	if err != nil {
		log.Fatal("json.Marshal()", err.Error())
	}
	log.Println("setting data to storage..")
	if err = setData(getMyPath() + "/blockman.dat", string(encrypt(js, genPwd(spwd)))); err != nil { return }
	js = nil
	log.Println("done!")
}

func loadInfo(spwd string) *BlockInfo {
	var bi *BlockInfo
	var err error
	var data string
	var btdata []byte
	log.Println("getting data from storage..")
	data, err = getData(getMyPath() + "/blockman.dat")
	if err != nil { return nil }
	bi = &BlockInfo{}
	btdata, err = decrypt([]byte(data), genPwd(spwd))
	if err != nil { bi = nil; return nil }
	if err = json.Unmarshal(btdata, bi); err != nil {
		log.Fatal("json.Unmarshal():", err.Error())
		bi = nil
		return nil
	}
	return bi
}

func printInfo(bi BlockInfo) {
	var err error
	var js []byte

	if len(bi.Wpass) > 1 {
		bi.Wpass = "***" + bi.Wpass[len(bi.Wpass)-1:]
	}
	if bi.Cli != nil && len(bi.Cli.Chsum) > 13 {
		bi.Cli.Chsum = "***" + bi.Cli.Chsum[len(bi.Cli.Chsum)-13:]
	}
	if bi.Daemon != nil && len(bi.Daemon.Chsum) > 13 {
		bi.Daemon.Chsum = "***" + bi.Daemon.Chsum[len(bi.Daemon.Chsum)-13:]
	}

	js, err = json.MarshalIndent(&bi, "", "\t")
	if err != nil {
		log.Fatal("json.Marshal()", err.Error())
	}
	log.Println("\n"+string(js))
}

func getStatus(bi *BlockInfo) (bexec, bstake bool) {
	var cmd *exec.Cmd
	var err error
	var stdout io.ReadCloser
	var stderr io.ReadCloser
	var bterr []byte

	if bi == nil {
		log.Println("blockinfo is nil! cann't check the status")
		return false, false
	}
	log.Println("checking the status..")
	cmd = exec.Command(bi.Path + bi.Cli.Name, "getstakingstatus")
	stdout, err = cmd.StdoutPipe()
	if err != nil {
		log.Fatal("cmd.StdoutPipe():", err.Error())
	}
	stderr, err = cmd.StderrPipe()
	if err != nil {
		log.Fatal("cmd.StderrPipe():", err.Error())
	}
	if err = cmd.Start(); err != nil {
		log.Fatal("cmd.Start():", err.Error())
	}
	bterr, err = ioutil.ReadAll(stderr); if err != nil {
		log.Fatal("ioutil.ReadAll():", err.Error())
	}
	if len(bterr) == 0 {
		log.Println(bi.Daemon.Name, "run")
		bexec = true
		st := &StakingStatus{}
		if err = json.NewDecoder(stdout).Decode(&st); err != nil {
			log.Println("json.NewDecoder():", err.Error())
		} else {
			log.Println("staking", st.Staking)
			bstake = st.Staking
		}
		st = nil
	} else {
		log.Print("getstakingstatus:", string(bterr))
	}
	if err = cmd.Wait(); err != nil {
		log.Println(err.Error())
	}
	cmd = nil
	bterr = nil
	return bexec, bstake
}

func recoverDaemon(bi *BlockInfo) (result bool) {
	var cmd *exec.Cmd
	var err error
	if bi == nil {
		log.Println("blockinfo is nil! cann't start recovering")
		return false
	}
	log.Println("initiate recover..")
	cmd = exec.Command(bi.Path + bi.Daemon.Name, "-daemon")
	err = cmd.Start()
	if err != nil {
		log.Fatal("cmdDaemon.Start():", err.Error())
	}
	err = cmd.Wait()
	cmd = nil
	if err != nil {
		log.Println(bi.Daemon.Name, err.Error())
	} else {
		result = true
		log.Println(bi.Daemon.Name, "started")
		log.Println("waiting 30 seconds to load index..")
		<-time.After(time.Second * 30)
		log.Println("done!")
	}
	return result
}

func unlockWallet(bi *BlockInfo) (result bool) {
	var cmd *exec.Cmd
	var err error
	var stdout io.ReadCloser
	var stderr io.ReadCloser
	var bterr []byte
	if bi == nil {
		log.Println("blockinfo is nil! cann't start unlocking")
		return false
	}
	log.Println("trying to unlock the wallet..")
	cmd = exec.Command(bi.Path + bi.Cli.Name, "walletpassphrase", bi.Wpass, "63072000", "true")
	stdout, err = cmd.StdoutPipe()
	if err != nil {
		log.Fatal("cmd.StdoutPipe():", err.Error())
	}
	stderr, err = cmd.StderrPipe()
	if err != nil {
		log.Fatal("cmd.StderrPipe():", err.Error())
	}
	if err = cmd.Start(); err != nil {
		log.Fatal("cmd.Start():", err.Error())
	}
	bterr, err = ioutil.ReadAll(stderr); if err != nil {
		log.Fatal("ioutil.ReadAll():", err.Error())
	}
	if len(bterr) == 0 {
		_, err = ioutil.ReadAll(stdout); if err != nil {
			log.Fatal("ioutil.ReadAll():", err.Error())
		}
		log.Println("wallet unlocked!")
	} else {
		uer := &UnlockError{}
		bterr = bytes.TrimPrefix(bterr, []byte("error: "))
		if err := json.Unmarshal(bterr, &uer); err != nil {
			log.Println("json.NewDecoder():", err.Error())
		} else {
			log.Println(uer.Message)
			if uer.Code == -17 { result = true }
		}
		uer = nil
	}
	bterr = nil
	if err = cmd.Wait(); err != nil {
		log.Println(err.Error())
	}
	cmd = nil
	return result
}

func main() {
	var bi *BlockInfo
	var err error
	var executed, staking bool

	pathPtr := flag.String("path", "", "set path to " + C_CLI_FNAME + ", " + C_DAEMON_FNAME)
	wpwdPtr := flag.String("wpwd", "", "set blocknet wallet passphrase")
	spwdPtr := flag.String("spwd", "", "set additional storage passphrase (if set, must use it to execute too)")
	prntPtr := flag.Bool("print", false, "print data from storage")
	estkPtr := flag.Bool("estake", false, "execute staking with daemon recovering")
	flag.Parse()

	if *pathPtr != "" || *wpwdPtr != "" {
		saveInfo(*pathPtr, *wpwdPtr, *spwdPtr)
		debug.FreeOSMemory()
	}
	if *prntPtr {
		if bi = loadInfo(*spwdPtr); bi == nil { return }
                printInfo(*bi)
		bi = nil
		debug.FreeOSMemory()
	}
	if *estkPtr {
		for {
			if bi = loadInfo(*spwdPtr); bi == nil { return }

		        if err = checkCliFile(bi); err != nil { return }
		        executed, staking = getStatus(bi)
			if !executed {
			        if err = checkDaemonFile(bi); err != nil { return }
				executed = recoverDaemon(bi)
			}
			if !staking {
			        if err = checkCliFile(bi); err != nil { return }
				staking = unlockWallet(bi)
			}
			bi = nil
			debug.FreeOSMemory()
			if !staking {
				log.Println("waiting 5 minutes to confirm..")
				<-time.After(time.Minute * 5)
			} else {
				log.Println("waiting 30 minutes before new check..")
				<-time.After(time.Minute * 30)
			}
		}
	}
}
