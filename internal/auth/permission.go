package auth

import (
	"github.com/Confialink/wallet-messages/internal/srvdiscovery"
	"github.com/Confialink/wallet-permissions/rpc/permissions"

	"context"
	"net/http"
)

const PermissionSendReplyMessage = Permission("send_reply_internal_messages")

type PermissionService struct {
}

func NewPermissionService() *PermissionService {
	return &PermissionService{}
}

//Check checks if specified user is granted permission to perform some action
func (p *PermissionService) Check(userId, actionKey string) (bool, error) {
	request := &permissions.PermissionReq{UserId: userId, ActionKey: actionKey}

	checker, err := p.checker()
	if nil != err {
		return false, err
	}

	response, err := checker.Check(context.Background(), request)
	if nil != err {
		return false, err
	}
	return response.IsAllowed, nil
}

func (p *PermissionService) checker() (permissions.PermissionChecker, error) {
	permissionsUrl, err := srvdiscovery.ResolveRPC(srvdiscovery.ServiceNamePermissions)
	if nil != err {
		return nil, err
	}
	checker := permissions.NewPermissionCheckerProtobufClient(permissionsUrl.String(), http.DefaultClient)
	return checker, nil
}
