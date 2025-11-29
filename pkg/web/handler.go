package web

import (
	"fmt"
	"log_parser/pkg/models"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Filter struct {
	Level     []string `json:"level"`
	Component []string `json:"component"`
	Host      []string `json:"host"`
	RequestID []string `json:"request_id"`
	StartTime string   `json:"startTime"`
	EndTime   string   `json:"endTime"`
}

type LogsResponse struct {
	Entries []models.Entry `json:"entries"`
	Total   int64          `json:"total_count"`
}

type Error struct {
	Error string `json:"error"`
}

func FilterPaginatedLogs(c *gin.Context) {

	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "100"))

	var filter Filter
	if err := c.ShouldBindJSON(&filter); err != nil {
		if err.Error() != "EOF" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return

		}
	}

	offset := page * pageSize

	var entries []models.Entry
	var total int64

	countQuery := DB.Model(&models.Entry{})
	db := DB.
		Model(&models.Entry{}).
		Preload("Level").
		Preload("Component").
		Preload("Host")

	if len(filter.Level) > 0 {
		db = db.Joins("Level").Where(`"Level"."level" IN ?`, filter.Level)
		countQuery = countQuery.Joins("Level").Where(`"Level"."level" IN ?`, filter.Level)
	}

	if len(filter.Host) > 0 {
		db = db.Joins("Host").Where(`"Host"."host" IN ?`, filter.Host)
		countQuery = countQuery.Joins("Host").Where(`"Host"."host" IN ?`, filter.Host)
	}
	if len(filter.Component) > 0 {
		db = db.Joins("Component").Where(`"Component"."component" IN ?`, filter.Component)
		countQuery = countQuery.Joins("Component").Where(`"Component"."component" IN ?`, filter.Component)
	}
	if len(filter.RequestID) > 0 {
		db = db.Where(`entries.request_id IN ?`, filter.RequestID)
		countQuery = countQuery.Where(`entries.request_id IN ?`, filter.RequestID)
	}

	if filter.StartTime != "" && filter.EndTime != "" {
		db = db.Where("time_stamp BETWEEN ? AND ?", filter.StartTime, filter.EndTime)
	} else if filter.StartTime != "" {
		db = db.Where("time_stamp >= ?", filter.StartTime)
	} else if filter.EndTime != "" {
		db = db.Where("time_stamp <= ?", filter.EndTime)
	}

	// count with same filters
	if err := countQuery.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err := db.Order("entries.id ASC").
		Offset(offset).
		Limit(pageSize).
		Find(&entries).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, &Error{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &LogsResponse{
		Entries: entries,
		Total:   total,
	})
	fmt.Println("page:", page, "pageSize:", pageSize, "offset:", offset)

}

func GetAllLogs(c *gin.Context) {
	fmt.Println("show all")
	entries, err := models.QueryDB(DB, []string{}) //empty filter to get all logs
	fmt.Println(entries[0].Component)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"entries": entries,
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
