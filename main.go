package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/dns"
	"github.com/cloudflare/cloudflare-go/v6/option"
	"github.com/robfig/cron/v3"
)

// Config struct per contenere i valori letti dall'ambiente
type Config struct {
	APIToken string
	ZoneID   string
	CronSpec string
}

// loadConfigFromEnv carica la configurazione dalle variabili d'ambiente
func loadConfigFromEnv() (*Config, error) {
	apiToken := os.Getenv("CF_API_TOKEN")
	if apiToken == "" {
		return nil, errors.New("la variabile d'ambiente CF_API_TOKEN non è impostata")
	}

	zoneID := os.Getenv("CF_ZONE_ID")
	if zoneID == "" {
		return nil, errors.New("la variabile d'ambiente CF_ZONE_ID non è impostata")
	}

	cronSpec := os.Getenv("CRON_SPEC")
	if cronSpec == "" {
		cronSpec = "@every 30m" // Valore predefinito
		log.Printf("ℹ️ CRON_SPEC non impostato, utilizzo il default: %s", cronSpec)
	}

	return &Config{
		APIToken: apiToken,
		ZoneID:   zoneID,
		CronSpec: cronSpec,
	}, nil
}

// getPublicIP ottiene l'IP pubblico corrente
func getPublicIP() (string, error) {
	resp, err := http.Get("https://api.ipify.org")
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)

	ip, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(ip), nil
}

// updateDNS esegue il ciclo di controllo e aggiornamento
func updateDNS(config *Config) {
	log.Println("🚀 Inizio controllo DDNS...")

	publicIP, err := getPublicIP()
	if err != nil {
		log.Printf("❌ Impossibile ottenere l'IP pubblico: %v", err)
		return
	}
	log.Printf("ℹ️ IP pubblico corrente rilevato: %s", publicIP)

	client := cloudflare.NewClient(option.WithAPIToken(config.APIToken))

	ctx := context.Background()

	records, err := client.DNS.Records.List(ctx, dns.RecordListParams{
		ZoneID: cloudflare.F(config.ZoneID),
		Type:   cloudflare.F(dns.RecordListParamsTypeA),
	})
	if err != nil {
		log.Printf("❌ Errore nel recuperare i record: %v", err)
	}

	for _, record := range records.Result {

		if record.Content == publicIP {
			log.Printf("✅ L'IP è già aggiornato (%s). Nessuna azione richiesta.", publicIP)
		} else {
			log.Printf("🔄 L'IP è obsoleto. Vecchio: %s, Nuovo: %s. Aggiornamento in corso...", record.Content, publicIP)

			updateParams := dns.RecordEditParams{
				ZoneID: cloudflare.F(config.ZoneID),
				Body: dns.RecordEditParamsBody{
					Content: cloudflare.F(publicIP),
				},
			}

			_, err := client.DNS.Records.Edit(ctx, record.ID, updateParams)
			if err != nil {
				log.Printf("❌ Errore durante l'aggiornamento del record '%s': %v", record.Name, err)
			} else {
				log.Printf("🎉 Record DNS '%s' aggiornato con successo a %s!", record.Name, publicIP)
			}
		}
	}

	log.Println("✨ Controllo DDNS completato.")
}

func main() {
	// Carica la configurazione una sola volta all'avvio
	config, err := loadConfigFromEnv()
	if err != nil {
		log.Fatalf("Impossibile avviare il servizio. Errore di configurazione: %v", err)
	}

	c := cron.New()

	// Usa una closure per passare la configurazione al job schedulato
	_, err = c.AddFunc(config.CronSpec, func() {
		updateDNS(config)
	})
	if err != nil {
		log.Fatalf("Errore nell'impostare il job cron: %v", err)
	}

	log.Println("✅ Servizio DDNS Cloudflare avviato.")
	log.Printf("🕐 Schedulazione impostata a: %s", config.CronSpec)

	// Esegui un primo controllo all'avvio per un feedback immediato
	go updateDNS(config)

	c.Start()

	// Mantieni il programma in esecuzione
	select {}
}
