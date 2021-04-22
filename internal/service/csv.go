package service

import (
	"fmt"
	"time"

	"github.com/Confialink/wallet-messages/internal/model"

	csvPkg "github.com/Confialink/wallet-pkg-utils/csv"
	"github.com/Confialink/wallet-pkg-utils/timefmt"

	"github.com/Confialink/wallet-messages/internal/service/settings"
	userpb "github.com/Confialink/wallet-users/rpc/proto/users"
)

// Csv
type Csv struct{}

func NewCsv() *Csv {
	return &Csv{}
}

func (s *Csv) MessagesToCsv(messages []*model.Message, currentUser *userpb.User) (*csvPkg.File, error) {
	if currentUser.RoleName == "admin" || currentUser.RoleName == "root" {
		return s.adminListToCsv(messages, currentUser.UID)
	}

	return s.userListToCsv(messages, currentUser.UID)
}

func (s *Csv) adminListToCsv(messages []*model.Message, currentUid string) (*csvPkg.File, error) {
	header := []string{"User Name", "Subject", "User Email", "Date"}

	timeSettings, err := settings.GetTimeSettings()
	if err != nil {
		return nil, err
	}

	receiver := func(message *model.Message) []string {
		formattedLastMessageCreatedAt := timefmt.FormatWithTime(message.LastMessageCreatedAt, timeSettings.DateFormat, timeSettings.Timezone)

		var username string
		var email string
		if (message.RecipientId != nil && *message.RecipientId == currentUid && message.Sender != nil) || (message.RecipientId == nil && message.Sender != nil) {
			username = *message.Sender.Username
			email = *message.Sender.Email
		} else if message.SenderId != nil && *message.SenderId == currentUid && message.Recipient != nil {
			username = *message.Recipient.Username
			email = *message.Recipient.Email
		} else {
			username = "Administrator"
			email, _ = settings.GetNoReplayEmail()
		}

		var subject string
		if message.Subject != nil {
			subject = *message.Subject
		}

		return []string{
			username,
			subject,
			email,
			formattedLastMessageCreatedAt,
		}
	}

	return s.messagesToCsv(messages, header, receiver)
}

func (s *Csv) userListToCsv(messages []*model.Message, currentUid string) (*csvPkg.File, error) {
	header := []string{"User", "Subject", "Date"}

	timeSettings, err := settings.GetTimeSettings()
	if err != nil {
		return nil, err
	}

	receiver := func(message *model.Message) []string {
		formattedLastMessageCreatedAt := timefmt.FormatWithTime(message.LastMessageCreatedAt, timeSettings.DateFormat, timeSettings.Timezone)

		var username string
		if message.RecipientId != nil && *message.RecipientId == currentUid && message.Sender != nil {
			username = *message.Sender.Username
		} else if message.SenderId != nil && *message.SenderId == currentUid && message.Recipient != nil {
			username = *message.Recipient.Username
		}

		var subject string
		if message.Subject != nil {
			subject = *message.Subject
		}

		return []string{
			username,
			subject,
			formattedLastMessageCreatedAt,
		}
	}

	return s.messagesToCsv(messages, header, receiver)
}

func (s *Csv) messagesToCsv(
	messages []*model.Message,
	header []string,
	dataReceiver func(message *model.Message) []string,
) (*csvPkg.File, error) {
	currentTime := time.Now()
	timeSettings, err := settings.GetTimeSettings()
	if err != nil {
		return nil, err
	}

	file := csvPkg.NewFile()
	formattedCurrentTime := timefmt.FormatFilenameWithTime(currentTime, timeSettings.Timezone)
	file.Name = fmt.Sprintf("messages-%s.csv", formattedCurrentTime)

	file.WriteRow(header)

	for _, v := range messages {
		record := dataReceiver(v)
		file.WriteRow(record)
	}

	return file, nil
}
