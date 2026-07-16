import { useEffect, useState } from 'react'
import { ActivityMap } from '../activity/ActivityMap'
import { ActivityTabs } from '../activity/ActivityTabs'
import { ActivityEditForm } from '../activity/ActivityEditForm'
import { LeafletMap } from '../activity/LeafletMap'
import { MetricGrid } from '../activity/MetricGrid'
import { SplitTable, StreamChart, computeSplits, parseStreamValues } from '../activity/StreamChart'
import { api, type ActivityDetail, type ActivityStream } from '../../lib/api'
import type { Theme } from '../../lib/theme'

export function WorkoutFactPanel({
  theme,
  activity,
  onUpdated,
}: {
  theme: Theme
  activity: ActivityDetail
  onUpdated?: (a: ActivityDetail) => void
}) {
  const [streams, setStreams] = useState<ActivityStream[]>([])
  const [tab, setTab] = useState('overview')
  const [current, setCurrent] = useState(activity)
  const [editing, setEditing] = useState(false)
  const [phoneEl, setPhoneEl] = useState<HTMLElement | null>(null)

  useEffect(() => {
    setCurrent(activity)
  }, [activity])

  useEffect(() => {
    void api.activityStreams(current.id).then(setStreams).catch(() => setStreams([]))
  }, [current.id])

  useEffect(() => {
    setPhoneEl(document.querySelector('.phone'))
  }, [])

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
    <div style={{ display: 'flex', flexDirection: 'column', gap: 12 }}>
      <div className="card">
        <div style={{ fontWeight: 700, fontSize: 16 }}>{current.title}</div>
        <div style={{ fontSize: 11, color: theme.dim }}>{current.when}{current.source ? ` · ${current.source}` : ''}</div>
      </div>
      <button
        data-testid="edit-workout-fact"
        className="btn"
        onClick={() => setEditing(true)}
        style={{ background: theme.card2, color: theme.accent, borderRadius: 12, padding: 12 }}
      >
        Редактировать
      </button>
      <ActivityTabs theme={theme} tab={tab} tabs={tabs} onChange={setTab} />

      {tab === 'overview' && (
        <>
          {current.polyline ? <LeafletMap polyline={current.polyline} accent={theme.accent} /> : <ActivityMap activity={current} accent={theme.accent} height={180} />}
          <MetricGrid theme={theme} items={[
            { label: 'км', value: current.dist },
            { label: 'время', value: current.time },
            { label: 'темп', value: current.pace },
            { label: 'пульс', value: current.hr || '—' },
            { label: 'набор', value: current.elevation ? `${Math.round(current.elevation)} м` : '—' },
            { label: 'kudos', value: String(current.kudos) },
          ]} />
        </>
      )}

      {tab === 'stats' && (
        <div className="card">
          <div style={{ fontWeight: 700, marginBottom: 10 }}>Сплиты</div>
          <SplitTable splits={splits} />
          {current.movingSeconds != null && (
            <div style={{ marginTop: 12, fontSize: 12, color: theme.dim }}>
              Moving: {Math.floor(current.movingSeconds / 60)}:{String(current.movingSeconds % 60).padStart(2, '0')}
              {current.elapsedSeconds ? ` · Elapsed: ${Math.floor(current.elapsedSeconds / 60)}:${String(current.elapsedSeconds % 60).padStart(2, '0')}` : ''}
            </div>
          )}
        </div>
      )}

      {tab === 'hr' && (
        <div className="card">
          <StreamChart values={hrStream ? parseStreamValues(hrStream.data) : []} color={theme.accent} label={`Средний ${current.hr}${current.maxHeartrate ? ` · Max ${current.maxHeartrate}` : ''}`} />
        </div>
      )}

      {tab === 'elevation' && (
        <div className="card">
          <StreamChart values={altStream ? parseStreamValues(altStream.data) : []} color="#6ecbff" label={`Набор ${current.elevation ? Math.round(current.elevation) : 0} м`} />
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
          activity={current}
          phoneEl={phoneEl}
          onClose={() => setEditing(false)}
          onSave={async (body) => {
            const updated = await api.updateActivity(current.id, body)
            setCurrent(updated)
            onUpdated?.(updated)
          }}
        />
      )}
    </div>
  )
}
