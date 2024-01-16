export function isUndefindOrEmpty(str: string | undefined): boolean {
  return str === undefined || str === null || str === ''
}

export function randomId(length: number = 7) {
  let result = ''
  const characters =
    'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789'
  const charactersLength = characters.length
  let counter = 0
  while (counter < length) {
    result += characters.charAt(Math.floor(Math.random() * charactersLength))
    counter += 1
  }
  return result
}

export function capitalize(str: string) {
  return str.charAt(0).toUpperCase() + str.slice(1)
}
