import { useState } from 'react'
import { Button } from '../../globalStyles'
import { save } from '../../services/gameService'
import { AppState, state } from '../../state'
import { ToolHeadBar } from './gameStyles'
import { createPortal } from 'react-dom'
import SaveModal from './SaveModal'

export default function ToolBar() {
  const [requestProcessing, setRequestProcessing] = useState(false)
  const [showSaveModal, setShowSaveModal] = useState(false)
  const quit = async (needSaving: boolean) => {
    setRequestProcessing(true)
    let success = true
    if (needSaving) {
      success = await save()
    }
    if (success) {
      state.value = {
        ...state.value,
        player: '',
        selectedGameId: '',
        adventureId: '',
        actualPosition: undefined,
      } as AppState
    }
    setRequestProcessing(false)
    setShowSaveModal(false)
  }

  const cancel = () => {
    setShowSaveModal(false)
  }
  return (
    <ToolHeadBar>
      {showSaveModal &&
        createPortal(
          <SaveModal onQuit={quit} onCancel={cancel} />,
          document.getElementById('root')!
        )}
      <Button
        onClick={() => setShowSaveModal(true)}
        disabled={requestProcessing}
      >
        Quit
      </Button>
    </ToolHeadBar>
  )
}
