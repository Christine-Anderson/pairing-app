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