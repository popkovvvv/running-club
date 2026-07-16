import type { CSSProperties, ReactNode } from 'react'

export function HeroPanel({
  children,
  imageUrl,
  className,
  style,
}: {
  children: ReactNode
  imageUrl?: string
  className?: string
  style?: CSSProperties
}) {
  const bgStyle: CSSProperties | undefined = imageUrl
    ? { ...style, ['--hero-image' as string]: `url(${imageUrl})` }
    : style

  return (
    <div className={className ? `hero-panel ${className}` : 'hero-panel'} style={bgStyle}>
      {children}
    </div>
  )
}
