import { useEffect, useState } from 'react'
import { ActivityMap } from '../components/activity/ActivityMap'
import { AccentPill } from '../components/ui/AccentPill'
import { HeroPanel } from '../components/ui/HeroPanel'
import { api } from '../lib/api'
import { nearestAnnounce } from '../lib/nearestAnnounce'
import { useApp } from '../lib/store'

export function HomeAthlete() {
  const { theme, club, announces, activities, setScreen, user, openOverlay, reloadAnnounces, reloadActivities, tabEpoch } = useApp()
  const [weekKm, setWeekKm] = useState('0')
  const [weekPlan, setWeekPlan] = useState('')
  const [weekRange, setWeekRange] = useState('')
  const next = nearestAnnounce(announces)

  useEffect(() => {
    void reloadAnnounces()
    void reloadActivities()
    void api.plan(0).then((p) => {
      setWeekKm(p.weekKm || '0')
      setWeekPlan(p.weekPlan || '')
      setWeekRange(p.weekRange || '')
    }).catch(() => {
      setWeekKm('0')
      setWeekPlan('')
      setWeekRange('')
    })
  }, [tabEpoch, reloadAnnounces, reloadActivities])

  return (
    <div className="fade" style={{ display: 'flex', flexDirection: 'column', gap: 14 }}>
      <HeroPanel>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', gap: 8, marginBottom: 10 }}>
          <div className="display-title" style={{ fontSize: 24 }}>Эта неделя</div>
          {weekRange ? <AccentPill style={{ fontSize: 10, padding: '4px 10px' }}>{weekRange}</AccentPill> : null}
        </div>
        <div style={{ fontFamily: theme.display, fontSize: 36, fontWeight: 800, fontStyle: 'italic', lineHeight: 1 }}>
          {weekKm}
          {weekPlan ? <span style={{ fontSize: 16, color: theme.dim, fontStyle: 'normal' }}> / {weekPlan}</span> : <span style={{ fontSize: 16, color: theme.dim, fontStyle: 'normal' }}> км</span>}
        </div>
        {user?.inClub && next && (
          <div style={{ display: 'flex', flexWrap: 'wrap', gap: 8, alignItems: 'center', marginTop: 14 }}>
            <AccentPill>{next.day}</AccentPill>
            <span style={{ fontSize: 13, fontWeight: 700 }}>{next.time}</span>
            <span style={{ fontSize: 12, color: theme.dim }}>{next.place}</span>
          </div>
        )}
      </HeroPanel>
      {user?.inClub && (
        <div className="card" data-testid="group-preview" onClick={() => setScreen('schedule')} style={{ cursor: 'pointer' }}>
          <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 12, alignItems: 'center', gap: 8 }}>
            <AccentPill style={{ fontSize: 10, padding: '4px 10px' }}>Групповые</AccentPill>
            <span style={{ fontSize: 12, fontWeight: 700, color: theme.accent }}>Расписание →</span>
          </div>
          {next ? (
            <>
              <div style={{ fontWeight: 700, fontSize: 14 }}>Ближайшая: {next.place}</div>
              <div style={{ fontSize: 11, color: theme.dim, marginTop: 4 }}>{next.day} · старт {next.time} · идут {next.going}</div>
            </>
          ) : (
            <div style={{ fontSize: 13, color: theme.dim }}>Пока нет анонсов — откройте расписание</div>
          )}
        </div>
      )}
      <div className="card" style={{ display: 'flex', alignItems: 'center', gap: 10, borderRadius: 14, padding: '11px 14px' }}>
        <span style={{ width: 9, height: 9, borderRadius: '50%', background: theme.good, boxShadow: `0 0 10px ${theme.good}` }} />
        <div style={{ flex: 1 }}>
          <div style={{ fontSize: 12, fontWeight: 700 }}>Активности</div>
          <div style={{ fontSize: 11, color: theme.dim }}>Клуб {club?.name || '—'}</div>
        </div>
      </div>
      <div className="display-title" style={{ fontSize: 16 }}>Лента</div>
      {activities.length === 0 && (
        <div className="card" style={{ fontSize: 13, color: theme.dim }}>Пока нет активностей</div>
      )}
      {activities.map((a) => (
        <div key={a.id} className="card" data-testid="activity-card" onClick={() => openOverlay({ type: 'activity', id: a.id })} style={{ cursor: 'pointer' }}>
          <div style={{ display: 'flex', justifyContent: 'space-between', gap: 8 }}>
            <div style={{ fontWeight: 700 }}>{a.title}</div>
            {a.source && <AccentPill style={{ fontSize: 9, padding: '3px 8px' }}>{a.source}</AccentPill>}
          </div>
          <div style={{ fontSize: 11, color: theme.dim }}>{a.when}</div>
          <div style={{ display: 'flex', gap: 14, marginTop: 10 }}>
            <div><div style={{ fontFamily: theme.display, fontWeight: 800 }}>{a.dist}</div><div style={{ fontSize: 10, color: theme.dim }}>км</div></div>
            <div><div style={{ fontFamily: theme.display, fontWeight: 800 }}>{a.time}</div><div style={{ fontSize: 10, color: theme.dim }}>время</div></div>
            <div><div style={{ fontFamily: theme.display, fontWeight: 800 }}>{a.pace}</div><div style={{ fontSize: 10, color: theme.dim }}>темп</div></div>
          </div>
          <div style={{ marginTop: 12 }}>
            <ActivityMap activity={a} accent={theme.accent} height={120} />
          </div>
        </div>
      ))}
    </div>
  )
}
