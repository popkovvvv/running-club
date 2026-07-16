import { useEffect, useState } from 'react'
import { ActivityMap } from '../activity/ActivityMap'
import { ActivityTabs } from '../activity/ActivityTabs'
import { LeafletMap } from '../activity/LeafletMap'
import { MetricGrid } from '../activity/MetricGrid'
import { SplitTable, StreamChart, computeSplits, parseStreamValues } from '../activity/StreamChart'
import { api, type ActivityDetail, type ActivityStream } from '../../lib/api'
import type { Theme } from '../../lib/theme'

export function WorkoutFactPanel({ theme, activity }: { theme: Theme; activity: ActivityDetail }) {
  const [streams, setStreams] = useState<ActivityStream[]>([])
  const [tab, setTab] = useState('overview')

  useEffect(() => {
    void api.activityStreams(activity.id).then(setStreams).catch(() => setStreams([]))
  }, [activity.id])

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
        <div style={{ fontWeight: 700, fontSize: 16 }}>{activity.title}</div>
        <div style={{ fontSize: 11, color: theme.dim }}>{activity.when}{activity.source ? ` · ${activity.source}` : ''}</div>
      </div>
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
    </div>
  )
}
