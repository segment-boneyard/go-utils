package utils

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/segmentio/pointer"
	"github.com/stretchr/testify/assert"
)

type Example struct {
	ID            string
	Success       bool
	Count1        int `json:"count1,omitempty"`
	Count2        int `json:"count2,omitempty"`
	Count3        int `json:"count3,omitempty"`
	CountPointer1 *int
	CountPointer2 *int `json:"count_pointer_2,omitempty"`
}

// func TestDateFormat(t *testing.T) {
// 	now := time.Now()
// 	midnight_before_now := MidnightBeforeOrEqual(now)
//
// 	t.Log("date is", midnight_before_now.Format("yyyy-mm-dd"))
// 	assert.Fail(t, "DDDD", "")
// }

func TestYesterday(t *testing.T) {
	now := time.Now()
	midnight_before_now := MidnightBeforeOrEqual(now)
	midnight_before_now_idempotent := MidnightBeforeOrEqual(midnight_before_now)
	assert.Equal(t, 0, midnight_before_now.Hour())

	assert.Equal(t, now.UTC().Day(), midnight_before_now.Day(), "Should still be same day")
	assert.Equal(t, midnight_before_now, midnight_before_now_idempotent,
		"Multiple calls to MidnightBeforeOrEqual should not change timestamp")

	t.Log("now", now, now.Day())
	t.Log("midnight_before_now", midnight_before_now, midnight_before_now.Day())
	t.Log("now_truncate2", midnight_before_now_idempotent)
}

func TestLast7(t *testing.T) {
	midnight := MidnightBeforeOrEqual(time.Now())
	yesterdayMidnight := DaysBefore(midnight, 1)
	weekAgoMidnight := DaysBefore(midnight, 7)

	t.Log("midnight", midnight)
	t.Log("yesterdayMidnight", yesterdayMidnight)

	assert.Equal(t, int64(7*24*60*60), midnight.Unix()-weekAgoMidnight.Unix())
	assert.Equal(t, int64(1*24*60*60), midnight.Unix()-yesterdayMidnight.Unix())
}

func TestLast30(t *testing.T) {

}

func TestTimeDelta(t *testing.T) {
	start := time.Now()
	end := start.Add(15 * time.Second)
	delta := end.Sub(start)
	t.Log("start", start)
	t.Log("end", end)
	t.Log("delta", delta)
	assert.Equal(t, float64(15), delta.Seconds())
}

func TestStringFormatting(t *testing.T) {
	percent := float32(22) / float32(177) * 100
	str := fmt.Sprintf("Percent: %.1f%%", percent)
	assert.Equal(t, "Percent: 12.4%", str)

	dur := time.Duration(Round(10.871874065, 1) * float64(time.Second))
	val := fmt.Sprintf("%v", dur)
	assert.Equal(t, "10.9s", val)
}

func TestFloatComparison(t *testing.T) {
	val := map[string]interface{}{
		"count":   float32(33),
		"count64": float64(99),
	}
	assert.Equal(t, float32(33), val["count"])
	assert.NotEqual(t, float32(99), val["count64"])
}

func TestJsonOmitEmptyFalsyValues(t *testing.T) {
	val := JsonStringIndent(Example{
		Success: true,
		Count1:  5,
		Count2:  0,
	})
	assert.Contains(t, val, "ID")
	assert.Contains(t, val, "Success")
	assert.Contains(t, val, "count1", "Expect count1 to be included")
	assert.NotContains(t, val, "count2", "Expect count2 to be omitted")
	assert.NotContains(t, val, "count3", "Expect count3 to be omitted")

	var valMap map[string]interface{}
	err := json.Unmarshal([]byte(val), &valMap)
	assert.Nil(t, err, "Got unexpected error")

	countPointer1, countPointer1Exists := valMap["CountPointer1"]
	assert.Nil(t, countPointer1, "CountPointer1 should be nil")
	assert.True(t, countPointer1Exists, "CountPointer1 should exist and be null")
	assert.Contains(t, valMap, "CountPointer1", "Should contain CountPointer1")

	countPointer2, countPointer2Exists := valMap["CountPointer2"]
	assert.Nil(t, countPointer2, "CountPointer2 should be nil when accessed")
	assert.False(t, countPointer2Exists, "CountPointer2 should have been omitted")
	// t.Log(val)
}

func TestJsonMarshallPointerValues(t *testing.T) {
	val := JsonStringIndent(Example{
		CountPointer1: pointer.Int(55),
		CountPointer2: pointer.Int(66),
	})
	// t.Log(val)
	assert.Contains(t, val, "count_pointer_2")

	var valMap map[string]interface{}
	err := json.Unmarshal([]byte(val), &valMap)
	// t.Logf("%#v", valMap)

	assert.Nil(t, err, "Got unexpected error")
	assert.Equal(t, float64(55), valMap["CountPointer1"])
	assert.Equal(t, 66, int(valMap["count_pointer_2"].(float64)))
}
