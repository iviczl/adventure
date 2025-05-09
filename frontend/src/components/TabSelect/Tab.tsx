import { SelectedTab, UnselectedTab } from './tabSelectStyles'

export default function Tab({
  id,
  selectedId,
  title,
  onClick,
}: {
  id: string
  selectedId: string
  title: string
  onClick: (id: string) => void
}) {
  const keyUpHandler = (e: React.KeyboardEvent<HTMLDivElement>): void => {
    if (e.key === 'Enter') {
      onClick(id)
    }
  }

  if (id === selectedId) {
    return (
      <SelectedTab tabIndex={0} onKeyUp={keyUpHandler}>
        {title}
      </SelectedTab>
    )
  }
  return (
    <UnselectedTab
      tabIndex={0}
      onClick={() => onClick(id)}
      onKeyUp={keyUpHandler}
    >
      {title}
    </UnselectedTab>
  )
}
