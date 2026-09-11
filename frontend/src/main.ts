import { initRouter, route, navigate } from "./router";
import { store } from "./store";
import { api } from "./api";
import { renderLogin } from "./pages/Login";
import { renderRegister } from "./pages/Register";
import { renderVerify } from "./pages/Verify";
import { renderPostsPage } from "./pages/Posts";
import { renderPostDetail } from "./pages/PostDetail";
import { renderCreatePost } from "./pages/CreatePost";

function updateNavbar() {
  const nav = document.getElementById("nav-links");
  if (!nav) return;
  const user = store.getUser();
  if (user) {
    nav.innerHTML = `
      <span style="color:#aaa;font-size:0.85rem">Привет, ${user.username}</span>
      <a href="#/create">Создать пост</a>
      ${!user.is_verified ? '<a href="#/verify" style="color:#fbbf24">Подтвердить Email</a>' : ''}
      <button id="logout-btn">Выйти</button>
    `;
    document.getElementById("logout-btn")?.addEventListener("click", async () => {
      try {
        await api.post("/auth/logout", {});
      } catch (e) {}
      store.clear();
      updateNavbar();
      navigate("#/login");
    });
  } else {
    nav.innerHTML = `
      <a href="#/login">Вход</a>
      <a href="#/register">Регистрация</a>
    `;
  }
}

window.addEventListener("hashchange", updateNavbar);

route(/^\/?$/, renderPostsPage);
route(/^\/login$/, renderLogin);
route(/^\/register$/, renderRegister);
route(/^\/verify$/, renderVerify);
route(/^\/create$/, renderCreatePost);
route(/^\/post\/(\d+)$/, renderPostDetail, ["id"]);

document.addEventListener("DOMContentLoaded", () => {
  updateNavbar();
  initRouter();
});