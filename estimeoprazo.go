package estimeoprazo

import (
	"encoding/json"
	"io"
	"math"
	"math/rand"
	"net/http"
	"sort"
	"strconv"
	"sync"
	"time"
)

type Config struct {
	MinTasks      float64
	MaxTasks      float64
	MinSplitTasks float64
	MaxSplitTasks float64
	MinTasksDone  float64
	MaxTasksDone  float64
	Simulations   int
}

type ForecastResult struct {
	Likelihood int
	Weeks      float64
}

func TotalStories(config Config) float64 {
	return math.Ceil(getRandBeetween(config.MinTasks, config.MaxTasks) * (getRand()*float64(config.MaxSplitTasks-config.MinSplitTasks) + float64(config.MinSplitTasks)))
}

func EndWeekStories(totalStories float64, config Config) float64 {
	return math.Max(0, totalStories-getRandBeetween(config.MinTasksDone, config.MaxTasksDone))
}

func getRand() float64 {
	rand.Seed(time.Now().UTC().UnixNano())

	return rand.Float64()
}

func WeeksToZero(config Config, wg *sync.WaitGroup, weeks *sort.Float64Slice) float64 {
	defer wg.Done()
	totalStories := TotalStories(config)
	week := 1
	endWeekStories := EndWeekStories(totalStories, config)
	for endWeekStories > 0 {
		week = week + 1
		endWeekStories = EndWeekStories(endWeekStories, config)
	}

	*weeks = append(*weeks, float64(week))

	return float64(week)
}

func getRandBeetween(min, max float64) float64 {
	rand.Seed(time.Now().UTC().UnixNano())

	return rand.Float64()*(max-min) + min
}

func Percentile(numbers sort.Float64Slice, l int, n int) float64 {
	i := l*n/100 - 1
	return numbers[i]
}

func Forecast(config Config) []ForecastResult {
	var wg sync.WaitGroup
	var weeks sort.Float64Slice
	for i := 0; i <= config.Simulations; i++ {
		wg.Add(1)
		go WeeksToZero(config, &wg, &weeks)
	}
	wg.Wait()
	sort.Sort(weeks)
	l := len(weeks)
	result := []ForecastResult{
		{5, Percentile(weeks, l, 5)},
		{10, Percentile(weeks, l, 10)},
		{15, Percentile(weeks, l, 15)},
		{20, Percentile(weeks, l, 20)},
		{25, Percentile(weeks, l, 25)},
		{30, Percentile(weeks, l, 30)},
		{35, Percentile(weeks, l, 35)},
		{40, Percentile(weeks, l, 40)},
		{45, Percentile(weeks, l, 45)},
		{50, Percentile(weeks, l, 50)},
		{55, Percentile(weeks, l, 55)},
		{60, Percentile(weeks, l, 60)},
		{65, Percentile(weeks, l, 55)},
		{70, Percentile(weeks, l, 70)},
		{75, Percentile(weeks, l, 75)},
		{80, Percentile(weeks, l, 80)},
		{85, Percentile(weeks, l, 85)},
		{90, Percentile(weeks, l, 90)},
		{95, Percentile(weeks, l, 95)},
		{100, Percentile(weeks, l, 100)},
	}
	return result
}

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var config = Config{}

	config.MinTasks = filterInputToFloat(r.FormValue("MinTasks"))
	config.MaxTasks = filterInputToFloat(r.FormValue("MaxTasks"))
	config.MinSplitTasks = filterInputToFloat(r.FormValue("MinSplitTasks"))
	config.MaxSplitTasks = filterInputToFloat(r.FormValue("MaxSplitTasks"))
	config.MinTasksDone = filterInputToFloat(r.FormValue("MinTasksDone"))
	config.MaxTasksDone = filterInputToFloat(r.FormValue("MaxTasksDone"))
	config.Simulations = 1000

	w.Header().Set("Content-Type", "application/json")
	result, _ := json.Marshal(Forecast(config))
	io.WriteString(w, string(result))
}

func filterInputToFloat(value string) float64 {
	result, err := strconv.ParseFloat(value, 64)
	if err != nil {
		panic("Error converting value")
	}
	return result
}

func init() {

	http.HandleFunc("/", HandleIndex)
}
