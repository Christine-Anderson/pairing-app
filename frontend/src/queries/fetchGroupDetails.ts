import { GroupDetailsAPIRequest, GroupDetailsAPIResponse } from "../types";

const fetchGroupDetails = async ({groupId, jwt}: GroupDetailsAPIRequest): Promise<GroupDetailsAPIResponse> => {
    // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
    const BASE_URL = import.meta.env.VITE_BASE_URL as string;

    const response = await fetch(`${BASE_URL}group-details/${groupId}?jwt=${jwt}`)
    const result = await response.json() as GroupDetailsAPIResponse;

    if (!response.ok) {
        throw new Error(`Error fetching group ${groupId}.`);
    }

    return result;
}

export default fetchGroupDetails;