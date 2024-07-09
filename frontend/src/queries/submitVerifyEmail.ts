import { VerifyEmailAPIRequest, VerifyEmailAPIResponse } from "../types";

const submitVerifyEmail = async ({email}: VerifyEmailAPIRequest): Promise<VerifyEmailAPIResponse> => {
    // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
    const BASE_URL = import.meta.env.VITE_BASE_URL as string;

    const requestBody = {
        email: email,
    }

    const response = await fetch(
        `${BASE_URL}verify-email`,
        {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(requestBody)
        }
    )
    const result = await response.json() as VerifyEmailAPIResponse;

    if (!response.ok) {
        throw new Error(`Error verifying email ${email}.`);
    }

    return result;
}

export default submitVerifyEmail;