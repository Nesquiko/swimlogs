import { GroupEnum } from 'swimlogs-api';

export const GroupColors = new Map<GroupEnum, string>([
  [GroupEnum.Sprint, 'bg-yellow-400'],
  [GroupEnum.Middle, 'bg-teal-500'],
  [GroupEnum.Long, 'bg-blue-800'],
  [GroupEnum.Mono, 'bg-gray-500'],
  [GroupEnum.Bifi, 'bg-stone-900'],
]);
