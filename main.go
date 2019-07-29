package main

import (
	"context"
	"fmt"
	"os"
	"time"

	ipfs "github.com/ipfs/go-ipfs-api"
)

const (
	PEER_ADDR     = "/ip4/140.116.245.242/tcp/4001/ipfs/QmVti8ZHZ2o1SbieAVz7o2Qzxj8RCanEahGUM1mnajaQTi"
	REMOTE_FOLDER = "QmeJudURYGSJADGHewzbuodKpfN9x1RLrwxNEeX7Vqamrb"
	LOCAL_FOLDER  = "./local_folder"
)

func SwarmConnect(sh *ipfs.Shell, addr ...string) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)

	defer cancel()

	err := sh.SwarmConnect(ctx, addr...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err)
		os.Exit(1)
	}
}

func List(sh *ipfs.Shell, path string) {
	lsLink, err := sh.List(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err)
		os.Exit(1)
	}

	for _, v := range lsLink {
		fmt.Println(v.Hash, v.Name)
	}
}

func AddDir(sh *ipfs.Shell, dir string) (string, error) {
	return sh.AddDir(dir)
}

func main() {
	sh := ipfs.NewLocalShell()

	SwarmConnect(sh, PEER_ADDR)
	List(sh, REMOTE_FOLDER)

	hash, err := AddDir(sh, LOCAL_FOLDER)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err)
		os.Exit(1)
	}

	List(sh, hash)
}
