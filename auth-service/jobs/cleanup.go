package jobs

import (
	"auth-service/repository"
	"log"
	"time"
)

func StartCleanupJob(repo repository.AuthRepository) {
	go func() {
		for {
			now := time.Now()
			// Run at midnight
			next := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
			duration := next.Sub(now)

			time.Sleep(duration)

			log.Println("🧹 Running cleanup job...")
			if err := repo.DeleteInactiveUsersOver30Days(); err != nil {
				log.Printf("❌ Cleanup job failed: %v", err)
			} else {
				log.Println("✅ Inactive users cleanup completed.")
			}
		}

	}()
}
