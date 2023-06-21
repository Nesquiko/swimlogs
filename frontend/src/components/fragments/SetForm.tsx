import { Component, For, Show } from 'solid-js'
import {
  InvalidTraining,
  InvalidTrainingSet,
  NewTrainingSet,
  StartingRuleType
} from '../../generated'
import { SmallIntMax } from '../../lib/consts'
import { useCreateTraining } from '../context/CreateTrainingContextProvider'

interface SetProps {
  set: NewTrainingSet
  blockNum: number
  deleteSet: () => void
}

export const Set: Component<SetProps> = (props) => {
  const [{ setTraining }, { invalidTraining, setInvalidTraining }, , , ,] =
    useCreateTraining()

  const getInvalidSet = (invalidTraining: InvalidTraining) => {
    return invalidTraining.blocks?.[props.blockNum]?.sets?.[props.set.num]
  }

  const isInvalid = (is: InvalidTrainingSet | undefined) => {
    return (
      is?.what !== undefined ||
      is?.repeat !== undefined ||
      is?.distance !== undefined ||
      is?.startingRule?.type !== undefined ||
      is?.startingRule?.seconds !== undefined
    )
  }

  return (
    <div
      classList={{
        'bg-sky-50': !isInvalid(
          invalidTraining.blocks?.[props.blockNum]?.sets?.[props.set.num]
        ),
        'bg-red-50': isInvalid(
          invalidTraining.blocks?.[props.blockNum]?.sets?.[props.set.num]
        )
      }}
      class="mx-auto my-2 rounded-lg border border-solid border-slate-300 p-2 shadow"
    >
      <div class="my-2 flex items-center">
        <span class="mr-4 text-lg">Set {props.set.num + 1}</span>
        <input
          type="number"
          placeholder="1"
          classList={{
            'border-red-500 text-red-500':
              getInvalidSet(invalidTraining)?.repeat !== undefined,
            'border-slate-300':
              getInvalidSet(invalidTraining)?.repeat === undefined
          }}
          class="w-1/6 rounded-md border p-2 text-center focus:border-blue-500 focus:outline-none focus:ring"
          value={props.set.repeat}
          onChange={(e) => {
            const val = e.target.value
            let repeat = parseInt(val)
            if (Number.isNaN(repeat) || repeat < 1 || repeat > SmallIntMax) {
              repeat = 0
            }
            setTraining(
              'blocks',
              (block) => block.num === props.blockNum,
              'sets',
              (set) => set.num === props.set.num,
              'repeat',
              repeat
            )
            setInvalidTraining(
              'blocks',
              (b) => b.num === props.blockNum,
              'sets',
              (s) => s.num === props.set.num,
              'repeat',
              repeat > 0
                ? undefined
                : 'Repeat must be a number between 1 and 32767'
            )
          }}
        />
        <span class="mx-4 text-lg">x</span>
        <input
          type="number"
          placeholder="400"
          classList={{
            'border-red-500 text-red-500':
              getInvalidSet(invalidTraining)?.distance !== undefined,
            'border-slate-300':
              getInvalidSet(invalidTraining)?.distance === undefined
          }}
          class="w-1/3 rounded-md border p-2 text-center focus:border-blue-500 focus:outline-none focus:ring"
          value={props.set.distance}
          onChange={(e) => {
            const val = e.target.value
            let distance = parseInt(val)
            if (
              Number.isNaN(distance) ||
              distance < 1 ||
              distance > SmallIntMax
            ) {
              distance = 0
            }
            setTraining(
              'blocks',
              (block) => block.num === props.blockNum,
              'sets',
              (set) => set.num === props.set.num,
              'distance',
              distance
            )
            setInvalidTraining(
              'blocks',
              (b) => b.num === props.blockNum,
              'sets',
              (s) => s.num === props.set.num,
              'distance',
              distance > 0
                ? undefined
                : 'Repeat must be a number between 1 and 32767'
            )
          }}
        />
        <i
          class="fa-solid fa-trash fa-xl ml-auto text-red-500"
          onClick={() => props.deleteSet()}
        ></i>
      </div>
      <textarea
        placeholder="Freestyle"
        maxlength="255"
        classList={{
          'border-red-500 text-red-500':
            getInvalidSet(invalidTraining)?.what !== undefined,
          'border-slate-300 text-black':
            getInvalidSet(invalidTraining)?.what === undefined
        }}
        class="w-full rounded-md border p-2 focus:border-sky-500 focus:outline-none focus:ring"
        value={props.set.what}
        onChange={(e) => {
          const what = e.target.value
          setTraining(
            'blocks',
            (block) => block.num === props.blockNum,
            'sets',
            (set) => set.num === props.set.num,
            'what',
            what.trim()
          )
          setInvalidTraining(
            'blocks',
            (b) => b.num === props.blockNum,
            'sets',
            (s) => s.num === props.set.num,
            'what',
            what.length > 0 ? undefined : "What can't be empty"
          )
        }}
      />
      <div class="flex items-center space-x-4">
        <label for="start">Starting</label>
        <select
          id="start"
          class="my-2 rounded-md border border-solid border-slate-300 bg-white p-2 focus:border-sky-500 focus:outline-none focus:ring"
          onChange={(e) => {
            const typ = e.target.value as StartingRuleType
            const startingRule = {
              type: typ,
              seconds: props.set.startingRule.seconds ?? 20
            }
            setTraining(
              'blocks',
              (block) => block.num === props.blockNum,
              'sets',
              (set) => set.num === props.set.num,
              'startingRule',
              startingRule
            )
          }}
        >
          <For each={Object.keys(StartingRuleType)}>
            {(typ) => (
              <option
                selected={typ === props.set.startingRule.type}
                value={typ}
              >
                {typ}
              </option>
            )}
          </For>
        </select>
        <Show when={props.set.startingRule.type !== StartingRuleType.None}>
          <label for="seconds">Seconds</label>
          <input
            id="seconds"
            type="number"
            placeholder="20"
            classList={{
              'border-red-500 text-red-500':
                getInvalidSet(invalidTraining)?.startingRule?.seconds !==
                undefined,
              'border-slate-300':
                getInvalidSet(invalidTraining)?.startingRule?.seconds ===
                undefined
            }}
            class="w-1/4 rounded-md border border-solid px-4 py-2 text-center focus:border-sky-500 focus:outline-none focus:ring"
            value={props.set.startingRule.seconds}
            onChange={(e) => {
              const val = e.target.value
              let seconds = parseInt(val)
              if (
                Number.isNaN(seconds) ||
                seconds < 1 ||
                seconds > SmallIntMax
              ) {
                seconds = 0
              }
              const startingRule = {
                type: props.set.startingRule.type,
                seconds: seconds
              }
              setTraining(
                'blocks',
                (block) => block.num === props.blockNum,
                'sets',
                (set) => set.num === props.set.num,
                'startingRule',
                startingRule
              )
              setInvalidTraining(
                'blocks',
                (b) => b.num === props.blockNum,
                'sets',
                (s) => s.num === props.set.num,
                'startingRule',
                'seconds',
                seconds > 0
                  ? undefined
                  : 'Seconds must be a number between 1 and 32767'
              )
            }}
          />
        </Show>
      </div>
    </div>
  )
}
