package xrp

import (
	"errors"
	"fmt"
	"math"
	"qc-defi-graphql-server/internal/models"
	"strings"
	"sync"
	"time"

	"gorm.io/gorm"
)

// Remove the import for 'qc-defi-graphql-server/internal/configs/log'
// Comment out or remove any code that references 'log.BackgroundJobsLogger' or any other symbol from this package

type histDataOpts struct {
	interval       string
	histTableName  string
	prcChangeDbCol string
	timeBackSec    int
}

func stringToSlices(slice []string, size int) [][]string {
	var sliceOfSlices [][]string

	for i := 0; i < len(slice); i += size {
		end := i + size

		if end > len(slice) {
			end = len(slice)
		}

		sliceOfSlices = append(sliceOfSlices, slice[i:end])
	}

	return sliceOfSlices
}
func intervalsToSlices(slice []*models.IntervalBase, size int) [][]*models.IntervalBase {
	var sliceOfSlices [][]*models.IntervalBase

	for i := 0; i < len(slice); i += size {
		end := i + size

		if end > len(slice) {
			end = len(slice)
		}

		sliceOfSlices = append(sliceOfSlices, slice[i:end])
	}

	return sliceOfSlices
}
func spotPricesToSlices(slice []*models.SpotPrice, size int) [][]*models.SpotPrice {
	var sliceOfSlices [][]*models.SpotPrice

	for i := 0; i < len(slice); i += size {
		end := i + size

		if end > len(slice) {
			end = len(slice)
		}

		sliceOfSlices = append(sliceOfSlices, slice[i:end])
	}

	return sliceOfSlices
}

func getCandlestickTimeFrame(currentTimeInUTC time.Time, interval string) (int64, int64) {
	var (
		openTimestamp  int64
		closeTimestamp int64
	)

	conf := map[string]int64{
		"1min": 60,
		"5min": 5 * 60,
		"1h":   60 * 60,
	}

	if c, ok := conf[interval]; ok {
		if interval == "1min" || interval == "5min" || interval == "1h" {
			current := currentTimeInUTC.Unix()
			openTimestamp = int64(math.Floor(float64(current/c)) * float64(c))
			closeTimestamp = openTimestamp + c
		}
	}

	if interval == "24h" {
		beginningOfDay := time.Date(currentTimeInUTC.Year(), currentTimeInUTC.Month(), currentTimeInUTC.Day(), 0, 0, 0, 0, time.UTC)
		openTimestamp = beginningOfDay.Unix()
		closeTimestamp = beginningOfDay.Unix() + 60*60*24
	}
	return openTimestamp, closeTimestamp
}

var conf = []histDataOpts{
	{interval: "5min", histTableName: "historical_data_5min", prcChangeDbCol: "change_5min", timeBackSec: 5 * 60},
	{interval: "1h", histTableName: "historical_data_1h", prcChangeDbCol: "change_1h", timeBackSec: 60 * 60},
	{interval: "24h", histTableName: "historical_data_24h", prcChangeDbCol: "change_24h", timeBackSec: 60 * 60 * 24},
}

func getOpts(interval string) (histDataOpts, error) {
	var c histDataOpts
	for _, v := range conf {
		if v.interval == interval {
			c = v
			return c, nil
		}
	}
	return c, errors.New("invalid interval")
}

type HistoricalDataService struct {
	deFiDb  *gorm.DB
	network *models.Network // network configuration
}

func NewHistoricalDataService(deFeDb *gorm.DB, network *models.Network) *HistoricalDataService {
	return &HistoricalDataService{deFiDb: deFeDb, network: network}
}

func (svc *HistoricalDataService) GetAllSpotPrices() map[string]*models.SpotPrice {
	var (
		pools []models.Pool
		r     = make(map[string]*models.SpotPrice)
	)

	err := svc.deFiDb.Select("id", "network", "dex", "address", "token_0_price_usd", "token_1_price_usd", "quote_exact_in", "reserve_ratio").Find(&pools)

	if err != nil {
		fmt.Println(err)
	}
	for _, p := range pools {
		key := p.Network + ":" + p.Dex + ":" + p.Address
		s := models.SpotPrice{Id: p.ID, Network: p.Network, Dex: p.Dex, PairAddress: p.Address, Token0PriceUSD: p.Token0PriceUSD, Token1PriceUSD: p.Token1PriceUSD, ReserveRatio: p.ReserveRatio, QuoteExactIn: p.QuoteExactIn}
		r[key] = &s
	}

	return r
}

