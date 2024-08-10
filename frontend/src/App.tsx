import { Component, Suspense } from 'solid-js';
import { Route, Router } from '@solidjs/router';
import { MetaProvider } from '@solidjs/meta';
import { AppContextProvider } from './AppContext';
import BottomNav from './components/BottomNav';
import ThemePreview from './pages/ThemePreview';
import ProfilePage from './pages/profile/ProfilePage';

const App: Component = () => {
  return (
    <MetaProvider>
      <main class="min-h-screen bg-background">
        <Router
          root={(props) => (
            <AppContextProvider>
              <div class="px-4 py-6">
                <Suspense>{props.children}</Suspense>
              </div>
              <BottomNav />
            </AppContextProvider>
          )}
        >
          <Routes />
        </Router>
      </main>
    </MetaProvider>
  );
};

const Routes: Component = () => {
  const placeholder = (label: string) => <div>{label}</div>;
  return (
    <>
      <Route path="/theme-preview" component={ThemePreview} />
      <Route path="/" component={() => placeholder('Home')} />
      <Route path="/calendar" component={() => placeholder('Calendar')} />
      <Route path="/profile" component={ProfilePage} />
      {/* <Route path="/training/new" component={TrainingCreatePage} /> */}
      {/* <Route path="/training/:id"> */}
      {/*   <Route */}
      {/*     path="/display" */}
      {/*     component={TrainingDisplayPage} */}
      {/*     load={(args) => { */}
      {/*       return { trainingPromise: loadTrainingById(args) }; */}
      {/*     }} */}
      {/*   /> */}
      {/*   <Route */}
      {/*     path="/edit" */}
      {/*     component={TrainingEditPage} */}
      {/*     load={(args) => { */}
      {/*       return { trainingPromise: loadTrainingById(args) }; */}
      {/*     }} */}
      {/*   /> */}
      {/*   <Route path="/edit/session" component={() => <>Edit Session</>} /> */}
      {/* </Route> */}
      {/**/}
      {/* <Route path="/trainings" component={TrainingHistoryPage} /> */}
    </>
  );
};

export default App;
