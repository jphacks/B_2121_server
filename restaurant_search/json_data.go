package restaurant_search

type hotpepperResponseJson struct {
	Results struct {
		ResultsAvailable int    `json:"results_available"`
		ResultsReturned  string `json:"results_returned"`
		ResultsStart     int    `json:"results_start"`
		Shop             []struct {
			Address string  `json:"address"`
			Id      string  `json:"id"`
			Lat     float64 `json:"lat"`
			Lng     float64 `json:"lng"`
			Name    string  `json:"name"`
			Photo   struct {
				Mobile struct {
					L string `json:"l"`
				} `json:"mobile"`
			} `json:"photo"`
			ShopDetailMemo string `json:"shop_detail_memo"`
			Urls           struct {
				Pc string `json:"pc"`
			} `json:"urls"`
		} `json:"shop"`
	} `json:"results"`
}
