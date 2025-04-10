import JSZip from "jszip"
import { makeRequest, fetchApiToJson } from "./apiFetch";

export class BackendFile {
  name: string;
  data: Blob;

  constructor(name: string, data: Blob) {
    this.name = name;
    this.data = data;
  }

  download() {
    const url = URL.createObjectURL(this.data);
    const a = document.createElement("a");
    a.href = url;
    a.download = this.name;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
  }
}

export async function apiFetchFiles(chatID: string, token: string | null): Promise<BackendFile[]> {
  const response = await makeRequest(`/files/${chatID}`, token);
  const blob = await response.blob();
  const zip = await JSZip.loadAsync(blob);
  const extractedFiles: BackendFile[] = [];

  for (const filename of Object.keys(zip.files)) {
    const fileData = zip.files[filename];
    if (!fileData.dir) {
      const fileBlob = await fileData.async("blob");
      extractedFiles.push(new BackendFile(filename, new Blob([fileBlob])));
    }
  }
  return extractedFiles;
}

export async function apiUploadFiles(
  file: File,
  chatID: string,
  token: string | null
): Promise<{ message: string; chatID: string; file: string; userId: string }> {
  const formData = new FormData();
  formData.append("file", file);
  formData.append("chatID", chatID);

  return fetchApiToJson<{ message: string; chatID: string; file: string; userId: string }>("/file", token, {
    method: "POST",
    body: formData,
  });
}

export async function apiDeleteFiles(
  filename: string,
  chatID: string,
  token: string | null
): Promise<{ message: string; chatID: string; file: string; userId: string }> {
  const formData = new FormData();
  formData.append("filename", filename);
  formData.append("chatID", chatID);

  return fetchApiToJson<{ message: string; chatID: string; file: string; userId: string }>("/file/delete", token, {
    method: "POST",
    body: formData,
  });
}
