import { createContext, ParentComponent, Signal, useContext } from 'solid-js'

const makeShownComponentContext = (currentComponentSignal: Signal<number>) => {
  return [currentComponentSignal[0], currentComponentSignal[1]] as const
}
type ShownComponentContextType = ReturnType<typeof makeShownComponentContext>

const ShownComponentContext = createContext<ShownComponentContextType>()

export const useShownComponent = () => useContext(ShownComponentContext)!

interface ShownComponentContextProps {
  currentComponentSignal: Signal<number>
}

export const ShownComponentContextProvider: ParentComponent<
  ShownComponentContextProps
> = (props) => {
  return (
    <ShownComponentContext.Provider
      value={makeShownComponentContext(props.currentComponentSignal)}
    >
      {props.children}
    </ShownComponentContext.Provider>
  )
}
