import { JoinGroupAPIRequest, JoinGroupAPIResponse } from "../types";

const submitJoinGroup = async ({name, email, groupId}: JoinGroupAPIRequest): Promise<JoinGroupAPIResponse> => {
    // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
    const BASE_URL = import.meta.env.VITE_BASE_URL as string;

    const requestBody = {
        name: name,
        email: email,
        groupId: groupId
    }

    const response = await fetch(
        `${BASE_URL}join-group`,
        {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(requestBody)
        }
    )
    const result = await response.json() as JoinGroupAPIResponse;

    if (!response.ok) {
        throw new Error(`Error joining group ${groupId}.`);
    }

    return result;
}

export default submitJoinGroup;