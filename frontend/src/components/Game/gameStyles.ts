import styled from 'styled-components'
import { devices } from '../../globalStyles'

export const Item = styled.option`
  margin: 0.5rem;
  padding: 0.5rem;
  border-radius: 0.8rem;

  @media ${devices.mobile} {
    background-color: #aa8866;
    color: white;
  }

  @media ${devices.desktop} {
    &:hover {
      background-color: #aa8866;
      color: white;
      cursor: pointer;
    }
  }
`
