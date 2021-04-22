package connection

import (
	"github.com/Confialink/wallet-messages/internal/srvdiscovery"
	pb "github.com/Confialink/wallet-settings/rpc/proto/settings"

	"net/http"
)

// GetSystemSettingsClient returns rpc client to settings service
func GetSystemSettingsClient() (pb.SettingsHandler, error) {
	settingsUrl, err := srvdiscovery.ResolveRPC(srvdiscovery.ServiceNameSettings)
	if err != nil {
		return nil, err
	}

	return pb.NewSettingsHandlerProtobufClient(settingsUrl.String(), http.DefaultClient), nil
}
