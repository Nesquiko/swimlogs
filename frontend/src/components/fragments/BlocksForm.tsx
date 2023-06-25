import { Trans } from '@mbarzda/solid-i18next'
import { batch, Component, For, Show } from 'solid-js'
import { produce } from 'solid-js/store'
import { NewBlock, StartingRuleType } from '../../generated'
import { cloneBlock } from '../../lib/clone'
import { mergeInvalidBlocks, validateBlocks } from '../../lib/validation'
import { useCreateTraining } from '../context/CreateTrainingContextProvider'
import { BlockForm } from './BlockForm'

export const BlocksForm: Component = () => {
  const [
    { training, setTraining },
    { invalidTraining, setInvalidTraining },
    ,
    state,
    [, setCurrentComponent]
  ] = useCreateTraining()
  const [currentBlock, setCurrentBlock] = state.currentBlock

  const sumbit = () => {
    let isValid = true

    const invalidBlocks = validateBlocks(training.blocks)
    if (invalidBlocks !== undefined) {
      isValid = false
      setInvalidTraining(
        'blocks',
        mergeInvalidBlocks(invalidTraining.blocks, invalidBlocks)
      )
    }

    if (!isValid) {
      return
    }
    setCurrentComponent((c) => c + 1)
  }

  const shownBlockButtons = () => {
    return training.blocks
      .slice(
        // if the current block is the last one, show the last 3 blocks
        (currentBlock() + 1) % training.blocks.length === 0
          ? Math.max(0, currentBlock() - 2)
          : Math.max(0, currentBlock() - 1),
        // if the current block is the first one, show the first 3 blocks
        currentBlock() === 0 ? currentBlock() + 3 : currentBlock() + 2
      )
      .map((b) => {
        const ib = invalidTraining.blocks?.find((ib) => ib.num === b.num)
        return {
          blockNum: b.num,
          isSelected: currentBlock() === b.num,
          isValid:
            ib?.name === undefined &&
            ib?.repeat === undefined &&
            ib?.sets?.every(
              (s) =>
                s.repeat === undefined &&
                s.what === undefined &&
                s.distance === undefined &&
                s.startingRule?.type === undefined &&
                s.startingRule?.seconds === undefined
            )
        }
      })
  }

  const addNewBlock = () => {
    const newB = {
      num: training.blocks.length,
      repeat: 1,
      name: '',
      totalDistance: 0,
      sets: [
        {
          num: 0,
          repeat: 1,
          distance: 100,
          what: '',
          startingRule: {
            type: StartingRuleType.None
          }
        }
      ]
    }
    addBlock(newB)
  }

  const duplicateCurrentBlock = () => {
    const newB = cloneBlock(training.blocks[currentBlock()])
    newB.num = training.blocks.length
    addBlock(newB)
  }

  const addBlock = (b: NewBlock) => {
    batch(() => {
      setTraining(
        'blocks',
        produce((blocks) => blocks.push(b))
      )
      setCurrentBlock(training.blocks.length - 1)
      setInvalidTraining(
        'blocks',
        produce((blocks) =>
          blocks?.push({
            num: b.num,
            sets: [{ num: b.sets[0].num, startingRule: {} }]
          })
        )
      )
    })
  }

  const deleteCurrentBlock = () => {
    batch(() => {
      setTraining('blocks', (blocks) =>
        blocks
          .filter((b) => b.num !== currentBlock())
          .map((b, i) => ({ ...b, num: i }))
      )
      setCurrentBlock((i) => Math.max(0, i - 1))
      setInvalidTraining('blocks', (blocks) =>
        blocks
          ?.filter((b) => b.num !== currentBlock())
          .map((b, i) => ({ ...b, num: i }))
      )
    })
  }

  return (
    <div>
      <Show
        when={training.blocks.length !== 0}
        fallback={
          <div class="m-4 flex items-center justify-start rounded bg-blue-200 p-4 text-xl font-bold">
            <Trans key="no.blocks.in.training" />
          </div>
        }
      >
        <BlockForm
          block={training.blocks[currentBlock()]}
          onDelete={() => {
            deleteCurrentBlock()
          }}
        />
      </Show>
      <div class="fixed bottom-16 mx-auto my-4 flex w-screen justify-center">
        <Show when={training.blocks.length !== 0}>
          <button
            class="mr-8 h-12 w-12 rounded-full bg-sky-500 text-2xl text-white shadow"
            onClick={() => duplicateCurrentBlock()}
          >
            <i class="fa-regular fa-copy"></i>
          </button>
        </Show>

        <For each={shownBlockButtons()}>
          {(blockButton) => {
            return (
              <button
                classList={{
                  // selected and valid
                  'bg-sky-500 text-white border-sky-500':
                    blockButton.isSelected && blockButton.isValid,
                  // not selected and valid
                  'bg-white text-sky-500 border-sky-500':
                    !blockButton.isSelected && blockButton.isValid,
                  // selected and invalid
                  'bg-red-500 text-black border-black':
                    blockButton.isSelected && !blockButton.isValid,
                  // not selected and invalid
                  'bg-red-500 text-white border-red-500':
                    !blockButton.isSelected && !blockButton.isValid
                }}
                class="mx-2 h-12 w-12 rounded-full border-2 border-solid text-2xl font-bold shadow"
                onClick={() => setCurrentBlock(blockButton.blockNum)}
              >
                {blockButton.blockNum + 1}
              </button>
            )
          }}
        </For>
        <button
          class="ml-8 h-12 w-12 rounded-full bg-green-500 text-2xl text-white shadow"
          onClick={() => addNewBlock()}
        >
          <i class="fa-sharp fa-solid fa-plus fa-lg"></i>
        </button>

        <button
          class="fixed bottom-0 right-4 mx-auto my-4 w-1/4 rounded border bg-purple-dark py-2 text-xl font-bold text-white"
          onClick={() => sumbit()}
        >
          <Trans key="next" />
        </button>
        <button
          class="fixed bottom-0 left-4 mx-auto my-4 w-1/4 rounded border bg-purple-dark py-2 text-xl font-bold text-white"
          onClick={() => setCurrentComponent((c) => c - 1)}
        >
          <Trans key="previous" />
        </button>
      </div>
    </div>
  )
}
