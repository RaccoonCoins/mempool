package mainold

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	_ "github.com/lib/pq"
	"github.com/tonkeeper/tonapi-go"
	"github.com/tonkeeper/tonapi-go/examples/sse/db"
)

func intPointer(x int) *int {
	return &x
}

const myAccount = "EQAEksSxKrjw9qav1tgEM5AKYkj-8XgmiPczEVyNUDYdaFNE"

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

func subscribeToMempool(token string) {

	streamingAPI := tonapi.NewStreamingAPI(tonapi.WithStreamingToken("AF2BN4CF7W243KAAAAANNF3BE5Z2UDOVWVWGECV7BX3V43RJO6AUTQXPJJH23YDJFKKAHEY"))
	//streamingAPI := tonapi.NewStreamingAPI(tonapi.WithStreamingToken(token))
	for {
		err := streamingAPI.SubscribeToMempool(context.Background(),
			[]string{"0:fb86103da27e2a43732988af8b2275d8a637549c7020be7f547bf4c784a4d9d3"},
			func(data tonapi.MempoolEventData) {
				//value, _ := json.Marshal(data)
				value, _ := json.Marshal(data)
				//decodedBytes, err1 := hex.DecodeString(string(value))
				fmt.Printf("mempool event: %#v\n", string(value))

			})
		if err != nil {
			fmt.Printf("mempool error: %v, reconnecting...\n", err)
		}
	}
}

func subscribeToTransactions(token string) {

	streamingAPI := tonapi.NewStreamingAPI(tonapi.WithStreamingToken("AF2BN4CF7W243KAAAAANNF3BE5Z2UDOVWVWGECV7BX3V43RJO6AUTQXPJJH23YDJFKKAHEY"))
	//streamingAPI := tonapi.NewStreamingAPI(tonapi.WithStreamingToken(token))

	for {
		err := streamingAPI.SubscribeToTransactions(context.Background(),
			[]string{"0:fb86103da27e2a43732988af8b2275d8a637549c7020be7f547bf4c784a4d9d3"},
			nil,
			func(data tonapi.TransactionEventData) {

				fmt.Printf("New tx with hash: %v\n", data.TxHash)
				//GetBlockchainTransaction(data.TxHash)
				ketnoidata(data.TxHash, 11111111)

			})
		if err != nil {
			fmt.Printf("tx error: %v, reconnecting...\n", err)
		}
	}
}

func GetBlockchainTransaction(myTransaction string) error {
	client, err := tonapi.New()
	if err != nil {
		return err
	}
	//account, err := client.GetAccount(context.Background(), tonapi.GetAccountParams{AccountID: myAccount})
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

func subscribeToTraces(token string) {
	streamingAPI := tonapi.NewStreamingAPI(tonapi.WithStreamingToken("AF2BN4CF7W243KAAAAANNF3BE5Z2UDOVWVWGECV7BX3V43RJO6AUTQXPJJH23YDJFKKAHEY"))
	//streamingAPI := tonapi.NewStreamingAPI(tonapi.WithStreamingToken(token))
	for {
		err := streamingAPI.SubscribeToTraces(context.Background(), []string{"0:fb86103da27e2a43732988af8b2275d8a637549c7020be7f547bf4c784a4d9d3"},
			func(data tonapi.TraceEventData) {

				fmt.Printf("New trace with hash: %v\n", data.Hash)
			})
		if err != nil {
			fmt.Printf("trace error: %v, reconnecting...\n", err)
		}
	}
}

func subscribeToBlocks(token string) {
	streamingAPI := tonapi.NewStreamingAPI(tonapi.WithStreamingToken("AF2BN4CF7W243KAAAAANNF3BE5Z2UDOVWVWGECV7BX3V43RJO6AUTQXPJJH23YDJFKKAHEY"))
	//streamingAPI := tonapi.NewStreamingAPI(tonapi.WithStreamingToken(token))
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
func ketnoidata(txname string, price int) {
	sql := &db.Sql{

		Host:     "localhost",
		Port:     5432,
		Username: "postgres",
		Password: "12345678",
		Dbname:   "golang",
	}

	sql.Connect()

	defer sql.Close()

	fmt.Println("Successfully connected!")
	sql.InsertData(txname, price)

}
func main() {

	token := "0:fb86103da27e2a43732988af8b2275d8a637549c7020be7f547bf4c784a4d9d3"

	go subscribeToTraces(token)
	go subscribeToMempool(token)
	go subscribeToTransactions(token)
	go subscribeToBlocks(token)
	go printAccountInformation(myAccount)
	//go GetBlockchainTransaction("e08ddf584a8908f1df4c379d3ce2c650eddd09aaa3f26d6798ffd9e0355be225")
	//go gobot()

	//gobot("asasas")
	//go ketnoidata("haimeo", 6666666)
	select {}
}
