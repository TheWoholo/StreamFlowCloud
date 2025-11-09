package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()

	// --- CORS for your video player frontend ---
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://98.70.25.253, http://98.70.25.253:5173, http://localhost:5173, http://98.70.25.253:8081",
		AllowMethods:     "GET,OPTIONS", // Playback is GET-only
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowCredentials: true,
	}))

	// --- THIS IS THE CRITICAL LOGIC ---
	// This is the path *inside the container* where the shared volume will be.
	uploadDir := "/app/uploads"

	// This route serves the master playlist (e.g., /hls/my-video_hls/index.m3u8)
	// and all the video segments (e.g., /hls/my-video_hls/index0.ts)
	app.Static("/hls", uploadDir)

	// This route serves the original MP4 file (e.g., /uploads/test.mp4)
	app.Static("/uploads", uploadDir)
	// --- END OF CRITICAL LOGIC ---

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	port := 8083
	fmt.Printf("🚀 Playback service running at http://98.70.25.253:%d\n", port)

	// Ensure the directory exists (it will be the mount point)
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		log.Printf("Warning: Could not create %s: %v", uploadDir, err)
	}

	log.Fatal(app.Listen(fmt.Sprintf(":%d", port)))
}
