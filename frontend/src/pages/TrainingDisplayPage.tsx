import { Trans, useTransContext } from '@mbarzda/solid-i18next';
import {
  createAsync,
  RouteSectionProps,
  useNavigate,
  useParams,
} from '@solidjs/router';
import { Component, createSignal, onCleanup, onMount, Show } from 'solid-js';
import { ResponseError, Training } from 'swimlogs-api';
import { showToast } from '../App';
import DismissibleModal from '../components/DismissibleModal';
import { ToastMode } from '../components/DismissibleToast';
import { clearHeaderMenu, setHeaderMenu } from '../components/Header';
import SetsPreview, { SkeletonSetsPreview } from '../components/SetsPreview';
import TrainingSummary from '../components/TraningSummary';
import { deleteTrainingById } from '../state/trainings';

interface TrainingDisplayPageProps {
  trainingPromise: Promise<Training>;
}

const TrainingDisplayPage: Component<
  RouteSectionProps<TrainingDisplayPageProps>
> = (props) => {
  const [t] = useTransContext();
  const navigate = useNavigate();
  const params = useParams();
  const training = createAsync(() =>
    props.data!.trainingPromise.catch((e: ResponseError) => {
      let msg = t('server.error');
      if (e.response?.status === 404) {
        msg = t('training.not.found');
      }
      showToast(msg, ToastMode.ERROR);
      navigate('/', { replace: true });
      return Promise.reject(e);
    })
  );
  const [openConfirmationModal, setOpenConfirmationModal] = createSignal(false);

  onMount(() => {
    setHeaderMenu({
      items: [
        {
          textKey: 'edit',
          icon: 'fa-pen',
          onClick: () => navigate(`/training/${params.id}/edit`),
        },
      ],
      lastItem: {
        textKey: 'delete',
        icon: 'fa-trash',
        textColorCls: 'text-red-500',
        onClick: () => setOpenConfirmationModal(true),
      },
    });
  });
  onCleanup(clearHeaderMenu);

  async function deleteTraining(id: string) {
    await deleteTrainingById(id)
      .then(() => showToast(t('training.deleted')))
      .catch((e: ResponseError) => {
        let msg = t('server.error');
        if (e.response?.status === 404) {
          msg = t('training.not.found');
        }
        showToast(msg, ToastMode.ERROR);
      })
      .finally(() => {
        navigate('/', { replace: true });
      });
  }

  return (
    <div class="px-4">
      <Show when={training()} fallback={<SkeletonSetsPreview />} keyed>
        {(training) => (
          <>
            <DismissibleModal
              open={openConfirmationModal()}
              setOpen={setOpenConfirmationModal}
              message={t('confirm.training.delete.message')}
              confirmLabel={t('confirm.delete.training')}
              cancelLabel={t('no.cancel')}
              onConfirm={() => deleteTraining(params.id)}
            />
            <TrainingSummary training={training} />
            <h1 class="py-2 text-2xl font-bold text-sky-900">
              <Trans key="sets" />
            </h1>
            <SetsPreview sets={training.sets} />
          </>
        )}
      </Show>
    </div>
  );
};

export default TrainingDisplayPage;
