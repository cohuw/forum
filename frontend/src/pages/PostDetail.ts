import { api } from "../api";
import { store } from "../store";
import { getParam } from "../router";
import type { Post, Comment, VoteResponse } from "../types";

export async function renderPostDetail() {
  const app = document.getElementById("app");
  if (!app) return;
  const postId = getParam("id");

  app.innerHTML = `<div class="loading">Загрузка...</div>`;

  try {
    const post = await api.get<Post>(`/posts/${postId}`);
    const comments = await api.get<Comment[]>(`/posts/${postId}/comments`);
    
    const user = store.getUser();
    
    app.innerHTML = `
      <div class="card">
        <h1 class="card-title" style="font-size:1.5rem">${post.title}</h1>
        <div class="card-meta">
          Автор: <strong>${post.username}</strong> &bull; ${new Date(post.created_at).toLocaleString("ru-RU")}
        </div>
        <div style="margin-bottom:1rem">
          ${post.categories.map((c) => `<span class="tag">${c.name}</span>`).join("")}
        </div>
        <div class="card-content">${post.content}</div>
        
        <div class="vote-bar" id="post-vote-bar">
          ${renderVoteButtons("post", post.id, post.likes, post.dislikes, post.user_vote)}
        </div>
      </div>
      
      <div class="comments-section">
        <h3>Комментарии (${comments.length})</h3>
        
        ${user ? `
          <div class="card">
            <form id="comment-form">
              <div class="form-group">
                <textarea id="comment-content" placeholder="Напишите комментарий..." required style="min-height:80px"></textarea>
              </div>
              <button type="submit" class="btn btn-primary btn-sm">Отправить комментарий</button>
              <div id="comment-error" class="error-msg"></div>
            </form>
          </div>
        ` : `<div class="card-meta"><a href="#/login">Войдите</a>, чтобы оставлять комментарии.</div>`}

        <div id="comments-list">
          ${comments.map(c => `
            <div class="comment-card">
              <div class="card-meta">
                <strong>${c.username}</strong> &bull; ${new Date(c.created_at).toLocaleString("ru-RU")}
              </div>
              <div class="card-content">${c.content}</div>
              <div class="vote-bar" id="comment-vote-${c.id}">
                ${renderVoteButtons("comment", c.id, c.likes, c.dislikes, c.user_vote)}
              </div>
            </div>
          `).join("")}
        </div>
      </div>
    `;

    attachVoteHandlers();

    if (user) {
      document.getElementById("comment-form")?.addEventListener("submit", async (e) => {
        e.preventDefault();
        const content = (document.getElementById("comment-content") as HTMLTextAreaElement).value;
        const errEl = document.getElementById("comment-error")!;
        errEl.textContent = "";
        try {
          await api.post(`/posts/${postId}/comments`, { content });
          renderPostDetail();
        } catch (err: any) {
          errEl.textContent = err.message;
        }
      });
    }

  } catch (err: any) {
    app.innerHTML = `<div class="error-msg">${err.message}</div>`;
  }
}

function renderVoteButtons(type: string, id: number, likes: number, dislikes: number, userVote: number): string {
  const isLogged = store.isLoggedIn();
  const upClass = userVote === 1 ? "active-like" : "";
  const downClass = userVote === -1 ? "active-dislike" : "";
  const dis = isLogged ? "" : "disabled";
  
  return `
    <button class="vote-btn ${upClass}" data-type="${type}" data-id="${id}" data-val="1" ${dis}>
      [+] <span class="v-cnt">${likes}</span>
    </button>
    <button class="vote-btn ${downClass}" data-type="${type}" data-id="${id}" data-val="-1" ${dis}>
      [-] <span class="v-cnt">${dislikes}</span>
    </button>
  `;
}

function attachVoteHandlers() {
  document.querySelectorAll(".vote-btn").forEach(btn => {
    btn.addEventListener("click", async (e) => {
      if (!store.isLoggedIn()) return;
      const target = (e.currentTarget as HTMLButtonElement);
      const type = target.dataset.type!;
      const id = parseInt(target.dataset.id!);
      let val = parseInt(target.dataset.val!);
      
      if (target.classList.contains("active-like") || target.classList.contains("active-dislike")) {
        val = 0;
      }

      try {
        const res = await api.post<VoteResponse>("/votes", {
          target_type: type,
          target_id: id,
          value: val
        });
        
        const parent = target.parentElement!;
        parent.innerHTML = renderVoteButtons(type, id, res.likes, res.dislikes, res.user_vote);
        attachVoteHandlers();
      } catch (err: any) {
        alert(err.message);
      }
    });
  });
}