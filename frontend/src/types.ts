export interface VerifyEmailAPIRequest {
    email: string;
}

export interface VerifyEmailAPIResponse {
    email: string;
}

export interface CreateGroupAPIRequest {
    name: string;
    email: string;
    groupName: string;
}
export interface CreateGroupAPIResponse {
    groupId: string;
    groupName: string;
}

export interface JoinGroupAPIRequest {
    name: string;
    email: string;
    groupId: string;
}
export interface JoinGroupAPIResponse {
    groupId: string;
    groupName: string;
}

export interface GroupDetailsAPIRequest {
    groupId: string;
    jwt: string;
}

export interface GroupMember {
    memberId: string;
    name: string;
    email: string;
}

export interface GroupDetailsAPIResponse {
    groupId: string;
    groupName: string;
    groupMembers: GroupMember[];
}