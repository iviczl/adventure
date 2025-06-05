import { useState } from 'react'
import { Button, Label } from '../../globalStyles'
import { save } from '../../services/gameService'
import { AppState, state } from '../../state'
import { AttributeList, ToolHeadBar } from './gameStyles'
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

  const attributeList = () => {
    return state.value.playerAttributes.map((attr) => (
      <>
        <div>
          <Label>{Object.keys(attr)[0]}:</Label>
        </div>
        <div>
          <Label>{attr[Object.keys(attr)[0]]}</Label>
        </div>
      </>
    ))
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
      <AttributeList>{attributeList()}</AttributeList>
      <Button
        onClick={() => setShowSaveModal(true)}
        disabled={requestProcessing}
      >
        Quit
      </Button>
    </ToolHeadBar>
  )
}
