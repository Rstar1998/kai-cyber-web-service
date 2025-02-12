package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"sync"
	"time"
)

// schema that should match the content of json file
type ScanResult struct {
	ScanResults struct {
		ScanID          string `json:"scan_id"`
		Timestamp       string `json:"timestamp"`
		ScanStatus      string `json:"scan_status"`
		ResourceType    string `json:"resource_type"`
		ResourceName    string `json:"resource_name"`
		Vulnerabilities []struct {
			ID             string   `json:"id"`
			Severity       string   `json:"severity"`
			CVSS           float64  `json:"cvss"`
			Status         string   `json:"status"`
			PackageName    string   `json:"package_name"`
			CurrentVersion string   `json:"current_version"`
			FixedVersion   string   `json:"fixed_version"`
			Description    string   `json:"description"`
			PublishedDate  string   `json:"published_date"`
			Link           string   `json:"link"`
			RiskFactors    []string `json:"risk_factors"`
		} `json:"vulnerabilities"`
	} `json:"scanResults"`
}

// add Vulnerabilities to db
func saveVulnerabilities(scanResults []ScanResult, filename string) {
	for _, scan := range scanResults {
		for _, vuln := range scan.ScanResults.Vulnerabilities {
			riskFactorsJSON, _ := json.Marshal(vuln.RiskFactors) // Convert []string to JSON string

			vulnerability := Vulnerability{
				ID:             vuln.ID,
				Severity:       vuln.Severity,
				CVSS:           vuln.CVSS,
				Status:         vuln.Status,
				PackageName:    vuln.PackageName,
				CurrentVersion: vuln.CurrentVersion,
				FixedVersion:   vuln.FixedVersion,
				Description:    vuln.Description,
				PublishedDate:  vuln.PublishedDate,
				Link:           vuln.Link,
				RiskFactors:    string(riskFactorsJSON), // Store as a JSON string
				SourceFile:     filename,
				ScanTime:       time.Now(),
			}
			db.Create(&vulnerability) //insert record : we can do batch insert as well
		}
	}
}

func fetchAndProcessFile(filename string, repo string, fileStatus map[string]string, wg *sync.WaitGroup) {
	defer wg.Done()

	owner, repo, err := ExtractOwnerRepo(repo) // get repo name an owner name
	if err != nil {
		fmt.Println("Error:", err)
		fileStatus[filename] = fmt.Sprintf("error: %s", err.Error())
		return
	}

	// create a raw file link
	url := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/refs/heads/main/%s", owner, repo, filename)

	resp, err := http.Get(url)

	if err != nil {
		fmt.Println("Failed to fetch file:", err)
		fileStatus[filename] = fmt.Sprintf("Failed to fetch file: %s", err.Error())
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body) // read json file
	if err != nil {
		fmt.Println("Failed to read file contents:", err)
		fileStatus[filename] = fmt.Sprintf("Failed to read file contents: %s", err.Error())
		return
	}

	var scanResults []ScanResult
	if err := json.Unmarshal(body, &scanResults); err != nil { // match json file with schema
		fmt.Println("Failed to parse JSON:", err)
		fileStatus[filename] = fmt.Sprintf("Failed to parse JSON: %s", err.Error())
		return
	}

	saveVulnerabilities(scanResults, filename) //  add Vulnerabilities to db
	fileStatus[filename] = "Sucess"

}

func ScanRepo(c *gin.Context) {

	var request struct { // request structure
		Repo  string   `json:"repo" binding:"required"`
		Files []string `json:"files" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil { // check if request structure matches the received on
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Exec("DELETE FROM vulnerabilities") // for each load in order to prevent duplicates I am deleteing old records and inserting the new ones

	// Map to store the status of each file
	fileStatus := make(map[string]string)

	var wg sync.WaitGroup //  to process all files concurrently
	for _, file := range request.Files {
		wg.Add(1)
		go fetchAndProcessFile(file, request.Repo, fileStatus, &wg)
	}
	wg.Wait()

	count, err := GetTotalVulnerabilities() // get total records inserted

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve total rows: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Scan complete", "total_rows": count, "file_status": fileStatus})

}

func QueryVulnerabilities(c *gin.Context) {

	var request struct { // request structure
		Filters struct {
			Severity string `json:"severity" binding:"required"`
		} `json:"filters" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil { // check if request structure matches the received on
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	type VulnerabilityOutput struct {
		ID             string   `json:"id"`
		Severity       string   `json:"severity"`
		CVSS           float64  `json:"cvss"`
		Status         string   `json:"status"`
		PackageName    string   `json:"package_name"`
		CurrentVersion string   `json:"current_version"`
		FixedVersion   string   `json:"fixed_version"`
		Description    string   `json:"description"`
		PublishedDate  string   `json:"published_date"`
		Link           string   `json:"link"`
		RiskFactors    []string `json:"risk_factors"` // Slice of strings
	}
	var outputVulnerabilities []VulnerabilityOutput

	var results []Vulnerability //  get Vulnerability based on severity
	db.Where("severity = ?", request.Filters.Severity).Find(&results)

	if len(results) == 0 {
		c.JSON(http.StatusOK, results)
		return 
	}

	for _, vuln := range results {
		var riskFactors []string
		err := json.Unmarshal([]byte(vuln.RiskFactors), &riskFactors) // Unmarshal JSON string
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		outputVuln := VulnerabilityOutput{
			ID:             vuln.ID,
			Severity:       vuln.Severity,
			CVSS:           vuln.CVSS,
			Status:         vuln.Status,
			PackageName:    vuln.PackageName,
			CurrentVersion: vuln.CurrentVersion,
			FixedVersion:   vuln.FixedVersion,
			Description:    vuln.Description,
			PublishedDate:  vuln.PublishedDate,
			Link:           vuln.Link,
			RiskFactors:    riskFactors, // Assign the slice
		}
		outputVulnerabilities = append(outputVulnerabilities, outputVuln)
	}
	
	c.JSON(http.StatusOK, outputVulnerabilities) // send results
}
