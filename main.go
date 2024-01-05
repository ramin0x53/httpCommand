package main

import (
	"bufio"
	"flag"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/gin-gonic/gin"
)

func ReadFile(path string) []string {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

func main() {
	configPath := flag.String("c", "", "configPath")
	host := flag.String("h", "", "host")
	flag.Parse()

	configLines := ReadFile(*configPath)
	httpPaths := make(map[string]string)
	for _, configLine := range configLines {
		parsed := strings.Split(configLine, ":")
		httpPaths[parsed[0]] = parsed[1]
	}

	server := gin.Default()

	for path, command := range httpPaths {
		cmd := command
		server.GET(path, func(c *gin.Context) {
			out, err := exec.Command("sh", "-c", cmd).Output()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			c.JSON(200, gin.H{
				"result": string(out),
			})
		})
	}

	server.Run(*host)
}
