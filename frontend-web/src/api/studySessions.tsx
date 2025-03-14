import { fetchApi, fetchApiToJson } from "../utils/apiFetch";

export interface StudySession {
    id: number;
    studyGroupId: number;
    creatorId: string;
    title: string;
    startTime: Date;
    durationMinutes: number;
    endTime: Date;
}

export interface StudySessionDetails {
    title: string;
    startTime: Date;
    durationMinutes: number;
}

export interface StudySessionAvailabilityRequest {
    id: number;
    studyGroupId: number;
    title: string;
    availabilityPeriodStart: Date;
    availabilityPeriodEnd: Date;

}

export interface StudySessionAvailabilityEntry {
    userId: string;
    availabilityStart: Date;
    availabilityEnd: Date;
}

export async function fetchStudySessions(token: string | null, studyGroupId: number): Promise<StudySession[]> {
    const studySessions = await fetchApiToJson<StudySessionResponse[]>(`/study-groups/${studyGroupId}/study-sessions`, token);

    return studySessions.map(toStudySession);
}

export async function createStudySession(token: string | null, studyGroupId: number, studySession: StudySessionDetails): Promise<StudySession> {
    const response = await fetchApiToJson<StudySessionResponse>(`/study-groups/${studyGroupId}/study-sessions`, token, {
        method: 'POST',
        body: JSON.stringify(studySession),
    });

    return toStudySession(response);
}

export async function updateStudySession(token: string | null, studySessionId: number, studySession: StudySessionDetails): Promise<StudySession> {
    const response = await fetchApiToJson<StudySessionResponse>(`/study-sessions/${studySessionId}`, token, {
        method: 'PUT',
        body: JSON.stringify(studySession),
    });

    return toStudySession(response);
}

export async function deleteStudySession(token: string | null, studySessionId: number): Promise<void> {
    await fetchApi(`/study-sessions/${studySessionId}`, token, {
        method: 'DELETE',
    });
}

interface StudySessionResponse {
    id: number;
    studyGroupId: number;
    creatorId: string;
    title: string;
    startTime: string;
    durationMinutes: number;
    endTime: string;
}

function toStudySession(response: StudySessionResponse): StudySession {
    return {
        ...response,
        startTime: new Date(response.startTime),
        endTime: new Date(response.endTime),
    };
}
