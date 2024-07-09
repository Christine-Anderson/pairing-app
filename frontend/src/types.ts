export interface VerifyEmailAPI {
    email: string;
}

export interface CreateGroupAPIRequest {
    name: string;
    email: string;
    groupName: string;
}

export interface JoinGroupAPIRequest {
    name: string;
    email: string;
    groupId: string;
}

export interface GroupAPIResponse {
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

export interface GenerateAssignmentsAPIRequest {
    groupId: string;
    jwt: string;
    restrictions: {
        [key: string]: string[];
    };
}