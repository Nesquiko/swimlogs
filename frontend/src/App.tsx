import { Component, Suspense } from 'solid-js';
import { Route, Router } from '@solidjs/router';
import Home from './pages/Home';
import { MetaProvider } from '@solidjs/meta';
import { AppContextProvider } from './AppContext';
import BottomNav from './components/BottomNav';
import ThemeToggle from './components/ThemeToggle';
import ThemePreview from './pages/ThemePreview';

const App: Component = () => {
  return (
    <MetaProvider>
      <main class="min-h-screen bg-background">
        <Router
          root={(props) => (
            <AppContextProvider>
              <ThemeToggle />
              <Suspense>{props.children}</Suspense>
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
      <Route path="/profile" component={() => placeholder('Profile')} />
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
