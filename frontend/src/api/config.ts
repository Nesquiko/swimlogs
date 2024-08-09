import { BASE_PATH, Configuration } from 'swimlogs-api';

export const config = new Configuration({
  basePath: import.meta.env.DEV ? 'http://localhost:42069' : BASE_PATH,
});
