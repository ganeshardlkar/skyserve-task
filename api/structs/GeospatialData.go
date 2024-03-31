package structs

type GeospatialData struct {
	ID     int    `json:"id"`
	UserID int    `json:"user_id"`
	Geom   string `json:"geom"`
}
