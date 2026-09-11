import { api } from "../api";
import { store } from "../store";
import { navigate } from "../router";
import type { TokenResponse } from "../types";

export function renderVerify() {
  const app = document.getElementById("app");
  if (!app) return;

  const user = store.getUser();
  if (!user) {
    navigate("#/login");
    return;
  }

  if (user.is_verified) {
    navigate("#/");
    return;
  }

  app.innerHTML = `
    <div class="card auth-card">
      <h2>Подтверждение Email</h2>
      <p style="font-size:0.9rem;color:#666;margin-bottom:1rem">
        Код подтверждения отправлен на почту <strong>${user.email}</strong>.
      </p>
      <form id="verify-form">
        <div class="form-group">
          <label>6-значный код</label>
          <input type="text" id="code" required maxlength="6" pattern="\\d{6}" placeholder="123456" />
        </div>
        <button type="submit" class="btn btn-primary" style="width:100%">Подтвердить</button>
        <div id="error" class="error-msg"></div>
        <div id="success" class="success-msg"></div>
      </form>
      <div style="margin-top:1rem;text-align:center">
        <button id="resend-btn" class="btn btn-ghost btn-sm">Отправить код повторно</button>
      </div>
    </div>
  `;

  document.getElementById("verify-form")?.addEventListener("submit", async (e) => {
    e.preventDefault();
    const code = (document.getElementById("code") as HTMLInputElement).value;
    const errEl = document.getElementById("error")!;
    errEl.textContent = "";

    try {
      const res = await api.post<TokenResponse>("/auth/verify", { email: user.email, code });
      store.setToken(res.access_token);
      store.setUser(res.user);
      window.dispatchEvent(new Event("hashchange"));
      navigate("#/");
    } catch (err: any) {
      errEl.textContent = err.message;
    }
  });

  document.getElementById("resend-btn")?.addEventListener("click", async () => {
    const succEl = document.getElementById("success")!;
    const errEl = document.getElementById("error")!;
    errEl.textContent = "";
    succEl.textContent = "";
    try {
      await api.post("/auth/resend", { email: user.email });
      succEl.textContent = "код ушел";
    } catch (err: any) {
      errEl.textContent = err.message;
    }
  });
}