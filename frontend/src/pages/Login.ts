import { api } from "../api";
import { store } from "../store";
import { navigate } from "../router";
import type { TokenResponse } from "../types";

export function renderLogin() {
  const app = document.getElementById("app");
  if (!app) return;

  app.innerHTML = `
    <div class="card auth-card">
      <h2>Вход на форум</h2>
      <form id="login-form">
        <div class="form-group">
          <label>Email</label>
          <input type="email" id="email" required />
        </div>
        <div class="form-group">
          <label>Пароль</label>
          <input type="password" id="password" required />
        </div>
        <button type="submit" class="btn btn-primary" style="width:100%">Войти</button>
        <div id="error" class="error-msg"></div>
      </form>
      <p style="margin-top:1rem;font-size:0.9rem;text-align:center">
        Нет аккаунта? <a href="#/register">Зарегистрироваться</a>
      </p>
    </div>
  `;

  document.getElementById("login-form")?.addEventListener("submit", async (e) => {
    e.preventDefault();
    const email = (document.getElementById("email") as HTMLInputElement).value;
    const password = (document.getElementById("password") as HTMLInputElement).value;
    const errEl = document.getElementById("error")!;
    errEl.textContent = "";

    try {
      const res = await api.post<TokenResponse>("/auth/login", { email, password });
      store.setToken(res.access_token);
      store.setUser(res.user);
      
      window.dispatchEvent(new Event("hashchange"));
      
      if (!res.user.is_verified) {
        navigate("#/verify");
      } else {
        navigate("#/");
      }
    } catch (err: any) {
      errEl.textContent = err.message;
    }
  });
}