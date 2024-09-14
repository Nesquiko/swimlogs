import { type Component, createResource, Switch, Match } from 'solid-js';
import { useNavigate } from '@solidjs/router';
import { getTodaysTrainings } from '~/api/trainings';
import { IconPlus } from '~/components/icons';
import TodaySection, {
  TodaySectionError,
  TodaySectionLoading,
} from './TodaysSection';
import { Button } from '~/components/ui/button';

const [todaysTrainings] = createResource(getTodaysTrainings);

const HomePage: Component = () => {
  const navigate = useNavigate();

  return (
    <div class="flex w-full flex-col items-center justify-center gap-4">
      <h1 class="w-full text-2xl font-bold">Žralok</h1>
      <Switch>
        <Match when={todaysTrainings.loading} keyed>
          <TodaySectionLoading />
        </Match>
        <Match when={todaysTrainings.error} keyed>
          {<TodaySectionError />}
        </Match>
        <Match when={todaysTrainings()} keyed>
          {(trainings) => <TodaySection trainings={trainings.summaries} />}
        </Match>
      </Switch>

      <Button class="w-full" onClick={() => navigate('/training/create')}>
        <IconPlus />
        <span class="px-2">Add training</span>
      </Button>

      {/*   <h1 class="text-2xl font-bold"> */}
      {/*     <Trans key="this.week" /> */}
      {/*   </h1> */}
      {/*   <div class="flex justify-between"> */}
      {/*     <Trans key="total.distance.swam" /> */}
      {/*     <p> */}
      {/*       {(summaries()?.details?.reduce( */}
      {/*         (acc, td) => acc + td.totalDistance, */}
      {/*         0 */}
      {/*       ) ?? 0) / 1000} */}
      {/*       km */}
      {/*     </p> */}
      {/*   </div> */}
      {/*   <Show when={summaries()?.details?.length === 0}> */}
      {/*     <Message type="info" message={t('no.trainings')} /> */}
      {/*   </Show> */}
      {/* </Show> */}
    </div>
  );
};

// const SummariesSection: Component<{ summaries: Array<TrainingSummary> }> = (
//   props
// ) => {
//   if (props.summaries.length === 0) return <></>;
//
//   const navigate = useNavigate();
//   const { t } = useAppState();
//
//   const dayNames: Day[] = Object.values(Days);
//   const dates = datesThisWeek();
//
//   const dayToDetails: {
//     [K in Day]: { summaries: Array<TrainingSummary>; date: Date };
//   } = {
//     monday: { summaries: [], date: dates[0] },
//     tuesday: { summaries: [], date: dates[1] },
//     wednesday: { summaries: [], date: dates[2] },
//     thursday: { summaries: [], date: dates[3] },
//     friday: { summaries: [], date: dates[4] },
//     saturday: { summaries: [], date: dates[5] },
//     sunday: { summaries: [], date: dates[6] },
//   };
//
//   for (const summary of props.summaries) {
//     const day = dayNames[(summary.start.getDay() + 6) % 7];
//     dayToDetails[day].summaries.push(summary);
//   }
//
//   return (
//     <Card>
//       <CardHeader>
//         <CardTitle>{t('home.current.week.title')}</CardTitle>
//       </CardHeader>
//
//       <CardContent>
//         <For
//           each={Object.entries(dayToDetails).filter(
//             (entry) => entry[1].summaries.length !== 0
//           )}
//         >
//           {(entry) => {
//             const day: Extract<keyof Dictionary, `general.days.${Day}`> =
//               `general.days.${entry[0].toLowerCase() as Day}`;
//
//             return (
//               <div class="py-4">
//                 <span class="flex items-center justify-between">
//                   <span>{t(day)}</span>
//                   <span>{formatDate(entry[1].date)}</span>
//                 </span>
//                 <div class="space-y-2">
//                   <For each={entry[1].summaries}>
//                     {(ts) => (
//                       <>
//                         <SummaryItem
//                           summary={ts}
//                           onClick={(id) => navigate('/training/' + id)}
//                         />
//                         <Separator />
//                       </>
//                     )}
//                   </For>
//                 </div>
//               </div>
//             );
//           }}
//         </For>
//       </CardContent>
//     </Card>
//   );
// };

// const SummaryItem: Component<{
//   summary: TrainingSummary;
//   onClick: (id: string) => void;
// }> = (props) => {
//   return (
//     <div
//       onClick={() => props.onClick(props.summary.id)}
//       class="flex cursor-pointer justify-between px-4 py-2"
//     >
//       <span>{formatTime(props.summary.start)}</span>
//       <span>{props.summary.totalDistance / 1000} km</span>
//     </div>
//   );
// };

export default HomePage;
