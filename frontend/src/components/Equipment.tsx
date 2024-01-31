import finsBlackSvg from '../assets/fins-black.svg'
import boardBlackSvg from '../assets/board-black.svg'
import paddlesBlackSvg from '../assets/paddles-black.svg'
import snorkelBlackSvg from '../assets/snorkel-black.svg'
import monofinBlackSvg from '../assets/monofin-black.svg'
import { EquipmentEnum } from 'swimlogs-api'

export const EquipmentIcons = new Map<Equipment, string>([
  [EquipmentEnum.Fins, finsBlackSvg],
  [EquipmentEnum.Board, boardBlackSvg],
  [EquipmentEnum.Paddles, paddlesBlackSvg],
  [EquipmentEnum.Snorkel, snorkelBlackSvg],
  [EquipmentEnum.Monofin, monofinBlackSvg],
])
