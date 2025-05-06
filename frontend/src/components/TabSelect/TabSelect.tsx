import React from 'react'
import { TabContainer } from './tabSelectStyles'

export default function TabSelect({ children }: { children: React.ReactNode }) {
  return (
    <TabContainer>
      {React.Children.map(children, (child) => {
        if (React.isValidElement(child)) {
          return React.cloneElement(child)
        }
        return child
      })}
    </TabContainer>
  )
}
