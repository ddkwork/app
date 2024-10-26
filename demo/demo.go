package main

import (
	"strconv"
	"time"

	"github.com/ddkwork/app/widget"
	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/unison"
	"github.com/ddkwork/unison/app"
)

func main() {
	app.Run("", func(w *unison.Window) {
		w.Content().AddChild(layout())
	})
}

type Data struct {
	Time time.Time
	// 收支类型 string
	Trade  TradeKind
	Seller string  // 商家姓名
	Name   string  // 挂礼人
	Amount float64 // 金额
	Note   string  // 备注
}

func layout() unison.Paneler {
	return widget.NewTableScroll(Data{}, widget.TableContext[Data]{
		ContextMenuItems: nil,
		MarshalRow: func(node *widget.Node[Data]) (cells []widget.CellData) {
			timeFmt := node.Data.Time.Format("2006-01-02")
			sumFmt := strconv.FormatFloat(node.Data.Amount, 'f', 2, 64)
			trade := node.Data.Trade.Tooltip()
			sum := 0.00
			if node.Container() {
				timeFmt = node.Sum()
				node.Walk(func(node *widget.Node[Data]) {
					sum += node.Data.Amount
				})
				sumFmt = strconv.FormatFloat(sum, 'f', 2, 64)
				trade = ""
			}
			return []widget.CellData{
				{Text: timeFmt},
				{Text: trade},
				{Text: node.Data.Seller},
				{Text: node.Data.Name},
				{Text: sumFmt},
				{Text: node.Data.Note},
			}
		},
		UnmarshalRow: func(node *widget.Node[Data], values []string) {
			node.Data = Data{
				Time:   time.Time{},
				Trade:  0,
				Seller: "",
				Name:   "",
				Amount: 0,
				Note:   "",
			}
		},
		SelectionChangedCallback: func(root *widget.Node[Data]) {
		},
		SetRootRowsCallBack: func(root *widget.Node[Data]) {
			containerSum := widget.NewContainerNode("all", Data{
				Time: mylog.Check2(time.Parse("2006-01-02", "2024-08-08")),
				// Trade:  MeatKind,
				// Seller: "杨学春",
				// Name:   "",
				// Amount: 19.4 * 45,
				// Note:   "二流",
			})

			container := widget.NewContainerNode("8号", Data{
				Time: mylog.Check2(time.Parse("2006-01-02", "2024-08-08")),
				// Trade:  MeatKind,
				// Seller: "杨学春",
				// Name:   "",
				// Amount: 19.4 * 45,
				// Note:   "二流",
			})
			container.AddChildByData(Data{
				Time:   mylog.Check2(time.Parse("2006-01-02", "2024-08-08")),
				Trade:  OtherKind,
				Seller: "",
				Name:   "",
				Amount: -16,
				Note:   "",
			})
			container.AddChildByData(Data{
				Time:   mylog.Check2(time.Parse("2006-01-02", "2024-08-08")),
				Trade:  OtherKind,
				Seller: "",
				Name:   "",
				Amount: -29,
				Note:   "",
			})
			container.AddChildByData(Data{
				Time:   mylog.Check2(time.Parse("2006-01-02", "2024-08-08")),
				Trade:  OtherKind,
				Seller: "",
				Name:   "",
				Amount: -20,
				Note:   "",
			})

			container.AddChildByData(Data{
				Time:   mylog.Check2(time.Parse("2006-01-02", "2024-08-08")),
				Trade:  OtherKind,
				Seller: "",
				Name:   "",
				Amount: -10,
				Note:   "",
			})

			container.AddChildByData(Data{
				Time:   mylog.Check2(time.Parse("2006-01-02", "2024-08-08")),
				Trade:  OtherKind,
				Seller: "",
				Name:   "",
				Amount: -5,
				Note:   "",
			})
			container.AddChildByData(Data{
				Time:   mylog.Check2(time.Parse("2006-01-02", "2024-08-08")),
				Trade:  OtherKind,
				Seller: "",
				Name:   "",
				Amount: -10,
				Note:   "",
			})
			container.AddChildByData(Data{
				Time:   mylog.Check2(time.Parse("2006-01-02", "2024-08-08")),
				Trade:  OtherKind,
				Seller: "",
				Name:   "",
				Amount: -40,
				Note:   "",
			})
			container.AddChildByData(Data{
				Time:   mylog.Check2(time.Parse("2006-01-02", "2024-08-08")),
				Trade:  OtherKind,
				Seller: "",
				Name:   "",
				Amount: -32,
				Note:   "",
			})
			container.AddChildByData(Data{
				Time:   mylog.Check2(time.Parse("2006-01-02", "2024-08-08")),
				Trade:  OtherKind,
				Seller: "",
				Name:   "",
				Amount: -500,
				Note:   "取款",
			})
			container.AddChildByData(Data{
				Time:   mylog.Check2(time.Parse("2006-01-02", "2024-08-08")),
				Trade:  OtherKind,
				Seller: "",
				Name:   "杨国科",
				Amount: -500,
				Note:   "挂礼",
			})

			container.AddChildByData(Data{
				Time:   mylog.Check2(time.Parse("2006-01-02", "2024-08-08")),
				Trade:  OtherKind,
				Seller: "",
				Name:   "",
				Amount: -18,
				Note:   "",
			})

			container.AddChildByData(Data{
				Time:   mylog.Check2(time.Parse("2006-01-02", "2024-08-08")),
				Trade:  OtherKind,
				Seller: "",
				Name:   "",
				Amount: -29,
				Note:   "",
			})
			container.AddChildByData(Data{
				Time:   mylog.Check2(time.Parse("2006-01-02", "2024-08-09")), // todo
				Trade:  OtherKind,
				Seller: "",
				Name:   "",
				Amount: -18,
				Note:   "",
			})
			container.AddChildByData(Data{
				Time:   mylog.Check2(time.Parse("2006-01-02", "2024-08-09")),
				Trade:  OtherKind,
				Seller: "",
				Name:   "",
				Amount: -2,
				Note:   "",
			})

			container.AddChildByData(Data{
				Time:   mylog.Check2(time.Parse("2006-01-02", "2024-08-09")),
				Trade:  OtherKind,
				Seller: "",
				Name:   "",
				Amount: -101,
				Note:   "",
			})

			container.AddChildByData(Data{
				Time:   mylog.Check2(time.Parse("2006-01-02", "2024-08-09")),
				Trade:  OtherKind,
				Seller: "",
				Name:   "",
				Amount: -94,
				Note:   "",
			})

			container.AddChildByData(Data{
				Time:   mylog.Check2(time.Parse("2006-01-02", "2024-08-09")),
				Trade:  OtherKind,
				Seller: "",
				Name:   "",
				Amount: -660,
				Note:   "喇叭匠",
			})

			container.AddChildByData(Data{
				Time:   mylog.Check2(time.Parse("2006-01-02", "2024-08-08")),
				Trade:  MeatKind,
				Seller: "杨学春",
				Name:   "",
				Amount: -19.4 * 45,
				Note:   "二流",
			})

			container.AddChildByData(Data{
				Time:   mylog.Check2(time.Parse("2006-01-02", "2024-08-09")),
				Trade:  OtherKind,
				Seller: "",
				Name:   "",
				Amount: -15,
				Note:   "",
			})

			container.AddChildByData(Data{
				Time:   mylog.Check2(time.Parse("2006-01-02", "2024-08-09")),
				Trade:  OtherKind,
				Seller: "",
				Name:   "",
				Amount: -24,
				Note:   "",
			})

			container.AddChildByData(Data{
				Time:   mylog.Check2(time.Parse("2006-01-02", "2024-08-09")),
				Trade:  OtherKind,
				Seller: "",
				Name:   "",
				Amount: -130,
				Note:   "",
			})

			//////////////////挂礼
			giftContainer := widget.NewContainerNode("gift", Data{
				Time: mylog.Check2(time.Parse("2006-01-02", "2024-08-09")),
				// Trade:  MeatKind,
				// Seller: "杨学春",
				// Name:   "",
				// Amount: 19.4 * 45,
				// Note:   "二流",
			})

			giftContainer.AddChildByData(Data{
				Time:   mylog.Check2(time.Parse("2006-01-02", "2024-08-09")),
				Trade:  OtherKind,
				Seller: "",
				Name:   "彭光坤",
				Amount: 100,
				Note:   "收礼",
			})
			giftContainer.AddChildByData(Data{
				Time:   mylog.Check2(time.Parse("2006-01-02", "2024-08-09")),
				Trade:  OtherKind,
				Seller: "",
				Name:   "彭光文",
				Amount: 100,
				Note:   "收礼",
			})
			giftContainer.AddChildByData(Data{
				Time:   mylog.Check2(time.Parse("2006-01-02", "2024-08-09")),
				Trade:  OtherKind,
				Seller: "",
				Name:   "罗得华",
				Amount: 100,
				Note:   "收礼",
			})
			giftContainer.AddChildByData(Data{
				Time:   mylog.Check2(time.Parse("2006-01-02", "2024-08-09")),
				Trade:  OtherKind,
				Seller: "",
				Name:   "杨国伟",
				Amount: 200,
				Note:   "收礼",
			})
			giftContainer.AddChildByData(Data{
				Time:   mylog.Check2(time.Parse("2006-01-02", "2024-08-09")),
				Trade:  OtherKind,
				Seller: "",
				Name:   "杨国友",
				Amount: 200,
				Note:   "收礼",
			})

			giftContainer.AddChildByData(Data{
				Time:   mylog.Check2(time.Parse("2006-01-02", "2024-08-09")),
				Trade:  OtherKind,
				Seller: "",
				Name:   "杨周泽",
				Amount: 150,
				Note:   "收礼",
			})
			giftContainer.AddChildByData(Data{
				Time:   mylog.Check2(time.Parse("2006-01-02", "2024-08-09")),
				Trade:  OtherKind,
				Seller: "",
				Name:   "杨国富",
				Amount: 100,
				Note:   "收礼",
			})
			giftContainer.AddChildByData(Data{
				Time:   mylog.Check2(time.Parse("2006-01-02", "2024-08-09")),
				Trade:  OtherKind,
				Seller: "",
				Name:   "杨红春",
				Amount: 500,
				Note:   "收礼",
			})
			giftContainer.AddChildByData(Data{
				Time:   mylog.Check2(time.Parse("2006-01-02", "2024-08-09")),
				Trade:  OtherKind,
				Seller: "",
				Name:   "杨国双",
				Amount: 200,
				Note:   "收礼",
			})
			giftContainer.AddChildByData(Data{
				Time:   mylog.Check2(time.Parse("2006-01-02", "2024-08-09")),
				Trade:  OtherKind,
				Seller: "",
				Name:   "杨国安",
				Amount: 200,
				Note:   "收礼",
			})
			//container.AddChildByData(Data{
			//	Time:   mylog.Check2(time.Parse("2006-01-02", "2024-08-09")),
			//	Trade:  OtherKind,
			//	Seller: "",
			//	Name:   "彭光学",
			//	Amount: 300,
			//	Note:   "挂礼，不是我收",
			//})
			//container.AddChildByData(Data{
			//	Time:   mylog.Check2(time.Parse("2006-01-02", "2024-08-09")),
			//	Trade:  OtherKind,
			//	Seller: "",
			//	Name:   "杨同泽",
			//	Amount: 200,
			//	Note:   "挂礼，不是我收",
			//})

			///////////zhifubao/
			container.AddChildByData(Data{
				Time:   mylog.Check2(time.Parse("2006-01-02", "2024-08-09")),
				Trade:  OtherKind,
				Seller: "",
				Name:   "",
				Amount: -180,
				Note:   "支付宝，买菜",
			})
			container.AddChildByData(Data{
				Time:   mylog.Check2(time.Parse("2006-01-02", "2024-08-09")),
				Trade:  OtherKind,
				Seller: "",
				Name:   "",
				Amount: -140,
				Note:   "支付宝，买菜",
			})

			containerSum.AddChild(container)
			containerSum.AddChild(giftContainer)
			root.AddChild(containerSum)
		},
		JsonName:   "demo",
		IsDocument: false,
	})
}
