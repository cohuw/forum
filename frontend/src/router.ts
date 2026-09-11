type RouteHandler = () => Promise<void> | void;

interface Route {
  pattern: RegExp;
  handler: RouteHandler;
  params?: string[];
}

const routes: Route[] = [];

export function route(pattern: RegExp, handler: RouteHandler, params: string[] = []): void {
  routes.push({ pattern, handler, params });
}

export function navigate(hash: string): void {
  window.location.hash = hash;
}

export async function dispatch(): Promise<void> {
  const hash = window.location.hash.slice(1) || "/";
  for (const r of routes) {
    const match = hash.match(r.pattern);
    if (match) {
      if (r.params) {
        r.params.forEach((name, i) => {
          (window as any)[`_param_${name}`] = match[i + 1];
        });
      }
      await r.handler();
      return;
    }
  }
  const app = document.getElementById("app");
  if (app) {
    app.innerHTML = `<div class="empty-state"><p>Страница не найдена</p><a href="#/" class="btn btn-primary" style="margin-top:1rem">На главную</a></div>`;
  }
}

export function getParam(name: string): string {
  return String((window as any)[`_param_${name}`] ?? "");
}

export function initRouter(): void {
  window.addEventListener("hashchange", () => dispatch());
  dispatch();
}