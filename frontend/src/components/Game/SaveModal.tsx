import { useEffect, useRef } from 'react'
import { Button, Dialog, Row } from '../../globalStyles'

export default function SaveModal({
  onQuit,
  onCancel,
}: {
  onQuit: (save: boolean) => void
  onCancel: () => void
}) {
  const dialogElement = useRef(null)

  useEffect(() => {
    const element = dialogElement.current as HTMLDialogElement | null
    if (element) {
      element.showModal()
    }
  }, [])

  return (
    <Dialog ref={dialogElement}>
      <Row $marginBottom='2rem'>
        You are about to quit the game. Do you want to save it for later?
      </Row>
      <Row $justifyContent='center'>
        <Button onClick={() => onQuit(true)}>Save Game</Button>
        <Button onClick={() => onQuit(false)}>Do not save</Button>
        <Button onClick={onCancel}>Cancel</Button>
      </Row>
    </Dialog>
  )
}
