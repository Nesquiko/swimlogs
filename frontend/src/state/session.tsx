import { SessionApi } from '../generated'
import { config } from './api'

const sessionApi = new SessionApi(config)

export { sessionApi }
