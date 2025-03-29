import { makeRequest } from '../utils/apiFetch';

interface ChatResponse {
    success: boolean;
    message: string;
}

interface ApiResponse {
    message?: string;
}

export async function sendChatMessage({
    token,
    chatID,
    message,
    memory,
}: {
        token: string | null;
        chatID: string;
        message: string;
        memory?: 'prevchats' | 'file' | 'all';
    }): Promise<ChatResponse> {
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

        const data = (await res.json()) as ApiResponse;

        return {
            success: true,
            message: data.message ?? 'No message returned',
        };
    } catch {
        return {
            success: false,
            message: 'Could not reach chatbot',
        };
    }
}
