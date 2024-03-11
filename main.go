package main

import (
  "context"
  "encoding/json"
  "fmt"
  "io"

  "net/http"
  "time"

  rpchttp "github.com/cometbft/cometbft/rpc/client/http"
  "github.com/cometbft/cometbft/types"
)

type StatusResponse struct {
  Jsonrpc string `json:"jsonrpc"`
  ID      int    `json:"id"`
  Result  struct {
    NodeInfo struct {
      ProtocolVersion struct {
        P2P   string `json:"p2p"`
        Block string `json:"block"`
        App   string `json:"app"`
      } `json:"protocol_version"`
      ID         string `json:"id"`
      ListenAddr string `json:"listen_addr"`
      Network    string `json:"network"`
      Version    string `json:"version"`
      Channels   string `json:"channels"`
      Moniker    string `json:"moniker"`
      Other      struct {
        TxIndex    string `json:"tx_index"`
        RPCAddress string `json:"rpc_address"`
      } `json:"other"`
    } `json:"node_info"`
    SyncInfo struct {
      LatestBlockHash     string `json:"latest_block_hash"`
      LatestAppHash       string `json:"latest_app_hash"`
      LatestBlockHeight   string `json:"latest_block_height"`
      LatestBlockTime     string `json:"latest_block_time"`
      EarliestBlockHash   string `json:"earliest_block_hash"`
      EarliestAppHash     string `json:"earliest_app_hash"`
      EarliestBlockHeight string `json:"earliest_block_height"`
      EarliestBlockTime   string `json:"earliest_block_time"`
      CatchingUp          bool   `json:"catching_up"`
    } `json:"sync_info"`
    ValidatorInfo struct {
      Address string `json:"address"`
      PubKey  struct {
        Type  string `json:"type"`
        Value string `json:"value"`
      } `json:"pub_key"`
      VotingPower string `json:"voting_power"`
    } `json:"validator_info"`
  } `json:"result"`
}

type ChainStatus struct {
  Network string
  Moniker string
}

func fetchChainStatus() ChainStatus {
  resp, err := http.Get("https://cosmos-rpc.polkachu.com/status")
  if err != nil {
    fmt.Printf("%s\n", err)
  }

  defer resp.Body.Close()

  body, err := io.ReadAll(resp.Body)
  if err != nil {
    fmt.Printf("%s\n", err)
  }

  var status StatusResponse
  if err := json.Unmarshal(body, &status); err != nil {
    fmt.Printf("%s\n", err)
  }

  return ChainStatus{
    Network: status.Result.NodeInfo.Network,
    Moniker: status.Result.NodeInfo.Moniker,
  }
}

func main() {
  status := fetchChainStatus()
  client, err := rpchttp.New("https://cosmos-rpc.polkachu.com", "/websocket")
  if err != nil {
    fmt.Println(err)
  }

  err = client.Start()

  if err != nil {
    fmt.Println(err)
    return
  }
  defer client.Stop()

  ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
  defer cancel()

  query := "tm.event = 'NewBlock'"
  txs, err := client.Subscribe(ctx, "test-client", query)
  if err != nil {
    // handle error
    fmt.Println(err)
  }
  // fmt.Println(txs)

  // for e := range txs {
  //   // fmt.Println("got ", e.Data.(types.EventDataTx))
  //   eventDataTx, _ := e.Data.(types.EventDataTx)
  //   fmt.Println((eventDataTx))

  // }
  for e := range txs {
    eventDataNewBlock, ok := e.Data.(types.EventDataNewBlock)
    if !ok {
      fmt.Println("Event data is not of type NewBlock")
      continue
    }

    // jsonData, err := json.MarshalIndent(eventDataTx, "", "  ")
    // if err != nil {
    //   fmt.Printf("%s\n", err)
    //   continue
    // }
    fmt.Println("Moniker:", status.Moniker)
    fmt.Println("Network:", status.Network)

    fmt.Printf("Block Height: %d\n", eventDataNewBlock.Block.Height)

  }

}
