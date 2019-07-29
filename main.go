package main

import (
	"context"
	"fmt"
	"os"

	ipfs "github.com/ipfs/go-ipfs-api"
)

const PEER_ADDR = "/ip4/140.116.245.242/tcp/4001/ipfs/QmVti8ZHZ2o1SbieAVz7o2Qzxj8RCanEahGUM1mnajaQTi"

func main() {
	fmt.Println("rehentai-content-server")

	ctx := context.Background()

	sh := ipfs.NewLocalShell()

	err := sh.SwarmConnect(ctx, PEER_ADDR)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err)
		os.Exit(1)
	}

	lsLink, err := sh.List("QmeJudURYGSJADGHewzbuodKpfN9x1RLrwxNEeX7Vqamrb")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err)
		os.Exit(1)
	}

	for _, v := range lsLink {
		fmt.Println(v.Hash, v.Name)
	}
}
