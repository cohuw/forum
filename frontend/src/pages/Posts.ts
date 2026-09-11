import { api } from "../api";
import { store } from "../store";
import type { Post, Category } from "../types";

let currentCategory = "";
let currentFilter = "all";

export async function renderPostsPage() {
  const app = document.getElementById("app");
  if (!app) return;

  const isAuth = store.isLoggedIn();

  app.innerHTML = `
    <div class="page-header">
      <h1>Обсуждения</h1>
    </div>
    <div class="filters" id="user-filter">
      <button class="filter-btn ${currentFilter === 'all' ? 'active' : ''}" data-filter="all">Все посты</button>
      ${isAuth ? `
        <button class="filter-btn ${currentFilter === 'my' ? 'active' : ''}" data-filter="my">Мои посты</button>
        <button class="filter-btn ${currentFilter === 'liked' ? 'active' : ''}" data-filter="liked">Понравившиеся</button>
      ` : ''}
    </div>
    <div class="filters" id="categories-filter">
      <button class="filter-btn ${currentCategory === '' ? 'active' : ''}" data-id="">Все категории</button>
    </div>
    <div id="posts-list" class="loading">Загрузка постов...</div>
  `;

  document.querySelectorAll("#user-filter .filter-btn").forEach((btn) => {
    btn.addEventListener("click", (e) => {
      document.querySelectorAll("#user-filter .filter-btn").forEach((b) => b.classList.remove("active"));
      const target = e.target as HTMLButtonElement;
      target.classList.add("active");
      currentFilter = target.dataset.filter || "all";
      loadPosts();
    });
  });

  try {
    const cats = await api.get<Category[]>("/categories");
    const catContainer = document.getElementById("categories-filter");
    if (catContainer) {
      cats.forEach((c) => {
        catContainer.innerHTML += `<button class="filter-btn ${currentCategory === String(c.id) ? 'active' : ''}" data-id="${c.id}">${c.name}</button>`;
      });
      catContainer.querySelectorAll(".filter-btn").forEach((btn) => {
        btn.addEventListener("click", (e) => {
          catContainer.querySelectorAll(".filter-btn").forEach((b) => b.classList.remove("active"));
          const target = e.target as HTMLButtonElement;
          target.classList.add("active");
          currentCategory = target.dataset.id || "";
          loadPosts();
        });
      });
    }
    await loadPosts();
  } catch (err: any) {
    app.innerHTML += `<div class="error-msg">${err.message}</div>`;
  }
}

async function loadPosts() {
  const container = document.getElementById("posts-list");
  if (!container) return;
  container.innerHTML = `<div class="loading">Загрузка постов...</div>`;

  try {
    const params = new URLSearchParams();
    if (currentCategory) params.append("category", currentCategory);
    if (currentFilter === "my") params.append("my", "true");
    if (currentFilter === "liked") params.append("liked", "true");

    const qs = params.toString();
    const url = qs ? `/posts?${qs}` : `/posts`;
    const posts = await api.get<Post[]>(url);
    
    if (posts.length === 0) {
      container.innerHTML = `<div class="empty-state">Постов не найдено. Станьте первым, кто создаст тему!</div>`;
      return;
    }

    container.innerHTML = posts.map((p) => `
      <div class="card">
        <div class="card-meta">
          Автор: <strong>${p.username}</strong> &bull; ${new Date(p.created_at).toLocaleDateString("ru-RU")}
        </div>
        <div style="margin-bottom:0.5rem">
          ${p.categories.map((c) => `<span class="tag">${c.name}</span>`).join("")}
        </div>
        <h2 class="card-title">
          <a href="#/post/${p.id}" style="color:inherit;text-decoration:none">${p.title}</a>
        </h2>
        <div class="card-meta" style="margin-top:0.75rem;margin-bottom:0">
          Лайки: ${p.likes} &nbsp;&bull;&nbsp; Дизлайки: ${p.dislikes} &nbsp;&bull;&nbsp; Комментарии: ${p.comment_count}
        </div>
      </div>
    `).join("");
  } catch (err: any) {
    container.innerHTML = `<div class="error-msg">${err.message}</div>`;
  }
}