package utils

type MarketDataRowType struct {
	AUCNG_DE             string
	PBLMNG_WHSAL_MRKT_NM string
	PBLMNG_WHSAL_MRKT_CD string
	CPR_NM               string
	CPR_CD               string
	RISENO               string
	ORGNO                string
	BIDTIME              string
	PRDLST_NM            string
	PRDLST_CD            string
	SPCIES_NM            string
	SPCIES_CD            string
	PRICE                int
	DELNGBUNDLE_QY       int
	STNDRD               string
	STNDRD_CD            string
	GRAD                 string
	GRAD_CD              string
	SANJI_CD             string
	SANJI_NM             string
	DELNG_QY             int
}

type MargetDataGrid struct {
	Row      []MarketDataRowType
	StartRow int
	EndRow   int
	TotalCnt int
}

type MarketDataType struct {
	Grid_20161221000000000429_1 MargetDataGrid
}
