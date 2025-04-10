import { makeRequest } from '../utils/apiFetch';

interface BotResponse {
    message: string;
}

export async function prompt({
    token,
    chatID,
    message,
    memory,
}: {
    token: string | null;
    chatID: string;
    message: string;
    memory?: string;
}): Promise<BotResponse> {
    const body = new URLSearchParams({
        text: message,
        chatID,
        memory: memory ?? '',
    });

    try {
        const res = await makeRequest(`/chatbot`, token, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded',
            },
            body,
        });
        const data = (await res.json()) as BotResponse;
        return data;
    } catch {
        return {
            message: 'Could not reach chatbot',
        };
    }
}
