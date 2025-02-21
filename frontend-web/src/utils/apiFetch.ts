const BASE_URL = "http://localhost:8080/api";

export async function fetchApiWithToken<T>(
    url: string,
    token: string | null,
    options: RequestInit = {}
): Promise<T> {
    const headers = new Headers(options.headers);
    if (token) {
        headers.append("Authorization", `Bearer ${token}`);
    }

    const response = await fetch(BASE_URL + url, {
        ...options,
        headers,
    });

    if (!response.ok) {
        throw new Error("HTTP error ${response.status}: ${response.statusText}");
    }

    return response.json() as Promise<T>;
};
