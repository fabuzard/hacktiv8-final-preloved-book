package job

import (
	"log"
	"main/model"
	"time"

	"gorm.io/gorm"
)

func UpdateStatus(db *gorm.DB) {
	now := time.Now()

	// Lakukan update pada baris yang sesuai
	result := db.Model(&model.Transaction{}).
		Where("status = ? AND expiration_date < ?", "pending", now).
		Update("status", "fail")

	if result.Error != nil {
		log.Println(now, " - Failed to update payments:", result.Error)
		return
	}

	log.Println(now, " -", result.RowsAffected, "payments updated to 'fail'")
}
