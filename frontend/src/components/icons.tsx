import { type IconProps as SolidIconProps } from 'solid-icons';
import { HiOutlineXMark } from 'solid-icons/hi';
import {
  BiRegularSun,
  BiSolidHome,
  BiRegularCalendar,
  BiRegularStopwatch,
} from 'solid-icons/bi';
import {
  IoMoonOutline,
  IoLanguageOutline,
  IoContrastOutline,
} from 'solid-icons/io';
import { Component } from 'solid-js';
import {
  FaRegularUser,
  FaSolidBug,
  FaRegularClock,
  FaSolidRepeat,
  FaSolidPlus,
  FaSolidArrowLeft,
  FaSolidMinus,
} from 'solid-icons/fa';

export type IconProps = SolidIconProps;

export const IconSun: Component<IconProps> = BiRegularSun;
export const IconMoon: Component<IconProps> = IoMoonOutline;
export const IconHome: Component<IconProps> = BiSolidHome;
export const IconCalendar: Component<IconProps> = BiRegularCalendar;
export const IconUser: Component<IconProps> = FaRegularUser;
export const IconLanguage: Component<IconProps> = IoLanguageOutline;
export const IconColorTheme: Component<IconProps> = IoContrastOutline;
export const IconBug: Component<IconProps> = FaSolidBug;
export const IconStopwatch: Component<IconProps> = BiRegularStopwatch;
export const IconClock: Component<IconProps> = FaRegularClock;
export const IconRepeat: Component<IconProps> = FaSolidRepeat;
export const IconX: Component<IconProps> = HiOutlineXMark;
export const IconPlus: Component<IconProps> = FaSolidPlus;
export const IconMinus: Component<IconProps> = FaSolidMinus;
export const IconArrowLeft: Component<IconProps> = FaSolidArrowLeft;
