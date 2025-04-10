import { fetchApiToJson } from "../utils/apiFetch";

export interface ChatMessage {
    id: number;
    userId: string;
    userName: string;
    text: string;
    timestamp: Date;
}

export async function fetchPastChatMessagesByGroupId(token: string | null, studyGroupId: number): Promise<ChatMessage[]> {
    const chatMessages = await fetchApiToJson<ChatMessageResponse[]>(`/chat/${studyGroupId}/past`, token);
    
    return chatMessages.map(toChatMessage);
}

interface ChatMessageResponse {
    id: number;
    userId: string;
    userName: string;
    text: string;
    timestamp: string;
}

function toChatMessage(response: ChatMessageResponse): ChatMessage {
    return {
        ...response,
        timestamp: new Date(response.timestamp)
    }
}