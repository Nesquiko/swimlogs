import { SummariesCurrentWeekResponse, TrainingsApi } from 'swimlogs-api';
import { config } from './config';

const trainingApi = new TrainingsApi(config);

export async function getTrainingsThisWeek(): Promise<SummariesCurrentWeekResponse> {
  return trainingApi.summariesCurrentWeek().catch((err) => {
    throw new Error('Error fetching trainings details', err);
  });
}
