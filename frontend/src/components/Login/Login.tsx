import { Signal } from '@preact/signals-react'
import { AppState } from '../../state'
import {
  Button,
  Container,
  Input,
  Label,
  MainHeading,
  Paragraph,
  Row,
} from '../../globalStyles'
import { ColumnContainer } from './loginStyles'
import { useState } from 'react'
import { login } from '../../services/gameService'

export default function Login({ state }: { state: Signal<AppState> }) {
  const [requestProcessing, setRequestProcessing] = useState(false)
  const [userName, setUserName] = useState(state.value.userName)
  const [password, setPassword] = useState('')
  const [loginFailed, setLoginFailed] = useState(false)

  const doLogin = async () => {
    setRequestProcessing(true)
    await login(userName, password)
    setRequestProcessing(false)
    if (state.value.userName === '' || state.value.token === undefined) {
      setLoginFailed(true)
    }
  }

  if (state.value.userName !== '') {
    return null
  }

  return (
    <Container>
      <MainHeading>Adventure games</MainHeading>
      <Paragraph>Log in to play.</Paragraph>
      <ColumnContainer>
        <Row>
          <Label $width='10rem'>User name:</Label>
          <Input
            $width='12rem'
            type='text'
            onChange={(e) => {
              setLoginFailed(false)
              setUserName(e.target.value)
            }}
            value={userName}
            autoFocus
          />
        </Row>
        <Row>
          <Label $width='10rem'>Password:</Label>
          <Input
            $width='12rem'
            type='text'
            onChange={(e) => {
              setLoginFailed(false)
              setPassword(e.target.value)
            }}
            value={password}
          />
        </Row>
      </ColumnContainer>
      <Row>
        <Button
          onClick={doLogin}
          disabled={!userName || !password || requestProcessing}
        >
          {requestProcessing ? 'Processing...' : 'Log in'}
        </Button>
      </Row>
      {loginFailed ? (
        <Row>
          <Label $color='red'>
            The login failed. Please, check if you gave the right credentials.
          </Label>
        </Row>
      ) : null}
    </Container>
  )
}
