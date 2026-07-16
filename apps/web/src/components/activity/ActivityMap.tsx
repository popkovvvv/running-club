import type { Activity, ActivityDetail } from '../../lib/api'

export function ActivityMap({ activity, accent, height = 160 }: { activity: Activity | ActivityDetail; accent: string; height?: number }) {
  const route = activity.route || ''
  if (!route) {
    return (
      <div style={{ height, borderRadius: 14, background: 'linear-gradient(160deg,#1a2330,#0f141c)', display: 'flex', alignItems: 'center', justifyContent: 'center', color: '#9aa2ab', fontSize: 12 }}>
        Нет маршрута
      </div>
    )
  }
  return (
    <div style={{ height, borderRadius: 14, overflow: 'hidden', background: 'linear-gradient(160deg,#1a2330,#0f141c)', position: 'relative' }}>
      <svg viewBox="0 0 300 140" style={{ width: '100%', height: '100%' }}>
        <path d={route} fill="none" stroke={accent} strokeWidth="4.5" strokeLinecap="round" strokeLinejoin="round" />
        {activity.sx > 0 && <circle cx={activity.sx} cy={activity.sy} r="5" fill={accent} />}
        {activity.ex > 0 && <circle cx={activity.ex} cy={activity.ey} r="5" fill="#fff" stroke={accent} strokeWidth="2" />}
      </svg>
    </div>
  )
}
