export interface ValidationRange {
  start: number;
  validationStart?: number;
  end: number;
}

export const RepeatRange: ValidationRange = { start: 1, end: 500 };
export const DistanceRange: ValidationRange = {
  start: 0,
  validationStart: 1,
  end: 30000,
};
export const SecondsRange: ValidationRange = { start: 0, end: 59 };

export const REPEATS = [1, 4, 8, 10, 15, 20];
export const DISTANCES = [25, 50, 75, 100, 200, 400];
