import { type IconProps as SolidIconProps } from 'solid-icons';
import { BiRegularSun, BiSolidHome, BiRegularCalendar } from 'solid-icons/bi';
import {
  IoMoonOutline,
  IoLanguageOutline,
  IoContrastOutline,
} from 'solid-icons/io';
import { Component } from 'solid-js';
import { FaRegularUser } from 'solid-icons/fa';

export type IconProps = SolidIconProps;

export const IconSun: Component<IconProps> = (props) => (
  <BiRegularSun {...props} />
);

export const IconMoon: Component<IconProps> = (props) => (
  <IoMoonOutline {...props} />
);

export const IconHome: Component<IconProps> = (props) => (
  <BiSolidHome {...props} />
);

export const IconCalendar: Component<IconProps> = (props) => (
  <BiRegularCalendar {...props} />
);

export const IconUser: Component<IconProps> = (props) => (
  <FaRegularUser {...props} />
);

export const IconLanguage: Component<IconProps> = (props) => (
  <IoLanguageOutline {...props} />
);

export const IconColorTheme: Component<IconProps> = (props) => (
  <IoContrastOutline {...props} />
);
