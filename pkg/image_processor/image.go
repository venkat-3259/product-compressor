package imageprocessor

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"zocket/app/models"
	"zocket/app/queries"
	"zocket/pkg/utils"

	"github.com/nfnt/resize"
	"github.com/streadway/amqp"
)

func InitChannel(imageProcessor *ImageProcessor) {

	defer log.Println("exiting from image processor")

	imageProcessor.Subscriber()
}

type ImageProcessor struct {
	Ctx     context.Context
	Channel *amqp.Channel
	Queries *queries.Queries
}

func (ch *ImageProcessor) ImageDownload(images models.ProductLinks) {

	var compressedPaths []string

	for _, link := range images.Links {

		// Download image
		req, err := utils.PrepareHTTPRequest(nil, http.MethodGet, link)
		if err != nil {
			log.Println("Failed to prepare request", err)
			return
		}

		StatusCode, imageBytes, err := utils.SendHTTPRequest(req)
		if err != nil {
			log.Println("Failed to send prepared request", err)
			return
		}

		log.Println("Status-code :", StatusCode)

		// Compress and save image locally
		compressedPath := saveCompressedImage(imageBytes)
		if compressedPath == "" {
			return
		}

		compressedPaths = append(compressedPaths, compressedPath)
	}
	// Send a message to RabbitMQ
	ch.sendMessage(images.ID, compressedPaths)

}

func saveCompressedImage(imageBytes []byte) string {
	// Generate a unique filename based on timestamp
	timestamp := time.Now().UnixNano()
	filename := fmt.Sprintf("compressed_image_%d.jpg", timestamp)

	// Specify the directory where you want to save the compressed images
	saveDir := "./compressed_images"

	// Create the directory if it doesn't exist
	if err := os.MkdirAll(saveDir, os.ModePerm); err != nil {
		return ""
	}

	// Specify the full path of the saved image
	fullPath := fmt.Sprintf("%s/%s", saveDir, filename)

	// Create a file for the compressed image
	outFile, err := os.Create(fullPath)
	if err != nil {
		return ""
	}
	defer outFile.Close()

	// Decode the original image
	originalImage, _, err := image.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		return ""
	}

	// Resize the image (you can adjust the dimensions as needed)
	resizedImage := resize.Resize(300, 0, originalImage, resize.Lanczos3)

	// Encode and save the resized image as JPEG
	err = jpeg.Encode(outFile, resizedImage, nil)
	if err != nil {
		return ""
	}

	return fullPath
}
func (ch *ImageProcessor) sendMessage(id int, imagePaths []string) {
	// Define the queue name
	queueName := "image_queue"

	// Declare a queue
	_, err := ch.Channel.QueueDeclare(
		queueName, // Queue name
		false,     // Durable
		false,     // Delete when unused
		false,     // Exclusive
		false,     // No-wait
		nil,       // Arguments
	)
	if err != nil {
		log.Println("Failed to declare a queue:", err)
		return
	}

	// Create a message containing the ID and image paths
	messageBody := fmt.Sprintf("%d:%s", id, strings.Join(imagePaths, ","))

	// Publish the message to the queue
	err = ch.Channel.Publish(
		"",        // Exchange
		queueName, // Routing key (queue name)
		false,     // Mandatory
		false,     // Immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(messageBody),
		})
	if err != nil {
		log.Println("Failed to publish a message:", err)
		return
	}
}

func (ch *ImageProcessor) Subscriber() {
	// Define the queue name to consume from
	queueName := "image_queue"

	// Declare the queue to make sure it exists
	_, err := ch.Channel.QueueDeclare(
		queueName, // Queue name
		false,     // Durable
		false,     // Delete when unused
		false,     // Exclusive
		false,     // No-wait
		nil,       // Arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %s", err)
	}

	// Consume messages from the queue
	msgs, err := ch.Channel.Consume(
		queueName, // Queue name
		"",        // Consumer name
		true,      // Auto-acknowledge (acknowledge messages automatically)
		false,     // Exclusive
		false,     // No-local
		false,     // No-wait
		nil,       // Arguments
	)
	if err != nil {
		log.Fatalf("Failed to start consumer: %s", err)
	}

	log.Println("Consumer started. Waiting for messages...")

	// Process messages in a loop
	for msg := range msgs {
		// Extract and process the message body
		messageBody := string(msg.Body)
		log.Printf("Received message: %s\n", messageBody)

		// Split the message to get the ID and image paths
		parts := strings.Split(messageBody, ":")
		if len(parts) != 2 {
			log.Println("Invalid message format")
			continue
		}

		id := parts[0]
		imagePaths := strings.Split(parts[1], ",")
		if len(imagePaths) > 0 {
			var image models.ProductImages

			image.ID, err = strconv.Atoi(id)
			if err != nil {
				log.Println("Failed to convert ID")
			}
			image.Paths = imagePaths

			err = ch.Queries.ProductQueries.AddCompressedImages(ch.Ctx, image)
			if err != nil {
				log.Println("Failed to put entry in database, ", err)
			}
		}
		log.Println("ID :", id)
		log.Println("Paths: ", imagePaths)
	}
}
