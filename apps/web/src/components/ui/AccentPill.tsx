import type { CSSProperties, ReactNode } from 'react'

export function AccentPill({
  children,
  className,
  style,
}: {
  children: ReactNode
  className?: string
  style?: CSSProperties
}) {
  return (
    <span className={className ? `accent-pill ${className}` : 'accent-pill'} style={style}>
      {children}
    </span>
  )
}
