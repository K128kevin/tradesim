package model

type Transaction struct {
	Symbol 			string `form:"Symbol" json:"Symbol" binding:"required"`
	Quantity 		float64 `form:"Quantity" json:"Quantity" binding:"required"`
	Fee 			float64 `form:"Fee" json:"Fee" binding:"required"`
}

type TransactionDetail struct {
	Symbol 			string `form:"Symbol" json:"Symbol" binding:"required"`
	Quantity 		float64 `form:"Quantity" json:"Quantity" binding:"required"`
	Fee 			float64 `form:"Fee" json:"Fee" binding:"required"`
	Action 			string `form:"Action" json:"Action" binding:"required"`
	Rate 			float64 `form:"Rate" json:"Rate" binding:"required"`
	Time 			string `form:"Time" json"Time" binding:"required"`
}

type Rate struct {
	Symbol 			string `json:"symbol" binding:"required"`
	Name 			string `json:"description" binding:"required"`
	Price 			float64 `json:"last" binding:"required"`
	Change 			float64 `json:"change" binding:"required"`
	Bid 			float64 `json:"bid" binding:"required"`
	Ask 			float64 `json:"ask" binding:"required"`
}