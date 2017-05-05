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