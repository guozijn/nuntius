package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/guozijn/nuntius"
	"github.com/guozijn/nuntius/provider/telegram"
	"github.com/prometheus/alertmanager/template"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	bind = flag.String(
		"bind",
		":9096",
		"Address to listen on for request.")
	config = flag.String(
		"config",
		"nuntius.yml",
		"Nuntius configuration file path.")
	printVn = flag.Bool("v", false, "Print current build and version tags.")
	build   string
	gitHash string
	version string
)

func main() {
	flag.Parse()

	if *printVn {
		fmt.Println("Version    :", version)
		fmt.Println("Git hash   :", gitHash)
		fmt.Println("Build date :", build)
		os.Exit(0)
	}

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if err := LoadConfig(*config); err != nil {
		log.Fatalf("Error loading configuration: %s", err)
	}

	http.HandleFunc("/alert", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		// https://godoc.org/github.com/prometheus/alertmanager/template#Data
		data := template.Data{}
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			errorHandler(w, http.StatusBadRequest, err, "?")
			return
		}

		receiverConf := receiverConfByReceiver(data.Receiver)
		if receiverConf == nil {
			errorHandler(w, http.StatusBadRequest, fmt.Errorf("Receiver missing: %s", data.Receiver), "?")
			return
		}
		provider, err := providerByName(receiverConf.Provider)
		if err != nil {
			errorHandler(w, http.StatusInternalServerError, err, receiverConf.Provider)
			return
		}

		var text string
		if receiverConf.Text != "" {
			text, err = tmpl.ExecuteTextString(receiverConf.Text, data)
			if err != nil {
				errorHandler(w, http.StatusInternalServerError, err, receiverConf.Provider)
				return
			}
		} else {
			if len(data.Alerts) > 1 {
				labelAlerts := map[string]template.Alerts{
					"Firing":   data.Alerts.Firing(),
					"Resolved": data.Alerts.Resolved(),
				}
				for label, alerts := range labelAlerts {
					if len(alerts) > 0 {
						text += label + ": \n"
						for _, alert := range alerts {
							text += alert.Labels["alertname"] + " @" + alert.Labels["instance"]
							if len(alert.Labels["exported_instance"]) > 0 {
								text += " (" + alert.Labels["exported_instance"] + ")"
							}
							text += "\n"
						}
					}
				}
			} else if len(data.Alerts) == 1 {
				alert := data.Alerts[0]
				tuples := []string{}
				for k, v := range alert.Labels {
					tuples = append(tuples, k+"= "+v)
				}
				text = strings.ToUpper(data.Status) + " \n" + strings.Join(tuples, "\n")
			} else {
				text = "Alert \n" + strings.Join(data.CommonLabels.Values(), " | ")
			}
		}

		message := nuntius.Message{
			To:   receiverConf.To,
			From: receiverConf.From,
			Text: text,
		}

		if err = provider.Send(message); err != nil {
			errorHandler(w, http.StatusBadRequest, err, receiverConf.Provider)
			return
		}

		requestTotal.WithLabelValues("200", receiverConf.Provider).Inc()
	})

	http.Handle("/metrics", promhttp.Handler())

	http.HandleFunc("/-/reload", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		if r.Method == "POST" {
			log.Println("Loading configuration file: ", *config)
			if err := LoadConfig(*config); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			http.Error(w, "Invalid request method.", http.StatusMethodNotAllowed)
		}
	})

	if os.Getenv("NUNTIUS_PORT") != "" {
		*bind = ":" + os.Getenv("NUNTIUS_PORT")
	}

	log.Fatal(http.ListenAndServe(*bind, nil))
}

// receiverConfByReceiver loops the receiver conf list and returns the first instance with that name
func receiverConfByReceiver(name string) *ReceiverConf {
	for i := range providerConfig.Receivers {
		rc := &providerConfig.Receivers[i]
		if rc.Name == name {
			return rc
		}
	}
	return nil
}

func providerByName(name string) (nuntius.Provider, error) {
	switch name {
	case "telegram":
		return telegram.NewTelegram(providerConfig.Providers.Telegram)
	}

	return nil, fmt.Errorf("%s: Unknown provider", name)
}

func errorHandler(w http.ResponseWriter, status int, err error, provider string) {
	w.WriteHeader(status)

	data := struct {
		Error   bool
		Status  int
		Message string
	}{
		true,
		status,
		err.Error(),
	}
	// respond json
	bytes, _ := json.Marshal(data)
	json := string(bytes[:])
	fmt.Fprint(w, json)

	log.Println("Error: " + json)
	requestTotal.WithLabelValues(strconv.FormatInt(int64(status), 10), provider).Inc()
}
