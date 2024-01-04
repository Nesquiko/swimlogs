import finsBlackSvg from '../assets/fins-black.svg'
import boardBlackSvg from '../assets/board-black.svg'
import paddlesBlackSvg from '../assets/paddles-black.svg'
import snorkelBlackSvg from '../assets/snorkel-black.svg'
import monofinBlackSvg from '../assets/monofin-black.svg'
import { Equipment } from '../generated'

export const EquipmentIcons = new Map<Equipment, string>([
  [Equipment.Fins, finsBlackSvg],
  [Equipment.Board, boardBlackSvg],
  [Equipment.Paddles, paddlesBlackSvg],
  [Equipment.Snorkel, snorkelBlackSvg],
  [Equipment.Monofin, monofinBlackSvg],
])
