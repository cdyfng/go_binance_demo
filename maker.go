package main

import (
	"fmt"
	//"testing"
	"strconv"
	time2 "time"
	"os"
	//"errors"
	//"github.com/stretchr/testify/assert"

	sdk "github.com/binance-chain/go-sdk/client"
	"github.com/binance-chain/go-sdk/client/query"
	//"github.com/binance-chain/go-sdk/common"
	"github.com/binance-chain/go-sdk/keys"
	"github.com/binance-chain/go-sdk/types"
	"github.com/binance-chain/go-sdk/types/msg"
	"github.com/cdyfng/go_binance_demo/util"
	//tx2 "github.com/binance-chain/go-sdk/types/tx"
	"github.com/binance-chain/go-sdk/client/transaction"
	"math/rand"
)

type order struct {
	Id         int64
	Symbol     string
	OpenTime   time2.Time
	Side       string
	Price      float64
	Quantity   float64
	StartDepth []string
}

func initDex() (sdk.DexClient, error) {
	priv := readBinance_private()
	keyManager, _ := keys.NewPrivateKeyManager(priv)
	testAccount1 := keyManager.GetAddr()
	fmt.Printf("testAccount1: %v \n", testAccount1)
	//-----   Init sdk  -------------
	return sdk.NewDexClient("https://testnet-dex.binance.org", types.TestNetwork, keyManager)
}

func isTradeable(ct sdk.DexClient, tradeSymbol string, nativeSymbol string) bool {
	//return false
	dq := query.NewDepthQuery(tradeSymbol, nativeSymbol)
	//fmt.Printf("Depth Query: %v \n", dq)
	//fmt.Printf("orders: %v \n", len(orders))
	fmt.Printf("ct: %v \n", ct)
	depth, err := ct.GetDepth(dq)
	if err != nil {
		fmt.Printf("depth 0 : %v \n", depth)
		return false
	}
	//assert.NoError(t, err)
	//asks := depth.Asks
	//bids := depth.Bids
	//disp :=
	fmt.Printf("depth 0 : %v \n", [...][]string{depth.Asks[0], depth.Bids[0], depth.Asks[1], depth.Bids[1]})
	//assert.True(t, depth.Height > 0)
	//fmt.Printf("tradeSymbol: %v nativeSymbol: %v\n", tradeSymbol, nativeSymbol)
	//fmt.Printf("nativeSymbol: %v \n", nativeSymbol)

	a, _ := strconv.ParseFloat(depth.Asks[0][0], 32)
	b, _ := strconv.ParseFloat(depth.Bids[0][0], 32)
	spread := a - b
	spread_percent := spread / a * 100
	fmt.Printf("spread: %v spread_percent: %v %%\n", spread, spread_percent)
	if spread_percent >= 0.001 {
		return true
	} else {
		return false
	}
}

func getDepth(ct sdk.DexClient, tradeSymbol string, nativeSymbol string) (map[string]string, error) {
	dq := query.NewDepthQuery(tradeSymbol, nativeSymbol)
	depth, error := ct.GetDepth(dq)
	//fmt.Printf("d: %v, %v", depth, error)
	if error == nil {
		return map[string]string{
			"ASK0": depth.Asks[0][0],
			"BID0": depth.Bids[0][0],
		}, nil
	}
	return nil, error
}

func main() {

	client, _ := initDex()
	fmt.Printf("client : %v \n", client)

	tradeSymbol := "BNB"
	nativeSymbol := "BTC.B-918"
	//var orders [2]order
	//assert.Equal(t, 1, len(markets))

	//-----  Get Depth  ----------
	// dq := query.NewDepthQuery(tradeSymbol, nativeSymbol)
	// fmt.Printf("Depth Query: %v \n", dq)

	side := "Buy"
	count := 1
	for {
		if rand.Intn(2)%2 == 0 {
				side = "Buy"
			}else{
				side = "Sell"
			}
		fmt.Printf("loop : %v side: %v\n", count, side)
		//得到depth
		//是否满足下单条件，不满足，退出延时
		//满足则根据条件下单          进入子进程 判断价格是否有效，有效等待，否则取消

		if isTradeable(client, tradeSymbol, nativeSymbol) {
			time2.Sleep(500 * time2.Millisecond)
			d, error := getDepth(client, tradeSymbol, nativeSymbol)
			if error != nil {
				time2.Sleep(500 * time2.Millisecond)
				continue
			}
			var createOrderResult *transaction.CreateOrderResult
			if side == "Sell" {
				fmt.Printf("Do trade sell \n")
				a, _ := strconv.ParseFloat(d["ASK0"], 64)
				myAsk := util.Times8(a) - 1
				fmt.Printf("myBid: %v \n", myAsk)
				createOrderResult, error = client.CreateOrder(tradeSymbol, nativeSymbol, msg.OrderSide.SELL, myAsk, 500000000, true)
			} else {
				fmt.Printf("Do trade buy \n")
				b, _ := strconv.ParseFloat(d["BID0"], 64)
				myBid := util.Times8(b) + 1
				fmt.Printf("myBid: %v \n", myBid)
				createOrderResult, error = client.CreateOrder(tradeSymbol, nativeSymbol, msg.OrderSide.BUY, myBid, 500000000, true)
			}

			fmt.Printf("createOrderResult: %v \n", createOrderResult)
			fmt.Printf("err: %v \n", error)
			if error != nil {
				continue
			}

			i := 0
			for {
				i++
				time2.Sleep(500 * time2.Millisecond)
				curOrder, err := client.GetOrder(createOrderResult.OrderId)
				if err != nil {
					fmt.Printf("break GetOrder: %v \n", curOrder)
					break
				}
				fmt.Printf("%v order: %v\n", i, curOrder)
				if curOrder.Status == "FullyFill" {
					fmt.Printf("stat: %v %v\n", query.OrderStatus.FULLY_FILLED, curOrder.Status)
					break
				}

				if curOrder.Status == "Ack" {
					time2.Sleep(500 * time2.Millisecond)
					d, error := getDepth(client, tradeSymbol, nativeSymbol)
					if error != nil {
						break
					}
					var cur_price string
					if side == "Sell" {
						cur_price = d["ASK0"]
					} else {
						cur_price = d["BID0"]
					}
					fmt.Printf(" %s %s to compare\n", cur_price, curOrder.Price)
					if cur_price != curOrder.Price {
						fmt.Printf(" %s %s not same\n", cur_price, curOrder.Price)
						time2.Sleep(500 * time2.Millisecond)
						cancleOrderResult, _ := client.CancelOrder(tradeSymbol, nativeSymbol, createOrderResult.OrderId, true)
						fmt.Printf("cancleOrderResult:  %v \n", cancleOrderResult)
						break
					}
				}
			}
		}

		time2.Sleep(1 * time2.Second)
		count++
		if count >= 1000000 {
			break
		}
	}
}



func readBinance_private() string {
	envName := "BINANCE_ACCOUNT_PRIVATE"
	account_1 := os.Getenv(envName)
	if account_1 == "" {
		panic(fmt.Errorf("BINANCE_ACCOUNT_PRIVATE environment variable %q must be set", envName))
	}

	return account_1
}
