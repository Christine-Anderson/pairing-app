package types

type GroupMember struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Group struct {
	GroupId      string        `json:"groupId"`
	GroupName    string        `json:"groupName"`
	GroupMembers []GroupMember `json:"groupMembers"`
}

type CreateGroupDetails struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	GroupName string `json:"groupName"`
}

type JoinGroupDetails struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	GroupId string `json:"groupId"`
}

func NewGroupMember(name string, email string) GroupMember {
	return GroupMember{
		Name:  name,
		Email: email,
	}
}

func NewGroup(groupId string, groupDetails CreateGroupDetails) Group {
	return Group{
		GroupId:   groupId,
		GroupName: groupDetails.GroupName,
		GroupMembers: []GroupMember{
			NewGroupMember(groupDetails.Name, groupDetails.Email),
		},
	}
}
