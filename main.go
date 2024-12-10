package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/0xjeffro/tx-parser/solana"
	"github.com/0xjeffro/tx-parser/solana/types"
	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
)

func main() {
	main2()
	// r := gin.Default()
	// r.POST("/solana", solanaHandler)
	// err := r.Run()
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }
}

type RpcRsp struct {
	JsonRpc string      `json:"jsonrpc"`
	Result  types.RawTx `json:"result"`
	Id      int         `json:"id"`
}

func main2() {
	// Fetch a raw transaction from the Solana RPC
	// url := "https://solana-rpc.publicnode.com"
	url := "https://api.mainnet-beta.solana.com"
	method := "POST"
	tx := "5RRSy3gqTBDsnWkEN4uj467vfn4geVbgef2b7wZCMhbaK5KTZdwhTPchzTtimcTnBTvbA8zbCaFo2sFyxvwDD33g"

	payload := strings.NewReader(fmt.Sprintf(`{
		"jsonrpc": "2.0",
		"id": 1,
		"method": "getTransaction",
		"params": [
		  "%s",
		  {
			"encoding": "json",
			"maxSupportedTransactionVersion": 0
		  }
		]
	}`, tx))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	// Parse the raw transaction
	var result RpcRsp
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println(err)
		return
	}
	// spew.Dump(result)
	parsed := solana.TxParser(result.Result)
	for _, v := range parsed.Actions {
		spew.Dump(v)

	}

	// parsedJson, err := json.Marshal(parsed)

	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// fmt.Println(string(parsedJson))
}

func solanaHandler(c *gin.Context) {
	bytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		Fail(c, err)
		return
	}

	res, err := solana.Parser(bytes)
	if err != nil {
		Fail(c, err)
		return
	}
	Success(c, res)
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(200, gin.H{
		"message": "success",
		"data":    data,
	})
}

func Fail(c *gin.Context, err error) {
	c.JSON(400, gin.H{
		"message": "error",
		"error":   err.Error(),
	})
}
