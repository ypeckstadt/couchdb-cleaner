package model


type Designs struct {
	TotalRows int `json:"total_rows"`
	Offset int `json:"offset"`
	Rows []Design `json:"rows"`
}
type Design struct {
	ID string `json:"id"`
	Key string `json:"key"`
	Value DesignValue `json:"value"`
}

type DesignValue struct {
	Revision string `json:"rev"`
}


