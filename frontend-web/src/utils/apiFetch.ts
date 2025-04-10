const BASE_URL = "http://localhost:8080/api";

export async function fetchApiToJson<T>(
    url: string,
    token: string | null,
    options: RequestInit = {}
): Promise<T> {
    const response = await makeRequest(url, token, options);
    return response.json() as Promise<T>;
};

export async function fetchApi(
    url: string,
    token: string | null,
    options: RequestInit = {}
): Promise<void> {
    await makeRequest(url, token, options);
}

export async function makeRequest(
    url: string,
    token: string | null,
    options: RequestInit = {}
): Promise<Response> {
    const headers = new Headers(options.headers);
    if (token) {
        headers.append("Authorization", `Bearer ${token}`);
    }

    const response = await fetch(BASE_URL + url, {
        ...options,
        headers,
    });

    if (!response.ok) {
        let serverMsg = "";
        try {
            serverMsg = await response.text();
        } catch {
            // fallback if reading the text body fails
            serverMsg = response.statusText;
        }
        throw new Error(`HTTP error ${response.status}: ${serverMsg}`);
    }

    return response;
}
