import { fetchApi, fetchApiToJson } from "../utils/apiFetch";

export interface StudyGroup {
    id: number;
    name: string;
    description: string;
    type: "public" | "closed" | "invite-only";
    moduleId: number;
    members: StudyGroupMember[];
}

export interface StudyGroupMember {
    is: string;
    name: string;
    role: "admin" | "member" | "invitee" | "requester";
}

export interface StudyGroupCreationDetails {
    name: string;
    description: string;
    type: "public" | "closed" | "invite-only";
    moduleId: number;
}

export async function fetchStudyGroups(token: string): Promise<StudyGroup[]> {
    const studyGroups = await fetchApiToJson<StudyGroup[]>("/study-groups", token);
    return studyGroups;
}

export async function createStudyGroup(token: string, studyGroup: StudyGroupCreationDetails): Promise<StudyGroup> {
    const createdStudyGroup = await fetchApiToJson<StudyGroup>("/study-groups", token, {
        method: "POST",
        body: JSON.stringify(studyGroup),
    });
    return createdStudyGroup;
}

export async function requestToJoinStudyGroup(token: string, studyGroupId: number): Promise<void> {
    await fetchApi(`/study-groups/${studyGroupId}/request-to-join`, token, {
        method: "POST",
    });
}

export async function removeMemberFromStudyGroup(token: string, studyGroupId: number, memberId: string): Promise<void> {
    await fetchApi(`/study-groups/${studyGroupId}/remove-member`, token, {
        method: "POST",
        body: JSON.stringify({ targetUserId: memberId }),
    });
}
