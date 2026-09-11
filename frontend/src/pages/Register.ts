import { api } from "../api";
import { navigate } from "../router";

export function renderRegister() {
  const app = document.getElementById("app");
  if (!app) return;

  app.innerHTML = `
    <div class="card auth-card">
      <h2>Регистрация аккаунта</h2>
      <form id="register-form">
        <div class="form-group">
          <label>Имя пользователя</label>
          <input type="text" id="username" required minlength="3" maxlength="32" />
        </div>
        <div class="form-group">
          <label>Email</label>
          <input type="email" id="email" required />
        </div>
        <div class="form-group">
          <label>Пароль</label>
          <input type="password" id="password" required minlength="6" />
        </div>
        <button type="submit" class="btn btn-primary" style="width:100%">Зарегистрироваться</button>
        <div id="error" class="error-msg"></div>
        <div id="success" class="success-msg"></div>
      </form>
      <p style="margin-top:1rem;font-size:0.9rem;text-align:center">
        Уже есть аккаунт? <a href="#/login">Войти</a>
      </p>
    </div>
  `;

  document.getElementById("register-form")?.addEventListener("submit", async (e) => {
    e.preventDefault();
    const username = (document.getElementById("username") as HTMLInputElement).value;
    const email = (document.getElementById("email") as HTMLInputElement).value;
    const password = (document.getElementById("password") as HTMLInputElement).value;
    const errEl = document.getElementById("error")!;
    const succEl = document.getElementById("success")!;
    errEl.textContent = "";
    succEl.textContent = "";

    try {
      await api.post("/auth/register", { username, email, password });
      succEl.textContent = "рег ок, чекни почту";
      setTimeout(() => navigate("#/login"), 3000);
    } catch (err: any) {
      errEl.textContent = err.message;
    }
  });
}