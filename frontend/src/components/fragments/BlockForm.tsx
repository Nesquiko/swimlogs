import { Trans, useTransContext } from '@mbarzda/solid-i18next'
import { batch, Component, createEffect, For } from 'solid-js'
import { produce } from 'solid-js/store'
import { NewBlock, NewTrainingSet, StartingRuleType } from '../../generated'
import { cloneSet } from '../../lib/clone'
import { SmallIntMax } from '../../lib/consts'
import { useCreateTraining } from '../context/CreateTrainingContextProvider'
import { Set } from './SetForm'

interface BlockFormProps {
  block: NewBlock
  onDelete: () => void
}

export const BlockForm: Component<BlockFormProps> = (props) => {
  const [{ setTraining }, { invalidTraining, setInvalidTraining }] =
    useCreateTraining()
  const [t] = useTransContext()

  createEffect(() => {
    const dist = props.block.sets.reduce(
      (acc, set) => acc + set.distance * set.repeat,
      0
    )

    setTraining(
      'blocks',
      (block) => block.num === props.block.num,
      'totalDistance',
      props.block.repeat * dist
    )
  })

  const addNewSet = () => {
    const newSet = {
      num: props.block.sets.length,
      repeat: 1,
      distance: 100,
      what: '',
      startingRule: {
        type: StartingRuleType.None
      }
    } as NewTrainingSet
    addSet(newSet)
  }

  const duplicateSet = (setNum: number) => {
    const newSet = cloneSet(props.block.sets[setNum])
    newSet.num = props.block.sets.length
    addSet(newSet)
  }

  const addSet = (s: NewTrainingSet) => {
    batch(() => {
      setTraining(
        'blocks',
        (block) => block.num === props.block.num,
        'sets',
        produce((sets) => sets.push(s))
      )
      setInvalidTraining(
        'blocks',
        (b) => b.num === props.block.num,
        'sets',
        produce((sets) => sets?.push({ num: s.num, startingRule: {} }))
      )
    })
  }

  const deleteSet = (setIndex: number) => {
    batch(() => {
      setTraining(
        'blocks',
        (block) => block.num === props.block.num,
        'sets',
        (sets) =>
          sets
            .filter((s) => s.num !== setIndex)
            .map((s, i) => ({ ...s, num: i }))
      )
      setInvalidTraining(
        'blocks',
        (b) => b.num === props.block.num,
        'sets',
        (sets) =>
          sets
            ?.filter((s) => s.num !== setIndex)
            .map((s, i) => ({ ...s, num: i }))
      )
    })
  }

  return (
    <div class="m-4">
      <div class="my-2 flex items-center justify-between">
        <label class="w-2/5 text-2xl" for="name">
          <Trans key="block" /> {props.block.num + 1}
        </label>
        <input
          id="name"
          placeholder={t('block.name.placeholder', 'Warm up')}
          maxlength="255"
          classList={{
            'border-red-500 text-red-500':
              invalidTraining.blocks?.[props.block.num]?.name !== undefined,
            'border-slate-300 text-black':
              invalidTraining.blocks?.[props.block.num]?.name === undefined
          }}
          class=" w-1/2 rounded-md border p-2 text-lg focus:border-sky-500 focus:outline-none focus:ring"
          value={props.block.name}
          onChange={(e) => {
            const val = e.target.value
            setTraining(
              'blocks',
              (block) => block.num === props.block.num,
              'name',
              val.trim()
            )
            setInvalidTraining(
              'blocks',
              (block) => block.num === props.block.num,
              'name',
              undefined
            )
          }}
        />
      </div>
      <div class="my-2 flex items-center justify-between">
        <label class="w-2/5 text-2xl" for="repeat">
          <Trans key="repeat" />
        </label>
        <input
          id="repeat"
          type="number"
          placeholder="1"
          classList={{
            'border-slate-300':
              invalidTraining.blocks?.[props.block.num].repeat === undefined,
            'border-red-500 text-red-500':
              invalidTraining.blocks?.[props.block.num].repeat !== undefined
          }}
          class="w-1/2 rounded-md border p-2 text-lg focus:border-sky-500 focus:outline-none focus:ring"
          value={props.block.repeat}
          onChange={(e) => {
            const val = e.target.value
            setInvalidTraining(
              'blocks',
              (block) => block.num === props.block.num,
              'repeat',
              undefined
            )
            const repeat = parseInt(val)
            if (Number.isNaN(repeat) || repeat < 1 || repeat > SmallIntMax) {
              setInvalidTraining(
                'blocks',
                (block) => block.num === props.block.num,
                'repeat',
                'Repeat must be a number between 1 and 32767'
              )
              return
            }
            setTraining(
              'blocks',
              (block) => block.num === props.block.num,
              'repeat',
              repeat
            )
          }}
        />
      </div>
      <div class="flex justify-between">
        <span class="my-2 text-2xl">
          <Trans key="total" />: <b>{props.block.totalDistance}m</b>
        </span>
        <button
          class="h-12 w-12 rounded-full bg-red-500 shadow"
          onClick={() => props.onDelete()}
        >
          <img src="/src/assets/bin.svg" width={48} height={48} />
        </button>
      </div>

      <For each={props.block.sets}>
        {(set) => {
          return (
            <Set
              set={set}
              blockNum={props.block.num}
              deleteSet={() => deleteSet(set.num)}
              duplicateSet={() => duplicateSet(set.num)}
            />
          )
        }}
      </For>
      <button
        class="float-right my-4 rounded bg-sky-500 p-2 font-bold text-white"
        onClick={() => addNewSet()}
      >
        <Trans key="add.set" />
      </button>
      {/* Add space at the bottom, so the buttons dont hide block form */}
      <div class="h-48 w-full"></div>
    </div>
  )
}
