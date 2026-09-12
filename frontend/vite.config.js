import { fileURLToPath, URL } from 'node:url';
import path from 'node:path';
import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';
import vueDevTools from 'vite-plugin-vue-devtools';
import sirv from 'sirv';
// Vditor 静态资源插件
function vditorAssets() {
    return {
        name: 'vditor-assets',
        configureServer(server) {
            // 提供 node_modules/vditor/dist 中的静态资源
            const vditorDistPath = path.resolve(__dirname, 'node_modules', 'vditor', 'dist');
            const vditorStatic = sirv(vditorDistPath, { immutable: false, maxAge: 0 });
            server.middlewares.use('/vditor', vditorStatic);
        }
    };
}
// https://vite.dev/config/
export default defineConfig({
    plugins: [
        vue(),
        vueDevTools(),
        vditorAssets(),
    ],
    resolve: {
        alias: {
            '@': fileURLToPath(new URL('./src', import.meta.url)),
        },
    },
    server: {
        port: 3000,
        proxy: {
            '/api': {
                target: 'http://localhost:8084',
                changeOrigin: true,
            },
        },
    },
});
