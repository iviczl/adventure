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
  if (id === selectedId) {
    return <SelectedTab tabIndex={0}>{title}</SelectedTab>
  }
  return (
    <UnselectedTab tabIndex={0} onClick={() => onClick(id)}>
      {title}
    </UnselectedTab>
  )
}
