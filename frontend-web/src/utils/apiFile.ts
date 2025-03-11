import JSZip from "jszip"
import { makeRequest } from "./apiFetch";

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
