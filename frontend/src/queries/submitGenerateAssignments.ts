import { GenerateAssignmentsAPIRequest, GroupAPIResponse } from "../types";

const submitGenerateAssignments = async ({groupId, jwt, restrictions}: GenerateAssignmentsAPIRequest): Promise<GroupAPIResponse> => {
    // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
    const BASE_URL = import.meta.env.VITE_BASE_URL as string;

    const requestBody = {
        restrictions: restrictions
    }

    const response = await fetch(
        `${BASE_URL}assign/${groupId}?jwt=${jwt}`,
        {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(requestBody)
        }
    )

    const result = await response.json() as GroupAPIResponse;

    if (!response.ok) {
        throw new Error(`Error generating assignments for group ${groupId}.`);
    }

    return result;
}

export default submitGenerateAssignments;