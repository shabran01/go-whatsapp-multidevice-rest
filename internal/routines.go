package internal

import (
	"github.com/robfig/cron/v3"

	"github.com/shabran01/go-whatsapp-multidevice-rest/pkg/log"
	pkgWhatsApp "github.com/shabran01/go-whatsapp-multidevice-rest/pkg/whatsapp"
)

func Routines(cron *cron.Cron) {
	log.Print(nil).Info("Running Routine Tasks")

	// Reduced frequency to avoid overwhelming the server
	cron.AddFunc("0 */5 * * * *", func() {
		// If WhatsAppClient Connection is more than 0
		if len(pkgWhatsApp.WhatsAppClient) > 0 {
			// Check Every Authenticated MSISDN
			for jid, client := range pkgWhatsApp.WhatsAppClient {
				// Get Real JID from Datastore
				realJID := client.Store.ID.User

				// Mask JID for Logging Information
				maskJID := realJID[0:len(realJID)-4] + "xxxx"

				// Print Log Show Information of Device Checking
				log.Print(nil).Info("Checking WhatsApp Client for " + maskJID)

				// Check WhatsAppClient Registered JID with Authenticated MSISDN
				if jid != realJID {
					// Print Log Show Information to Force Log-out Device
					log.Print(nil).Info("Logging out WhatsApp Client for " + maskJID + " Due to Missmatch Authentication")

					// Logout WhatsAppClient Device
					_ = pkgWhatsApp.WhatsAppLogout(jid)
					delete(pkgWhatsApp.WhatsAppClient, jid)
				} else {
					// Check if client is still connected and healthy
					if !client.IsConnected() || !client.IsLoggedIn() {
						log.Print(nil).Warn("WhatsApp Client for " + maskJID + " is disconnected, attempting reconnection")
						
						go func(clientJID string) {
							err := pkgWhatsApp.WhatsAppReconnect(clientJID)
							if err != nil {
								log.Print(nil).Error("Failed to reconnect " + maskJID + ": " + err.Error())
							}
						}(jid)
					}
				}
			}
		}
	})
	
	// Add health check routine
	cron.AddFunc("0 */10 * * * *", func() {
		log.Print(nil).Info("Running WhatsApp Client Health Check")
		
		for jid, client := range pkgWhatsApp.WhatsAppClient {
			if client == nil {
				continue
			}
			
			maskJID := jid[0:len(jid)-4] + "xxxx"
			
			// Check client health
			if !client.IsConnected() {
				log.Print(nil).Warn("Health Check: Client " + maskJID + " is not connected")
			} else if !client.IsLoggedIn() {
				log.Print(nil).Warn("Health Check: Client " + maskJID + " is not logged in")
			} else {
				log.Print(nil).Info("Health Check: Client " + maskJID + " is healthy")
			}
		}
	})

	cron.Start()
}
