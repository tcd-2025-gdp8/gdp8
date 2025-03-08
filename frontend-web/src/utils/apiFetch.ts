const BASE_URL = "http://localhost:8080/api";

export async function uploadFile(
  file: File,
  chatID: string,
  userId: string,
  token: string | null
): Promise<{ message: string; chatID: string; file: string; userId: string }> {
  const formData = new FormData();
  formData.append("file", file);
  formData.append("chatID", chatID);
  formData.append("userID", userId);

  return fetchApiToJson<{ message: string; chatID: string; file: string; userId: string }>("/file", token, {
    method: "POST",
    body: formData,
  });
}

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

async function makeRequest(
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
        throw new Error(`HTTP error ${response.status}: ${response.statusText}`);
    }

    return response;
}
