package internal

import (
	"time"

	"github.com/shabran01/go-whatsapp-multidevice-rest/pkg/log"
	pkgWhatsApp "github.com/shabran01/go-whatsapp-multidevice-rest/pkg/whatsapp"
)

func Startup() {
	log.Print(nil).Info("Running Startup Tasks")
	
	// Add startup delay to ensure database is ready
	time.Sleep(2 * time.Second)

	// Load All WhatsApp Client Devices from Datastore
	devices, err := pkgWhatsApp.WhatsAppDatastore.GetAllDevices()
	if err != nil {
		log.Print(nil).Error("Failed to Load WhatsApp Client Devices from Datastore: " + err.Error())
		return
	}

	// Do Reconnect for Every Device in Datastore
	for _, device := range devices {
		// Get JID from Datastore
		jid := pkgWhatsApp.WhatsAppDecomposeJID(device.ID.User)

		// Mask JID for Logging Information
		maskJID := jid[0:len(jid)-4] + "xxxx"

		// Print Restore Log
		log.Print(nil).Info("Restoring WhatsApp Client for " + maskJID)

		// Initialize WhatsApp Client
		pkgWhatsApp.WhatsAppInitClient(device, jid)

		// Reconnect WhatsApp Client WebSocket with delay
		go func(deviceJID string) {
			// Stagger reconnections to avoid overwhelming the server
			time.Sleep(time.Duration(len(pkgWhatsApp.WhatsAppClient)*2) * time.Second)
			
			err := pkgWhatsApp.WhatsAppReconnect(deviceJID)
			if err != nil {
				log.Print(nil).Error("Failed to reconnect " + pkgWhatsApp.WhatsAppDecomposeJID(deviceJID)[0:len(pkgWhatsApp.WhatsAppDecomposeJID(deviceJID))-4] + "xxxx: " + err.Error())
			}
		}(jid)
	}
}
