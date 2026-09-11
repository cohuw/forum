import { store } from "./store";
import { navigate } from "./router";
import type { ApiError } from "./types";

const BASE_URL = "/api";

async function request<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
  const headers = new Headers(options.headers || {});
  if (!headers.has("Content-Type") && !(options.body instanceof FormData)) {
    headers.set("Content-Type", "application/json");
  }

  const token = store.getToken();
  if (token) {
    headers.set("Authorization", `Bearer ${token}`);
  }

  options.credentials = "include";

  const res = await fetch(`${BASE_URL}${endpoint}`, { ...options, headers });
  
  if (!res.ok) {
    let errBody: any = null;
    try {
      errBody = await res.json();
    } catch {}

    const code = errBody?.code || errBody?.["код"];
    if (res.status === 401 && (code === "token_expired" || code === "токен_истек")) {
      const refreshed = await fetch(`${BASE_URL}/auth/refresh`, { method: "GET", credentials: "include" });
      if (refreshed.ok) {
        const data = await refreshed.json();
        store.setToken(data.access_token);
        store.setUser(data.user);
        
        const retryHeaders = new Headers(options.headers || {});
        retryHeaders.set("Authorization", `Bearer ${data.access_token}`);
        if (!retryHeaders.has("Content-Type")) retryHeaders.set("Content-Type", "application/json");
        
        const retryRes = await fetch(`${BASE_URL}${endpoint}`, { ...options, headers: retryHeaders });
        if (!retryRes.ok) {
           let rErr: any = null;
           try { rErr = await retryRes.json(); } catch {}
           throw new Error(rErr?.["ошибка"] || rErr?.error || `Ошибка: статус ${retryRes.status}`);
        }
        const text = await retryRes.text();
        return text ? JSON.parse(text) : ({} as T);
      } else {
        store.clear();
        navigate("#/login");
        throw new Error("ерр сессия истекла");
      }
    }

    throw new Error(errBody?.["ошибка"] || errBody?.error || `Ошибка: статус ${res.status}`);
  }

  const text = await res.text();
  return text ? JSON.parse(text) : ({} as T);
}

export const api = {
  get: <T>(endpoint: string) => request<T>(endpoint, { method: "GET" }),
  post: <T>(endpoint: string, body: any) => request<T>(endpoint, { method: "POST", body: JSON.stringify(body) }),
  delete: <T>(endpoint: string) => request<T>(endpoint, { method: "DELETE" }),
};
