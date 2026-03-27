package model

type SellerStats struct {
	TotalListings int     `json:"total_listings"`
	TotalSold     int     `json:"total_sold"`
	TotalRevenue  float64 `json:"total_revenue"`
}
