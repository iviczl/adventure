import { useState } from 'react'
import { Button } from '../../globalStyles'
import { save } from '../../services/gameService'
import { AppState, state } from '../../state'
import { ToolHeadBar } from './gameStyles'

export default function ToolBar() {
  const [requestProcessing, setRequestProcessing] = useState(false)
  const quit = async () => {
    setRequestProcessing(true)
    const success = await save()
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
  }

  return (
    <ToolHeadBar>
      <Button onClick={quit} disabled={requestProcessing}>
        Quit
      </Button>
    </ToolHeadBar>
  )
}
