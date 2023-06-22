import { NewBlock, NewTrainingSet } from '../generated'

export function cloneBlock(b: NewBlock) {
  const newBlock = {} as NewBlock
  newBlock.num = b.num
  newBlock.name = b.name
  newBlock.repeat = b.repeat
  newBlock.repeat = b.repeat
  newBlock.sets = new Array(b.sets.length)

  b.sets.forEach((s, i) => {
    const newSet = {} as NewTrainingSet
    newSet.repeat = s.repeat
    newSet.num = s.num
    newSet.what = s.what
    newSet.distance = s.distance
    newSet.startingRule = Object.assign({}, s.startingRule)

    newBlock.sets[i] = newSet
  })

  return newBlock
}
