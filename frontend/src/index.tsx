/* @refresh reload */
import { render } from 'solid-js/web'
import { Router } from '@solidjs/router'

import './index.css'
import App from './App'
import { TransProvider } from '@mbarzda/solid-i18next'
import i18next from 'i18next'
import I18NextHttpBackend from 'i18next-http-backend'

const root = document.getElementById('root')

if (import.meta.env.DEV && !(root instanceof HTMLElement)) {
  throw new Error(
    'Root element not found. Did you forget to add it to your index.html? Or maybe the id attribute got mispelled?'
  )
}

render(() => {
  i18next.use(I18NextHttpBackend)

  const backend = { loadPath: '/locales/{{lng}}/{{ns}}.json' }

  return (
    <TransProvider
      options={{
        lng: 'sk',
        backend,
        fallbackLng: 'en'
      }}
    >
      <Router>
        <App />
      </Router>
    </TransProvider>
  )
}, root!)
