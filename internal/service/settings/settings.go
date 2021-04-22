package settings

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	pb "github.com/Confialink/wallet-settings/rpc/proto/settings"

	"github.com/Confialink/wallet-messages/internal/service/settings/connection"
)

type MessageParams struct {
	IncomingMessageUrl string
	ProfileSettingsUrl string
}

// TimeSettings struct has timezone and date format
type TimeSettings struct {
	Timezone       string
	DateFormat     string
	TimeFormat     string
	DateTimeFormat string
}

func GetMessageParamsByMessageId(id uint64) (*MessageParams, error) {
	client, err := connection.GetSystemSettingsClient()
	if err != nil {
		return nil, err
	}

	response, err := client.List(context.Background(), &pb.Request{Path: "regional/general/%"})
	if err != nil {
		return nil, err
	}

	url := getSettingValue(response.Settings, "regional/general/site_url")
	messagePath := getSettingValue(response.Settings, "regional/general/site_incoming_message_path")
	settingsPath := getSettingValue(response.Settings, "regional/general/site_my_profile_settings_path")
	params := &MessageParams{
		IncomingMessageUrl: strings.Replace(url+messagePath, "{id}", strconv.Itoa(int(id)), -1),
		ProfileSettingsUrl: url + settingsPath,
	}
	return params, nil
}

// GetTimeSettings returns new TimeSettings from settings service or err if can not get it
func GetTimeSettings() (*TimeSettings, error) {
	timeSettings := TimeSettings{}
	client, err := connection.GetSystemSettingsClient()
	if err != nil {
		return &timeSettings, err
	}

	response, err := client.List(context.Background(), &pb.Request{Path: "regional/general/%"})
	if err != nil {
		return &timeSettings, err
	}

	timeSettings.Timezone = getSettingValue(response.Settings, "regional/general/default_timezone")
	timeSettings.DateFormat = getSettingValue(response.Settings, "regional/general/default_date_format")
	timeSettings.TimeFormat = getSettingValue(response.Settings, "regional/general/default_time_format")
	timeSettings.DateTimeFormat = fmt.Sprintf("%s %s", timeSettings.DateFormat, timeSettings.TimeFormat)
	return &timeSettings, nil
}

// GetNoReplayEmail returns no-replay email or err if can not get it
func GetNoReplayEmail() (string, error) {
	client, err := connection.GetSystemSettingsClient()
	if err != nil {
		return "", err
	}

	response, err := client.List(context.Background(), &pb.Request{Path: "regional/general/no-replay-email"})
	if err != nil {
		return "", err
	}

	return getSettingValue(response.Settings, "regional/general/no-replay-email"), nil
}

func getSettingValue(settings []*pb.Setting, path string) string {
	for _, v := range settings {
		if v.Path == path {
			return v.Value
		}
	}
	return ""
}
