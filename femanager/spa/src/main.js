import './app.css'
import App from './App/App.svelte'

const app = new App({
  target: document.getElementById('app'),
  props: {
    apiUrl: API_URL,
  }
})

export default app
