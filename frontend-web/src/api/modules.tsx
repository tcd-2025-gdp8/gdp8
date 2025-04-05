import { fetchApi, fetchApiToJson } from "../utils/apiFetch";

export interface Module {
    id: number;
    code: string;
    name: string;
}

export interface ModuleCreationDetails {
    code: string;
    name: string;
}

export interface StudyGroupStatsDTO {
    id: number;
    name: string;
    members: number;
    totalHours: number;
    weeklyHours: number[];
}

export async function fetchAllModules(token: string | null): Promise<Module[]> {
    const modules = await fetchApiToJson<Module[]>("/modules", token);
    return modules;
}

export async function createModule(token: string | null, module: ModuleCreationDetails): Promise<void> {
    await fetchApi("/modules", token, {
        method: "POST",
        body: JSON.stringify(module),
    });
}

export async function fetchStudyGroupStats(token: string | null): Promise<Record<string, StudyGroupStatsDTO[] | null>> {
    const stats = await fetchApiToJson<Record<string, StudyGroupStatsDTO[] | null>>("/study-group-stats", token);
    return stats;
}
