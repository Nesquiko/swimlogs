import { type Component } from 'solid-js';
import { useAppState } from '~/AppContext';
import { Separator } from '~/components/ui/separator';
import AppearanceSection from './AppearanceSection';

const ProfilePage: Component = () => {
  const { t } = useAppState();
  return (
    <div class="grid gap-4">
      <p class="text-3xl tracking-tight">{t('profile.headline')}</p>
      <Separator />
      <AppearanceSection />
    </div>
  );
};

export default ProfilePage;
