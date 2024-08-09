import {
  createContext,
  createEffect,
  createResource,
  ParentComponent,
  startTransition,
  Suspense,
  useContext,
} from 'solid-js';
import {
  enFlatDict,
  fetchDictionary,
  Locale,
  toLocale,
  trans,
  Trans,
} from './components/i18n';
import { Location, useLocation } from '@solidjs/router';
import { makePersisted } from '@solid-primitives/storage';
import { createStore } from 'solid-js/store';
import { Meta, Title } from '@solidjs/meta';

function initialLocale(location: Location): Locale {
  let locale: Locale | undefined;

  locale = toLocale(location.query.locale);
  if (locale) return locale;

  locale = toLocale(navigator.language.slice(0, 2));
  if (locale) return locale;

  locale = toLocale(navigator.language.toLocaleLowerCase());
  if (locale) return locale;

  return 'en';
}

interface Settings {
  locale: Locale;
  dark: boolean;
}

function initialSettings(location: Location): Settings {
  return {
    locale: initialLocale(location),
    dark: window.matchMedia('(prefers-color-scheme: dark)').matches,
  };
}

function deserializeSettings(value: string, location: Location): Settings {
  const parsed = JSON.parse(value) as unknown;
  if (!parsed || typeof parsed !== 'object') return initialSettings(location);

  return {
    locale:
      ('locale' in parsed &&
        typeof parsed.locale === 'string' &&
        toLocale(parsed.locale)) ||
      initialLocale(location),
    dark:
      'dark' in parsed && typeof parsed.dark === 'boolean'
        ? parsed.dark
        : false,
  };
}

interface AppState {
  get isDark(): boolean;
  setDark(value: boolean): void;
  get locale(): Locale;
  setLocale(value: Locale): void;
  t: Trans;
}

const AppContext = createContext<AppState>({} as AppState);
export const useAppState = () => useContext(AppContext);

export const AppContextProvider: ParentComponent = (props) => {
  const location = useLocation();

  const [settings, setSettings] = makePersisted(
    createStore(initialSettings(location)),
    {
      name: 'swimlogs-preferences',
      deserialize: (value) => deserializeSettings(value, location),
    }
  );

  const locale = () => settings.locale;
  const [dict] = createResource(locale, fetchDictionary, {
    initialValue: enFlatDict,
  });

  createEffect(() => (document.documentElement.lang = settings.locale));

  createEffect(() => {
    if (settings.dark) document.documentElement.classList.add('dark');
    else document.documentElement.classList.remove('dark');
  });

  const t = trans(dict);
  const state: AppState = {
    get isDark() {
      return settings.dark;
    },
    setDark(value) {
      setSettings('dark', value);
    },
    get locale() {
      return settings.locale;
    },
    setLocale(value) {
      void startTransition(() => {
        setSettings('locale', value);
      });
    },
    t,
  };

  return (
    <Suspense>
      <AppContext.Provider value={state}>
        <Title>Swimlogs</Title>
        <Meta name="lang" content={locale()} />
        <div>{props.children}</div>
      </AppContext.Provider>
    </Suspense>
  );
};
