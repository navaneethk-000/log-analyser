package web

import (
	"fmt"
	"log_parser/pkg/models"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetAllLogs(c *gin.Context) {
	fmt.Println("show all")
	entries, err := models.QueryDB(DB, []string{}) //empty filter to get all logs
	fmt.Println(entries[0].Component)
	if err != nil {
		c.HTML(500, "index.html", gin.H{"error": err.Error()})
		return
	}

	c.HTML(200, "index.html", gin.H{
		"entries": entries,
		"count":   len(entries),
	})
}

func ExecuteFilterQuery(c *gin.Context) {
	queryParts := []string{}

	c.Request.ParseForm()
	formData := c.Request.PostForm

	result := make(map[string][]string)

	for key, values := range formData {
		if len(values) == 0 || values[0] == "" {
			continue
		}

		val := values[0]

		if key == "filter" {

			parts := strings.Split(val, " ")
			for _, p := range parts {
				if strings.TrimSpace(p) != "" {
					queryParts = append(queryParts, p)
				}
			}
		} else {
			// dropdowns
			queryParts = append(queryParts, fmt.Sprintf("%s=%s", key, val))
		}

	}

	fmt.Println(result)

	for key, vals := range result {
		if len(vals) > 0 {
			// join multiple values into one comma string: "INFO,ERROR"
			joined := strings.Join(vals, ",")
			queryParts = append(queryParts, fmt.Sprintf("%s=%s", key, joined))
		}
	}

	// Execute query
	entries, err := models.QueryDB(DB, queryParts)
	if err != nil {
		c.JSON(500, err)
		return
	}

	c.JSON(200, gin.H{
		"entries": entries,
		"count":   len(entries),
	})
}
