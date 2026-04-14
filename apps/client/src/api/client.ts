const GATEWAY_URL = "http://127.0.0.1:8000";

async function request<T>(path: string, options?: RequestInit): Promise<T> {
  const res = await fetch(`${GATEWAY_URL}${path}`, {
    headers: { "Content-Type": "application/json" },
    ...options,
  });
  if (!res.ok) {
    const err = await res.json().catch(() => ({ detail: res.statusText }));
    throw new Error(err.detail || err.message || "Request failed");
  }
  return res.json();
}

export interface NoteData {
  id: string;
  title: string;
  author: string;
  author_id: string;
  type: string;
  liked: number;
  collected: number;
  commented: number;
  shared: number;
  url: string;
  tags: string[];
  time: string;
  ip_location: string;
}

export interface NotesListResponse {
  notes: NoteData[];
  count: number;
  partial: boolean;
  errors: string[];
}

export interface SingleNoteResponse {
  note: NoteData;
  partial: boolean;
}

export async function getHealth() {
  return request<{ status: string; spider_ready: boolean }>("/health");
}

export async function getCookies() {
  return request<{ configured: boolean; has_content: boolean; domain: string }>(
    "/api/v1/cookies"
  );
}

export async function setCookies(cookies: string) {
  return request<{ message: string }>("/api/v1/cookies", {
    method: "POST",
    body: JSON.stringify({ cookies }),
  });
}

export async function deleteCookies() {
  return request<{ message: string }>("/api/v1/cookies", {
    method: "DELETE",
  });
}

export async function collectSingleNote(noteUrl: string) {
  return request<SingleNoteResponse>("/api/v1/notes/single", {
    method: "POST",
    body: JSON.stringify({ note_url: noteUrl }),
  });
}

export async function collectUserNotes(userUrl: string) {
  return request<NotesListResponse>("/api/v1/notes/user", {
    method: "POST",
    body: JSON.stringify({ user_url: userUrl }),
  });
}

export interface SearchParams {
  keyword: string;
  count: number;
  sort_type?: number;
  note_type?: number;
  note_time?: number;
}

export async function searchNotes(params: SearchParams) {
  return request<NotesListResponse>("/api/v1/notes/search", {
    method: "POST",
    body: JSON.stringify(params),
  });
}

export async function exportNotes(notes: NoteData[], filename = "xhs_notes.xlsx") {
  const res = await fetch(`${GATEWAY_URL}/api/v1/notes/export?filename=${encodeURIComponent(filename)}`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ notes }),
  });
  if (!res.ok) {
    const err = await res.json().catch(() => ({ detail: res.statusText }));
    throw new Error(err.detail || err.message || "Export failed");
  }
  const blob = await res.blob();
  const url = URL.createObjectURL(blob);
  const a = document.createElement("a");
  a.href = url;
  a.download = filename;
  a.click();
  URL.revokeObjectURL(url);
}
