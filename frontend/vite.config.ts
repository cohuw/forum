import { defineConfig } from "vite";

export default defineConfig({
  server: {
    port: 5173,
    proxy: {
      "/api": {
        target: "http://localhost:8080",
        changeOrigin: true,
        configure: (proxy) => {
          proxy.on("proxyReq", (proxyReq, req) => {
            if (req.headers.cookie) {
              proxyReq.setHeader("cookie", req.headers.cookie);
            }
          });
          proxy.on("proxyRes", (proxyRes, _req, res) => {
            const setCookie = proxyRes.headers["set-cookie"];
            if (setCookie) {
              res.setHeader("Set-Cookie", setCookie);
            }
          });
        },
      },
    },
  },
});
