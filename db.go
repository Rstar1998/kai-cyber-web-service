package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"time"
)

// sqlite table schema where we store all Vulnerabilities 
// primary key is combo of id and SourceFile since different files can have same Vulnerabilities
type Vulnerability struct {
	ID             string    `gorm:"primaryKey" json:"id"`
	Severity       string    `json:"severity"`
	CVSS           float64   `json:"cvss"`
	Status         string    `json:"status"`
	PackageName    string    `json:"package_name"`
	CurrentVersion string    `json:"current_version"`
	FixedVersion   string    `json:"fixed_version"`
	Description    string    `json:"description"`
	PublishedDate  string    `json:"published_date"`
	Link           string    `json:"link"`
	RiskFactors    string    `json:"risk_factors"` // Store as a JSON string
	SourceFile     string    `gorm:"primaryKey" json:"source_file"`
	ScanTime       time.Time `json:"scan_time"`
}

var db *gorm.DB

// create sqlite db named vulnerabilities.db and create Vulnerability table
func InitDB() {
	var err error
	db, err = gorm.Open(sqlite.Open("vulnerabilities.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}
	if err := db.AutoMigrate(&Vulnerability{}); err != nil {
		panic("failed to migrate database: " + err.Error())
	}
}

// get total no of rows in Vulnerability table
func GetTotalVulnerabilities() (int64, error) {
	var count int64
	if err := db.Model(&Vulnerability{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
