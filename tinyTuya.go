package main

import (
	"log"
	"net/http"
	"text/template"
	//"github.com/py60800/tuya"
	"encoding/json"
	"flag"
	"io/ioutil"
	"strconv"
	"time"
	"tuya"
)

var tmpl *template.Template
var dm *tuya.DeviceManager
var devList []tuya.Device

func initDevList(dm *tuya.DeviceManager) {
	keys := dm.DeviceKeys()
	devList = make([]tuya.Device, 0)
	for _, k := range keys {
		d, _ := dm.GetDevice(k)
		devList = append(devList, d)
	}
}

type devDescr struct {
	Name    string
}

func setHandler(w http.ResponseWriter, r *http.Request) {
	sw := r.FormValue("switch")
	val := r.FormValue("set")
	log.Println("Set:", sw, val)
	s, ok := dm.GetDevice(sw)
	if ok {
		p := s.(tuya.Switch)
		if val == "on" || val == "true" {
			p.Set(true)
		} else {
			p.Set(false)
		}
	} else {
		log.Println("Unknown device :", sw)
	}
}
func getHandler(w http.ResponseWriter, r *http.Request) {
	//log.Println("Get")

	// The code hereafter is not optimized
	// if there were many devices or many users, an intermediate coroutine
	// should be used to avoid the burden and the load of
	// subscription/subscription

	skeys := make([]int64, 0)

	// Subscribe for events from these devices
	wait := r.FormValue("w")
	if wait != "false" {
		syncChannel := tuya.MakeSyncChannel()
		for _, b := range devList {
			skeys = append(skeys, b.Subscribe(syncChannel))
		}

		// Wait until data update or timeout
		select {
		case <-syncChannel:
		case <-time.After(time.Second * 15):
		}

		// cancel subscriptions
		for i := range skeys {
			devList[i].Unsubscribe(skeys[i])
		}
	}
	// build the response
	result := make(map[string]interface{})
	for _, b := range devList {
		s, ok := b.(tuya.Switch)
		if ok {
			st, err := s.Status()
			t := make(map[string]interface{})
			t["Value"] = st
			if err == nil {
				t["Status"] = "OK"
			} else {
				t["Status"] = err.Error()
			}
			result[b.Name()] = t
		}
	}

	// send the response
	json, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}
func homeHandler(w http.ResponseWriter, r *http.Request) {

	tmpl.ExecuteTemplate(w, "header.tmpl", nil)
	for _, b := range devList {
		d := devDescr{b.Name()}
		if b.Type() == "Switch" {
			tmpl.ExecuteTemplate(w, "Switch.tmpl", d)
		}
	}
	tmpl.ExecuteTemplate(w, "footer.tmpl", nil)

}
func main() {
	pconfig := flag.String("c", "", "configuration file")
	port := flag.Int("p", 8080, "Port number")
	flag.Parse()
	if len(*pconfig) == 0 {
		flag.PrintDefaults()
		return
	}
	dm = tuya.NewDeviceManager(getConfig(*pconfig))
	initDevList(dm)

	var err error
	tmpl, err = template.ParseGlob("tmpl/*.tmpl")
	if err != nil {
		log.Fatal("Template error:", err)
	}
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/set", setHandler)
	http.HandleFunc("/get", getHandler)
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static", fs))
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*port), nil))
}

func getConfig(tuyaconf string) string {
	b, err := ioutil.ReadFile(tuyaconf)
	if err != nil {
		log.Fatal("Cannot read:", tuyaconf)
	}
	return string(b)
}
