import { fetchApi, fetchApiToJson } from "../utils/apiFetch";
import { StudyGroup } from "./studyGroups";

function formatTime(dateObj: Date): string {
    return dateObj.toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" });
  }

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
    availabilityEntries: StudySessionAvailabilityEntry[];
}

export interface StudySessionAvailabilityEntry {
    userId: string;
    availabilityStart: Date;
    availabilityEnd: Date;
}

export interface StudySessionAvailabilityRequestDetails {
    title: string;
    availabilityPeriodStart: Date;
    availabilityPeriodEnd: Date;
}

export interface StudySessionAvailabilityEntryDetails {
    availabilityStart: Date;
    availabilityEnd: Date;
}

export async function fetchStudySessions(token: string | null, studyGroupId: number): Promise<StudySession[]> {
    const studySessions = await fetchApiToJson<StudySessionResponse[]>(`/study-groups/${studyGroupId}/study-sessions`, token);
    return studySessions.map(toStudySession);
}

export async function fetchAllUserStudySessions(token: string | null, userGroups: StudyGroup[]): Promise<Session[]> {
    if (!token || userGroups.length === 0) return [];

    const sessionsByGroup = await Promise.all(
        userGroups.map(async (group) => {
            const sessions = await fetchStudySessions(token, group.id);
            return sessions.map((session) => ({
                id: session.id,
                name: session.title,
                date: session.startTime,
                earliestTime: formatTime(session.startTime),
                latestTime: formatTime(session.endTime),
            }));
        })
    );

    return sessionsByGroup.flat();
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
    await fetchApi(`/study-sessions/${studySessionId}`, token, { method: 'DELETE' });
}

export async function fetchCurrentAvailabilityRequestsByGroupId(token: string | null, studyGroupId: number): Promise<StudySessionAvailabilityRequest[]> {
    const response = await fetchApiToJson<StudySessionAvailabilityRequestResponse[]>(`/study-groups/${studyGroupId}/availability-requests`, token);
    return response.map(toStudySessionAvailabilityRequest);
}

export async function createAvailabilityRequest(token: string | null, studyGroupId: number, availabilityRequest: StudySessionAvailabilityRequestDetails): Promise<void> {
    await fetchApi(`/study-groups/${studyGroupId}/availability-requests`, token, {
        method: 'POST',
        body: JSON.stringify(availabilityRequest),
    });
}

export async function deleteAvailabilityRequest(token: string | null, availabilityRequestId: number): Promise<void> {
    await fetchApi(`/availability-requests/${availabilityRequestId}`, token, { method: 'DELETE' });
}

export async function upsertAvailabilityEntries(
    token: string | null,
    availabilityRequestId: number,
    availabilityEntries: StudySessionAvailabilityEntryDetails[]
): Promise<void> {
    await fetchApi(`/availability-requests/${availabilityRequestId}/entries`, token, {
        method: 'PUT',
        body: JSON.stringify({ entries: availabilityEntries }),
    });
}

function toStudySession(response: StudySessionResponse): StudySession {
    return {
        ...response,
        startTime: new Date(response.startTime),
        endTime: new Date(response.endTime),
    };
}

function toStudySessionAvailabilityRequest(response: StudySessionAvailabilityRequestResponse): StudySessionAvailabilityRequest {
    return {
        ...response,
        availabilityPeriodStart: new Date(response.availabilityPeriodStart),
        availabilityPeriodEnd: new Date(response.availabilityPeriodEnd),
        availabilityEntries: response.entries.map(toStudySessionAvailabilityEntry),
    };
}

function toStudySessionAvailabilityEntry(response: StudySessionAvailabilityEntryResponse): StudySessionAvailabilityEntry {
    return {
        ...response,
        availabilityStart: new Date(response.availabilityStart),
        availabilityEnd: new Date(response.availabilityEnd),
    };
}

export interface Session {
    id: number;
    name: string;
    date: Date | null;
    earliestTime: string;
    latestTime: string;
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

interface StudySessionAvailabilityRequestResponse {
    id: number;
    studyGroupId: number;
    title: string;
    availabilityPeriodStart: string;
    availabilityPeriodEnd: string;
    entries: StudySessionAvailabilityEntryResponse[];
}

interface StudySessionAvailabilityEntryResponse {
    userId: string;
    availabilityStart: string;
    availabilityEnd: string;
}