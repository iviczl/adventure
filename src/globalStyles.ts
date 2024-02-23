import styled, { createGlobalStyle } from 'styled-components'

const GlobalStyle = createGlobalStyle`
font-family: Inter, system-ui, Avenir, Helvetica, Arial, sans-serif;
line-height: 1.5;
font-weight: 400;

color: #331111;
background-color: #242424;

font-synthesis: none;
text-rendering: optimizeLegibility;
-webkit-font-smoothing: antialiased;
-moz-osx-font-smoothing: grayscale;
padding: auto;
`
export default GlobalStyle

export const MainHeading = styled.h1`
  font-size: clamp(3rem, 6vw, 5rem);
  margin-bottom: 1rem;
  margin-top: 2rem;
`
export const MainContainer = styled.div`
  margin: 0;
  background-color: #cc9977;
  display: flex;
  place-items: center;
  flex-flow: column;
  min-width: 320px;
  min-height: 90vh;
  border: solid 1px #331111;
  border-radius: clamp(5px, 20px, 25px);
`
export const Container = styled.div`
  margin: 10vh 10vh;
  padding: 3rem;
  display: flex;
  flex-flow: column;
  place-items: center;
  min-width: 320px;
  min-height: 10vh;
  border: solid 1px #331111;
  border-radius: clamp(5px, 20px, 25px);
  background-color: #cc9977;
`
export const Button = styled.button`
  border-radius: 8px;
  border: 1px solid transparent;
  padding: 0.6em 1.2em;
  font-size: 1em;
  font-weight: 500;
  font-family: inherit;
  background-color: #eeccbb;
  color: #331111;
  border: solid 1px #331111;

  cursor: pointer;
  transition: border-color 0.25s;
  &:disabled {
    background-color: #ffddcc;
    color: #bbcccc;
    border: solid 1px #ffddcc;
    cursor: not-allowed;
  }

  &:not(:disabled):hover {
    border-color: white;
  }

  &:focus,
  &:focus-visible {
    outline: 4px auto -webkit-focus-ring-color;
  }
}
`
export const Paragraph = styled.p`
  margin: 1rem 0 0.5rem 0;
`

export const Label = styled.label`
  font-size: 1rem;
  font-weight: 500;
  font-family: inherit;
  line-height: 0.8rem; //calc(2.2rem + 2px);
  padding: 0.5rem 0;
`
export const Select = styled.select`
  border-radius: 8px;
  margin-left: ${(props) => props.$marginLeft || '1rem'};
  padding: 0.6rem .6rem .6rem .6rem;
  font-size: 1rem;
  font-weight: 500;
  font-family: inherit;
  background-color: white;
  color: #331111;
  border: solid 1px #331111;
  max-width: 10rem;

  cursor: pointer;
  transition: border-color 0.25s;
  &:disabled {
    color: #bbcccc;
    border: solid 1px #ffddcc;
    cursor: not-allowed;
  }

  &:focus,
  &:focus-visible {
    outline: 4px auto -webkit-focus-ring-color;
  }
}
`
export const Input = styled.input`
  border-radius: 8px;
  margin-left: ${(props) => props.$marginLeft || '1rem'};

  padding: 0.6rem .6rem;
  font-size: 1rem;
  font-weight: 500;
  font-family: inherit;
  background-color: white;
  color: #331111;
  border: solid 1px #331111;
  max-width: 10rem;

  cursor: pointer;
  transition: border-color 0.25s;
  &:disabled {
    color: #bbcccc;
    border: solid 1px #ffddcc;
    cursor: not-allowed;
  }

  &:focus,
  &:focus-visible {
    outline: 4px auto -webkit-focus-ring-color;
  }
}
`