package service

import (
	"context"
	"fmt"
	"time"

	"github.com/finalops/analytics-service/internal/client"
	"github.com/finalops/analytics-service/internal/model"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type AnalyticsService struct {
	workoutClient *client.WorkoutClient
	redis         *redis.Client
	cacheTTL      time.Duration
}

func NewAnalyticsService(wc *client.WorkoutClient, rdb *redis.Client, cacheTTL time.Duration) *AnalyticsService {
	return &AnalyticsService{workoutClient: wc, redis: rdb, cacheTTL: cacheTTL}
}

func (s *AnalyticsService) GetCalendar(ctx context.Context, userID uuid.UUID, year, month int) ([]model.CalendarDay, error) {
	cacheKey := fmt.Sprintf("analytics:calendar:%s:%d:%02d", userID, year, month)

	if cached, err := s.redis.Get(ctx, cacheKey).Bytes(); err == nil {
		var days []model.CalendarDay
		if err := jsonUnmarshal(cached, &days); err == nil {
			return days, nil
		}
	}

	start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, -1)
	from := start.Format("2006-01-02")
	to := end.Format("2006-01-02")

	sessions, err := s.workoutClient.GetSessionsInRange(ctx, userID, from, to)
	if err != nil {
		return nil, err
	}

	sessionMap := make(map[string]model.SessionRangeItem)
	for _, s := range sessions {
		sessionMap[s.SessionDate] = s
	}

	var days []model.CalendarDay
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format("2006-01-02")
		day := model.CalendarDay{Date: dateStr}
		if sess, ok := sessionMap[dateStr]; ok {
			day.HasSession = true
			day.SessionID = sess.ID
			day.DurationMinutes = sess.DurationMinutes
			day.CaloriesBurned = sess.CaloriesBurned
		}
		days = append(days, day)
	}

	if b, err := jsonMarshal(days); err == nil {
		_ = s.redis.Set(ctx, cacheKey, b, s.cacheTTL).Err()
	}
	return days, nil
}

func (s *AnalyticsService) GetProgression(ctx context.Context, userID, exerciseID uuid.UUID) ([]model.ProgressionPoint, error) {
	cacheKey := fmt.Sprintf("analytics:progression:%s:%s", userID, exerciseID)

	if cached, err := s.redis.Get(ctx, cacheKey).Bytes(); err == nil {
		var points []model.ProgressionPoint
		if err := jsonUnmarshal(cached, &points); err == nil {
			return points, nil
		}
	}

	history, err := s.workoutClient.GetExerciseHistory(ctx, userID, exerciseID)
	if err != nil {
		return nil, err
	}

	byDate := make(map[string][]model.ExerciseSetHistory)
	for _, h := range history {
		byDate[h.SessionDate] = append(byDate[h.SessionDate], h)
	}

	var points []model.ProgressionPoint
	for date, sets := range byDate {
		point := model.ProgressionPoint{Date: date}
		var totalVolume float64
		var rpeSum float64
		var rpeCount int

		for _, set := range sets {
			if set.WeightKg != nil && set.Reps != nil {
				vol := *set.WeightKg * float64(*set.Reps)
				totalVolume += vol
				if point.MaxWeight == nil || *set.WeightKg > *point.MaxWeight {
					w := *set.WeightKg
					point.MaxWeight = &w
				}
			}
			if set.Reps != nil {
				if point.MaxReps == nil || *set.Reps > *point.MaxReps {
					r := *set.Reps
					point.MaxReps = &r
				}
			}
			if set.RPE != nil {
				rpeSum += *set.RPE
				rpeCount++
			}
		}
		point.TotalVolume = totalVolume
		if rpeCount > 0 {
			avg := rpeSum / float64(rpeCount)
			point.AvgRPE = &avg
		}
		points = append(points, point)
	}

	sortProgressionByDate(points)

	if b, err := jsonMarshal(points); err == nil {
		_ = s.redis.Set(ctx, cacheKey, b, s.cacheTTL*2).Err()
	}
	return points, nil
}

func (s *AnalyticsService) GetWeeklyStats(ctx context.Context, userID uuid.UUID) (*model.WeeklyStats, error) {
	now := time.Now().UTC()
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	weekStart := now.AddDate(0, 0, -(weekday - 1))
	weekEnd := weekStart.AddDate(0, 0, 6)

	sessions, err := s.workoutClient.GetSessionsInRange(ctx, userID,
		weekStart.Format("2006-01-02"), weekEnd.Format("2006-01-02"))
	if err != nil {
		return nil, err
	}

	stats := &model.WeeklyStats{
		WeekStart: weekStart.Format("2006-01-02"),
		WeekEnd:   weekEnd.Format("2006-01-02"),
	}
	for _, sess := range sessions {
		stats.TotalSessions++
		if sess.DurationMinutes != nil {
			stats.TotalDuration += int(*sess.DurationMinutes)
		}
		if sess.CaloriesBurned != nil {
			stats.TotalCalories += int(*sess.CaloriesBurned)
		}
	}
	return stats, nil
}

func (s *AnalyticsService) GetMonthlyStats(ctx context.Context, userID uuid.UUID, year, month int) (*model.MonthlyStats, error) {
	cacheKey := fmt.Sprintf("analytics:stats:monthly:%s:%d:%02d", userID, year, month)

	if cached, err := s.redis.Get(ctx, cacheKey).Bytes(); err == nil {
		var stats model.MonthlyStats
		if err := jsonUnmarshal(cached, &stats); err == nil {
			return &stats, nil
		}
	}

	start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, -1)

	sessions, err := s.workoutClient.GetSessionsInRange(ctx, userID,
		start.Format("2006-01-02"), end.Format("2006-01-02"))
	if err != nil {
		return nil, err
	}

	stats := &model.MonthlyStats{Year: year, Month: month}
	for _, sess := range sessions {
		stats.TotalSessions++
		stats.TrainingDays++
		if sess.DurationMinutes != nil {
			stats.TotalDuration += int(*sess.DurationMinutes)
		}
		if sess.CaloriesBurned != nil {
			stats.TotalCalories += int(*sess.CaloriesBurned)
		}
	}

	if b, err := jsonMarshal(stats); err == nil {
		_ = s.redis.Set(ctx, cacheKey, b, s.cacheTTL).Err()
	}
	return stats, nil
}
