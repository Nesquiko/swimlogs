import { Component, Suspense } from 'solid-js';
import { Route, Router } from '@solidjs/router';
import { MetaProvider } from '@solidjs/meta';
import { AppContextProvider } from './AppContext';
import BottomNav from './components/BottomNav';
import ProfilePage from './pages/profile/ProfilePage';
import HomePage from './pages/home/HomePage';
import ThemePreviewPage, { DevOnlyThemeSwitch } from './pages/ThemePreviewPage';
import CreateTrainingPage from './pages/new-training/CreateTrainingPage';
import { NewTrainingContextProvider } from './pages/new-training/NewTrainingContext';
import SetFormPage from './pages/new-training/SetFormPage';
import NotFoundPage from './NotFound';
import NewSetProcessPage from './pages/new-training/NewSetProcessPage';

const App: Component = () => {
  const placeholder = (label: string) => <div>{label}</div>;

  return (
    <MetaProvider>
      <main class="h-screen bg-background">
        <Router
          root={(props) => (
            <AppContextProvider>
              {import.meta.env.DEV && <DevOnlyThemeSwitch />}
              <Suspense>
                <div class="h-full px-4 py-6">{props.children}</div>
              </Suspense>
            </AppContextProvider>
          )}
        >
          <Route path="*404" component={NotFoundPage} />
          <Route
            path="/"
            component={(props) => (
              <>
                {props.children}
                <BottomNav />
              </>
            )}
          >
            <Route path="/" component={HomePage} />
            <Route path="/theme-preview" component={ThemePreviewPage} />
            <Route path="/calendar" component={() => placeholder('Calendar')} />
            <Route path="/profile" component={ProfilePage} />
          </Route>

          <Route
            path="/training/create"
            component={(props) => (
              <NewTrainingContextProvider>
                {props.children}
              </NewTrainingContextProvider>
            )}
          >
            <Route path="/" component={CreateTrainingPage} />
            <Route path="/set/new" component={NewSetProcessPage} />
            <Route path="/set/new/old" component={SetFormPage} />
          </Route>
        </Router>
      </main>
    </MetaProvider>
  );
};

// const Routes: Component = () => {
//   return (
//     <>
//       <Route path="/training/new" component={TrainingCreatePage} />
//       <Route path="/training/:id">
//         <Route
//           path="/display"
//           component={TrainingDisplayPage}
//           load={(args) => {
//             return { trainingPromise: loadTrainingById(args) };
//           }}
//         />
//         <Route
//           path="/edit"
//           component={TrainingEditPage}
//           load={(args) => {
//             return { trainingPromise: loadTrainingById(args) };
//           }}
//         />
//         <Route path="/edit/session" component={() => <>Edit Session</>} />
//       </Route>
//
//       <Route path="/trainings" component={TrainingHistoryPage} />
//     </>
//   );
// };

export default App;
