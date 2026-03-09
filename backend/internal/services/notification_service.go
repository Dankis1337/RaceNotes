package services

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strings"
	"time"

	"github.com/levferril/racenotes/backend/internal/models"
	"gorm.io/gorm"
)

type NotificationService struct {
	DB *gorm.DB
}

func NewNotificationService(db *gorm.DB) *NotificationService {
	return &NotificationService{DB: db}
}

// StartReminderCron runs a goroutine that checks for tomorrow's races at 20:00 daily.
func (s *NotificationService) StartReminderCron() {
	smtpHost := os.Getenv("SMTP_HOST")
	if smtpHost == "" {
		log.Println("[Notifications] SMTP_HOST not configured, email reminders disabled")
		return
	}

	go func() {
		for {
			now := time.Now()
			// Schedule next run at 20:00
			next := time.Date(now.Year(), now.Month(), now.Day(), 20, 0, 0, 0, now.Location())
			if now.After(next) {
				next = next.Add(24 * time.Hour)
			}
			time.Sleep(time.Until(next))

			s.sendTomorrowReminders()
		}
	}()

	log.Println("[Notifications] Email reminder cron started (daily at 20:00)")
}

func (s *NotificationService) sendTomorrowReminders() {
	tomorrow := time.Now().Add(24 * time.Hour).Format("2006-01-02")

	var races []models.Race
	s.DB.Where("date = ? AND is_completed = false", tomorrow).
		Preload("Setup").
		Find(&races)

	if len(races) == 0 {
		return
	}

	// Group by user
	userRaces := make(map[uint][]models.Race)
	for _, r := range races {
		userRaces[r.UserID] = append(userRaces[r.UserID], r)
	}

	for userID, raceList := range userRaces {
		var user models.User
		if err := s.DB.First(&user, userID).Error; err != nil {
			continue
		}
		if user.Email == "" {
			continue
		}

		body := buildReminderEmail(user.Name, raceList)
		if err := sendEmail(user.Email, "RaceNotes: Tomorrow's Race Reminder", body); err != nil {
			log.Printf("[Notifications] Failed to send email to %s: %v", user.Email, err)
		} else {
			log.Printf("[Notifications] Reminder sent to %s for %d race(s)", user.Email, len(raceList))
		}
	}
}

func buildReminderEmail(userName string, races []models.Race) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Hi %s!\n\nYou have %d race(s) tomorrow:\n\n", userName, len(races)))

	for _, r := range races {
		sb.WriteString(fmt.Sprintf("--- %s ---\n", r.Name))
		sb.WriteString(fmt.Sprintf("Type: %s\n", r.Type))

		if r.Setup != nil {
			sb.WriteString(fmt.Sprintf("Setup: %s (%s, %s)\n", r.Setup.Name, r.Setup.BikeName, r.Setup.Tires))
		} else if r.BikeName != nil {
			sb.WriteString(fmt.Sprintf("Bike: %s\n", *r.BikeName))
			if r.Tires != nil {
				sb.WriteString(fmt.Sprintf("Tires: %s\n", *r.Tires))
			}
		}

		if r.TirePressureFront != nil {
			sb.WriteString(fmt.Sprintf("Pressure: Front %.2f bar", *r.TirePressureFront))
			if r.TirePressureRear != nil {
				sb.WriteString(fmt.Sprintf(", Rear %.2f bar", *r.TirePressureRear))
			}
			sb.WriteString("\n")
		}

		if r.Conditions != nil {
			sb.WriteString(fmt.Sprintf("Weather: %s", *r.Conditions))
			if r.Temperature != nil {
				sb.WriteString(fmt.Sprintf(", %d°C", *r.Temperature))
			}
			sb.WriteString("\n")
		}

		if r.NutritionPlan != nil {
			sb.WriteString(fmt.Sprintf("Nutrition: %s\n", *r.NutritionPlan))
		}

		sb.WriteString("\n")
	}

	sb.WriteString("Good luck!\n— RaceNotes")
	return sb.String()
}

func sendEmail(to, subject, body string) error {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASSWORD")
	fromEmail := os.Getenv("SMTP_FROM")

	if smtpPort == "" {
		smtpPort = "587"
	}
	if fromEmail == "" {
		fromEmail = smtpUser
	}

	msg := fmt.Sprintf("From: RaceNotes <%s>\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		fromEmail, to, subject, body)

	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)
	return smtp.SendMail(smtpHost+":"+smtpPort, auth, fromEmail, []string{to}, []byte(msg))
}
