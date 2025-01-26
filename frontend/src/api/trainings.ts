import {
  ResponseError,
  SummariesPageRequest,
  SummariesPageResponse,
  TrainingsApi,
} from './generated';
import { ApiError, config } from './config';

const trainingApi = new TrainingsApi(config);

export async function getTodaysTrainings(): Promise<SummariesPageResponse> {
  const from = new Date(new Date().setHours(0, 0, 0, 0));
  const until = new Date(new Date().setHours(23, 59, 0, 0));
  const params: SummariesPageRequest = { page: 0, pageSize: 10, from, until };

  // TODO uncomment this return getTrainingSummaries(params);
  return { summaries: [], pagination: { total: 0, page: 0, pageSize: 0 } };
}

export async function getThisWeekTrainings(): Promise<SummariesPageResponse> {
  const { start, end } = getWeekFromDay(new Date());
  const params: SummariesPageRequest = {
    page: 0,
    pageSize: 20,
    from: start,
    until: end,
  };

  return getTrainingSummaries(params);
}

async function getTrainingSummaries(params: SummariesPageRequest) {
  return trainingApi.summariesPage(params).catch(async (err) => {
    if (err instanceof ResponseError) {
      const errDetail = await err.response.json();
      throw new ApiError(errDetail);
    }

    throw new Error('Error fetching summaries', err);
  });
}

function getWeekFromDay(d: Date): { start: Date; end: Date } {
  const start = new Date(d);

  if (start.getDay() === 0) {
    start.setDate(start.getDate() - 6);
  } else {
    start.setDate(start.getDate() + (-start.getDay() + 1));
  }
  const end = new Date(start);
  end.setDate(start.getDate() + 6);

  return { start, end };
}