func (svc *HistoricalDataService) getCurrentCandle(table string, openTime int64, closeTime int64) map[string]*models.IntervalBase {
	var (
		candles  []*models.IntervalBase
		tokenMap = make(map[string]*models.IntervalBase)
	)
	svc.deFiDb.Table(table).Where("open_time = ? and close_time = ?", openTime, closeTime).Find(&candles)
	for _, t := range candles {
		tokenMap[strings.ToLower(t.Network+":"+t.Dex+":"+t.Address)] = t
	}
	return tokenMap
}

func (svc *HistoricalDataService) saveNewCandles(table string, prices map[string]*models.SpotPrice, openTime int64, closeTime int64) error {
	var (
		slices [][]*models.SpotPrice
		slice  []*models.SpotPrice
		size   = 1000
		i      = 0
	)

	for _, value := range prices {
		slice = append(slice, value)
		i++
		if i == size {
			slices = append(slices, slice)
			slice = make([]*models.SpotPrice, 0, size)
			i = 0
		}
	}
	if i != 0 {
		slices = append(slices, slice)
	}

	tx := svc.deFiDb.Begin()
	for _, slice := range slices {
		var (
			valueStrings []string
			valueArgs    []interface{}
		)

		for _, t := range slice {
			valueStrings = append(valueStrings, "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
			valueArgs = append(valueArgs, openTime)
			valueArgs = append(valueArgs, closeTime)
			valueArgs = append(valueArgs, t.Network)
			valueArgs = append(valueArgs, t.Dex)
			valueArgs = append(valueArgs, t.PairAddress)
			valueArgs = append(valueArgs, t.Token0PriceUSD)
			valueArgs = append(valueArgs, t.Token0PriceUSD)
			valueArgs = append(valueArgs, t.Token0PriceUSD)
			valueArgs = append(valueArgs, t.Token0PriceUSD)
			valueArgs = append(valueArgs, time.Now().Unix())
		}

		stmt := fmt.Sprintf(`INSERT INTO `+table+` (open_time, close_time, network, dex, address, open, close, high,low, created_at) VALUES %s`, strings.Join(valueStrings, ","))

		err := svc.deFiDb.Transaction(func(tx *gorm.DB) error {
			var err error
			for i := 0; i < 10; i++ {
				err = tx.Exec(stmt, valueArgs...).Error
				if err != nil {
					fmt.Printf("create new candles restarting transaction try: %d\n", i)
					time.Sleep(time.Duration(i) * time.Second)
					if i == 9 {
						tx.Rollback()
						break
					}
					continue
				}
				break
			}
			return err
		})

		if err != nil {
			tx.Rollback()
			return err
		}

	}

	err := tx.Commit().Error
	if err != nil {
		return err
	}
	fmt.Printf("created new candles table: %s count: %d\n", table, len(prices))
	return nil
}

func (svc *HistoricalDataService) updateCandles(table string, currentData map[string]*models.IntervalBase, prices map[string]*models.SpotPrice) error {
	var (
		res []*models.IntervalBase
		tx  = svc.deFiDb.Begin()
	)
	for key, current := range currentData {
		if price, ok := prices[key]; ok {
			current.SetHigh(price.Token0PriceUSD)
			current.SetLow(price.Token0PriceUSD)
			current.Close = price.Token0PriceUSD
			current.UpdatedAt = time.Now().Unix()
		}
		res = append(res, current)
	}

	slices := intervalsToSlices(res, 1000)
	tx = tx.Set("gorm:query_option", "FOR UPDATE") // isolation level

	for _, slice := range slices {
		var (
			valueStrings []string
			valueArgs    []interface{}
		)

		for _, t := range slice {
			valueStrings = append(valueStrings, "(?, ?, ?, ?, ?, ?)")
			valueArgs = append(valueArgs, t.ID)
			valueArgs = append(valueArgs, t.OpenTime)
			valueArgs = append(valueArgs, t.High)
			valueArgs = append(valueArgs, t.Low)
			valueArgs = append(valueArgs, t.Close)
			valueArgs = append(valueArgs, t.UpdatedAt)
		}

		stmt := fmt.Sprintf(`INSERT INTO `+table+` (id, open_time, high,low,close,updated_at) VALUES %s
															on duplicate key update
														    high=values(high),
														    low=values(low),
														    close=values(close),
														    updated_at=values(updated_at)
			                                             `, strings.Join(valueStrings, ","))
		//err := tx.Exec(stmt, valueArgs...).Error

		err := svc.deFiDb.Transaction(func(tx *gorm.DB) error {
			var err error
			for i := 0; i < 10; i++ {
				err = tx.Exec(stmt, valueArgs...).Error
				if err != nil {
					fmt.Printf("update candles restarting transaction try: %d\n", i)
					time.Sleep(time.Duration(i) * time.Second)
					if i == 9 {
						tx.Rollback()
						break
					}
					continue
				}
				break
			}
			return err
		})

		if err != nil {
			return err
		}
	}

	err := tx.Commit().Error
	if err != nil {
		return err
	}

	fmt.Printf("updated candles table: %s count: %d\n", table, len(currentData))
	return nil
}

func (svc *HistoricalDataService) ProcessPtcChange(interval string, prices map[string]*models.SpotPrice, historicalTale string, prcChangeDbCol string, timeBackSec int) error {
	type tempStruct struct {
		ID      uint
		Network string
		Address string
		Dex     string
		Close   float64
	}

	toSlices := func(slice []*tempStruct, size int) [][]*tempStruct {
		var sliceOfSlices [][]*tempStruct

		for i := 0; i < len(slice); i += size {
			end := i + size

			if end > len(slice) {
				end = len(slice)
			}

			sliceOfSlices = append(sliceOfSlices, slice[i:end])
		}

		return sliceOfSlices
	}

	var (
		result              []*tempStruct
		targetTable         = "pools"
		timeBack            = time.Now().UTC().Add(time.Duration(-1*timeBackSec) * time.Second)
		tx                  = svc.deFiDb.Begin()
		openTime, closeTime = getCandlestickTimeFrame(timeBack, interval)
	)

	svc.deFiDb.Table("pools AS pl").Select(`pl.id, pl.network, pl.dex, pl.address, hd.close`).
		Joins(`left join `+historicalTale+` hd on pl.network = hd.network and pl.address=hd.address and pl.dex=hd.dex`).
		Where("hd.open_time = ? and hd.close_time = ? ", openTime, closeTime).Scan(&result)

	slices := toSlices(result, 1000)
	tx = tx.Set("gorm:query_option", "FOR UPDATE") // isolation level

	for _, slice := range slices {
		var (
			valueStrings []string
			valueArgs    []interface{}
		)
		for _, historical := range slice {
			key := strings.ToLower(historical.Network + ":" + historical.Dex + ":" + historical.Address)
			if c, ok := prices[key]; ok {
				var change float64
				if historical.Close != 0 {
					change = (c.Token0PriceUSD - historical.Close) / historical.Close
				}

				valueStrings = append(valueStrings, "(?, ?)")
				valueArgs = append(valueArgs, historical.ID)
				valueArgs = append(valueArgs, change)
			}
		}

		stmt := fmt.Sprintf(`INSERT INTO `+targetTable+` (id, `+prcChangeDbCol+`) VALUES %s
															on duplicate key update
														    `+prcChangeDbCol+`=values(`+prcChangeDbCol+`)
			                                             	`, strings.Join(valueStrings, ","))

		err := svc.deFiDb.Transaction(func(tx *gorm.DB) error {
			var err error
			for i := 0; i < 10; i++ {
				err = tx.Exec(stmt, valueArgs...).Error
				if err != nil {
					fmt.Printf("percent change restarting transaction try: %d\n", i)
					time.Sleep(time.Duration(i) * time.Second)
					if i == 9 {
						tx.Rollback()
						break
					}
					continue
				}
				break
			}
			return err
		})

		if err != nil {
			return err
		}
	}

	err := tx.Commit().Error
	if err != nil {
		return err
	}

	fmt.Printf("updated ptc change db-col: %s count: %d\n", prcChangeDbCol, len(result))
	return nil
}

func (svc *HistoricalDataService) Process(interval string) error {

	conf, err := getOpts(interval)

	if err != nil {
		return err
	}

	var (
		openTime, closeTime = getCandlestickTimeFrame(time.Now().UTC(), interval)
		prices              = svc.GetAllSpotPrices()
		currentData         = svc.getCurrentCandle(conf.histTableName, openTime, closeTime)
	)

	if len(currentData) == 0 {
		err := svc.saveNewCandles(conf.histTableName, prices, openTime, closeTime)
		if err != nil {
			fmt.Printf("save new candles error table: %s interval: %s\n", conf.histTableName, interval)
		}

	} else {
		err := svc.updateCandles(conf.histTableName, currentData, prices)
		if err != nil {
			fmt.Printf("update candles error table: %s interval: %s\n", conf.histTableName, interval)
		}
	}

	err = svc.ProcessPtcChange(interval, prices, conf.histTableName, conf.prcChangeDbCol, conf.timeBackSec)
	if err != nil {
		fmt.Printf("update ptc-change error log: %s interval: %s\n", err.Error(), interval)
	}

	return nil
}

func (svc *HistoricalDataService) RunForAllInterval() {
	var wg = &sync.WaitGroup{}
	for {
		if time.Now().Unix()%60 == 0 {
			for _, i := range conf {
				wg.Add(1)
				go func(interval string) {
					defer wg.Done()
					_ = svc.Process(interval)
				}(i.interval)
			}
		}
		time.Sleep(time.Second * 1)
	}
}

func (svc *HistoricalDataService) DeleteData() {

	var (
		timeNow  = time.Now().Unix()
		timeBack = timeNow - 60*60*24*7
	)
	fmt.Println(timeNow, timeBack)

	ToSlices := func(slice []models.Interval5minCandle, size int) [][]models.Interval5minCandle {
		var sliceOfSlices [][]models.Interval5minCandle

		for i := 0; i < len(slice); i += size {
			end := i + size

			if end > len(slice) {
				end = len(slice)
			}

			sliceOfSlices = append(sliceOfSlices, slice[i:end])
		}

		return sliceOfSlices
	}
	for {

		var candles []models.Interval5minCandle

		if err := svc.deFiDb.Limit(1000000).Select("id", "open_time").Where("open_time < ?", timeBack).Find(&candles).Error; err != nil {
			continue
		}
		if len(candles) < 1000 {
			break
		}
		slices := ToSlices(candles, 10000)
		fmt.Println("TOTAL COUNT", len(candles), "SLICE", len(slices))
		for i, slice := range slices {
			var IDS []uint
			for _, elem := range slice {
				IDS = append(IDS, elem.ID)
			}
			svc.deFiDb.Where("id in ?", IDS).Delete(&models.Interval5minCandle{})
			fmt.Println("SLICE", i+1, "DELETED", time.Now(), len(IDS))
		}
	}

}

// ----------------------------------------------------------

/*
func (svc *HistoricalDataService) GetAllPrices() map[string]*model.SpotPrice {
	var (
		pricesMap = make(map[string]*model.SpotPrice)
		//wg        = &sync.WaitGroup{}
		//mutex     = &sync.RWMutex{}
		ctx = context.Background()
	)

	keys := svc.redisDb.Keys(context.Background(), "*").Val()
	//keySlices := stringToSlices(keys, 500)
	//fmt.Println(keys)
	pipe := svc.redisDb.Pipeline()
	for _, key := range keys {
		pipe.HGetAll(ctx, key).Val()

		//wg.Add(1)
		//go func(keys []string) {
		//	defer wg.Done()
		//	data := svc.redisDb.HGetAll(context.Background(), keys[0]).Val()
		//	fmt.Println(data)
		//	for _, d := range data {
		//		var p = model.SpotPrice{}
		//		err := json.Unmarshal([]byte(d), &p)
		//		if err != nil {
		//			continue
		//		}
		//		mutex.Lock()
		//		pricesMap[strings.ToLower(p.Network+"_"+p.PairAddress)] = &p
		//		mutex.Unlock()
		//	}
		//}(i)
	}
	//wg.Wait()

	// Execute the pipeline and retrieve the results
	results, err := pipe.Exec(ctx)
	if err != nil {
		return nil
	}
	//Process the retrieved results
	for i, result := range results {
		if result.Err() != nil {
			continue
		}
		//// Get the fields and values from the result
		fieldsAndValues, err := result.(*redis.StringStringMapCmd).Result()
		if err != nil {
			continue
		}
		jsonData, err := json.Marshal(fieldsAndValues)
		if err != nil {
			fmt.Println(err)
			continue
		}
		var p = model.SpotPrice{}
		err = json.Unmarshal(jsonData, &p)
		if err != nil {
			fmt.Println(err)
			continue
		}
		pricesMap[keys[i]] = &p
	}

	return pricesMap
}
*/
