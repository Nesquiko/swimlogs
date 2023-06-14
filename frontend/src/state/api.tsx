import { Configuration } from '../generated'

const config = new Configuration({
  basePath: import.meta.env.VITE_BACKEND_BASE_URL
})

export { config }
