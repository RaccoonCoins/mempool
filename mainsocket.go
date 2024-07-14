package main

import (
	"context"
	"fmt"
	"log"

	"github.com/tonkeeper/tonapi-go"
)

func main() {

	accounts := []string{"0:f51f88aa88ecea8034ced5b94a89d31b9d9904fb8bb3cec2e47fc468b0a97c00"}

	// When working with tonapi.io, you should consider getting an API key at https://tonconsole.com/
	// because tonapi.io has per-ip limits for sse and websocket connections.
	//
	token := "AF2BN4CF7W243KAAAAANNF3BE5Z2UDOVWVWGECV7BX3V43RJO6AUTQXPJJH23YDJFKKAHEY"

	streamingAPI := tonapi.NewStreamingAPI(tonapi.WithStreamingToken(token))

	err := streamingAPI.WebsocketHandleRequests(context.Background(), func(ws tonapi.Websocket) error {
		ws.SetMempoolHandler(func(data tonapi.MempoolEventData) {
			fmt.Printf("new mempool event\n")
		})
		ws.SetTransactionHandler(func(data tonapi.TransactionEventData) {
			fmt.Printf("New tx with hash: %v\n", data.TxHash)
		})
		ws.SetTraceHandler(func(data tonapi.TraceEventData) {
			fmt.Printf("New trace with hash: %v\n", data.Hash)
		})
		ws.SetBlockHandler(func(data tonapi.BlockEventData) {
			fmt.Printf("New block: (%v,%v,%v)\n", data.Workchain, data.Shard, data.Seqno)
		})

		if err := ws.SubscribeToMempool(nil); err != nil {

			return err
		}
		if err := ws.SubscribeToTransactions(accounts, nil); err != nil {

			return err
		}
		if err := ws.SubscribeToTraces(accounts); err != nil {
			return err
		}
		masterchain := -1
		if err := ws.SubscribeToBlocks(&masterchain); err != nil {
			return err
		}
		// It is possible to run a loop updating subscription on the go:
		//
		// subscribeCh := make(chan []string) // channel to send accounts to subscribe.
		// for {
		// 	select {
		//	case accounts := <-subscribeCh:
		//		if err := ws.SubscribeToTransactions(accounts); err != nil {
		//			return err
		//		}
		//		if err := ws.SubscribeToTraces(accounts); err != nil {
		//			return err
		//		}
		//	case <-ctx.Done():
		//		return nil
		//	}
		//}
		return nil
	})
	if err != nil {
		log.Fatalf("connection error: %v", err)
	}
}
