package seed

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
	"github.com/nikpopkov/running-club/api/internal/pkg/password"
)

func Run(ctx context.Context, pool *pgxpool.Pool) error {
	hash, err := password.Hash("password")
	if err != nil {
		return fmt.Errorf("password.Hash: %w", err)
	}
	now := time.Now().UTC()
	coachID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	athleteID := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	clubID := uuid.MustParse("33333333-3333-3333-3333-333333333333")

	_, err = pool.Exec(ctx, `
		INSERT INTO users (id,name,email,password_hash,role,created_at) VALUES
			($1,'Главный тренер','coach@pulse.run',$2,'coach',$3),
			($4,'Никита Попков','nikita@pulse.run',$2,'athlete',$3)
		ON CONFLICT (email) DO UPDATE SET password_hash=EXCLUDED.password_hash, name=EXCLUDED.name`,
		coachID, hash, now, athleteID)
	if err != nil {
		return fmt.Errorf("insert users: %w", err)
	}

	var n int
	if err := pool.QueryRow(ctx, `SELECT COUNT(*) FROM clubs`).Scan(&n); err != nil {
		return fmt.Errorf("count clubs: %w", err)
	}
	if n > 0 {
		_, _ = pool.Exec(ctx, `UPDATE memberships SET status='active', updated_at=NOW() WHERE user_id=$1 AND club_id=$2`, athleteID, clubID)
		return nil
	}
	_, err = pool.Exec(ctx, `INSERT INTO clubs (id,name,invite_code,accent_hex,coach_id,created_at) VALUES ($1,'PULSE','PULSE-7K42','#ff5c22',$2,$3)`,
		clubID, coachID, now)
	if err != nil {
		return fmt.Errorf("insert club: %w", err)
	}
	_, err = pool.Exec(ctx, `INSERT INTO memberships (id,user_id,club_id,status,created_at,updated_at) VALUES ($1,$2,$3,'active',$4,$4)`,
		uuid.New(), athleteID, clubID, now)
	if err != nil {
		return fmt.Errorf("insert membership: %w", err)
	}

	july21 := time.Date(2026, 7, 21, 0, 0, 0, 0, time.UTC)
	july23 := time.Date(2026, 7, 23, 0, 0, 0, 0, time.UTC)
	announces := []*model.Announce{
		model.NewAnnounce(clubID, "Стадион «Зина»", "Вт, 21 июля", "19:50", "Основная группа", "ОФП + темповая работа. Приходите на разминку к 19:30.", &july21),
		model.NewAnnounce(clubID, "ЛЭМЗ", "Чт, 23 июля", "19:50", "Основная группа", "Интервалы 5×800 м через 400 м трусцой. Форма по погоде.", &july23),
	}
	announces[0].GoingCount = 8
	announces[1].GoingCount = 6
	for _, a := range announces {
		_, err = pool.Exec(ctx, `INSERT INTO announces (id,club_id,place,day_label,time,group_name,note,starts_on,going_count,created_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
			a.ID, a.ClubID, a.Place, a.DayLabel, a.Time, a.GroupName, a.Note, a.StartsOn, a.GoingCount, a.CreatedAt)
		if err != nil {
			return fmt.Errorf("insert announce: %w", err)
		}
	}

	planDays := []struct {
		day, tag, title, dist, dur string
		km                         float64
	}{
		{"Пн", "Восстановление", "Зарядка + лёгкий бег", "5 км", "~35 мин", 5},
		{"Вт", "Интервалы", "Стадион «Зина» · отрезки", "8 км", "~1:05", 8},
		{"Ср", "Кросс", "Восстановительный кросс", "6 км", "~46 мин", 6},
		{"Чт", "Темповая", "ЛЭМЗ · темповый бег", "7 км", "~50 мин", 7},
		{"Пт", "Отдых", "Восстановление", "—", "—", 0},
		{"Сб", "Длительный", "Длительный кросс", "8 км", "~1:00", 8},
		{"Вс", "Отдых", "Полный отдых", "—", "—", 0},
	}
	for _, d := range planDays {
		_, err = pool.Exec(ctx, `INSERT INTO workouts (id,club_id,user_id,kind,day_label,tag,title,dist_km,duration,pace,hr,week_index,created_at)
			VALUES ($1,$2,$3,'plan',$4,$5,$6,$7,$8,'','',1,$9)`,
			uuid.New(), clubID, athleteID, d.day, d.tag, d.title, d.km, d.dur, now)
		if err != nil {
			return fmt.Errorf("insert workout: %w", err)
		}
	}

	activities := []*model.Activity{
		model.NewActivity(athleteID, "Вечерний кросс", "сегодня, 20:14", 8.2, "1:00:14", "7:20", 143, 12, 3, "M18 104 C 60 40, 96 128, 148 74 S 224 34, 284 92", 18, 104, 284, 92),
		model.NewActivity(athleteID, "Групповая · ЛЭМЗ", "вчера, 19:48", 6.1, "46:30", "7:38", 151, 9, 2, "M28 30 C 84 82, 58 122, 132 100 S 214 122, 272 44", 28, 30, 272, 44),
		model.NewActivity(athleteID, "Восстановительный бег", "2 дня назад", 6.0, "47:12", "7:52", 138, 6, 1, "M24 98 C 72 58, 122 112, 162 60 C 202 20, 244 84, 278 50", 24, 98, 278, 50),
	}
	activities[1].CreatedAt = now.Add(-time.Hour)
	activities[2].CreatedAt = now.Add(-2 * time.Hour)
	for _, a := range activities {
		_, err = pool.Exec(ctx, `INSERT INTO activities (id,user_id,title,when_label,dist_km,duration,pace,hr,kudos,comments,route_svg,start_x,start_y,end_x,end_y,created_at)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)`,
			a.ID, a.UserID, a.Title, a.WhenLabel, a.DistKm, a.Duration, a.Pace, a.HR, a.Kudos, a.Comments, a.RouteSVG, a.StartX, a.StartY, a.EndX, a.EndY, a.CreatedAt)
		if err != nil {
			return fmt.Errorf("insert activity: %w", err)
		}
	}

	months := []model.MonthStat{
		model.NewMonthStat("Июль", 108, 18, "7:20", "6.4"),
		model.NewMonthStat("Июнь", 96, 16, "7:28", "6.1"),
		model.NewMonthStat("Май", 88, 15, "7:35", "5.8"),
	}
	for _, m := range months {
		_, err = pool.Exec(ctx, `INSERT INTO month_stats (id,user_id,month,km,tr,pace,diff) VALUES ($1,$2,$3,$4,$5,$6,$7)`,
			uuid.New(), athleteID, m.Month, m.Km, m.Tr, m.Pace, m.Diff)
		if err != nil {
			return fmt.Errorf("insert month: %w", err)
		}
	}

	prs := []*model.PR{
		model.NewPR(athleteID, "1 км", "5:38", "мар 2026"),
		model.NewPR(athleteID, "5 км", "32:40", "апр 2026"),
		model.NewPR(athleteID, "10 км", "1:08:12", "май 2026"),
		model.NewPR(athleteID, "21.1 км", "2:29:40", "апр 2026"),
	}
	for _, p := range prs {
		_, err = pool.Exec(ctx, `INSERT INTO prs (id,user_id,distance,time,date_label,pending) VALUES ($1,$2,$3,$4,$5,false)`,
			p.ID, p.UserID, p.Distance, p.Time, p.DateLabel)
		if err != nil {
			return fmt.Errorf("insert pr: %w", err)
		}
	}

	races := []*model.Race{
		model.NewRace(athleteID, "Забег «Белые ночи»", "19 июл 2026", "10 км", "Цель 1:06", 10),
		model.NewRace(athleteID, "Осенний гром 21K", "13 сен 2026", "21.1 км", "Цель 2:20", 66),
		model.NewRace(athleteID, "Московский марафон", "20 сен 2026", "42.2 км", "Дебют", 73),
	}
	for _, race := range races {
		_, err = pool.Exec(ctx, `INSERT INTO races (id,user_id,name,date_label,dist,goal,days_left,finished,result) VALUES ($1,$2,$3,$4,$5,$6,$7,false,'')`,
			race.ID, race.UserID, race.Name, race.DateLabel, race.Dist, race.Goal, race.DaysLeft)
		if err != nil {
			return fmt.Errorf("insert race: %w", err)
		}
	}
	return nil
}
