import styled from 'styled-components'
import { devices } from '../../globalStyles'

export const Item = styled.button`
  margin: 0.5rem;
  padding: 0.5rem;
  border: none;
  border-radius: 0.8rem;
  max-width: 100%;
  text-wrap: wrap;

  @media ${devices.mobile} {
    background-color: #aa8866;
    color: white;
  }

  @media ${devices.desktop} {
    background-color: #eeccbb;
    color: #331111;
    &:hover {
      background-color: #aa8866;
      color: white;
      cursor: pointer;
    }
  }
`
