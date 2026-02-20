package admin_dashboard

type AdminDashboardSummary struct {
	TotalUsers        int64            `json:"total_users"`
	ActiveUsers       int64            `json:"active_users"`
	TotalTransactions int64            `json:"total_transactions"`
	TotalVolume       float64          `json:"total_volume"`
	UserGrowth        []UserGrowthData `json:"user_growth"`
}

type UserGrowthData struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}
