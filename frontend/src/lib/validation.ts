import {
  InvalidBlock,
  InvalidTraining,
  InvalidTrainingSet,
  NewBlock,
  NewTraining,
  NewTrainingSet,
  StartingRule,
  StartingRuleType
} from '../generated'

export const mergeInvalidBlocks = (
  invalidBlocks1: InvalidBlock[] | undefined,
  invalidBlocks2: InvalidBlock[] | undefined
): InvalidBlock[] => {
  if (invalidBlocks1 === undefined) return invalidBlocks2 || []
  if (invalidBlocks2 === undefined) return invalidBlocks1

  const merged: InvalidBlock[] = []

  invalidBlocks1.forEach((ib) => {
    const other = invalidBlocks2.find((ib2) => ib2.num === ib.num)
    if (other === undefined) {
      merged.push(ib)
      return
    }
    const mergedBlock = { ...ib, ...other }
    mergedBlock.sets = []

    ib.sets?.forEach((is) => {
      const otherSet = other.sets?.find((is2) => is2.num === is.num)
      if (otherSet === undefined) {
        mergedBlock.sets?.push(is)
        return
      }
      const newSet = { ...is, ...otherSet }
      newSet.startingRule = { ...is.startingRule, ...otherSet.startingRule }
      mergedBlock.sets?.push(newSet)
    })

    merged.push(mergedBlock)
  })

  return merged
}

export const validateTraining = (
  t: NewTraining
): InvalidTraining | undefined => {
  const invalidTraining: InvalidTraining = {}

  if (t.durationMin === undefined)
    invalidTraining.durationMin = "Duration can't be empty"
  else if (t.durationMin < 0)
    invalidTraining.durationMin = 'Duration must be positive'

  const invalidBlocks = validateBlocks(t.blocks)
  if (invalidBlocks) invalidTraining.blocks = invalidBlocks

  if (Object.keys(invalidTraining).length === 0) return undefined

  return invalidTraining
}

export function validateBlocks(
  blocks: Array<NewBlock>
): Array<InvalidBlock> | undefined {
  const invalidBlocks: Array<InvalidBlock> = []

  blocks.forEach((block, index) => {
    invalidBlocks[index] = {}

    if (block.name === '') invalidBlocks[index].name = "Name can't be empty"
    else if (block.name.length > 255)
      invalidBlocks[index].name = "Name can't be longer than 255 characters"

    if (block.repeat <= 0)
      invalidBlocks[index].repeat = 'Repeat must be positive'

    const invalidSets = validateSets(block.sets)
    if (invalidSets) invalidBlocks[index].sets = invalidSets

    if (Object.keys(invalidBlocks[index]).length === 0) {
      delete invalidBlocks[index]
    } else {
      invalidBlocks[index].num = block.num
    }
  })

  if (Object.keys(invalidBlocks).length === 0) return undefined

  return invalidBlocks.filter((block) => block !== undefined)
}

function validateSets(
  sets: Array<NewTrainingSet>
): Array<InvalidTrainingSet> | undefined {
  const invalidSets: Array<InvalidTrainingSet> = []

  sets.forEach((set, index) => {
    invalidSets[index] = {}
    if (set.repeat <= 0) invalidSets[index].repeat = 'Repeat must be positive'

    if (set.what === '') invalidSets[index].what = "What can't be empty"

    if (set.distance <= 0)
      invalidSets[index].distance = 'Distance must be positive'

    if (!isStartingRuleValid(set.startingRule)) {
      invalidSets[index].startingRule = { seconds: 'Seconds must be positive' }
    }

    if (Object.keys(invalidSets[index]).length === 0) {
      delete invalidSets[index]
    } else {
      invalidSets[index].num = set.num
    }
  })

  if (Object.keys(invalidSets).length === 0) return undefined

  return invalidSets.filter((set) => set !== undefined)
}

function isStartingRuleValid(rule: StartingRule): boolean {
  if (rule.type === StartingRuleType.None) {
    return true
  }

  // else it is Interval or Pause
  if (rule.seconds === undefined) {
    return false
  }

  return rule.seconds > 0
}
