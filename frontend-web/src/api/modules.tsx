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

export async function fetchAllModules(token: string): Promise<Module[]> {
    const modules = await fetchApiToJson<Module[]>("/modules", token);
    return modules;
}

export async function createModule(token: string, module: ModuleCreationDetails): Promise<void> {
    await fetchApi("/modules", token, {
        method: "POST",
        body: JSON.stringify(module),
    });
}
