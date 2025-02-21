import { fetchApiWithToken } from "../utils/apiFetch";

interface StudyGroup {
    id: number;
    name: string;
    description: string;
    type: "public" | "closed" | "invite-only";
    moduleId: number;
    members: StudyGroupMember[];
}

interface StudyGroupMember {
    is: string;
    name: string;
    role: "admin" | "member" | "invitee" | "requester";
}

interface StudyGroupCreationDetails {
    name: string;
    description: string;
    type: "public" | "closed" | "invite-only";
}

export async function fetchStudyGroups(token: string): Promise<StudyGroup[]> {
    const response = await fetchApiWithToken<StudyGroup[]>("/study-groups", token);
    // TODO: Remove this once the backend is updated
    response.forEach((group: StudyGroup) => {
        group.moduleId = 0;
    });
    return response;
}

export async function createStudyGroup(token: string, studyGroup: StudyGroupCreationDetails): Promise<StudyGroup> {
    const response = await fetchApiWithToken<StudyGroup>("/study-groups", token, {
        method: "POST",
        body: JSON.stringify(studyGroup),
    });
    return response;
}
