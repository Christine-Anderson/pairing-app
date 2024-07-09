package types

type GroupMember struct {
	MemberId string `json:"memberId"`
	Name     string `json:"name"`
	Email    string `json:"email"`
}

type Group struct {
	GroupId      string        `json:"groupId"`
	GroupName    string        `json:"groupName"`
	GroupMembers []GroupMember `json:"groupMembers"`
}

type VerifyEmailDetails struct {
	Email string `json:"email"`
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

type AssignmentRestrictions struct {
	Restrictions map[string][]string `json:"restrictions"`
}

func NewGroupMember(memberId string, name string, email string) GroupMember {
	return GroupMember{
		MemberId: memberId,
		Name:     name,
		Email:    email,
	}
}

func NewGroup(groupId string, memberId string, groupDetails CreateGroupDetails) Group {
	return Group{
		GroupId:   groupId,
		GroupName: groupDetails.GroupName,
		GroupMembers: []GroupMember{
			NewGroupMember(memberId, groupDetails.Name, groupDetails.Email),
		},
	}
}
