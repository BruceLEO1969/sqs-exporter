package collector

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	cmap "github.com/orcaman/concurrent-map"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	visibleMessageGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "sqs_messages_visible",
		Help: "Type: Gauge, The number of available messages in queue(s). Use the name of the queue as the label to get the messages of a specific queue e.g `sqs_messages_visible{queue_name=\"<QUEUE NAME>\"}`.",
	}, []string{"queue_name"})
	delayedMessageGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "sqs_messages_delayed",
		Help: "Type: Gauge, The number of messages waiting to be added into queue(s). Use the name of the queue as the label to get the messages of a specific queue e.g `sqs_messages_delayed{queue_name=\"<QUEUE NAME>\"}`.",
	}, []string{"queue_name"})
	invisibleMessageGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "sqs_messages_invisible",
		Help: "Type: Gauge, The number of messages in flight in queue(s). Use the name of the queue as the label to get the messages of a specific queue e.g `sqs_messages_invisible{queue_name=\"<QUEUE NAME>\"}`.",
	}, []string{"queue_name"})
)

func init() {
	prometheus.MustRegister(visibleMessageGauge)
	prometheus.MustRegister(delayedMessageGauge)
	prometheus.MustRegister(invisibleMessageGauge)
}

// MonitorSQS Retrieves the attributes of all allowed queues from SQS and appends the metrics
func MonitorSQS() error {
	queues, _, err := getQueues()
	if err != nil {
		return fmt.Errorf("[MONITORING ERROR]: Error occurred while retrieve queues info from SQS: %v", err)
	}

	queues.IterCb(func(key string, v interface{}) {
		attr, ok := v.(*sqs.GetQueueAttributesOutput)

		if !ok {
			return
		}

		msgAvailable, msgError := strconv.ParseFloat(*attr.Attributes["ApproximateNumberOfMessages"], 64)
		msgDelayed, delayError := strconv.ParseFloat(*attr.Attributes["ApproximateNumberOfMessagesDelayed"], 64)
		msgNotVisible, invisibleError := strconv.ParseFloat(*attr.Attributes["ApproximateNumberOfMessagesNotVisible"], 64)

		if msgError != nil {
			fmt.Errorf("Error in converting ApproximateNumberOfMessages: %v", msgError)
			return
		}
		visibleMessageGauge.WithLabelValues(key).Set(msgAvailable)

		if delayError != nil {
			fmt.Errorf("Error in converting ApproximateNumberOfMessagesDelayed: %v", delayError)
			return
		}
		delayedMessageGauge.WithLabelValues(key).Set(msgDelayed)

		if invisibleError != nil {
			fmt.Errorf("Error in converting ApproximateNumberOfMessagesNotVisible: %v", invisibleError)
			return
		}
		invisibleMessageGauge.WithLabelValues(key).Set(msgNotVisible)

		fmt.Printf("sqs_messages_visible{queue_name=\"%s\"} %+v\n", key, msgAvailable)
		fmt.Printf("sqs_messages_delayed{queue_name=\"%s\"} %+v\n", key, msgDelayed)
		fmt.Printf("sqs_messages_invisible{queue_name=\"%s\"} %+v\n", key, msgNotVisible)
	})

	return nil
}

func getQueueName(url string) (queueName string) {
	queue := strings.Split(url, "/")
	queueName = queue[len(queue)-1]
	return
}

func getQueues() (queues cmap.ConcurrentMap, tags cmap.ConcurrentMap, err error) {
	sess := session.Must(session.NewSession())
	client := sqs.New(sess)
	result, err := client.ListQueues(nil)

	if err != nil {
		return nil, nil, err
	}

	//fmt.Println(result)
	if result.QueueUrls == nil {
		err = fmt.Errorf("SQS did not return any QueueUrls")
		return nil, nil, err
	}

	queues = cmap.New()
	tags = cmap.New()

	wg := sync.WaitGroup{}

	for _, urls := range result.QueueUrls {
		urls := urls
		wg.Add(1)
		go func() {
			defer wg.Done()
			params := &sqs.GetQueueAttributesInput{
				QueueUrl: aws.String(*urls),
				AttributeNames: []*string{
					aws.String("ApproximateNumberOfMessages"),
					aws.String("ApproximateNumberOfMessagesDelayed"),
					aws.String("ApproximateNumberOfMessagesNotVisible"),
				},
			}

			tagsParams := &sqs.ListQueueTagsInput{
				QueueUrl: aws.String(*urls),
			}

			resp, err := client.GetQueueAttributes(params)
			if err != nil {
				fmt.Errorf("GetQueueAttributes Error: %v", err)
				return
			}

			tagsResp, err := client.ListQueueTags(tagsParams)
			if err != nil {
				fmt.Errorf("ListQueueTags Error: %v", err)
				return
			}

			queueName := getQueueName(*urls)

			queues.Set(queueName, resp)
			tags.Set(queueName, tagsResp)
		}()
	}
	wg.Wait()

	return queues, tags, nil
}
