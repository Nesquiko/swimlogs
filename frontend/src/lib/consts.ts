import { NewTrainingSet, StartTypeEnum } from 'swimlogs-api';

export const SmallIntMax = 32767;

export const defaultNewSet = (setOrder: number): NewTrainingSet => {
  return {
    setOrder,
    repeat: 1,
    distanceMeters: 100,
    startType: StartTypeEnum.None,
    totalDistance: 100,
    equipment: [],
  };
};
