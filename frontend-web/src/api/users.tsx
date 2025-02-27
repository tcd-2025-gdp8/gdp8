import { Module } from "./modules";
import { fetchApi, fetchApiToJson } from "../utils/apiFetch";

export interface UserRegisterDetails {
    id: string;
    name: string;
}

export interface UserModulesSelection {
    selectedModules: number[];
}

export async function registerUser(token: string | null, user: UserRegisterDetails): Promise<void> {
    const userWithModules = { ...user, modules: [] };
    await fetchApi("/user", token, {
        method: "POST",
        body: JSON.stringify(userWithModules),
    });
}

export async function getUserModules(token: string | null, userId: string): Promise<Module[]> {
    const modules = await fetchApiToJson<Module[]>(`/user/${userId}/modules`, token);
    return modules;
}

export async function updateUserModules(token: string | null, userId: string, modules: UserModulesSelection): Promise<void> {
    await fetchApi(`/user/${userId}/modules`, token, {
        method: "POST",
        body: JSON.stringify(modules),
    });
}
