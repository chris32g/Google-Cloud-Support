package main

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"time"

	"cloud.google.com/go/datastore"
	"contrib.go.opencensus.io/exporter/stackdriver"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
	"google.golang.org/api/option"
	"google.golang.org/appengine"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

func main() {
	projectId := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectId == "" {
		log.Fatal("Missing env var GOOGLE_CLOUD_PROJECT")
	}

	// Setup stackdriver tracing
	if !appengine.IsDevAppServer() {
		traceExporter, err := stackdriver.NewExporter(stackdriver.Options{})
		if err != nil {
			log.Fatalf("Failed to create the Stackdriver exporter: %v", err)
		}
		trace.RegisterExporter(traceExporter)
		trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
		defer traceExporter.Flush()
	}

	// Setup datastore client
	ctx := context.Background()
	options := []option.ClientOption{
		option.WithGRPCConnectionPool(runtime.GOMAXPROCS(0)),
		option.WithGRPCDialOption(grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time: 5 * time.Minute,
		})),
	}
	datastoreClient, err := datastore.NewClient(ctx, projectId, options...)
	if err != nil {
		log.Fatal(errors.Wrap(err, "NewDatastoreDB: Failed to create datastore client"))
	}

	// Setup HTTP routes
	http.Handle("/", GetIndex())
	http.Handle("/animals/random", GetAnimal(datastoreClient))

	// Run appengine main loop
	appengine.Main()
}

func GetIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; encoding=utf-8")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("Pong."))
		if err != nil {
			log.Printf("Error: Failed to write response: %s", err)
		}
	}
}

func GetAnimal(datastoreClient *datastore.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		animalName := animalNames[rand.Intn(len(animalNames) - 1)]

		animal := new(Animal)
		key := datastore.NameKey("Animal", animalName, nil)
		err := datastoreClient.Get(r.Context(), key, animal)
		if err == datastore.ErrNoSuchEntity {
			err = nil

			// Create animal, since it does not exist, yet
			animal.Name = animalName

			_, err := datastoreClient.Put(r.Context(), key, animal)
			if err != nil {
				log.Printf("Error: Failed to put animal: %s", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Header().Set("created-animal", "true")
		}
		if err != nil {
			log.Printf("Error: Failed to get animal: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		data, err := json.Marshal(animal)
		if err != nil {
			log.Printf("Error: Failed to marshal animal: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; encoding=utf-8")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(data)
		if err != nil {
			log.Printf("Error: Failed to write response: %s", err)
		}
	}
}

type Animal struct {
	Name string
}

var animalNames = []string{
	"Ace",
	"Alpha",
	"Apollo",
	"Archer",
	"Bandit",
	"Blaze",
	"Bolt",
	"Chess",
	"Cobra",
	"Copper",
	"Dirt",
	"Enigma",
	"Flame",
	"Flare",
	"Frisk",
	"Gunner",
	"Harley",
	"Hunter",
	"Ice",
	"Juno",
	"Kong",
	"Maze",
	"Mocha",
	"Pip",
	"Pistol",
	"Quartz",
	"Quiz",
	"Raven",
	"Red",
	"Rex",
	"Rhino",
	"Rider",
	"Ripley",
	"Ripper",
	"River",
	"Saber",
	"Shadow",
	"Shank",
	"Shot",
	"Solo",
	"Spike",
	"Stinger",
	"Storm",
	"Tango",
	"Thunder",
	"Titus",
	"Tricks",
	"Valley",
	"Vice",
	"Wolf",
}
