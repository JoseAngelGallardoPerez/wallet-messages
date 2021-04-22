package service

import (
	"github.com/Confialink/wallet-messages/internal/srvdiscovery"
	pb "github.com/Confialink/wallet-users/rpc/proto/users"

	"context"
	"net/http"
)

type UserService struct {
}

func NewUserService() *UserService {
	return &UserService{}
}

// GetByUsername returns User by passed username
func (u *UserService) GetByUsername(username string) (*pb.User, error) {
	req := pb.Request{Username: username}
	client, err := u.getClient()
	if err != nil {
		return nil, err
	}
	resp, err := client.GetByUsername(context.Background(), &req)
	if err != nil {
		return nil, err
	}
	return resp.User, nil
}

func (u *UserService) GetByUID(uid string) (*pb.User, error) {
	req := pb.Request{UID: uid}
	client, err := u.getClient()
	if err != nil {
		return nil, err
	}
	resp, err := client.GetByUID(context.Background(), &req)
	if err != nil {
		return nil, err
	}
	return resp.User, nil
}

func (u *UserService) GetByUIDs(uids []string) ([]*pb.User, error) {
	req := pb.Request{UIDs: uids}
	client, err := u.getClient()
	if err != nil {
		return nil, err
	}
	resp, err := client.GetByUIDs(context.Background(), &req)
	if err != nil {
		return nil, err
	}
	return resp.Users, nil
}

func (u *UserService) GetByGroup(id uint64) ([]*pb.User, error) {
	req := pb.Request{GroupId: id}
	client, err := u.getClient()
	if err != nil {
		return nil, err
	}
	resp, err := client.GetByUserGroupId(context.Background(), &req)
	if err != nil {
		return nil, err
	}
	return resp.Users, nil
}

func (u *UserService) GetAll() ([]*pb.User, error) {
	req := pb.Request{}
	client, err := u.getClient()
	if err != nil {
		return nil, err
	}
	resp, err := client.GetAll(context.Background(), &req)
	if err != nil {
		return nil, err
	}
	return resp.Users, nil
}

func (u *UserService) getClient() (pb.UserHandler, error) {
	usersUrl, err := srvdiscovery.ResolveRPC(srvdiscovery.ServiceNameUsers)
	if nil != err {
		return nil, err
	}

	return pb.NewUserHandlerProtobufClient(usersUrl.String(), http.DefaultClient), nil
}
