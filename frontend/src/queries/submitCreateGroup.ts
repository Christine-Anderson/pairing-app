import { CreateGroupAPIRequest, CreateGroupAPIResponse } from "../types";

const submitCreateGroup = async ({name, email, groupName}: CreateGroupAPIRequest): Promise<CreateGroupAPIResponse> => {
    // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
    const BASE_URL = import.meta.env.VITE_BASE_URL as string;

    const requestBody = {
        name: name,
        email: email,
        groupName: groupName
    }

    const response = await fetch(
        `${BASE_URL}create-group`,
        {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(requestBody)
        }
    )
    const result = await response.json() as CreateGroupAPIResponse;

    if (!response.ok) {
        throw new Error(`Error creating group ${groupName}.`);
    }

    return result;
}

export default submitCreateGroup;