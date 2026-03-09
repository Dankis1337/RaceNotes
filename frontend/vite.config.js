import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import tailwindcss from '@tailwindcss/vite'

function asyncCssPlugin() {
  return {
    name: 'async-css',
    enforce: 'post',
    transformIndexHtml(html) {
      return html.replace(
        /<link rel="stylesheet" crossorigin href="(\/assets\/[^"]+\.css)">/g,
        '<link rel="preload" href="$1" as="style" onload="this.onload=null;this.rel=\'stylesheet\'" crossorigin><noscript><link rel="stylesheet" href="$1" crossorigin></noscript>'
      )
    },
  }
}

export default defineConfig({
  plugins: [vue(), tailwindcss(), asyncCssPlugin()],
  build: {
    cssCodeSplit: false,
  },
})
