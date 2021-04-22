package model

type User struct {
	UID       *string `json:"uid"`
	Email     *string `json:"email"`
	Username  *string `json:"username"`
	FirstName *string `json:"firstName"`
	LastName  *string `json:"lastName"`
	RoleName  *string `json:"roleName"`
	GroupId   *uint64 `json:"groupId"`
}
