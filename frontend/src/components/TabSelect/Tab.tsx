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
    return <SelectedTab>{title}</SelectedTab>
  }
  return <UnselectedTab onClick={() => onClick(id)}>{title}</UnselectedTab>
}
