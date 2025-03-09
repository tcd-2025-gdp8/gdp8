import { fetchApi, fetchApiToJson } from "../utils/apiFetch";

export interface StudyGroup {
    id: number;
    name: string;
    description: string;
    type: "public" | "closed" | "invite-only";
    moduleId: number;
    maxMembers: number;
    members: StudyGroupMember[];
}

export interface StudyGroupMember {
    id: string;
    name: string;
    role: "admin" | "member" | "invitee" | "requester";
}

export interface StudyGroupCreationDetails {
    name: string;
    description: string;
    type: "public" | "closed" | "invite-only";
    moduleId: number;
    maxMembers: number;
}

export async function fetchStudyGroups(token: string | null): Promise<StudyGroup[]> {
    const studyGroups = await fetchApiToJson<StudyGroup[]>("/study-groups", token);
    return studyGroups;
}

export async function createStudyGroup(token: string | null, studyGroup: StudyGroupCreationDetails): Promise<StudyGroup> {
    const createdStudyGroup = await fetchApiToJson<StudyGroup>("/study-groups", token, {
        method: "POST",
        body: JSON.stringify(studyGroup),
    });
    return createdStudyGroup;
}

export async function requestToJoinStudyGroup(token: string | null, studyGroupId: number): Promise<void> {
    await fetchApi(`/study-groups/${studyGroupId}/request-to-join`, token, {
        method: "POST",
    });
}

export async function removeMemberFromStudyGroup(token: string | null, studyGroupId: number, memberId: string): Promise<void> {
    await fetchApi(`/study-groups/${studyGroupId}/remove-member`, token, {
        method: "POST",
        body: JSON.stringify({ targetUserId: memberId }),
    });
}

export async function fetchStudyGroupById(token: string | null, studyGroupId: number): Promise<StudyGroup> {
    return await fetchApiToJson<StudyGroup>(`/api/study-groups/${studyGroupId}`, token);
}

export async function updateStudyGroup(
    token: string | null,
    studyGroupId: number,
    updates: Partial<StudyGroupCreationDetails>
): Promise<StudyGroup> {
    return await fetchApiToJson<StudyGroup>(`/api/study-groups/${studyGroupId}`, token, {
        method: "PUT", // Adjust method if needed by your backend
        body: JSON.stringify(updates),
    });
}

export async function deleteStudyGroup(token: string | null, studyGroupId: number): Promise<void> {
    await fetchApi(`/api/study-groups/${studyGroupId}`, token, {
        method: "DELETE",
    });
}
