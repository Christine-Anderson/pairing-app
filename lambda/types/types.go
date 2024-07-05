package types

type CreateGroupDetails struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	GroupName string `json:"groupName"`
}

type GroupMember struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Group struct {
	GroupId      string        `json:"groupId"`
	GroupName    string        `json:"groupName"`
	GroupMembers []GroupMember `json:"groupMembers"`
}

func NewGroup(groupId string, groupDetails CreateGroupDetails) Group {
	return Group{
		GroupId:   groupId,
		GroupName: groupDetails.GroupName,
		GroupMembers: []GroupMember{
			{Name: groupDetails.Name, Email: groupDetails.Email},
		},
	}
}
