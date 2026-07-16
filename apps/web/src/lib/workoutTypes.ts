export type WorkoutType =
  | 'easy'
  | 'long'
  | 'tempo'
  | 'interval'
  | 'fartlek'
  | 'recovery'
  | 'hills'
  | 'race'
  | 'cross'
  | 'rest'

export type SegmentKind =
  | 'warmup'
  | 'work'
  | 'interval'
  | 'recovery'
  | 'cooldown'
  | 'drill'
  | 'stride'

export const WORKOUT_TYPES: { id: WorkoutType; label: string }[] = [
  { id: 'easy', label: 'Лёгкий' },
  { id: 'long', label: 'Длинный' },
  { id: 'tempo', label: 'Темповый' },
  { id: 'interval', label: 'Интервалы' },
  { id: 'fartlek', label: 'Фартlek' },
  { id: 'recovery', label: 'Восстановление' },
  { id: 'hills', label: 'Горки' },
  { id: 'race', label: 'Старт' },
  { id: 'cross', label: 'Кросс' },
  { id: 'rest', label: 'Отдых' },
]

export const SEGMENT_KINDS: { id: SegmentKind; label: string }[] = [
  { id: 'warmup', label: 'Разминка' },
  { id: 'work', label: 'Основная' },
  { id: 'interval', label: 'Интервал' },
  { id: 'recovery', label: 'Восстановление' },
  { id: 'cooldown', label: 'Заминка' },
  { id: 'drill', label: 'Упражнения' },
  { id: 'stride', label: 'Ускорения' },
]

export function workoutTypeLabel(id: string) {
  return WORKOUT_TYPES.find((t) => t.id === id)?.label || id
}

export function segmentKindLabel(id: string) {
  return SEGMENT_KINDS.find((k) => k.id === id)?.label || id
}

export const WORKOUT_PRESETS: Record<WorkoutType, { kind: SegmentKind; title: string; distKm: number; pace: string }[]> = {
  easy: [
    { kind: 'warmup', title: 'Лёгкий бег', distKm: 2, pace: '7:40' },
    { kind: 'work', title: 'Основной кросс', distKm: 4, pace: '7:20' },
    { kind: 'cooldown', title: 'Заминка', distKm: 1, pace: '8:00' },
  ],
  long: [
    { kind: 'warmup', title: 'Разминка', distKm: 2, pace: '7:30' },
    { kind: 'work', title: 'Длительный', distKm: 12, pace: '7:00' },
    { kind: 'cooldown', title: 'Заминка', distKm: 1, pace: '8:00' },
  ],
  tempo: [
    { kind: 'warmup', title: 'Разминка', distKm: 2, pace: '7:30' },
    { kind: 'work', title: 'Темповый блок', distKm: 6, pace: '5:30' },
    { kind: 'cooldown', title: 'Заминка', distKm: 1.5, pace: '8:00' },
  ],
  interval: [
    { kind: 'warmup', title: 'Разминка', distKm: 2, pace: '7:30' },
    { kind: 'interval', title: '5×800', distKm: 4, pace: '4:30' },
    { kind: 'recovery', title: 'Трусца', distKm: 2, pace: '7:00' },
    { kind: 'cooldown', title: 'Заминка', distKm: 1, pace: '8:00' },
  ],
  fartlek: [
    { kind: 'warmup', title: 'Разминка', distKm: 2, pace: '7:30' },
    { kind: 'work', title: 'Фартlek', distKm: 5, pace: '6:00' },
    { kind: 'cooldown', title: 'Заминка', distKm: 1, pace: '8:00' },
  ],
  recovery: [{ kind: 'work', title: 'Восстановительный', distKm: 4, pace: '8:00' }],
  hills: [
    { kind: 'warmup', title: 'Разминка', distKm: 2, pace: '7:30' },
    { kind: 'work', title: 'Горки 8×200', distKm: 3, pace: '5:00' },
    { kind: 'cooldown', title: 'Заминка', distKm: 1, pace: '8:00' },
  ],
  race: [{ kind: 'work', title: 'Стартовый темп', distKm: 10, pace: '5:00' }],
  cross: [{ kind: 'drill', title: 'ОФП', distKm: 0, pace: '—' }],
  rest: [],
}
