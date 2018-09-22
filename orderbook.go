package bean

import (
	"fmt"
	"math"
	"sort"
	"time"
)

// OrderBook: an orderbook from exchange
type OrderBook struct {
	Bids []Order
	Asks []Order
}

// timed order book
type OrderBookT struct {
	Time time.Time
	OB   OrderBook
}

type OrderBookTS []OrderBookT

type Order struct {
	Price  float64
	Amount float64
}

func (ob OrderBook) Mid() float64 {
	if ob.Valid() {
		return (ob.Bids[0].Price + ob.Asks[0].Price) / 2.0
	} else {
		return math.NaN()
	}
}

func (ob OrderBook) Valid() bool {
	return len(ob.Bids) > 0 && len(ob.Asks) > 0
}

// filter out orders with amount less than the Coin minimum trading amount
// assuming ob is sorted
func Denoise(pair Pair, ob OrderBook) OrderBook {
	var bids []Order
	var asks []Order
	minimumAmount := pair.MinimumTradingAmount()
	for i, b := range ob.Bids {
		if b.Amount < minimumAmount {
			if i+1 < len(ob.Bids) {
				ob.Bids[i+1].Amount += b.Amount
			}
		} else {
			bids = append(bids, b)
		}
	}
	for i, a := range ob.Asks {
		if a.Amount < minimumAmount {
			if i+1 < len(ob.Asks) {
				ob.Asks[i+1].Amount += a.Amount
			}
		} else {
			asks = append(asks, a)
		}
	}
	return OrderBook{Bids: bids, Asks: asks}
}

func (ob OrderBook) Sort() OrderBook {
	// asks in ascending order
	sort.Slice(ob.Asks, func(i, j int) bool { return ob.Asks[i].Price < ob.Asks[j].Price })
	// bids in descending order
	sort.Slice(ob.Bids, func(i, j int) bool { return ob.Bids[i].Price > ob.Bids[j].Price })
	return ob
}

// orderbook display
func (obts OrderBookTS) ShowBrief() {
	for _, ob := range obts {
		if ob.OB.Valid() {
			fmt.Println(ob.Time.Local().Format("Jan _2 15:04:05"), "depth:", len(ob.OB.Asks), "bestBid:", ob.OB.Bids[0].Price, "bestAsk:", ob.OB.Asks[0].Price)
		} else {
			fmt.Println(ob.Time.Local(), len(ob.OB.Asks))
		}
	}
}

func (obts OrderBookTS) Sort() OrderBookTS {
	sort.Slice(obts, func(i, j int) bool { return obts[i].Time.Before(obts[j].Time) })
	return obts
}

////////////////////////////////////////////////////////
// FIXME: move below functions to other module

type OrderState string

const (
	ALIVE     OrderState = "ALIVE"
	FILLED    OrderState = "FILLED"
	CANCELLED OrderState = "CANCELLED"
)

// Status of the placed order,
type OrderStatus struct {
	FilledAmount float64
	LeftAmount   float64
	PlacedPrice  float64 // initial price
	Price        float64 // filled price, if not applicable then placed price
	State        OrderState
}
