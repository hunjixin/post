package main

import (
	"context"
	"crypto/rand"
	"flag"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/spacemeshos/post/config"
	"github.com/spacemeshos/post/initialization"
	"github.com/spacemeshos/post/proving"
	"github.com/spacemeshos/post/shared"
	zap "go.uber.org/zap"
)

var dataDir string
var nonces uint
var thread uint

func main() {
	flag.StringVar(&dataDir, "dataDir", "", "data directory")
	flag.UintVar(&nonces, "nonces", 288, "")
	flag.UintVar(&thread, "thread", 0, "")
	flag.Parse()

	fmt.Println(run(context.Background()))
}

func run(ctx context.Context) error {
	dirs := strings.Split(dataDir, ":")
	wg := sync.WaitGroup{}
	for _, dir := range dirs {
		wg.Add(1)
		go func(dir string) {
			startT := time.Now()
			defer wg.Done()

			meta, err := initialization.LoadMetadata(dir)
			if err != nil {
				log.Println("failed to load metadata", dir, err)
				return
			}

			myLog := zap.L()
			mainnetCfg := config.MainnetConfig()
			var challenge shared.Challenge
			rand.Read(challenge)
			proof, _, err := proving.Generate(ctx, challenge, mainnetCfg, myLog, proving.WithDataSource(mainnetCfg, meta.NodeId, meta.CommitmentAtxId, dir), proving.WithNonces(nonces), proving.WithThreads(10))
			if err != nil {
				log.Fatalln("proof generation error", err)
				return
			}
			fmt.Println(proof)
			fmt.Println("%s time %s", dir, time.Now().Sub(startT).String())
		}(dir)
	}
	wg.Wait()
	return nil
}
