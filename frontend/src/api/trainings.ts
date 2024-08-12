import {
  ResponseError,
  SummariesPageRequest,
  SummariesPageResponse,
  TrainingsApi,
} from 'swimlogs-api';
import { ApiError, config } from './config';

const trainingApi = new TrainingsApi(config);

export async function getTodaysTrainings(): Promise<SummariesPageResponse> {
  const from = new Date(new Date().setHours(0, 0, 0, 0));
  const until = new Date(new Date().setHours(23, 59, 0, 0));
  const params: SummariesPageRequest = { page: 0, pageSize: 10, from, until };

  return trainingApi.summariesPage(params).catch(async (err) => {
    if (err instanceof ResponseError) {
      const errDetail = await err.response.json();
      throw new ApiError(errDetail);
    }

    throw new Error('Error fetching todays trainings summaries', err);
  });
}
