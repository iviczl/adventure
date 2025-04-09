import styled, { createGlobalStyle } from 'styled-components'

export const devices = {
  mobile: `(max-width: 1023px)`,
  desktop: `(min-width: 1024px)`,
}

const GlobalStyle = createGlobalStyle`
*,
*:before,
*:after {
  box-sizing: inherit;
  color: #331111;
  font-family: Inter, system-ui, Avenir, Helvetica, Arial, sans-serif;
  line-height: 1.5;
  font-weight: 400;
  
  font-synthesis: none;
  text-rendering: optimizeLegibility;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  padding: auto;
}

#root {
  box-sizing: border-box;
}
`
export default GlobalStyle

export const MainHeading = styled.h1`
  font-size: clamp(1.5rem, 6vw, 5rem);
  margin-bottom: 1rem;
  margin-top: 1rem;
`
export const SubHeading = styled.h5`
  font-size: 1rem;
  font-weight: 500;
  padding: 0;
  margin-bottom: 0;
  margin-top: 0;
`
export const MainContainer = styled.div`
  box-sizing: inherit;
  margin: auto;
  background-color: #cc9977;
  display: flex;
  place-items: center;
  flex-flow: column;
  min-width: 320px;
  max-width: 600px;
  height: 100%;
  max-height: 100%
  border: solid 1px #331111;
  border-radius: clamp(5px, 20px, 25px);
  padding: 1rem;
`
export const Container = styled.div`
  box-sizing: inherit;
  margin-left: auto;
  margin-right: auto;
  padding: 2rem;
  display: flex;
  flex-flow: column;
  place-items: center;
  min-width: 320px;
  max-width: 800px;
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
  margin-left: ${(props) => props.$marginLeft || '1rem'};
  font-weight: 500;
  font-family: inherit;
  line-height: 0.8rem; //calc(2.2rem + 2px);
  padding: 0.5rem 0;
  width: ${(props) => props.$width || 'unset'};
`
export const Option = styled.option`
  padding: 0.6rem 0.6rem 0.6rem 0.6rem;
  font-size: inherit;
  font-weight: 500;
  font-family: inherit;
  background-color: white;
  color: #331111;
  height: 3rem;
`
export const Select = styled.select`
  border-radius: 8px;
  margin-left: ${(props) => props.$marginLeft || '1rem'};
  padding: 0.6rem .6rem .6rem .6rem;
  font-size: inherit;
  font-weight: 500;
  font-family: inherit;
  background-color: white;
  color: #331111;
  border: solid 1px #331111;
  height: 3rem;
  max-width: 12rem;
  width: ${(props) => props.$width || 'unset'};
  // appearance: none;

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
  max-width: 12rem;
  width: ${(props) => props.$width || 'unset'};
  height: 3rem;

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
