package main

import (
        "context"
        "fmt"
		"os"
        "sync"
		"flag"
		"log"
		"path"
		
		"io/ioutil"

		"golang.org/x/oauth2"
        "cloud.google.com/go/pubsub"
		"google.golang.org/api/option"
		"google.golang.org/api/iterator"

)

const DEFAULT_SUBSCRIPTION_NAME string = "pubsub-monitor"

var (
	accessTokenPtr *string
    projectPtr *string
	logDirPtr *string
	publishPathPtr *string
 
    mu sync.Mutex
	counter int = 1
)

type StaticTokenSource struct {
	AccessToken string
}

func (mts StaticTokenSource) Token()  (*oauth2.Token, error) {
	re := oauth2.Token{}
	re.AccessToken = mts.AccessToken
	return &re, nil
}

func obtain_pubsub_client() (*pubsub.Client, error) {
	opts := []option.ClientOption{}
	
	static_token := *accessTokenPtr
	if static_token != "" {
		var myTokenSource = StaticTokenSource{}
		myTokenSource.AccessToken = static_token
		opts = append(opts, option.WithTokenSource(myTokenSource))
	}
	
	ctx := context.Background()
	pubsub_client, err := pubsub.NewClient(ctx, *projectPtr, opts...)
	if err != nil {
			return nil, fmt.Errorf("pubsub.NewClient: %v", err)
	}

	return pubsub_client, nil
}

func watch_topic(topic_name string, topic *pubsub.Topic) {
	pubsub_client, err := obtain_pubsub_client()
	if err != nil {
		log.Fatalf("Error while obtaining pubsub client: %s", err)
	}
	
	subscription_name := DEFAULT_SUBSCRIPTION_NAME+"--"+topic_name

	ctx:= context.Background()
	sub := pubsub_client.Subscription(subscription_name)
	isSubExist, err := sub.Exists(ctx)
	if err != nil {
		log.Fatalf("Error while working with subscription", err)
	}
	
	if ! isSubExist {
		sub, err = pubsub_client.CreateSubscription(context.Background(), subscription_name, pubsub.SubscriptionConfig{Topic: topic})
		if err != nil {
			log.Fatalf("Error while creating subscription: %s", err)
		}
		log.Printf("Subscription created: %s", subscription_name)
	} else {
		log.Printf("Reusing existing subscription %s", subscription_name)
	}
	
	
	// https://stackoverflow.com/questions/62635263/gcp-pub-sub-using-goroutines-to-make-multiple-subscriber-running-in-one-applica
	sub.ReceiveSettings.NumGoroutines = 1
	
	err = sub.Receive(context.Background(), func(ctx context.Context, msg *pubsub.Message) {
			mu.Lock()
			defer mu.Unlock()
			log.Printf("Got message in topic %s with %d bytes\n", topic_name, len(msg.Data))
			
			counter += 1

			bin_file := fmt.Sprintf("%04d-%s.bin", counter, topic_name)
			full_bin_file:= path.Join(*logDirPtr, bin_file)

			err := ioutil.WriteFile(full_bin_file, msg.Data, 0o644)

			if err != nil {
				log.Fatal(err)
			}

			msg.Ack()
	})
	if err != nil {
			log.Fatalf("Receive: %v", err)
	}		
}

func pullMsgs(topics []string) error {
	pubsub_client, err := obtain_pubsub_client()
        if err != nil {
                return fmt.Errorf("pubsub.NewClient: %v", err)
        }
        defer pubsub_client.Close()
		
		if len(topics) == 0 {
			log.Printf("Fetching a list of topics in %s", *projectPtr)
		    topic_iterator:= pubsub_client.Topics(context.Background())
			for {
				topic, err := topic_iterator.Next()
				if err == iterator.Done {
					break
				}
				if err != nil {
					return err
				}
				topics= append(topics, topic.ID())
			}			
			log.Printf("%d topics found", len(topics))
		}
		
		ctx := context.Background()
		
		var publish_msg []byte
		if *publishPathPtr != "" {
			publish_msg, err = ioutil.ReadFile(*publishPathPtr)
			if err != nil {
				log.Fatal(err)
			}
		}
		
		for _, topic_name:= range topics {
			log.Printf("Subscribing to topic %s", topic_name)
			topic:= pubsub_client.Topic(topic_name)
			if topic == nil {
				log.Fatalf("Topic %s not found", topic_name)
			}

		    if *publishPathPtr != "" {
			    log.Printf("Publishing message to topic: %s", topic_name)
				res := topic.Publish(ctx, &pubsub.Message{Data: publish_msg})
				_, err := res.Get(ctx)
				if err != nil {
					log.Fatalf("Error while publishing message to topic %s: %s", topic_name, err)
				}
			} else {
				go watch_topic(topic_name, topic)
			}
		}
		
		if *publishPathPtr != "" {
		   log.Printf("Publishing complete")
		   os.Exit(0)
		}
		
		log.Printf("Subscriptions complete, receiving messages for all topics")
		
        return nil
}

func main() {

   accessTokenPtr = flag.String("access-token", "", "access token to use $(curl --silent http://metadata.google.internal/computeMetadata/v1/instance/service-accounts/default/token -H Metadata-Flavor:Google | jq -r .access_token) or $(gcloud auth print-access-token)")
   projectPtr = flag.String("project", "", "project to inspect")
   logDirPtr = flag.String("log-dir", "", "directory to save the messages to")
   publishPathPtr = flag.String("publish-msg-path", "", "publishes the content of the specified file to the topics")
   
   flag.Parse()
   
   if *projectPtr == "" {
	  log.Fatal("Error: you must provide the project ID")
   }
   if *publishPathPtr == "" && *logDirPtr == "" {
	  log.Fatal("Error: you must specify the logdir")
   }
   
   tail := flag.Args() // the list of topics to limit the monitoring to

   err := pullMsgs(tail)
   if err != nil {
       log.Fatal(err)
   }
   
   // wait forever
   select {}
   
}
