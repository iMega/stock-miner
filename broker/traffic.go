package broker

import "github.com/imega/stock-miner/domain"

type StockItemTraffic struct {
	ApprovedCh          chan domain.Message
	PriceReceiptCh      chan domain.PriceReceiptMessage
	PriceReceiptStoreCh chan domain.PriceReceiptMessage
	SellCh              chan domain.Slot
	ConfirmSellCh       chan domain.Message
	BuyCh               chan domain.Slot
	ConfirmBuyCh        chan domain.Message
	OperationCh         chan domain.Message
}

func NewStockItemTraffic() StockItemTraffic {
	return StockItemTraffic{
		ApprovedCh:          make(chan domain.Message),
		PriceReceiptCh:      make(chan domain.PriceReceiptMessage),
		PriceReceiptStoreCh: make(chan domain.PriceReceiptMessage),
		SellCh:              make(chan domain.Slot),
		ConfirmSellCh:       make(chan domain.Message),
		BuyCh:               make(chan domain.Slot),
		ConfirmBuyCh:        make(chan domain.Message),
		OperationCh:         make(chan domain.Message),
	}
}
