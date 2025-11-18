package web

import (
	"fmt"
	"log_parser/pkg/models"
	"strings"

	"github.com/gin-gonic/gin"
)

func showAllLogs(c *gin.Context) {
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
		if len(values) > 0 && values[0] != "" {
			// result[key] = strings.Split(values[0], ",")
			result[key] = values
		}
	}
	fmt.Println(result)

	for key, vals := range result {
		if len(vals) > 0 {
			// join multiple values into one comma string: "INFO,DEBUG"
			joined := strings.Join(vals, ",")
			queryParts = append(queryParts, fmt.Sprintf("%s=%s", key, joined))
		}
	}

	// Execute query
	entries, err := models.QueryDB(DB, queryParts)
	if err != nil {
		c.HTML(500, "index.html", gin.H{"error": err.Error()})
		return
	}

	c.HTML(200, "index.html", gin.H{
		"entries": entries,
		"count":   len(entries),
	})
}
