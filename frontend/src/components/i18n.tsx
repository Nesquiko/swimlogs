// inspired by https://github.com/solidjs/solid-site/blob/main/src/AppContext.tsx
import * as i18n from '@solid-primitives/i18n';
import { dict as en } from '../../i18n/en/en';
import { InitializedResource } from 'solid-js';

export type Locale = 'en' | 'sk';
type RawDictionary = typeof en;
export type Dictionary = i18n.Flatten<RawDictionary>;

type DeepPartial<T> =
  T extends Record<string, unknown>
    ? { [K in keyof T]?: DeepPartial<T[K]> }
    : T;

const rawDictMap: Record<
  Locale,
  () => Promise<{ dict: DeepPartial<RawDictionary> }>
> = {
  // eslint-disable-next-line @typescript-eslint/no-unsafe-return, @typescript-eslint/no-explicit-any
  en: () => null as any, // en is loaded by default
  sk: () => import('../../i18n/sk/sk'),
};

export const enFlatDict: Dictionary = i18n.flatten(en);

export async function fetchDictionary(locale: Locale): Promise<Dictionary> {
  if (locale === 'en') return enFlatDict;

  const { dict } = await rawDictMap[locale]();
  const flatDict = i18n.flatten(dict) as RawDictionary;
  // the en is there to complete if some key is missing
  return { ...enFlatDict, ...flatDict };
}

export const toLocale = (string: string): Locale | undefined =>
  string in rawDictMap ? (string as Locale) : undefined;

export type Trans = i18n.Translator<Dictionary>;

export function trans(dict: InitializedResource<Dictionary>) {
  return i18n.translator(dict, i18n.resolveTemplate);
}
