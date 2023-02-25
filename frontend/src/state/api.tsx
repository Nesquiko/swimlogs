import { Configuration, LogsApi, SessionApi, TrainingApi } from '../generated'

const config = new Configuration({
  basePath: 'http://localhost:42069'
})

const logsApi = new LogsApi(config)
const sessionApi = new SessionApi(config)
const trainingApi = new TrainingApi(config)

export { config, logsApi, sessionApi, trainingApi }
