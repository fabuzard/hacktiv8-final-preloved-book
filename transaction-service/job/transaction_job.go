package job

import (
	"log"
	"main/model"
	"time"

	"gorm.io/gorm"
)

func UpdateStatus(db *gorm.DB) {
	expirationThreshold := time.Now().Add(-30 * time.Minute) //30 minute expiration threshold

	// Step 1: Update expired pending transactions to "fail"
	updateResult := db.Model(&model.Transaction{}).
		Where("status = ? AND expiration_date < ?", "pending", expirationThreshold).
		Update("status", "fail")

	if updateResult.Error != nil {
		log.Println("Failed to update expired transactions:", updateResult.Error)
		return
	}

	log.Println(updateResult.RowsAffected, "transactions marked as 'fail'")

	// Step 2: Soft-delete the failed + expired transactions
	deleteResult := db.Where("status = ? AND expiration_date < ?", "fail", expirationThreshold).
		Delete(&model.Transaction{})

	if deleteResult.Error != nil {
		log.Println("Failed to delete failed expired transactions:", deleteResult.Error)
		return
	}

	log.Println(deleteResult.RowsAffected, "failed expired transactions deleted")
}
