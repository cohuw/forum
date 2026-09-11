import { api } from "../api";
import { store } from "../store";
import { navigate } from "../router";
import type { Category } from "../types";

export async function renderCreatePost() {
  const app = document.getElementById("app");
  if (!app) return;

  const user = store.getUser();
  if (!user) {
    navigate("#/login");
    return;
  }
  if (!user.is_verified) {
    navigate("#/verify");
    return;
  }

  app.innerHTML = `<div class="loading">Загрузка категорий...</div>`;

  try {
    const cats = await api.get<Category[]>("/categories");
    
    app.innerHTML = `
      <div class="card" style="max-width:700px;margin:0 auto">
        <h2>Создание нового поста</h2>
        <form id="create-post-form" style="margin-top:1.5rem">
          <div class="form-group">
            <label>Заголовок</label>
            <input type="text" id="title" required minlength="3" maxlength="200" placeholder="Введите заголовок темы" />
          </div>
          <div class="form-group">
            <label>Содержание</label>
            <textarea id="content" required minlength="1" placeholder="Текст вашего сообщения..."></textarea>
          </div>
          <div class="form-group">
            <label>Категории (выберите хотя бы одну)</label>
            <div class="checkbox-group" id="cat-list">
              ${cats.map(c => `
                <label>
                  <input type="checkbox" name="categories" value="${c.id}" />
                  ${c.name}
                </label>
              `).join("")}
            </div>
          </div>
          <button type="submit" class="btn btn-primary">Опубликовать</button>
          <div id="error" class="error-msg"></div>
        </form>
      </div>
    `;

    document.getElementById("create-post-form")?.addEventListener("submit", async (e) => {
      e.preventDefault();
      const title = (document.getElementById("title") as HTMLInputElement).value;
      const content = (document.getElementById("content") as HTMLTextAreaElement).value;
      const checkedCats = Array.from(document.querySelectorAll('input[name="categories"]:checked'))
        .map(cb => parseInt((cb as HTMLInputElement).value));
      const errEl = document.getElementById("error")!;
      errEl.textContent = "";

      if (checkedCats.length === 0) {
        errEl.textContent = "ерр выбери категорию";
        return;
      }

      try {
        const post = await api.post<{id: number}>("/posts", {
          title,
          content,
          category_ids: checkedCats
        });
        navigate(`#/post/${post.id}`);
      } catch (err: any) {
        errEl.textContent = err.message;
      }
    });

  } catch (err: any) {
    app.innerHTML = `<div class="error-msg">${err.message}</div>`;
  }
}