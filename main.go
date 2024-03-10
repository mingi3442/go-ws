package main

import (
  "context"
  "fmt"
  "time"

  rpchttp "github.com/cometbft/cometbft/rpc/client/http"
  "github.com/cometbft/cometbft/types"
)

func main() {
  client, err := rpchttp.New("https://cosmos-rpc.polkachu.com", "/websocket")
  if err != nil {
    fmt.Println(err)
  }
  err = client.Start()
  if err != nil {
    // handle error
    fmt.Println(err)
    return
  }
  defer client.Stop()

  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
  defer cancel()

  query := "tm.event = 'Tx' AND tx.height = 3"
  txs, err := client.Subscribe(ctx, "test-client", query)
  if err != nil {
    // handle error
    fmt.Println(err)
  }

  for e := range txs {
    fmt.Println("got ", e.Data.(types.EventDataTx))
  }

}
