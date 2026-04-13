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

export async function getHealth() {
  return request<{ status: string; spider_ready: boolean }>("/health");
}

export async function getCookies() {
  return request<{ configured: boolean }>("/api/v1/cookies");
}

export async function setCookies(cookies: string) {
  return request<{ message: string }>("/api/v1/cookies", {
    method: "POST",
    body: JSON.stringify({ cookies }),
  });
}

export async function collectSingleNote(noteUrl: string) {
  return request<any>("/api/v1/notes/single", {
    method: "POST",
    body: JSON.stringify({ note_url: noteUrl }),
  });
}

export async function collectUserNotes(userUrl: string) {
  return request<{ notes: string[]; count: number }>("/api/v1/notes/user", {
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
  return request<{ notes: string[]; count: number }>("/api/v1/notes/search", {
    method: "POST",
    body: JSON.stringify(params),
  });
}
