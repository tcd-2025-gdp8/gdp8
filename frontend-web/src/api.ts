// src/api.ts
export const fetchWithToken = async <T>(
    token: string | null,
    url: string,
    options: RequestInit = {}
): Promise<T> => {
    const headers = new Headers(options.headers);
    if (token) {
        headers.append("Authorization", `Bearer ${token}`);
    }

    const response = await fetch(url, {
        ...options,
        headers,
    });

    if (!response.ok) {
        throw new Error("Network response was not ok");
    }

    return response.json() as Promise<T>;
};
