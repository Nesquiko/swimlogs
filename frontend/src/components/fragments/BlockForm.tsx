import { batch, Component, createEffect, For } from 'solid-js'
import { produce } from 'solid-js/store'
import { NewBlock, NewTrainingSet, StartingRuleType } from '../../generated'
import { SmallIntMax } from '../../lib/consts'
import { useCreateTraining } from '../context/CreateTrainingContextProvider'
import { Set } from './SetForm'

interface BlockFormProps {
  block: NewBlock
}

export const BlockForm: Component<BlockFormProps> = (props) => {
  const [{ setTraining }, { invalidTraining, setInvalidTraining }] =
    useCreateTraining()

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

  const addSet = () => {
    const newSet = {
      num: props.block.sets.length,
      repeat: 1,
      distance: 100,
      what: '',
      startingRule: {
        type: StartingRuleType.None
      }
    } as NewTrainingSet

    batch(() => {
      setTraining(
        'blocks',
        (block) => block.num === props.block.num,
        'sets',
        produce((sets) => sets.push(newSet))
      )
      setInvalidTraining(
        'blocks',
        (b) => b.num === props.block.num,
        'sets',
        produce((sets) => sets?.push({ num: newSet.num, startingRule: {} }))
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
      <div class="my-2 flex items-center space-x-4">
        <label class="w-1/5 text-2xl" for="name">
          Block {props.block.num + 1}
        </label>
        <input
          id="name"
          placeholder="Warm up"
          maxlength="255"
          classList={{
            'border-red-500 text-red-500':
              invalidTraining.blocks?.[props.block.num]?.name !== undefined,
            'border-slate-300 text-black':
              invalidTraining.blocks?.[props.block.num]?.name === undefined
          }}
          class=" rounded-md border p-2 text-lg focus:border-sky-500 focus:outline-none focus:ring"
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
      <div class="my-2 flex items-center space-x-4">
        <label class="w-1/5 text-2xl" for="repeat">
          Repeat
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
          class="rounded-md border p-2 text-lg focus:border-sky-500 focus:outline-none focus:ring"
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
      <span class="text-2xl">
        Distance in this block <b>{props.block.totalDistance}m</b>
      </span>

      <h1 class="my-4 text-2xl">Sets</h1>
      <For each={props.block.sets}>
        {(set) => {
          return (
            <Set
              set={set}
              blockNum={props.block.num}
              deleteSet={() => deleteSet(set.num)}
            />
          )
        }}
      </For>
      <button
        class="float-right rounded bg-sky-500 p-2 font-bold text-white"
        onClick={() => addSet()}
      >
        Add Set
      </button>
      {/* Add space at the bottom, so the buttons dont hide block form */}
      <div class="h-48 w-full"></div>
    </div>
  )
}
