package http

type RegisterRebate struct {
}

type RegisterTransaction struct {
}

type DateRange struct {
	StartDate string `json:"start_date" binding:"required"`
	EndDate   string `json:"end_date" binding:"required"`
}
