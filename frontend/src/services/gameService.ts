import { Position } from '../types/position'
import { state } from '../state'
import { GameInfo } from '../types/gameInfo'
import { LoginResponse } from '../types/loginResponse'

let abortController: AbortController
const portExpression = import.meta.env.VITE_SERVICE_PORT
  ? `:${import.meta.env.VITE_SERVICE_PORT}`
  : ''
const apiBasePath = import.meta.env.VITE_SERVICE_HOST + portExpression

export async function login(userName: string, password: string) {
  abortController = new AbortController()
  const result = await doFetch(`${apiBasePath}/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ userName, password }),
  })
  if (assertError(result)) {
    return
  }
  const response = (await result.response) as LoginResponse
  state.value = {
    ...state.value,
    userName,
    token: response.token,
  }
}

export async function getGames() {
  abortController = new AbortController()
  const result = await doFetch(`${apiBasePath}/games`)

  if (assertError(result)) {
    return
  }
  const list = (await result.response) as GameInfo[]
  if (list) {
    state.value = {
      ...state.value,
      games: list,
      selectedGameId: '',
    }
  }
}

export async function startGame(gameId: string, player: string) {
  abortController = new AbortController()
  const result = await doFetch(`${apiBasePath}/new`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ gameId, player }),
  })

  if (assertError(result)) {
    return
  }
  // {gameId: string, position: Position}
  const position = (await result.response) as Position
  state.value = {
    ...state.value,
    selectedGameId: gameId,
    player: player,
    adventureId: position.adventureId,
    actualPosition: position,
  }
}

export async function takeAction(
  actionId: string,
  player: string,
  adventureId: string
) {
  abortController = new AbortController()
  const result = await doFetch(`${apiBasePath}/do`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ actionCode: actionId, player, adventureId }),
  })

  if (assertError(result)) {
    return
  }
  const position = (await result.response) as Position
  state.value = {
    ...state.value,
    actualPosition: position,
  }
}

function assertError(result: { response: Promise<Position>; error: unknown }) {
  if (result.error) {
    console.log(result.error)
    return true
  }
  return false
}

export async function doFetch(url: string, options = {} as RequestInit) {
  let response = null
  let error = null
  try {
    if (state.value.token) {
      options.headers = {
        ...options.headers,
        Authorization: state.value.token,
      }
    }
    const signal = abortController.signal
    const res = await fetch(url, { ...options, credentials: 'include', signal })
    response = await res.json()
  } catch (problem) {
    abortController.abort()
    error = problem
  }
  return { response, error }
}
