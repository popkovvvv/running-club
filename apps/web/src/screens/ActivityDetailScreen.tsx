import { useEffect, useState } from 'react'
import { api, type ActivityDetail, type ActivityStream } from '../lib/api'
import { useApp } from '../lib/store'
import { ActivityMap } from '../components/activity/ActivityMap'
import { ActivityTabs } from '../components/activity/ActivityTabs'
import { ActivityEditForm } from '../components/activity/ActivityEditForm'
import { MetricGrid } from '../components/activity/MetricGrid'
import { SplitTable, StreamChart, computeSplits, parseStreamValues } from '../components/activity/StreamChart'
import { LeafletMap } from '../components/activity/LeafletMap'

export function ActivityDetailScreen({ id }: { id: string }) {
  const { theme, closeOverlay, reloadActivities } = useApp()
  const [activity, setActivity] = useState<ActivityDetail | null>(null)
  const [streams, setStreams] = useState<ActivityStream[]>([])
  const [tab, setTab] = useState('overview')
  const [editing, setEditing] = useState(false)
  const [phoneEl, setPhoneEl] = useState<HTMLElement | null>(null)

  useEffect(() => {
    void api.activity(id).then(setActivity).catch(() => setActivity(null))
    void api.activityStreams(id).then(setStreams).catch(() => setStreams([]))
  }, [id])

  useEffect(() => {
    setPhoneEl(document.querySelector('.phone'))
  }, [])

  if (!activity) {
    return <div className="card" style={{ color: theme.dim }}>Загрузка…</div>
  }

  const hrStream = streams.find((s) => s.type === 'heartrate')
  const altStream = streams.find((s) => s.type === 'altitude')
  const cadenceStream = streams.find((s) => s.type === 'cadence')
  const latlngStream = streams.find((s) => s.type === 'latlng')
  const timeStream = streams.find((s) => s.type === 'time')
  const splits = latlngStream && timeStream ? computeSplits(latlngStream.data, timeStream.data) : []

  const tabs = [
    { id: 'overview', label: 'Обзор' },
    { id: 'stats', label: 'Статистика' },
    { id: 'hr', label: 'Пульс' },
    { id: 'elevation', label: 'Высота' },
    ...(cadenceStream ? [{ id: 'cadence', label: 'Каденс' }] : []),
  ]

  return (
    <div className="fade" style={{ display: 'flex', flexDirection: 'column', gap: 12 }}>
      <button className="btn" onClick={closeOverlay} style={{ alignSelf: 'flex-start', background: theme.card2, color: theme.dim, borderRadius: 10, padding: '8px 12px' }}>← Назад</button>
      <div className="card">
        <div style={{ fontWeight: 700, fontSize: 18 }}>{activity.title}</div>
        <div style={{ fontSize: 11, color: theme.dim }}>{activity.when}{activity.source ? ` · ${activity.source}` : ''}</div>
      </div>
      <button
        data-testid="edit-activity"
        className="btn"
        onClick={() => setEditing(true)}
        style={{ background: theme.card2, color: theme.accent, borderRadius: 12, padding: 12 }}
      >
        Редактировать
      </button>
      <ActivityTabs theme={theme} tab={tab} tabs={tabs} onChange={setTab} />

      {tab === 'overview' && (
        <>
          {activity.polyline ? <LeafletMap polyline={activity.polyline} accent={theme.accent} /> : <ActivityMap activity={activity} accent={theme.accent} height={180} />}
          <MetricGrid theme={theme} items={[
            { label: 'км', value: activity.dist },
            { label: 'время', value: activity.time },
            { label: 'темп', value: activity.pace },
            { label: 'пульс', value: activity.hr || '—' },
            { label: 'набор', value: activity.elevation ? `${Math.round(activity.elevation)} м` : '—' },
            { label: 'kudos', value: String(activity.kudos) },
          ]} />
        </>
      )}

      {tab === 'stats' && (
        <div className="card">
          <div style={{ fontWeight: 700, marginBottom: 10 }}>Сплиты</div>
          <SplitTable splits={splits} />
          {activity.movingSeconds != null && (
            <div style={{ marginTop: 12, fontSize: 12, color: theme.dim }}>
              Moving: {Math.floor(activity.movingSeconds / 60)}:{String(activity.movingSeconds % 60).padStart(2, '0')}
              {activity.elapsedSeconds ? ` · Elapsed: ${Math.floor(activity.elapsedSeconds / 60)}:${String(activity.elapsedSeconds % 60).padStart(2, '0')}` : ''}
            </div>
          )}
        </div>
      )}

      {tab === 'hr' && (
        <div className="card">
          <StreamChart values={hrStream ? parseStreamValues(hrStream.data) : []} color={theme.accent} label={`Средний ${activity.hr}${activity.maxHeartrate ? ` · Max ${activity.maxHeartrate}` : ''}`} />
        </div>
      )}

      {tab === 'elevation' && (
        <div className="card">
          <StreamChart values={altStream ? parseStreamValues(altStream.data) : []} color="#6ecbff" label={`Набор ${activity.elevation ? Math.round(activity.elevation) : 0} м`} />
        </div>
      )}

      {tab === 'cadence' && (
        <div className="card">
          <StreamChart values={cadenceStream ? parseStreamValues(cadenceStream.data) : []} color="#b4ff6e" label="Каденс" />
        </div>
      )}

      {editing && phoneEl && (
        <ActivityEditForm
          theme={theme}
          activity={activity}
          phoneEl={phoneEl}
          onClose={() => setEditing(false)}
          onSave={async (body) => {
            const updated = await api.updateActivity(id, body)
            setActivity(updated)
            await reloadActivities()
          }}
        />
      )}
    </div>
  )
}
