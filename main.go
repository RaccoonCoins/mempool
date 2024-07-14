package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	_ "github.com/lib/pq"
	"github.com/tonkeeper/tonapi-go"
	"github.com/tonkeeper/tonapi-go/examples/sse/db"
)

func intPointer(x int) *int {
	return &x
}

const TokenAddress = "0:2212fce3578671c9c4c8577a658a5c444d071dbeb9ce47b9850ef0047be9e342"

func printAccountInformation(myAccount string) error {
	client, err := tonapi.New()
	if err != nil {
		return err
	}

	holder, err := client.GetJettonInfo(context.Background(), tonapi.GetJettonInfoParams{AccountID: myAccount})
	if err != nil {
		return err
	}
	fmt.Printf("Mintable : %v\n", holder.Mintable)
	fmt.Printf("TotalSupply : %v\n", holder.TotalSupply)
	fmt.Printf("holder: %v\n", holder.HoldersCount)

	return nil
}

func subscribeToMempool(token string, apiTon string) {

	streamingAPI := tonapi.NewStreamingAPI(tonapi.WithStreamingToken(apiTon))
	for {
		err := streamingAPI.SubscribeToMempool(context.Background(),
			[]string{token},
			func(data tonapi.MempoolEventData) {
				value, _ := json.Marshal(data.BOC)
				//decodedBytes, err1 := hex.DecodeString(string(value))
				fmt.Printf("mempool event: %#v\n", string(value))

			})
		if err != nil {
			fmt.Printf("mempool error: %v, reconnecting...\n", err)

		}
	}
}

func subscribeToTransactions(token string, apiTon string) {

	streamingAPI := tonapi.NewStreamingAPI(tonapi.WithStreamingToken(apiTon))
	for {
		err := streamingAPI.SubscribeToTransactions(context.Background(),
			[]string{token},
			nil,
			func(data tonapi.TransactionEventData) {
				fmt.Printf("New tx with hash: %v\n", data.TxHash)
			})
		if err != nil {
			fmt.Printf("tx error: %v, reconnecting...\n", err)
			time.Sleep(5 * time.Second) // wait for 5 seconds before retrying
		}
	}
}

func GetBlockchainTransaction(myTransaction string) error {
	client, err := tonapi.New()
	if err != nil {
		return err
	}
	transaction_acc, err := client.GetBlockchainTransaction(context.Background(), tonapi.GetBlockchainTransactionParams{TransactionID: myTransaction})
	if err != nil {
		return err
	}
	fmt.Printf("transaction  hash : %v\n", transaction_acc.Hash)
	fmt.Printf("transaction IN_Msgs WAllet : %v\n", transaction_acc.InMsg.Value.Source.Value.IsWallet)
	if transaction_acc.InMsg.Value.Source.Value.IsWallet {

		fmt.Printf("Gia tri Ton mua token : %v\n", transaction_acc.InMsg.Value.Value/1000000000)
		//fmt.Printf("transaction out_Msgs : %v\n", transaction_acc.OutMsgs[0].DecodedBody)
		//fmt.Printf("transaction Out_Msgs : %s\n", transaction_acc.OutMsgs)

		var data map[string]interface{}
		if err := json.Unmarshal([]byte(transaction_acc.OutMsgs[0].DecodedBody), &data); err != nil {
			fmt.Println("Error parsing JSON:", err)
			return nil
		}

		// Truy cập giá trị min_out
		minOut, ok := data["amount"].(string)
		if !ok {
			fmt.Println("Không thể truy cập giá trị min_out")
			return nil
		}
		// Chuyển chuỗi thành số thực
		minOutInt, err := strconv.Atoi(minOut)
		if err != nil {
			fmt.Println("Lỗi khi chuyển đổi giá trị minOut sang số nguyên:", err)
			return nil
		}
		result := minOutInt / 1000000000
		fmt.Println("Mua Token :", result)
		// Chuyển đổi số nguyên thành chuỗi
		//myStr_amount_token := strconv.Itoa(result)
		//fmt.Println("kkkkkkkkk", result)
		//go gobot(myStr_amount_token)

		//go ketnoidata(transaction_acc.Hash, result)

	} else {

		fmt.Println("Khong phai lenh mua la lenh ban")
	}

	return nil
}

func subscribeToTraces(token string, apiTon string) {
	streamingAPI := tonapi.NewStreamingAPI(tonapi.WithStreamingToken(apiTon))
	for {
		err := streamingAPI.SubscribeToTraces(context.Background(), []string{token},
			func(data tonapi.TraceEventData) {

				fmt.Printf("New trace with AccountIDs: %v\n", data)

				// Chuyển đổi data.AccountIDs thành chuỗi
				/*var accountIDsStr string
				for _, id := range data.AccountIDs {
					accountIDsStr += id.String() + " " // Giả sử ton.AccountID có phương thức String()
				}

				// Cắt chuỗi sau mỗi dấu khoảng trắng
				parts := strings.Split(accountIDsStr, " ")

				for _, part := range parts {
					fmt.Println(part)
				}*/

				fmt.Printf("New trace with hash: %v\n", data.Hash)
				ketnoidata(data.Hash)
			})
		if err != nil {
			fmt.Printf("trace error: %v, reconnecting...\n", err)
		}
	}
}

func subscribeToBlocks(token string, apiTon string) {
	streamingAPI := tonapi.NewStreamingAPI(tonapi.WithStreamingToken(apiTon))
	for {
		err := streamingAPI.SubscribeToBlocks(context.Background(), intPointer(-1),
			func(data tonapi.BlockEventData) {
				fmt.Printf("New block: (%v,%v,%v)\n", data.Workchain, data.Shard, data.Seqno)

			})
		if err != nil {
			fmt.Printf("block error: %v, reconnecting...\n", err)
		}
	}
}
func ketnoidata(txname string) {

	sql := &db.Sql{

		Host:     "dpg-cq811988fa8c738b8us0-a",
		Port:     5432,
		Username: "robotdog",
		Password: "R8WI48T1yxmwlec0WSdyLc3e7zphfzDi",
		Dbname:   "golang",
	}

	sql.Connect()

	defer sql.Close()

	fmt.Println("Successfully connected!")
	err := sql.InsertData(txname)

	if err != nil {
		fmt.Println("Failed to insert data:", err)
	} else {
		fmt.Println("Data inserted successfully!")
	}

}
func main() {

	token := "0:2212fce3578671c9c4c8577a658a5c444d071dbeb9ce47b9850ef0047be9e342,0:948673a596fc6d2c3123dfff64d3231e4d0e16bfaac506d5ebdb09f32b3b9c65"
	apiTon := "AF2BN4CF7W243KAAAAANNF3BE5Z2UDOVWVWGECV7BX3V43RJO6AUTQXPJJH23YDJFKKAHEY"
	//go subscribeToTransactions(token, apiTon)
	go subscribeToTraces(token, apiTon)
	//go subscribeToMempool(token, apiTon)

	go subscribeToBlocks(token, apiTon)
	//go printAccountInformation(myAccount)
	//go GetBlockchainTransaction("e08ddf584a8908f1df4c379d3ce2c650eddd09aaa3f26d6798ffd9e0355be225")
	//go gobot()

	//gobot("asasas")
	//go ketnoidata("haimeo", 9999999)
	select {}
}
