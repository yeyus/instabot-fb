package main

import (
	"log"
	"os"

	"crypto/tls"
	"net/http"

	"github.com/abhinavdahiya/go-messenger-bot"
	"github.com/joho/godotenv"
	"github.com/yeyus/instabot-fb/forecast"
	"github.com/yeyus/instabot-fb/handlers"
	"github.com/yeyus/witgo/v1/witgo"
	"golang.org/x/crypto/acme/autocert"
)

var (
	port              string
	domain            string
	fbPageAccessToken string
	fbVerifyToken     string
	fbSecretToken     string
	witaiToken        string
	generateCerts     bool
)

func main() {
	loadConfigs()

	bot := mbotapi.NewBotAPI(fbPageAccessToken, fbVerifyToken, fbSecretToken)
	// DEBUG:
	// bot.Debug = true

	// instantiate wit.ai
	client := witgo.NewClient(witaiToken)
	// DEBUG:
	// client.HttpClient = witgo.NewLoggingHttpClient(os.Stderr, client.HttpClient)

	handler := handlers.NewMessengerHandler(bot)
	handler.Actions["getForecast"] = func(session *witgo.Session, entities witgo.EntityMap) (*witgo.Session, error) {
		log.Printf("getForecast was called")
		if location, err := entities.FirstEntityValue("location"); err != nil {
			session.Context.Set("missingLocation", true)
		} else {
			f, err := forecast.GetForecast(location)
			if err != nil {
				session.Context.Merge(witgo.Context{
					"location": location,
					"notFound": "notFound",
				})
				return session, err
			}

			session.Context.Merge(witgo.Context{
				"location":    location,
				"forecast":    f["channel.item.condition.text"],
				"temperature": f["channel.item.condition.temp"],
			})
		}

		return session, nil
	}
	wg := witgo.NewWitgo(client, handler)

	handler.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("instabot-fb says hello world!"))
	})

	http.Handle("/", handler.Mux)

	go func() {
		err := wg.Process(handler)
		if err != nil {
			panic(err)
		}
	}()

	if generateCerts {
		certManager := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(domain),
			Cache:      autocert.DirCache("certs"),
		}

		server := &http.Server{
			Addr: ":" + port,
			TLSConfig: &tls.Config{
				GetCertificate: certManager.GetCertificate,
			},
		}

		log.Printf("Initiating https server for domain %s and port %s", domain, port)
		log.Printf("Endpoints are automatically signed with Let's Encrypt certificates")
		log.Printf("Go to FB's developer console and update your webhook to point to:")
		log.Printf("https://" + domain + ":" + port + "/webhook")
		log.Fatal(server.ListenAndServeTLS("", ""))
	} else {
		log.Printf("Initiating http server for domaing %s and port %s", domain, port)
		log.Printf("you will need to route calls through a proxy such as nginx as endpoints must be signed to use Messenger API.")
		log.Printf("Go to FB's developer console and update your webhook to point to:")
		log.Printf("http://" + domain + ":" + port + "/webhook")
		log.Fatal(http.ListenAndServe(":"+port, nil))
	}

}

func loadConfigs() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port = os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	domain = os.Getenv("DOMAIN")
	if domain == "" {
		log.Fatal("$DOMAIN must be set")
	}

	fbPageAccessToken = os.Getenv("FB_ACCESS_TOKEN")
	if fbPageAccessToken == "" {
		log.Fatal("$FB_ACCESS_TOKEN must be set")
	}

	fbVerifyToken = os.Getenv("FB_VERIFY_TOKEN")
	if fbVerifyToken == "" {
		log.Fatal("$FB_VERIFY_TOKEN must be set")
	}

	fbSecretToken = os.Getenv("FB_SECRET_TOKEN")
	if fbSecretToken == "" {
		log.Fatal("$FB_SECRET_TOKEN must be set")
	}

	witaiToken = os.Getenv("WITAI_TOKEN")
	if witaiToken == "" {
		log.Fatal("$WITAI_TOKEN must be set")
	}

	if os.Getenv("GENERATE_CERTS") == "TRUE" {
		generateCerts = true
	}
}
