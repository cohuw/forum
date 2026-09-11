import type { User } from "./types";

const TOKEN_KEY = "access_token";
const USER_KEY = "forum_user";

export const store = {
  getToken(): string | null {
    return sessionStorage.getItem(TOKEN_KEY);
  },
  setToken(token: string): void {
    sessionStorage.setItem(TOKEN_KEY, token);
  },
  getUser(): User | null {
    const raw = sessionStorage.getItem(USER_KEY);
    if (!raw) return null;
    try {
      return JSON.parse(raw) as User;
    } catch {
      return null;
    }
  },
  setUser(user: User): void {
    sessionStorage.setItem(USER_KEY, JSON.stringify(user));
  },
  isLoggedIn(): boolean {
    return !!this.getToken();
  },
  clear(): void {
    sessionStorage.removeItem(TOKEN_KEY);
    sessionStorage.removeItem(USER_KEY);
  },
};
