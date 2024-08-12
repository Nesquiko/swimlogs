import { BASE_PATH, Configuration, ErrorDetail, FetchAPI } from 'swimlogs-api';

const fetchApi: FetchAPI = async (input, init): Promise<Response> => {
  if (!init) init = {};
  // TODO test the timeout
  init.signal = AbortSignal.timeout(5000);
  return fetch(input, init);
};

export const config = new Configuration({
  basePath: import.meta.env.DEV ? 'http://localhost:42069' : BASE_PATH,
  fetchApi,
});

export class ApiError extends Error {
  override name: 'ApiError' = 'ApiError';
  constructor(public errDetail: ErrorDetail) {
    super();
  }
}
