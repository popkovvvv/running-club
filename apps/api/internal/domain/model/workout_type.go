package model

type WorkoutType string

const (
	WorkoutTypeEasy      WorkoutType = "easy"
	WorkoutTypeLong      WorkoutType = "long"
	WorkoutTypeTempo     WorkoutType = "tempo"
	WorkoutTypeInterval  WorkoutType = "interval"
	WorkoutTypeFartlek   WorkoutType = "fartlek"
	WorkoutTypeRecovery  WorkoutType = "recovery"
	WorkoutTypeHills     WorkoutType = "hills"
	WorkoutTypeRace      WorkoutType = "race"
	WorkoutTypeCross     WorkoutType = "cross"
	WorkoutTypeRest      WorkoutType = "rest"
)

type WorkoutStatus string

const (
	WorkoutStatusPlanned   WorkoutStatus = "planned"
	WorkoutStatusCompleted WorkoutStatus = "completed"
	WorkoutStatusSkipped   WorkoutStatus = "skipped"
)

type SegmentKind string

const (
	SegmentWarmup   SegmentKind = "warmup"
	SegmentWork     SegmentKind = "work"
	SegmentInterval SegmentKind = "interval"
	SegmentRecovery SegmentKind = "recovery"
	SegmentCooldown SegmentKind = "cooldown"
	SegmentDrill    SegmentKind = "drill"
	SegmentStride   SegmentKind = "stride"
)

func ValidWorkoutType(v string) bool {
	switch WorkoutType(v) {
	case WorkoutTypeEasy, WorkoutTypeLong, WorkoutTypeTempo, WorkoutTypeInterval,
		WorkoutTypeFartlek, WorkoutTypeRecovery, WorkoutTypeHills, WorkoutTypeRace,
		WorkoutTypeCross, WorkoutTypeRest:
		return true
	default:
		return false
	}
}
