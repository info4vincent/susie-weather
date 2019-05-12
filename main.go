//see uri = https://developer.accuweather.com/apis
package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
//	"strconv"
	"strings"
	"time"
	"encoding/json"
	"github.com/bitly/go-simplejson"
)

type ValueType struct {
    Value float64
    Unit string
    UnitType int64
}


func getContent() {
	apiKey := "{YOUR_ACCUWEATHER_DEVLOPER_ACCOUNT_API_KEY}"
	baseurl := "http://dataservice.accuweather.com/forecasts/v1/daily/1day/249209?apikey={APIKEY_ACCUWEATHER}&language=%20nl-nl%20&details=true&metric=true"

	yqlURL := strings.Replace(baseurl, "{APIKEY_ACCUWEATHER}", apiKey,1)
	fmt.Println(yqlURL)

	res, _ := http.Get(yqlURL)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}

	js, err := simplejson.NewJson(body)

	defer res.Body.Close()

	b, _ := js.EncodePretty()
	fmt.Printf("%s\n", b)
	fmt.Println(js)

	fmt.Println("--------------- v ---------------")
	timeNow := time.Now()
	v := js.GetPath().MustMap()
	a := v["DailyForecasts"].([]interface{})
	first := a[0].(map[string]interface{})

	day := strings.ToLower(time.Now().Weekday().String())
	date := strings.ToLower(first["Date"].(string))
	realFeelTemperature := first["RealFeelTemperature"].(map[string]interface{})
	lowTemperatureS := realFeelTemperature["Minimum"].(map[string]interface{})
	lowTemperature,_ := lowTemperatureS["Value"].(json.Number).Float64()
	highTemperature,_ := realFeelTemperature["Maximum"].(map[string]interface {})["Value"].(json.Number).Float64()
	dayWeather := first["Day"].(map[string]interface{})
	weatherText := strings.ToLower(dayWeather["LongPhrase"].(string))

	vs := js.Get("DailyForecasts").GetIndex(0).GetPath("Day","Wind").MustMap()
	speedStr := vs["Speed"].(map[string]interface{})
	speed,_ := speedStr["Value"].(json.Number).Float64()
	beaufort := int(speedInBeaufort(speed))

	fmt.Println("a:day", day)
	fmt.Println("a:date", date)
	fmt.Println("a:low", lowTemperature)
	fmt.Println("a:high", highTemperature)
	fmt.Println("a:text", weatherText)
	fmt.Println("a:speed", beaufort)

	fmt.Println("uur:", timeNow.Hour())
	fmt.Println("minuut", timeNow.Minute())

//	weerText := map[string]string{
		// "tornado":               "tornado",
		// "tropical storm":        "tropische_storm",
		// "hurricane":             "orkaan",
		// "severe thunderstorms":  "zware_onweersbuien",
		// "thunderstorms":         "onweersbuien",
		// "mixed rain and snow":   "regen en sneeuw",
		// "mixed rain and sleet":  "regen en natte sneeuw",
		// "mixed snow and sleet":  "sneeuw en natte sneeuw",
		// "freezing drizzle":      "ijzel",
		// "drizzle":               "motregen",
		// "freezing rain":         "ijskoude_regen",
		// "showers":               "stortbuien",
		// "snow flurries":         "sneeuwvlagen",
		// "light snow showers":    "lichte_sneeuwbuien",
		// "blowing snow":          "sneeuwbuien",
		// "snow":                  "sneeuw",
		// "hail":                  "hagel",
		// "sleet":                 "ijzel",
		// "dust":                  "stof",
		// "foggy":                 "mistig",
		// "haze":                  "nevel",
		// "smoky":                 "grijs",
		// "blustery":              "heftig",
		// "windy":                 "winderig",
		// "cold":                  "koud",
		// "cloudy":                "bewolkt",
		// "mostly cloudy":         "meest bewolkt",
		// "mostly cloudy (night)": "meest_bewolkt_vannacht",
		// "mostly cloudy (day)":   "meest_bewolkt_vandaag",
		// "partly cloudy (night)": "half bewolkt vannacht",
		// "partly cloudy (day)":   "half bewolkt vandaag",
		// "clear (night)":         "helder_nacht",
		// "sunny":                 "zonnig",
		// "fair (night)":          "mooi_vannacht",
		// "fair (day)":            "mooi_vandaag",
		// "mixed rain and hail":   "gemengde_regen_en_hagel",
		// "hot": "heet",
		// "isolated thunderstorms":  "plaatselijke_onweersbuien",
		// "scattered thunderstorms": "onweersbuien",
		// "scattered showers":       "verspreide_buien",
		// "heavy snow":              "zware_sneeuwbuien",
		// "scattered snow showers":  "verspreide_sneeuwbuien",
		// "partly cloudy":           "bewolkt",
		// "thundershowers":          "onweersbuien",
		// "snow showers":            "sneeuwbuien",
		// "isolated thundershowers": "plaatselijke_onweersbuien",
		// "not available":           "onbekend"}

	dagText := map[string]string{
		"monday": "maandag",
		"tuesday": "dinsdag",
		"wednesday": "woensdag",
		"thursday": "donderdag",
		"friday": "vrijdag",
		"saturday": "zaterdag",
		"sunday": "zondag"}

	dayMsg := fmt.Sprintf("goedemorgen het_is vandaag %s %d %s %d de temperatuur is nu %.0f graden en_loopt_op_tot %.0f graden het_weer_is %s wind %d beaufort", dagText[day], timeNow.Day(), strings.ToLower(timeNow.Month().String()), timeNow.Year(), lowTemperature, highTemperature, weatherText, beaufort)

	fmt.Printf("%s", dayMsg)
	stringSlice := strings.Split(dayMsg, " ")
	fmt.Println()

	f, err := os.Create("/var/lib/mpd/playlists/msgtoday.m3u")
	check(err)
	w := bufio.NewWriter(f)
	_, err = w.WriteString("# Playlist for susie\n")
	w.WriteString("#" + dayMsg + "\n")
	w.WriteString("susie/susie.mp3\n")

	for _, sound := range stringSlice {
		line := fmt.Sprintf("susie/%s.mp3\n", sound)
		w.WriteString(line)
		fmt.Println(line)
	}
	w.Flush()
	return
}

func speedInBeaufort(inKmPHr float64) (beaufort float64) {

	beaufort = math.Sqrt(inKmPHr)

	if beaufort < 7 {
		beaufort--
	}

	return
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	getContent()
	return
}

